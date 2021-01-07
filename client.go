package client

import (
	"strings"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	Client        mqtt.Client
	SubscCh       chan mqtt.Message
	subscTopics   []string
	subscRadiusKm float64
}

func NewClient(c mqtt.Client, subscRadiusKm float64) *Client {
	return &Client{Client: c, subscRadiusKm: subscRadiusKm}
}

func (c *Client) UpdateSubscribe(lat, lng float64, qos byte) {
	circle := capOnEarth(s2.PointFromLatLng(s2.LatLngFromDegrees(lat, lng)), c.subscRadiusKm)
	rc := &s2.RegionCoverer{MaxLevel: 30, MaxCells: 4}
	cells := rc.Covering(circle)
	newTopics := make([]string, len(cells))
	for i, c := range cells {
		newTopics[i] = celID2TopicName(c) + "/#"
	}

	unsubscTopics, subscTopics := extractUnduplicateTopics(c.subscTopics, newTopics)

	var callback mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		c.SubscCh <- msg
	}
	for _, t := range subscTopics {
		if t != "" {
			c.Client.Subscribe(t, qos, callback)
			c.Client.Publish("/api/register", 1, false, t)
		}
	}

	for _, t := range unsubscTopics {
		if t != "" {
			c.Client.Unsubscribe(t)
			c.Client.Publish("/api/unregister", 1, false, t)
		}
	}

	c.subscTopics = newTopics
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

func (c *Client) Publish(lat, lng float64, qos byte, retained bool, payload interface{}) mqtt.Token {
	topic := celID2TopicName(s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng)))
	return c.Client.Publish(topic, qos, retained, payload)
}

func celID2TopicName(id s2.CellID) string {
	idString := strings.Replace(id.String(), "/", "", 1)
	return strings.Replace(idString, "", "/", len(idString))
}

func capOnEarth(center s2.Point, radiusKm float64) s2.Cap {
	const earthRadiusKm = 6371.01
	ratio := (radiusKm / earthRadiusKm)
	return s2.CapFromCenterAngle(center, s1.Angle(ratio))
}
