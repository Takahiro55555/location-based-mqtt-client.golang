package client

type CoordinateXYZ interface {
	UpdateXYZ(x, y, z float64)
	UpdateXY(x, y float64)
	GetXYZ() (float64, float64, float64)
	GetXY() (float64, float64)
}

type coordinateXYZ struct {
	x float64
	y float64
	z float64
}

func NewCoordinateXYZ(x, y, z float64) CoordinateXYZ {
	return &coordinateXYZ{x, y, z}
}

func (c *coordinateXYZ) UpdateXYZ(x, y, z float64) {
	c.x = x
	c.y = y
	c.z = z
}

func (c *coordinateXYZ) UpdateXY(x, y float64) {
	c.x = x
	c.y = y
}

func (c *coordinateXYZ) GetXYZ() (float64, float64, float64) {
	return c.x, c.y, c.z
}

func (c *coordinateXYZ) GetXY() (float64, float64) {
	return c.x, c.y
}

func NewCoordinateLLA(latitude, longitude, altitude float64) CoordinateLLA {
	return &coordinateLLA{latitude: latitude, longitude: longitude, altitude: altitude}
}

type CoordinateLLA interface {
	UpdateLLA(latitude, longitude, altitude float64)
	UpdateLL(latitude, longitude float64)
	GetLLA() (float64, float64, float64)
	GetLL() (float64, float64)
}

type coordinateLLA struct {
	latitude  float64 // 緯度
	longitude float64 // 経度
	altitude  float64 // 高度
}

func (c *coordinateLLA) UpdateLLA(latitude, longitude, altitude float64) {
	c.latitude = latitude
	c.longitude = longitude
	c.altitude = altitude
}
func (c *coordinateLLA) UpdateLL(latitude, longitude float64) {
	c.latitude = latitude
	c.longitude = longitude
}
func (c *coordinateLLA) GetLLA() (float64, float64, float64) {
	return c.latitude, c.longitude, c.altitude
}
func (c *coordinateLLA) GetLL() (float64, float64) {
	return c.latitude, c.longitude
}
