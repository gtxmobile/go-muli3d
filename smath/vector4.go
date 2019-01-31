package smath

import "math"

type Vec_2 struct {

}

func Interp(x1 float64,x2 float64,t float64) float64 {
	return x1+(x2-x1)*t
}

type Vector_t struct {
	X	float64
	Y 	float64
	Z 	float64
	W 	float64
}
func Vector_length(v Vector_t) float64{
	sq := v.X * v.X + v.Y*v.Y + v.Z*v.Z
	return math.Sqrt((sq))
}

func Vector_add(z *Vector_t,x Vector_t, y Vector_t){
	z.X = x.X + y.X
	z.Y = x.Y + y.Y
	z.Z = x.Z + y.Z
	z.W = 1.0
}

func (x *Vector_t)Add_apply(y Vector_t){
	x.X = x.X + y.X
	x.Y = x.Y + y.Y
	x.Z = x.Z + y.Z
	x.W = 1.0
}

func (x *Vector_t)Sub_apply(y Vector_t){
	x.X = x.X - y.X
	x.Y = x.Y - y.Y
	x.Z = x.Z - y.Z
	x.W = 1.0
}

func (x *Vector_t)Sub_d(d float64){
	x.X = x.X - d
	x.Y = x.Y - d
	x.Z = x.Z - d
	x.W = 1.0
}
func Vector_sub(z *Vector_t,x Vector_t, y Vector_t){
	z.X = x.X - y.X
	z.Y = x.Y - y.Y
	z.Z = x.Z - y.Z
	z.W = 1.0
}

func V_sub(x Vector_t, y Vector_t) Vector_t{
	var z Vector_t
	z.X = x.X - y.X
	z.Y = x.Y - y.Y
	z.Z = x.Z - y.Z
	z.W = 1.0
	return z
}

func V_add(x Vector_t, y Vector_t) Vector_t{
	var z Vector_t
	z.X = x.X + y.X
	z.Y = x.Y + y.Y
	z.Z = x.Z + y.Z
	z.W = 1.0
	return z
}

func (z *Vector_t)Divide(d float64){
	z.X /= d
	z.Y /=	d
	z.Z /=	d
	z.W /=	d
}
func (z *Vector_t)Multiply(d float64)*Vector_t{
	z.X *= 	d
	z.Y *=	d
	z.Z *=	d
	z.W *=	d
	return z
}


func (x Vector_t) Matrix_apply(m Matrix_t)(Vector_t){
	X,Y,Z,W := x.X,x.Y,x.Z,x.W
	X = X * m.M[0][0] + Y * m.M[1][0]+Z * m.M[2][0] + W * m.M[3][0]
	Y = X * m.M[0][1] + Y * m.M[1][1]+Z * m.M[2][1] + W * m.M[3][1]
	Z = X * m.M[0][2] + Y * m.M[1][2]+Z * m.M[2][2] + W * m.M[3][2]
	W = X * m.M[0][3] + Y * m.M[1][3]+Z * m.M[2][3] + W * m.M[3][3]
	return Vector_t{X,Y,Z,W}
}

func Vector_dotproduct(x *Vector_t,y *Vector_t) float64 {
	return x.X * y.X + x.Y * y.Y + x.Z *y.Z
}

func Vector_crossproduct(x Vector_t,y Vector_t)Vector_t{
	z := Vector_t{}
	z.X = x.Y * y.Z - x.Z*y.Y
	z.Y = x.Z * y.X - x.X*y.Z
	z.Z = x.X * y.Y - x.Y*y.X
	z.W = 1.0
	return z
}

// 矢量插值，t取值 [0, 1]
func Vector_interp(z *Vector_t, x1 *Vector_t,x2 *Vector_t, t float64){
	z.X = Interp(x1.X,x2.X,t)
	z.Y = Interp(x1.Y,x2.Y,t)
	z.Z = Interp(x1.Z,x2.Z,t)
	z.W = 1.0
}

func (v *Vector_t) Normalize(){
	length := Vector_length(v)
	if length != 0.0 {
		inv := 1.0/length
		v.X *= inv
		v.Y *= inv
		v.Z *= inv
	}
}