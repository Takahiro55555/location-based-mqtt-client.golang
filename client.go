package client

import (
	"strings"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client interface {
}

type client struct {
	client        mqtt.Client
	subscCh       chan<- mqtt.Message
	subscTopic    string
	subscRadiusKm float64
}

func NewClient(c mqtt.Client, subscRadiusKm float64) *client {
	return &client{client: c, subscTopic: "", subscRadiusKm: subscRadiusKm}
}

func (c *client) UpdateSubscribe(lat, lng float64, qos byte) {
	circle := capOnEarth(s2.PointFromLatLng(s2.LatLngFromDegrees(lat, lng)), 4)
	rc := &s2.RegionCoverer{MaxLevel: 30, MaxCells: 1}
	newTopic := celID2TopicName(rc.Covering(circle)[0]) + "/#"

	var callback mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		c.subscCh <- msg
	}
	c.client.Subscribe(newTopic, qos, callback)

	if c.subscTopic != newTopic || c.subscTopic != "" {
		// NOTE: 古い Topic の Unsubscribe と GateWay への Unsibscribe 通知
		c.client.Unsubscribe(c.subscTopic)
		c.client.Publish("/api/unregister", 1, false, c.subscTopic)
	}
	c.subscTopic = newTopic
}

func (c *client) Publish(lat, lng float64, qos byte, retained bool, payload interface{}) mqtt.Token {
	topic := celID2TopicName(s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng)))
	return c.client.Publish(topic, qos, retained, payload)
}

func celID2TopicName(id s2.CellID) string {
	idString := strings.Replace(id.String(), "/", "", 1)
	return strings.Replace(idString, "", "/", -1)
}

func capOnEarth(center s2.Point, radiusKm float64) s2.Cap {
	const earthRadiusKm = 6371.01
	ratio := (radiusKm / earthRadiusKm)
	return s2.CapFromCenterAngle(center, s1.Angle(ratio))
}
