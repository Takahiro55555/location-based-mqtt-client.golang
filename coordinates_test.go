package client_test

import (
	"client"
	"testing"
)

func TestGet(t *testing.T) {
	type args struct {
		latitudeOrX  float64
		longitudeOrY float64
		altitudeOrZ  float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 01",
			args: args{
				latitudeOrX:  10.,
				longitudeOrY: 20.,
				altitudeOrZ:  30.,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewCoordinateXYZ(tt.args.latitudeOrX, tt.args.longitudeOrY, tt.args.altitudeOrZ)
			x, y, z := c.Get()
			if x != tt.args.latitudeOrX || y != tt.args.longitudeOrY || z != tt.args.altitudeOrZ {
				t.Fatalf("Not match input and output.")
			}

			c = client.NewCoordinateLLA(tt.args.latitudeOrX, tt.args.longitudeOrY, tt.args.altitudeOrZ)
			la, lo, al := c.Get()
			if la != tt.args.latitudeOrX || lo != tt.args.longitudeOrY || al != tt.args.altitudeOrZ {
				t.Fatalf("Not match input and output.")
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		latitudeOrX  float64
		longitudeOrY float64
		altitudeOrZ  float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 01",
			args: args{
				latitudeOrX:  10.,
				longitudeOrY: 20.,
				altitudeOrZ:  30.,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewCoordinateXYZ(0, 0, 0)
			c.Update(tt.args.latitudeOrX, tt.args.longitudeOrY, tt.args.altitudeOrZ)
			x, y, z := c.Get()
			if x != tt.args.latitudeOrX || y != tt.args.longitudeOrY || z != tt.args.altitudeOrZ {
				t.Fatalf("Not match input and output.")
			}

			c = client.NewCoordinateLLA(0, 0, 0)
			c.Update(tt.args.latitudeOrX, tt.args.longitudeOrY, tt.args.altitudeOrZ)
			la, lo, al := c.Get()
			if la != tt.args.latitudeOrX || lo != tt.args.longitudeOrY || al != tt.args.altitudeOrZ {
				t.Fatalf("Not match input and output.")
			}
		})
	}
}
