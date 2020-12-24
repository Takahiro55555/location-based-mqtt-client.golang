package client_test

import (
	"client"
	"testing"
)

func TestGetLLA(t *testing.T) {
	type args struct {
		latitude  float64
		longitude float64
		altitude  float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 01",
			args: args{
				latitude:  10.,
				longitude: 20.,
				altitude:  30.,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewCoordinateLLA(tt.args.latitude, tt.args.longitude, tt.args.altitude)
			la, lo, al := c.GetLLA()
			if la != tt.args.latitude || lo != tt.args.longitude || al != tt.args.altitude {
				t.Fatalf("Not match input and output.")
			}
		})
	}
}

func TestUpdateLLA(t *testing.T) {
	type args struct {
		latitude  float64
		longitude float64
		altitude  float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 01",
			args: args{
				latitude:  10.,
				longitude: 20.,
				altitude:  30.,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewCoordinateLLA(0, 0, 0)
			c.UpdateLLA(tt.args.latitude, tt.args.longitude, tt.args.altitude)
			la, lo, al := c.GetLLA()
			if la != tt.args.latitude || lo != tt.args.longitude || al != tt.args.altitude {
				t.Fatalf("Not match input and output.")
			}
		})
	}
}

func TestGetLL(t *testing.T) {
	type args struct {
		latitude  float64
		longitude float64
		altitude  float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 01",
			args: args{
				latitude:  10.,
				longitude: 20.,
				altitude:  30.,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewCoordinateLLA(tt.args.latitude, tt.args.longitude, tt.args.altitude)
			la, lo := c.GetLL()
			if la != tt.args.latitude || lo != tt.args.longitude {
				t.Fatalf("Not match input and output.")
			}
		})
	}
}

func TestUpdateLL(t *testing.T) {
	type args struct {
		latitude  float64
		longitude float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 01",
			args: args{
				latitude:  10.,
				longitude: 20.,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewCoordinateLLA(0, 0, 0)
			c.UpdateLL(tt.args.latitude, tt.args.longitude)
			la, lo := c.GetLL()
			if la != tt.args.latitude || lo != tt.args.longitude {
				t.Fatalf("Not match input and output.")
			}
		})
	}
}

func TestGetXYZ(t *testing.T) {
	type args struct {
		x float64
		y float64
		z float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 01",
			args: args{
				x: 10.,
				y: 20.,
				z: 30.,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewCoordinateXYZ(tt.args.x, tt.args.y, tt.args.z)
			x, y, z := c.GetXYZ()
			if x != tt.args.x || y != tt.args.y || z != tt.args.z {
				t.Fatalf("Not match input and output.")
			}
		})
	}
}

func TestUpdateXYZ(t *testing.T) {
	type args struct {
		x float64
		y float64
		z float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 01",
			args: args{
				x: 10.,
				y: 20.,
				z: 30.,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewCoordinateXYZ(0, 0, 0)
			c.UpdateXYZ(tt.args.x, tt.args.y, tt.args.z)
			x, y, z := c.GetXYZ()
			if x != tt.args.x || y != tt.args.y || z != tt.args.z {
				t.Fatalf("Not match input and output.")
			}
		})
	}
}

func TestGetXY(t *testing.T) {
	type args struct {
		x float64
		y float64
		z float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 01",
			args: args{
				x: 10.,
				y: 20.,
				z: 30.,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewCoordinateXYZ(tt.args.x, tt.args.y, tt.args.z)
			x, y := c.GetXY()
			if x != tt.args.x || y != tt.args.y {
				t.Fatalf("Not match input and output.")
			}
		})
	}
}

func TestUpdateXY(t *testing.T) {
	type args struct {
		x float64
		y float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 01",
			args: args{
				x: 10.,
				y: 20.,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := client.NewCoordinateXYZ(0, 0, 0)
			c.UpdateXY(tt.args.x, tt.args.y)
			x, y := c.GetXY()
			if x != tt.args.x || y != tt.args.y {
				t.Fatalf("Not match input and output.")
			}
		})
	}
}
