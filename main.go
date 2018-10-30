package main

import (
	"fmt"
	"math"
)

type Matrix_t struct{
	M 	[4][4]float64
}
type Vector_t struct {
	X	float64
	Y 	float64
	Z 	float64
	W 	float64
}
func CMID(x int,min int, max int) int {
	if x<min{
		return min
	} else{
		if x>max{
			return max
		}else {
			return x
		}
	}
}

func interp(x1 float64,x2 float64,t float64) float64 {
	return x1+(x2-x1)*t
}

func vector_length(v *Vector_t) float64{
	sq := v.X * v.X + v.Y*v.Y + v.Z*v.Z
	return math.Sqrt((sq))
}

func init(){
	fmt.Println("init")
}
func main(){
	fmt.Println("main")
}

func vector_add(z *Vector_t,x *Vector_t, y *Vector_t){
	z.X = x.X + y.X
	z.Y = x.Y + y.Y
	z.Z = x.Z + y.Z
	z.W = 1.0
}

func vector_sub(z *Vector_t,x *Vector_t, y *Vector_t){
	z.X = x.X - y.X
	z.Y = x.Y - y.Y
	z.Z = x.Z - y.Z
	z.W = 1.0
}

func vector_dotproduct(x *Vector_t,y *Vector_t) float64 {
	return x.X * y.X + x.Y * y.Y + x.Z *y.Z
}

func vector_crossproduct(z *Vector_t,x *Vector_t,y*Vector_t){
	z.X = x.Y * y.Z - x.Z*y.Y
	z.Y = x.Z * y.X - x.X*y.Z
	z.Z = x.X * y.Y - x.Y*y.X
	z.W = 1.0
}

// 矢量插值，t取值 [0, 1]
func vector_interp(z *Vector_t, x1 *Vector_t,x2 *Vector_t, t float64){
	z.X = interp(x1.X,x2.X,t)
	z.Y = interp(x1.Y,x2.Y,t)
	z.Z = interp(x1.Z,x2.Z,t)
	z.W = 1.0
}

func vector_normalize(v *Vector_t){
	length := vector_length(v)
	if length != 0.0 {
		inv := 1.0/length
		v.X *= inv
		v.Y *= inv
		v.Z *= inv
	}
}

func matrix_add(c *Matrix_t,a *Matrix_t, b *Matrix_t){
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			c.M[i][j] = a.M[i][j] + b.M[i][j]
		}
	}
}

func matrix_sub(c *Matrix_t,a *Matrix_t, b *Matrix_t){
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			c.M[i][j] = a.M[i][j] - b.M[i][j]
		}
	}
}

func matrix_mul(c *Matrix_t,a *Matrix_t, b *Matrix_t){
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			c.M[j][i] = a.M[j][0] * b.M[0][i] +
				a.M[j][1] * b.M[1][i] +
				a.M[j][2] * b.M[2][i] +
				a.M[j][3] * b.M[3][i]
		}
	}
}

func matrix_scale(c *Matrix_t,a *Matrix_t, f float64){
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			c.M[i][j] = a.M[i][j] * f
		}
	}
}

func matrix_apply(y *Vector_t, x *Vector_t,m *Matrix_t){
	X,Y,Z,W := x.X,x.Y,x.Z,x.W
	y.X = X * m.M[0][0] + Y * m.M[1][0]+Z * m.M[2][0] + W * m.M[3][0]
	y.Y = X * m.M[0][1] + Y * m.M[1][1]+Z * m.M[2][1] + W * m.M[3][1]
	y.Z = X * m.M[0][2] + Y * m.M[1][2]+Z * m.M[2][2] + W * m.M[3][2]
	y.W = X * m.M[0][3] + Y * m.M[1][3]+Z * m.M[2][3] + W * m.M[3][3]

}

func matrix_set_identiry(m *Matrix_t){
	m.M[0][0] ,m.M[1][1] , m.M[2][2] , m.M[3][3] = 1.0,1.0,1.0,1.0
	m.M[0][1] , m.M[0][2] , m.M[0][3] = 0,0,0
	m.M[1][0] , m.M[1][2] , m.M[1][3] = 0,0,0
	m.M[2][0] ,m.M[2][1] , m.M[2][3] = 0,0,0
	m.M[3][0] ,m.M[3][1] , m.M[3][2]  = 0,0,0
}

func matrix_set_zero(m *Matrix_t){
	m.M[0][0] ,m.M[1][1] , m.M[2][2] , m.M[3][3] = 0,0,0,0
	m.M[0][1] , m.M[0][2] , m.M[0][3] = 0,0,0
	m.M[1][0] , m.M[1][2] , m.M[1][3] = 0,0,0
	m.M[2][0] ,m.M[2][1] , m.M[2][3] = 0,0,0
	m.M[3][0] ,m.M[3][1] , m.M[3][2]  = 0,0,0
}

//平移
func matrix_set_translate(m *Matrix_t, x float64,y float64, z float64){
	matrix_set_identiry(m)
	m.M[3][0] = x
	m.M[3][1] = y
	m.M[3][2] = z

}
//平移
func matrix_set_scale(m *Matrix_t, x float64,y float64, z float64){
	matrix_set_identiry(m)
	m.M[0][0] = x
	m.M[1][1] = y
	m.M[2][2] = z

}
//从四元数构造旋转矩阵
func matrix_set_rotate(m *Matrix_t, x float64,y float64, z float64, theta float64){
	qsin := math.Sin(theta * 0.5)
	qcos := math.Cos(theta * 0.5)
	vec := Vector_t{x,y,z,1.0	}
	w := qcos
	vector_normalize(&vec)
	x = vec.X * qsin
	y = vec.Y * qsin
	z = vec.Z * qsin

	x2 := x * x
	y2 := y * y
	z2 := z * z

	xy := x * y
	xz := x * z
	yz := y * z

	wx := w * x
	wy := w * y
	wz := w * z
	m.M[0][0] = 1 - 2*y2 - 2*z2
	m.M[1][0] = 2*xy - 2*wz
	m.M[2][0] = 2*xz + 2*wy

	m.M[0][1] = 2*xy + 2*wz
	m.M[1][1] = 1-2*x2 -2*z2
	m.M[2][1] = 2*yz - 2*wx

	m.M[0][2] = 2*xz - 2*wy
	m.M[1][2] = 2*yz + 2*wx
	m.M[2][2] = 1-2*x2 - 2*y2

	m.M[0][3] , m.M[1][3],m.M[2][3] = 0,0,0
	m.M[3][0] , m.M[3][1],m.M[3][2] = 0,0,0
	m.M[3][3] = 1
	}

func matrix_set_lookat(m *Matrix_t,eye *Vector_t,at *Vector_t,up *Vector_t) {
	var xaxis,yaxis,zaxis Vector_t
	vector_sub(&zaxis,at,eye)
	vector_normalize(&zaxis)
	vector_crossproduct(&xaxis,up,&zaxis)
	vector_normalize(&xaxis)
	vector_crossproduct(&yaxis,&zaxis,&xaxis)

	m.M[0][0] = xaxis.X
	m.M[1][0] = xaxis.Y
	m.M[2][0] = xaxis.Z
	m.M[3][0] = -vector_dotproduct(&xaxis,eye)

	m.M[0][1] = yaxis.X
	m.M[1][1] = yaxis.Y
	m.M[2][1] = yaxis.Z
	m.M[3][1] = -vector_dotproduct(&yaxis,eye)

	m.M[0][2] = zaxis.X
	m.M[1][2] = zaxis.Y
	m.M[2][2] = zaxis.Z
	m.M[3][2] = -vector_dotproduct(&zaxis,eye)

	m.M[0][3],m.M[1][3],m.M[2][3] = 0,0,0
	m.M[3][3] = 1

}

func matrix_set_perspective(m *Matrix_t,fovy float64,aspect float64,zn float64,zf float64){
	fax := 1/math.Tan(fovy * 0.5)
	matrix_set_zero(m)
	m.M[0][0] = fax / aspect
	m.M[1][1] = fax
	m.M[2][2] = zf/(zf - zn)
	m.M[3][2] = - zn * zf /(zf-zn)
	m.M[2][3] = 1

}

//坐标变换
type Transform_t struct {
	world 	Matrix_t
	view	Matrix_t
	projection 	Matrix_t
	transform 	Matrix_t
	w	float64
	h	float64
}

//矩阵更新

func transform_update(ts *Transform_t){
	var m Matrix_t
	matrix_mul(&m, &ts.world,&ts.view)
	matrix_mul(&ts.transform,&m,&ts.projection)

}