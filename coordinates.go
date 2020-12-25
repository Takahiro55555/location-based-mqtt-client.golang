package client

type Coordinate interface {
	Update(a, b, c float64)
	Get() (float64, float64, float64)
}

type coordinateXYZ struct {
	x float64
	y float64
	z float64
}

func NewCoordinateXYZ(x, y, z float64) Coordinate {
	return &coordinateXYZ{x, y, z}
}

func (c *coordinateXYZ) Update(x, y, z float64) {
	c.x = x
	c.y = y
	c.z = z
}

func (c *coordinateXYZ) Get() (float64, float64, float64) {
	return c.x, c.y, c.z
}

func NewCoordinateLLA(latitude, longitude, altitude float64) Coordinate {
	return &coordinateLLA{latitude: latitude, longitude: longitude, altitude: altitude}
}

type coordinateLLA struct {
	latitude  float64 // 緯度
	longitude float64 // 経度
	altitude  float64 // 高度
}

func (c *coordinateLLA) Update(latitude, longitude, altitude float64) {
	c.latitude = latitude
	c.longitude = longitude
	c.altitude = altitude
}

func (c *coordinateLLA) Get() (float64, float64, float64) {
	return c.latitude, c.longitude, c.altitude
}
