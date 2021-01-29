package client

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	log "github.com/sirupsen/logrus"
)

type BrokerInfo struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

type GatewayBrokerInfo struct {
	Topic      string     `json:"topic"`
	BrokerInfo BrokerInfo `json:"broker_info"`
}

type Client struct {
	Client             mqtt.Client
	subscTopics        []string
	subscRadiusKm      float64
	publishTopicPrefix string
	managerHost        string
	managerPort        uint16
	gatewaysInfo       []GatewayBrokerInfo // This data structure is not optimized, so may update.
}

// Manager をにいったん接続し、Gateway ブローカの担当エリアと接続情報を取得してから最適な Gateway ブローカに接続する
func Connect(mHost string, mPort uint16, lat, lng float64, subscRadiusKm float64, timeout uint) (*Client, error) {
	log.WithFields(log.Fields{
		"mHost":         mHost,
		"mPort":         mPort,
		"lat":           lat,
		"lng":           lng,
		"subscRadiusKm": subscRadiusKm,
	}).Trace("Connecting managr brokr...")
	managerBroker := fmt.Sprintf("tcp://%v:%v", mHost, mPort)
	managerOpts := mqtt.NewClientOptions()
	managerOpts.AddBroker(managerBroker)
	managerClient := mqtt.NewClient(managerOpts)

	// connect to manager broker
	if token := managerClient.Connect(); token.Wait() && token.Error() != nil {
		log.WithFields(log.Fields{
			"mHost": mHost,
			"mPort": mPort,
			"err":   token.Error(),
		}).Debug("MQTT connection error (manager broker)")
		return nil, token.Error()
	}
	// Do not maintain connection with manager broker
	defer managerClient.Disconnect(100)

	gatewayInfoTopic := "/api/gateway/info/all"
	managerCh := make(chan mqtt.Message)
	var managerCallback mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
		managerCh <- m
	}
	if token := managerClient.Subscribe(gatewayInfoTopic, 2, managerCallback); token.Wait() && token.Error() != nil {
		log.WithFields(log.Fields{
			"mHost": mHost,
			"mPort": mPort,
			"topic": gatewayInfoTopic,
			"err":   token.Error(),
		}).Debug("MQTT subcribe error (manager broker)")
		return nil, token.Error()
	}

	var gatewaysInfo []GatewayBrokerInfo
	for {
		select {
		case m := <-managerCh:
			// JSONデコード
			if err := json.Unmarshal(m.Payload(), &gatewaysInfo); err != nil {
				log.WithFields(log.Fields{
					"data": string(m.Payload()),
					"err":  err,
				}).Debug("JSON decode error (gateway brokers infomation)")
			}
			log.WithFields(log.Fields{
				"data": string(m.Payload()),
			}).Trace("Get manager data (gateway brokers infomation)")
			break

		case <-time.After(time.Millisecond * time.Duration(timeout)):
			return nil, TimeoutError{Msg: "Timeout occured (gateway brokers information)"}
		}
		break
	}

	if len(gatewaysInfo) == 0 {
		log.WithFields(log.Fields{
			"gatewaysInfo": gatewaysInfo,
		}).Debug("Gateway info erro (info size is zero )")
		return nil, GatewayInfoError{Msg: "Gateway info erro (info size is zero )"}
	}

	// Shuffle
	rand.Shuffle(len(gatewaysInfo), func(i, j int) {
		gatewaysInfo[i], gatewaysInfo[j] = gatewaysInfo[j], gatewaysInfo[i]
	})

	// Unstable sorting
	sort.Slice(gatewaysInfo, func(i, j int) bool {
		return len(gatewaysInfo[i].Topic) > len(gatewaysInfo[j].Topic)
	})

	currentTopic := CelID2TopicName(s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng)))
	var gateway GatewayBrokerInfo
	for _, gateway = range gatewaysInfo {
		if strings.HasPrefix(currentTopic, gateway.Topic) {
			break
		}
	}

	log.WithFields(log.Fields{
		"gatewayHost":  gateway.BrokerInfo.Host,
		"gatewayPort":  gateway.BrokerInfo.Port,
		"gatewayTopic": gateway.Topic,
	}).Trace("Connecting gateway brokr...")
	gatewayBroker := fmt.Sprintf("tcp://%v:%v", gateway.BrokerInfo.Host, gateway.BrokerInfo.Port)
	gatewayOpts := mqtt.NewClientOptions()
	gatewayOpts.AddBroker(gatewayBroker)
	gatewayClient := mqtt.NewClient(gatewayOpts)

	// connect to manager broker
	if token := gatewayClient.Connect(); token.Wait() && token.Error() != nil {
		log.WithFields(log.Fields{
			"gatewayHost": gateway.BrokerInfo.Host,
			"gatewayPort": gateway.BrokerInfo.Port,
			"err":         token.Error(),
		}).Debug("MQTT connection error (gateway broker)")
		return nil, token.Error()
	}

	return NewClient(gatewayClient, subscRadiusKm), nil
}

// NewClient function is deprecated.
// Manager を経由せず直接 GateWay ブローカに接続（後方互換のため残している）
func NewClient(c mqtt.Client, subscRadiusKm float64) *Client {
	return &Client{Client: c, subscRadiusKm: subscRadiusKm, publishTopicPrefix: "/forward"}
}

func (c *Client) UpdateSubscribe(lat, lng float64, qos byte, callback mqtt.MessageHandler) error {
	circle := capOnEarth(s2.PointFromLatLng(s2.LatLngFromDegrees(lat, lng)), c.subscRadiusKm)
	rc := &s2.RegionCoverer{MaxLevel: 20, MaxCells: 4}
	cells := rc.Covering(circle)
	newTopics := make([]string, len(cells))
	for i, c := range cells {
		newTopics[i] = CelID2TopicName(c) + "/#"
	}

	unsubscTopics, subscTopics := extractUnduplicateTopics(c.subscTopics, newTopics)

	for _, t := range subscTopics {
		if t != "" {
			if token := c.Client.Subscribe(t, qos, callback); token.Wait() && token.Error() != nil {
				return token.Error()
			}
			if token := c.Client.Publish("/api/register", 0, false, t); token.Wait() && token.Error() != nil {
				return token.Error()
			}
		}
	}
	for _, t := range unsubscTopics {
		if t != "" {
			if token := c.Client.Unsubscribe(t); token.Wait() && token.Error() != nil {
				return token.Error()
			}
			if token := c.Client.Publish("/api/unregister", 0, false, t); token.Wait() && token.Error() != nil {
				return token.Error()
			}
		}
	}
	c.subscTopics = newTopics
	return nil
}

func (c *Client) Unsubscribe() error {
	for i, t := range c.subscTopics {
		if t != "" {
			if token := c.Client.Unsubscribe(t); token.Wait() && token.Error() != nil {
				return token.Error()
			}
			if token := c.Client.Publish("/api/unregister", 0, false, t); token.Wait() && token.Error() != nil {
				return token.Error()
			}
			c.subscTopics[i] = ""
		}
	}
	return nil
}

func extractUnduplicateTopics(currentSubscTopics, newSubscTopics []string) ([]string, []string) {
	unsubscTopics := make([]string, len(currentSubscTopics))
	copy(unsubscTopics, currentSubscTopics)
	subscTopics := make([]string, len(newSubscTopics))
	copy(subscTopics, newSubscTopics)
	for i, ct := range currentSubscTopics {
		for j, nt := range newSubscTopics {
			if ct == nt {
				unsubscTopics[i] = ""
				subscTopics[j] = ""
				break
			}
		}
	}
	return unsubscTopics, subscTopics
}

func (c *Client) Publish(lat, lng float64, qos byte, retained bool, payload interface{}) error {
	topic := c.publishTopicPrefix + CelID2TopicName(s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng)))
	if token := c.Client.Publish(topic, qos, retained, payload); token.WaitTimeout(100) && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func CelID2TopicName(id s2.CellID) string {
	idString := strings.Replace(id.String(), "/", "", 1)
	return strings.Replace(idString, "", "/", len(idString))
}

func capOnEarth(center s2.Point, radiusKm float64) s2.Cap {
	const earthRadiusKm = 6371.01
	ratio := (radiusKm / earthRadiusKm)
	return s2.CapFromCenterAngle(center, s1.Angle(ratio))
}

func TopicName2Token(topic string) (string, error) {
	tmp := strings.Replace(topic, "/", "", -1)
	var token uint64
	switch string(tmp[0]) {
	case "0":
		token = 0b0000000000000000000000000000000000000000000000000000000000000000
	case "1":
		token = 0b0010000000000000000000000000000000000000000000000000000000000000
	case "2":
		token = 0b0100000000000000000000000000000000000000000000000000000000000000
	case "3":
		token = 0b0110000000000000000000000000000000000000000000000000000000000000
	case "4":
		token = 0b1000000000000000000000000000000000000000000000000000000000000000
	case "5":
		token = 0b1010000000000000000000000000000000000000000000000000000000000000
	default:
		return "", TopicNameError{fmt.Sprintf("Invalid topic name (inputed topic name: %v)", topic)}
	}
	maskTail := uint64(0b0001000000000000000000000000000000000000000000000000000000000000)
	masks := [3]uint64{
		0b0000100000000000000000000000000000000000000000000000000000000000,
		0b0001000000000000000000000000000000000000000000000000000000000000,
		0b0001100000000000000000000000000000000000000000000000000000000000,
	}
	for _, v := range tmp[1:] {
		switch string(v) {
		case "0":
			// 何もしない
		case "1":
			token = token | masks[0]
		case "2":
			token = token | masks[1]
		case "3":
			token = token | masks[2]
		default:
			return "", TopicNameError{fmt.Sprintf("Invalid topic name (inputed topic name: %v)", topic)}
		}

		for j := 0; j < 3; j++ {
			masks[j] = masks[j] >> 2
		}
		maskTail = maskTail >> 2
	}
	tokenString := uint2Token(token | maskTail)
	tokenLen := 1
	tokenLen += int(len(tmp) / 2)
	return tokenString[:tokenLen], nil
}

func (c *Client) Disconnect(quiesce uint) {
	c.Client.Disconnect(quiesce)
}

func uint2Token(ui uint64) string {
	token := ""
	mask := uint64(0b1111000000000000000000000000000000000000000000000000000000000000)

	for i := 0; i < 16; i++ {
		tmp := (ui & mask)
		for j := i + 1; j < 16; j++ {
			tmp = tmp >> 4
		}
		token += fmt.Sprintf("%x", tmp)
		mask = mask >> 4
	}
	return token
}

//////////////           以下、エラー 関連                 //////////////
type TopicNameError struct {
	Msg string
}

func (e TopicNameError) Error() string {
	return fmt.Sprintf("Error: %v", e.Msg)
}

type TimeoutError struct {
	Msg string
}

func (e TimeoutError) Error() string {
	return fmt.Sprintf("Error: %v", e.Msg)
}

type GatewayInfoError struct {
	Msg string
}

func (e GatewayInfoError) Error() string {
	return fmt.Sprintf("Error: %v", e.Msg)
}
