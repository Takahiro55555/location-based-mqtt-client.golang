package client

import "math"

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

//////////////  　  以下 平面直角座標 構造体関連  　 //////////////

// 参考：緯度経度と平面直角座標の相互変換をPythonで実装する
//   URL：https://qiita.com/sw1227/items/e7a590994ad7dcd0e8ab

type RectangularPlaneCoordinate interface {
	// TranslateXYZ 関数
	// 緯度・経度を平面直角座標系に変換する。高度はそのまま
	// c: 平面直角座標系に変換したい座標
	TranslateXYZ(c CoordinateLLA) CoordinateXYZ
}

type rectangularPlaneCoordinate struct {
	origin CoordinateLLA

	// 定数 (a, F: 世界測地系-測地基準系1980（GRS80）楕円体)
	m0 float64
	a  float64
	F  float64
}

// origin: 平面直角座標系原点の緯度・経度[度]
func NewRectangularPlaneCoordinate(origin CoordinateLLA) RectangularPlaneCoordinate {
	return &rectangularPlaneCoordinate{
		// 定数 (a, F: 世界測地系-測地基準系1980（GRS80）楕円体)
		origin: origin, m0: 0.9999, a: 6378137, F: 298.257222101,
	}
}

func (r *rectangularPlaneCoordinate) TranslateXYZ(c CoordinateLLA) CoordinateXYZ {
	// HACK: 文字数が多すぎる行がある。要リファクタリング
	// latitude  ==> phi    緯度
	// longitude ==> lambda 経度

	// 緯度経度・平面直角座標系原点の緯度経度をラジアンに変換
	la, lo := c.GetLL()
	phiRad := rad(la)
	lambdaRad := rad(lo)

	la, lo = r.origin.GetLL()
	phi0Rad := rad(la)
	lambda0Rad := rad(lo)

	// (1) n, A_i, alpha_iの計算
	n := 1.0 / (2*r.F - 1)
	aArray := calcAArray(n)
	alphaArray := calcAlphaArray(n)

	// (2), S, Aの計算
	A_ := ((r.m0 * r.a) / (1. + n)) * aArray[0]                                                                                                                                            // [m]
	S_ := ((r.m0 * r.a) / (1. + n)) * (aArray[0]*phi0Rad + fVecDot(aArray[1:], fVecApply([]float64{1, 2, 3, 4, 5}, func(i int, v float64) float64 { return math.Sin(2. * phi0Rad * v) }))) // [m]

	// (3) lambda_c, lambda_sの計算
	lambda_c := math.Cos(lambdaRad - lambda0Rad)
	lambda_s := math.Sin(lambdaRad - lambda0Rad)

	// (4) t, t_の計算
	t := math.Sinh(math.Atanh(math.Sin(phiRad)) - ((2.*math.Sqrt(n))/(1.+n))*math.Atanh(((2.*math.Sqrt(n))/(1.+n))*math.Sin(phiRad)))
	t_ := math.Sqrt(1. + t*t)

	// (5) xi', eta'の計算
	xi2 := math.Atan(t / lambda_c) // [rad]
	eta2 := math.Atanh(lambda_s / t_)

	// (6) x, yの計算
	x := A_*(xi2+fVecSum(fVecMulti(alphaArray[1:], fVecMulti(fVecApply([]float64{1, 2, 3, 4, 5}, func(i int, v float64) float64 { return math.Sin(2 * xi2 * v) }), fVecApply([]float64{1, 2, 3, 4, 5}, func(i int, v float64) float64 { return math.Cosh(2 * eta2 * v) }))))) - S_ // [m]
	y := A_ * (eta2 + fVecSum(fVecMulti(alphaArray[1:], fVecMulti(fVecApply([]float64{1, 2, 3, 4, 5}, func(i int, v float64) float64 { return math.Cos(2 * xi2 * v) }), fVecApply([]float64{1, 2, 3, 4, 5}, func(i int, v float64) float64 { return math.Sinh(2 * eta2 * v) }))))) // [m]

	// return CoordinateXYZ # [m]
	_, _, al := c.GetLLA()
	cXYZ := NewCoordinateXYZ(x, y, al)
	return cXYZ
}

func calcAArray(n float64) [6]float64 {
	var a [6]float64
	a[0] = 1 + math.Pow(n, 2)/4. + math.Pow(n, 4)/64.
	a[1] = -(3. / 2) * (n - math.Pow(n, 3)/8. - math.Pow(n, 5)/64.)
	a[2] = (15. / 16) * (math.Pow(n, 2) - math.Pow(n, 4)/4.)
	a[3] = -(35. / 48) * (math.Pow(n, 3) - (5./16)*math.Pow(n, 5))
	a[4] = (315. / 512) * math.Pow(n, 4)
	a[5] = -(693. / 1280) * math.Pow(n, 5)
	return a
}

func calcAlphaArray(n float64) [6]float64 {
	var a [6]float64
	a[0] = 0 // dummy
	a[1] = (1./2)*n - (2./3)*math.Pow(n, 2) + (5./16)*math.Pow(n, 3) + (41./180)*math.Pow(n, 4) - (127./288)*math.Pow(n, 5)
	a[2] = (13./48)*math.Pow(n, 2) - (3./5)*math.Pow(n, 3) + (557./1440)*math.Pow(n, 4) + (281./630)*math.Pow(n, 5)
	a[3] = (61./240)*math.Pow(n, 3) - (103./140)*math.Pow(n, 4) + (15061./26880)*math.Pow(n, 5)
	a[4] = (49561./161280)*math.Pow(n, 4) - (179./168)*math.Pow(n, 5)
	a[5] = (34729. / 80640) * math.Pow(n, 5)
	return a
}

func rad(ang float64) float64 {
	return ang / 360 * (2.0 * math.Pi)
}

func fVecDot(a, b []float64) float64 {
	if len(a) != len(b) {
		panic("Array length not mutch")
	}
	result := 0.
	for i := 0; i < len(a); i++ {
		result += (a[i] * b[i])
	}
	return result
}

func fVecApply(a []float64, f func(i int, v float64) float64) []float64 {
	result := make([]float64, len(a))
	for i, v := range a {
		result[i] = f(i, v)
	}
	return result
}

func fVecSum(a []float64) float64 {
	result := 0.
	for _, v := range a {
		result += v
	}
	return result
}

func fVecMulti(a, b []float64) []float64 {
	if len(a) != len(b) {
		panic("Array length not mutch")
	}
	result := make([]float64, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = (a[i] * b[i])
	}
	return result
}

//////////////  　  以上 平面直角座標 構造体関連  　 //////////////
