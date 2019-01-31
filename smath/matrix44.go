package smath

import "math"

type Matrix_t struct{
	M 	[4][4]float64
}

//type Matrix_t struct {
//	M11,M12,M13,M14 float64
//	M21,M22,M23,M24 float64
//	M31,M32,M33,M34 float64
//	M41,M42,M43,M44 float64
//}
// matrix-inverse functions ---------------------------------------------------
func fDeterminant2x2( i_fA, i_fB, i_fC, i_fD float64 ) (float64){

	return i_fA * i_fD - i_fB * i_fC
}

func fDeterminant3x3(  i_fA1,  i_fA2,  i_fA3, i_fB1,  i_fB2,  i_fB3, i_fC1,  i_fC2,  i_fC3 float64)(float64){
	// src: http://www.acm.org/pubs/tog/GraphicsGems/gems/MatrixInvert.c
	return i_fA1 * fDeterminant2x2( i_fB2, i_fB3, i_fC2, i_fC3 )-
		i_fB1 * fDeterminant2x2( i_fA2, i_fA3, i_fC2, i_fC3 )+
		i_fC1 * fDeterminant2x2( i_fA2, i_fA3, i_fB2, i_fB3 )

}

func fMinorDeterminant(i_matMatrix *Matrix_t, i_iRow, i_iColumn uint32 )(float64) {
	// src: http://www.codeproject.com/csharp/Matrix.asp
	var fMat3x3 [3][3]float64

	for r,m := 0,0; r < 4; r++{
		if uint32(r) == i_iRow {continue}
		for c , n := 0,0; c < 4; c++ {
			if	uint32(c) == i_iColumn{
				continue
			}
			fMat3x3[m][n] = i_matMatrix.M[r][c]
			n++
		}
		m++
	}
	return fDeterminant3x3( fMat3x3[0][0], fMat3x3[0][1], fMat3x3[0][2],
		fMat3x3[1][0], fMat3x3[1][1], fMat3x3[1][2], fMat3x3[2][0], fMat3x3[2][1], fMat3x3[2][2] )
}


func (m *Matrix_t) matAdjoint ( t *Matrix_t ) (Matrix_t){
	var matReturn Matrix_t
	for r:=0 ; r<4; r++ {
		for c := 0; c < 4; c++ {
			matReturn.M[c][r] = math.Pow(-1.0,float64(r+c)) * fMinorDeterminant(t,uint32(r),uint32(c))
		}
	}
	return matReturn
}

func (m *Matrix_t)determinant()(float64){
	// src: http://www.acm.org/pubs/tog/GraphicsGems/gems/MatrixInvert.c
	return (m.M[0][0] * fDeterminant3x3( m.M[1][1], m.M[2][1], m.M[3][1], m.M[1][2], m.M[2][2], m.M[3][2],m.M[1][3], m.M[2][3], m.M[3][3])-
		m.M[0][1] * fDeterminant3x3( m.M[1][0], m.M[2][0], m.M[3][0], m.M[1][2], m.M[2][2], m.M[3][2], m.M[1][3], m.M[2][3], m.M[3][3])+
		m.M[0][2] * fDeterminant3x3( m.M[1][0], m.M[2][0], m.M[3][0], m.M[1][1], m.M[2][1], m.M[3][1], m.M[1][3], m.M[2][3], m.M[3][3])-
		m.M[0][3] * fDeterminant3x3( m.M[1][0], m.M[2][0], m.M[3][0], m.M[1][1], m.M[2][1], m.M[3][1], m.M[1][2], m.M[2][2], m.M[3][2]))

}

func (a *Matrix_t)add( b *Matrix_t){
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			a.M[i][j] = a.M[i][j] + b.M[i][j]
		}
	}
}

func (a *Matrix_t)sub(b *Matrix_t){
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			a.M[i][j] = a.M[i][j] - b.M[i][j]
		}
	}
}
func (a *Matrix_t)mul( b *Matrix_t){
	var c Matrix_t
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			c.M[j][i] = a.M[j][0] * b.M[0][i] +
				a.M[j][1] * b.M[1][i] +
				a.M[j][2] * b.M[2][i] +
				a.M[j][3] * b.M[3][i]
		}
	}
	*a = c
}
//
func (a *Matrix_t) scale(f float64){
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			a.M[i][j] = a.M[i][j] * f
		}
	}
}
// /运算
func (a *Matrix_t) divide(f float64){
	a.scale(1/f)
}
func (m *Matrix_t) fan()(Matrix_t){
	fDeterminant := m.determinant()
	if math.Abs(fDeterminant) < 0.00000005{
		return *m
	}
	return m.matAdjoint(m).matrix_divide(fDeterminant)
}

func (a *Matrix_t)matrix_add( b *Matrix_t)(Matrix_t){
	var c Matrix_t
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			c.M[i][j] = a.M[i][j] + b.M[i][j]
		}
	}
	return c
}
func (a *Matrix_t)matrix_sub(b *Matrix_t)(Matrix_t){
	var c Matrix_t
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			c.M[i][j] = a.M[i][j] - b.M[i][j]
		}
	}
	return c
}
func (a *Matrix_t)matrix_mul( b *Matrix_t)(Matrix_t){
	var c Matrix_t
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			c.M[j][i] = a.M[j][0] * b.M[0][i] +
				a.M[j][1] * b.M[1][i] +
				a.M[j][2] * b.M[2][i] +
				a.M[j][3] * b.M[3][i]
		}
	}
	return c
}
//
func (a *Matrix_t) matrix_scale(f float64)(Matrix_t){
	var c Matrix_t
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			c.M[i][j] = a.M[i][j] * f
		}
	}
	return c
}

// /运算
func (a *Matrix_t) matrix_divide(f float64)(Matrix_t){
	return a.matrix_scale(1/f)
}

func (a *Matrix_t) matrix_transpose()(Matrix_t){
	var o_mat Matrix_t
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			o_mat.M[i][j] = a.M[j][i]
		}
	}
	return o_mat
}


func (m *Matrix_t) Set_identity(){
	m.M[0][0] ,m.M[1][1] , m.M[2][2] , m.M[3][3] = 1.0,1.0,1.0,1.0
	m.M[0][1] , m.M[0][2] , m.M[0][3] = 0,0,0
	m.M[1][0] , m.M[1][2] , m.M[1][3] = 0,0,0
	m.M[2][0] ,m.M[2][1] , m.M[2][3] = 0,0,0
	m.M[3][0] ,m.M[3][1] , m.M[3][2]  = 0,0,0
}

func (m *Matrix_t) set_zero(){
	m.M[0][0] ,m.M[1][1] , m.M[2][2] , m.M[3][3] = 0,0,0,0
	m.M[0][1] , m.M[0][2] , m.M[0][3] = 0,0,0
	m.M[1][0] , m.M[1][2] , m.M[1][3] = 0,0,0
	m.M[2][0] ,m.M[2][1] , m.M[2][3] = 0,0,0
	m.M[3][0] ,m.M[3][1] , m.M[3][2]  = 0,0,0
}

//平移
func (m *Matrix_t)set_translation( x float64,y float64, z float64){
	m.Set_identity()
	m.M[3][0] = x
	m.M[3][1] = y
	m.M[3][2] = z

}

func (m *Matrix_t)set_translation_vector(v Vector_t){
	m.set_translation(v.X,v.Y,v.Z)
}




//缩放
func (m *Matrix_t)set_scale(x ,y, z float64){
	m.Set_identity()
	m.M[0][0] = x
	m.M[1][1] = y
	m.M[2][2] = z
}

func (m *Matrix_t)set_scale_vector(v Vector_t){
	m.set_scale(v.X,v.Y,v.Z)
}
//从四元数构造旋转矩阵
func (m *Matrix_t) set_rotate( x ,y,z, theta float64){
	qsin := math.Sin(theta * 0.5)
	qcos := math.Cos(theta * 0.5)
	vec := Vector_t{x,y,z,1.0	}
	w := qcos
	vec.Normalize()
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

func (m *Matrix_t) matrix_rotate_x(theta float64) (Matrix_t) {
	var o_matMatOut Matrix_t
	fSin := math.Sin( theta )
	fCos := math.Cos( theta )
	o_matMatOut.M[0][0] = 1.0; o_matMatOut.M[0][1] = 0.0; o_matMatOut.M[0][2] = 0.0; o_matMatOut.M[0][3] = 0.0
	o_matMatOut.M[1][0] = 0.0; o_matMatOut.M[1][1] = fCos; o_matMatOut.M[1][2] = fSin; o_matMatOut.M[1][3] = 0.0
	o_matMatOut.M[2][0] = 0.0; o_matMatOut.M[2][1] = -fSin; o_matMatOut.M[2][2] = fCos; o_matMatOut.M[2][3] = 0.0
	o_matMatOut.M[3][0] = 0.0; o_matMatOut.M[3][1] = 0.0; o_matMatOut.M[3][2] = 0.0; o_matMatOut.M[3][3] = 1.0
	return o_matMatOut
}

func (m *Matrix_t) matrix_rotate_y(theta float64) (Matrix_t){
	var o_matMatOut Matrix_t
	fSin := math.Sin( theta )
	fCos := math.Cos( theta )
	o_matMatOut.M[0][0] = fCos; o_matMatOut.M[0][1] = 0.0; o_matMatOut.M[0][2] = -fSin; o_matMatOut.M[0][3] = 0.0
	o_matMatOut.M[1][0] = 0.0; o_matMatOut.M[1][1] = 1.0; o_matMatOut.M[1][2] = 0.0; o_matMatOut.M[1][3] = 0.0
	o_matMatOut.M[2][0] = fSin; o_matMatOut.M[2][1] = 0.0; o_matMatOut.M[2][2] = fCos; o_matMatOut.M[2][3] = 0.0
	o_matMatOut.M[3][0] = 0.0; o_matMatOut.M[3][1] = 0.0; o_matMatOut.M[3][2] = 0.0; o_matMatOut.M[3][3] = 1.0
	return o_matMatOut
}

func (m *Matrix_t) matrix_rotate_z(theta float64) (Matrix_t){
	var o_matMatOut Matrix_t
	fSin := math.Sin( theta )
	fCos := math.Cos( theta )
	o_matMatOut.M[0][0] =  fCos; o_matMatOut.M[0][1] = fSin; o_matMatOut.M[0][2] = 0.0; o_matMatOut.M[0][3] = 0.0
	o_matMatOut.M[1][0] = -fSin; o_matMatOut.M[1][1] = fCos; o_matMatOut.M[1][2] = 0.0; o_matMatOut.M[1][3] = 0.0
	o_matMatOut.M[2][0] = 0.0; o_matMatOut.M[2][1] = 0.0; o_matMatOut.M[2][2] = 1.0; o_matMatOut.M[2][3] = 0.0
	o_matMatOut.M[3][0] = 0.0; o_matMatOut.M[3][1] = 0.0; o_matMatOut.M[3][2] = 0.0; o_matMatOut.M[3][3] = 1.0
	return o_matMatOut
}

func (m *Matrix_t) matrix_rotate_yaw_pitch_roll(yaw,pitch,roll float64)(Matrix_t){
	m_yaw := m.matrix_rotate_x(yaw)
	m_pitch := m.matrix_rotate_y(pitch)
	m_roll := m.matrix_rotate_z(roll)
	return m.matrix_mul(&m_yaw).matrix_mul(&m_pitch).matrix_mul(&m_roll)

}

func (m *Matrix_t) matrix_rotate_yaw_pitch_roll_vector(v Vector_t)(Matrix_t){

	return m.matrix_rotate_yaw_pitch_roll(v.X,v.Y,v.Z)
}

func (m *Matrix_t) matMatrix44RotationAxis( i_vAxis Vector_t, i_fRot float64 )(Matrix_t){
	// make sure incoming axis is normalized!
	// http://www.euclideanspace.com/maths/algebra/matrix/orthogonal/rotation/openforum.htm
	var matMatOut Matrix_t
	fSin := math.Sin( i_fRot )
	fCos := math.Cos( i_fRot )
	fInvCos := 1.0 - fCos

	matMatOut.M[0][0]= fInvCos * i_vAxis.X * i_vAxis.X + fCos;
	matMatOut.M[0][1]= fInvCos * i_vAxis.X* i_vAxis.Y - i_vAxis.Z * fSin;
	matMatOut.M[0][2]= fInvCos * i_vAxis.X* i_vAxis.Z + i_vAxis.Y * fSin;
	matMatOut.M[0][3]= 0.0;

	matMatOut.M[1][0] = fInvCos * i_vAxis.X * i_vAxis.X + i_vAxis.Z * fSin;
	matMatOut.M[1][1] = fInvCos * i_vAxis.Y * i_vAxis.Y + fCos;
	matMatOut.M[1][2] = fInvCos * i_vAxis.Y * i_vAxis.Z - i_vAxis.X * fSin;
	matMatOut.M[1][3] = 0.0;

	matMatOut.M[2][0] = fInvCos * i_vAxis.X * i_vAxis.Z- i_vAxis.Y * fSin;
	matMatOut.M[2][1] = fInvCos * i_vAxis.Y * i_vAxis.Z+ i_vAxis.X * fSin;
	matMatOut.M[2][2] = fInvCos * i_vAxis.Z * i_vAxis.Z+ fCos;
	matMatOut.M[2][3] = 0.0;

	matMatOut.M[3][0] = 0.0;
	matMatOut.M[3][1]= 0.0
	matMatOut.M[3][2]= 0.0
	matMatOut.M[3][3]= 1.0
	return matMatOut

}

func (m *Matrix_t) set_lookat(eye ,at ,up *Vector_t) {
	var xaxis,yaxis,zaxis Vector_t
	Vector_sub(&zaxis,at,eye)
	zaxis.Normalize()
	Vector_crossproduct(&xaxis,up,&zaxis)
	xaxis.Normalize()
	Vector_crossproduct(&yaxis,&zaxis,&xaxis)

	m.M[0][0] = xaxis.X
	m.M[1][0] = xaxis.Y
	m.M[2][0] = xaxis.Z
	m.M[3][0] = -Vector_dotproduct(&xaxis,eye)

	m.M[0][1] = yaxis.X
	m.M[1][1] = yaxis.Y
	m.M[2][1] = yaxis.Z
	m.M[3][1] = -Vector_dotproduct(&yaxis,eye)

	m.M[0][2] = zaxis.X
	m.M[1][2] = zaxis.Y
	m.M[2][2] = zaxis.Z
	m.M[3][2] = -Vector_dotproduct(&zaxis,eye)

	m.M[0][3],m.M[1][3],m.M[2][3] = 0,0,0
	m.M[3][3] = 1

}

func (m *Matrix_t) set_lookat_rh(eye ,at ,up *Vector_t) {
	var xaxis,yaxis,zaxis Vector_t
	Vector_sub(&zaxis,eye,at)
	zaxis.Normalize()
	Vector_crossproduct(&xaxis,up,&zaxis)
	xaxis.Normalize()
	Vector_crossproduct(&yaxis,&zaxis,&xaxis)

	m.M[0][0] = xaxis.X
	m.M[1][0] = xaxis.Y
	m.M[2][0] = xaxis.Z
	m.M[3][0] = -Vector_dotproduct(&xaxis,eye)

	m.M[0][1] = yaxis.X
	m.M[1][1] = yaxis.Y
	m.M[2][1] = yaxis.Z
	m.M[3][1] = -Vector_dotproduct(&yaxis,eye)

	m.M[0][2] = zaxis.X
	m.M[1][2] = zaxis.Y
	m.M[2][2] = zaxis.Z
	m.M[3][2] = -Vector_dotproduct(&zaxis,eye)

	m.M[0][3],m.M[1][3],m.M[2][3] = 0,0,0
	m.M[3][3] = 1

}

func (m *Matrix_t)matMatrix44OrthoOffCenterLH( i_fLeft, i_fRight,i_fBottom,i_fTop, i_fZNear,i_fZFar float64)(Matrix_t){
	var matMatOut Matrix_t

	matMatOut.M[0][0] = 2.0/ (i_fRight - i_fLeft); matMatOut.M[0][1] = 0.0; matMatOut.M[0][2] = 0.0; matMatOut.M[0][3] = 0.0
	matMatOut.M[1][0] = 0.0; matMatOut.M[1][1] =  2.0 / (i_fTop - i_fBottom); matMatOut.M[1][2] = 0.0; matMatOut.M[1][3] = 0.0
	matMatOut.M[2][0] = 0.0; matMatOut.M[2][1] = 0.0; matMatOut.M[2][2] = 1.0 / (i_fZFar - i_fZNear); matMatOut.M[2][3] = 0.0
	matMatOut.M[3][0] = (i_fLeft + i_fRight) / (i_fLeft - i_fRight); matMatOut.M[3][1] = (i_fBottom + i_fTop) / (i_fBottom - i_fTop); matMatOut.M[3][2] = i_fZNear / (i_fZNear - i_fZFar); matMatOut.M[3][3] = 1.0

	return matMatOut
}

func (m *Matrix_t)matMatrix44OrthoOffCenterRH( i_fLeft, i_fRight,i_fBottom,i_fTop, i_fZNear,i_fZFar float64)(Matrix_t){

	var matMatOut Matrix_t

	matMatOut.M[0][0] = 2.0/ (i_fRight - i_fLeft); matMatOut.M[0][1] = 0.0; matMatOut.M[0][2] = 0.0; matMatOut.M[0][3] = 0.0
	matMatOut.M[1][0] = 0.0; matMatOut.M[1][1] =  2.0 / (i_fTop - i_fBottom); matMatOut.M[1][2] = 0.0; matMatOut.M[1][3] = 0.0
	matMatOut.M[2][0] = 0.0; matMatOut.M[2][1] = 0.0; matMatOut.M[2][2] = 1.0 / (i_fZNear - i_fZFar); matMatOut.M[2][3] = 0.0
	matMatOut.M[3][0] = (i_fLeft + i_fRight) / (i_fLeft - i_fRight); matMatOut.M[3][1] = (i_fBottom + i_fTop) / (i_fBottom - i_fTop); matMatOut.M[3][2] = i_fZNear / (i_fZNear - i_fZFar); matMatOut.M[3][3] = 1.0

	return matMatOut

}

func (m *Matrix_t)matMatrix44OrthoLH( i_fWidth, i_fHeight,  i_fZNear, i_fZFar float64 )(Matrix_t){
	return m.matMatrix44OrthoOffCenterLH( -i_fWidth * 0.5, i_fWidth * 0.5, -i_fHeight * 0.5, i_fHeight * 0.5, i_fZNear, i_fZFar )
}

func (m *Matrix_t)matMatrix44OrthoRH( i_fWidth, i_fHeight, i_fZNear,i_fZFar  float64 )(Matrix_t) {
	return m.matMatrix44OrthoOffCenterRH( -i_fWidth * 0.5, i_fWidth * 0.5, -i_fHeight * 0.5, i_fHeight * 0.5, i_fZNear, i_fZFar )
}


func (m *Matrix_t) set_perspective(fovy ,aspect ,zn ,zf float64){
	fax := 1/math.Tan(fovy * 0.5)
	m.set_zero()
	m.M[0][0] = fax / aspect
	m.M[1][1] = fax
	m.M[2][2] = zf/(zf - zn)
	m.M[3][2] = - zn * zf /(zf-zn)
	m.M[2][3] = 1
}

func (m *Matrix_t) set_perspectiveFovLH(fovy ,aspect ,zn ,zf float64){
	fViewSpaceHeight := 1/math.Tan(fovy * 0.5)
	fViewSpaceWidth := fViewSpaceHeight/aspect
	m.set_zero()
	m.M[0][0] = fViewSpaceWidth
	m.M[1][1] = fViewSpaceHeight
	m.M[2][2] = zf/(zf - zn)
	m.M[3][2] = - zn * zf /(zf-zn)
	m.M[2][3] = 1
}

func (m *Matrix_t) set_perspectiveFovRH(fovy ,aspect ,zn ,zf float64){
	fViewSpaceHeight := 1/math.Tan(fovy * 0.5)
	fViewSpaceWidth := fViewSpaceHeight/aspect
	m.set_zero()
	m.M[0][0] = fViewSpaceWidth
	m.M[1][1] = fViewSpaceHeight
	m.M[2][2] = zf/(zn - zf)
	m.M[3][2] = - zn * zf /(zf-zn)
	m.M[2][3] = 1
}
func (m *Matrix_t) set_perspectiveLH(fWidth, fHeight,zn ,zf float64){
	m.set_zero()
	m.M[0][0] = 2.0 * zn / fWidth
	m.M[1][1] = 2.0 * zn / fHeight
	m.M[2][2] = zf/(zf - zn)
	m.M[3][2] = - zn * zf /(zf-zn)
	m.M[2][3] = 1
}

func (m *Matrix_t) set_perspectiveRH(fWidth ,fHeight ,zn ,zf float64){
	m.set_zero()
	m.M[0][0] = 2.0 * zn / fWidth
	m.M[1][1] = 2.0 * zn / fHeight
	m.M[2][2] = zf/(zn - zf)
	m.M[3][2] = - zn * zf /(zf-zn)
	m.M[2][3] = 1
}

func (m *Matrix_t) set_perspectiveOffCenterLH(i_fLeft, i_fRight,i_fBottom,i_fTop, i_fZNear,i_fZFar float64){

	m.set_zero()
	m.M[0][0] = 2.0 * i_fZNear / (i_fRight - i_fLeft)
	m.M[1][1] = 2.0 * i_fZNear / (i_fTop - i_fBottom)
	m.M[2][0] =  (i_fLeft + i_fRight) / (i_fLeft - i_fRight)
	m.M[2][1] =  (i_fBottom + i_fTop) / (i_fBottom - i_fTop)
	m.M[2][2] = i_fZFar / (i_fZFar - i_fZNear)
	m.M[3][2] =  i_fZNear * i_fZFar / (i_fZNear - i_fZFar)
	m.M[2][3] = 1

}

func (m *Matrix_t) set_perspectiveOffCenterRH(i_fLeft, i_fRight,i_fBottom,i_fTop, i_fZNear,i_fZFar float64){
	m.set_zero()
	m.M[0][0] = 2.0 * i_fZNear / (i_fRight - i_fLeft)
	m.M[1][1] = 2.0 * i_fZNear / (i_fTop - i_fBottom)
	m.M[2][0] =  (i_fLeft + i_fRight) / (i_fLeft - i_fRight)
	m.M[2][1] =  (i_fBottom + i_fTop) / (i_fBottom - i_fTop)
	m.M[2][2] = i_fZFar / (i_fZNear-i_fZFar)
	m.M[3][2] =  i_fZNear * i_fZFar / (i_fZNear - i_fZFar)
	m.M[2][3] = -1
}

func (m *Matrix_t)getViewport( i_iX,i_iY, i_iWidth, i_iHeight uint32, i_fZNear,i_fZFar float64)(Matrix_t){
	var matMatOut Matrix_t
	matMatOut.set_zero()
	matMatOut.M[0][0] = float64(i_iWidth) * 0.5
	matMatOut.M[1][1] = float64(i_iHeight) * -0.5
	matMatOut.M[2][2] = i_fZFar - i_fZNear
	matMatOut.M[3][0] = float64(i_iX) + float64(i_iWidth) * 0.5
	matMatOut.M[3][1] = float64(i_iY) + float64(i_iHeight) * 0.5
	matMatOut.M[3][2] = i_fZNear
	matMatOut.M[3][3] = 1.0
	return matMatOut
}