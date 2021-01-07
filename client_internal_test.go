package client

import (
	"reflect"
	"testing"

	"github.com/golang/geo/s2"
)

func TestExtractUnduplicateTopics(t *testing.T) {
	type args struct {
		currentTopics []string
		newTopics     []string
	}
	type want struct {
		unsubscribeTopics []string
		subscribeTopics   []string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test 01",
			args: args{
				currentTopics: []string{"/0/#"},
				newTopics:     []string{"/1/#"},
			},
			want: want{
				unsubscribeTopics: []string{"/0/#"},
				subscribeTopics:   []string{"/1/#"},
			},
		},
		{
			name: "Test 02",
			args: args{
				currentTopics: []string{"/0/#", "/1/#"},
				newTopics:     []string{"/1/#", "/2/#"},
			},
			want: want{
				unsubscribeTopics: []string{"/0/#", ""},
				subscribeTopics:   []string{"", "/2/#"},
			},
		},
		{
			name: "Test 03",
			args: args{
				currentTopics: []string{"/1/#"},
				newTopics:     []string{"/1/#"},
			},
			want: want{
				unsubscribeTopics: []string{""},
				subscribeTopics:   []string{""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ut, st := extractUnduplicateTopics(tt.args.currentTopics, tt.args.newTopics)
			if !reflect.DeepEqual(ut, tt.want.unsubscribeTopics) {
				t.Errorf("UnsubscribeTopics Want: %v, Result: %v", tt.want.unsubscribeTopics, ut)
			}
			if !reflect.DeepEqual(st, tt.want.subscribeTopics) {
				t.Errorf("SubscribeTopics Want: %v, Result: %v", tt.want.subscribeTopics, st)
			}
		})
	}
}

func TestCelID2TopicName(t *testing.T) {
	type args struct {
		lat float64
		lng float64
	}
	type want struct {
		topic string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test 01",
			args: args{
				lat: 1,
				lng: 1,
			},
			want: want{
				topic: "/0/2/0/0/0/0/0/2/2/0/2/0/0/2/2/2/2/0/0/1/0/0/0/0/2/1/1/3/2/2/2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			topic := celID2TopicName(s2.CellIDFromLatLng(s2.LatLngFromDegrees(tt.args.lat, tt.args.lng)))
			if topic != tt.want.topic {
				t.Errorf("UnsubscribeTopics Want: %v, Result: %v", tt.want.topic, topic)
			}
		})
	}
}
