package smath
import (
	"math"
)

type Quaternion_t struct {
	X	float64
	Y 	float64
	Z 	float64
	W 	float64
}


//return quaternion( -x, -y, -z, w ); }

func (q *Quaternion_t) add_with( i_vVal Quaternion_t) {
	q.X += i_vVal.X; q.Y += i_vVal.Y;
	q.Z += i_vVal.Z; q.W += i_vVal.W;

}

func (q *Quaternion_t) sub_with( i_vVal Quaternion_t) {
	q.X -= i_vVal.X; q.Y -= i_vVal.Y;
	q.Z -= i_vVal.Z; q.W -= i_vVal.W;

}

func (q *Quaternion_t) multi_with( i_qVal Quaternion_t) *Quaternion_t {
	var qResult Quaternion_t
	qResult.X = q.W * i_qVal.X + q.X * i_qVal.W + q.Y * i_qVal.Z - q.Z * i_qVal.Y;
	qResult.Y = q.W * i_qVal.Y - q.X * i_qVal.Z + q.Y * i_qVal.W + q.Z * i_qVal.X;
	qResult.Z = q.W * i_qVal.Z + q.X * i_qVal.Y - q.Y * i_qVal.X + q.Z * i_qVal.W;
	qResult.W = q.W * i_qVal.W - q.X * i_qVal.X - q.Y * i_qVal.Y - q.Z * i_qVal.Z;
	*q = qResult;

	return q
}

func (q *Quaternion_t) divid_with(i_fVal float64){
	fInvVal := 1.0 / i_fVal;
	q.X *= fInvVal; q.Y *= fInvVal;
	q.Z *= fInvVal; q.W *= fInvVal;
}

func (q *Quaternion_t) dot_multi(i_fVal float64){
	q.X *= i_fVal; q.Y *= i_fVal;
	q.Z *= i_fVal; q.W *= i_fVal;
}



func (q *Quaternion_t) add( i_vVal Quaternion_t) Quaternion_t{
	return Quaternion_t{q.X + i_vVal.X, q.Y + i_vVal.Y, q.Z + i_vVal.Z, q.W + i_vVal.W }
}

func (q *Quaternion_t) sub( i_vVal Quaternion_t) Quaternion_t{
	return Quaternion_t{ q.X - i_vVal.X, q.Y - i_vVal.Y, q.Z - i_vVal.Z, q.W - i_vVal.W };
}

func (q *Quaternion_t) multi(i_qVal Quaternion_t) Quaternion_t{
	var qResult Quaternion_t
	qResult.X = q.W * i_qVal.X + q.X * i_qVal.W + q.Y * i_qVal.Z - q.Z * i_qVal.Y;
	qResult.Y = q.W * i_qVal.Y - q.X * i_qVal.Z + q.Y * i_qVal.W + q.Z * i_qVal.X;
	qResult.Z = q.W * i_qVal.Z + q.X * i_qVal.Y - q.Y * i_qVal.X + q.Z * i_qVal.W;
	qResult.W = q.W * i_qVal.W - q.X * i_qVal.X - q.Y * i_qVal.Y - q.Z * i_qVal.Z;
	return qResult;
}

func (q *Quaternion_t) multi_i( i_fVal float64) Quaternion_t{
	return Quaternion_t{ q.X * i_fVal, q.Y * i_fVal, q.Z * i_fVal, q.W * i_fVal }
}

func (q Quaternion_t) divid_i( i_fVal float64) Quaternion_t{
	fInv := 1.0 / i_fVal
	return Quaternion_t{ q.X * fInv, q.Y * fInv, q.Z * fInv, q.W * fInv }
}

func (q Quaternion_t) length() float64{
	return math.Sqrt( q.X * q.X + q.Y * q.Y + q.Z * q.Z + q.W * q.W )
}

func (q Quaternion_t) lengthsq() float64{
	return q.X * q.X + q.Y * q.Y + q.Z * q.Z + q.W * q.W
}

func (q *Quaternion_t) normalize() {
	fLength := q.length()
	//if( fLength >= FLT_EPSILON )
	{
		fInvLength := 1.0 / fLength;
		q.X *= fInvLength; q.Y *= fInvLength;
		q.Z *= fInvLength; q.W *= fInvLength;
	}
	//return q
}

func (q Quaternion_t) qQuaternionIdentity( o_qQuatOut *Quaternion_t ) *Quaternion_t {
	o_qQuatOut.X = 0
	o_qQuatOut.Y = 0
	o_qQuatOut.Z = 0
	o_qQuatOut.W = 1.0
	return o_qQuatOut
}

func (q Quaternion_t) qQuaternionRotationMatrix(o_qQuatOut *Quaternion_t,i_matMatrix *Matrix_t) *Quaternion_t{
	// http://www.gamasutra.com/features/19980703/quaternions_01.htm
	//TODO

	fDiagonal := i_matMatrix._11 + i_matMatrix._22 + i_matMatrix._33;
	if fDiagonal > 0.0{
		s := math.Sqrt( fDiagonal + 1.0 )
		o_qQuatOut.W = s / 2.0;
		s = 0.5 / s;
		o_qQuatOut.X = ( i_matMatrix._32 - i_matMatrix._23 ) * s;
		o_qQuatOut.X = ( i_matMatrix._13 - i_matMatrix._31 ) * s;
		o_qQuatOut.X = ( i_matMatrix._21 - i_matMatrix._12 ) * s;
	} else	{
		iNext:= [3]int{1, 2, 0}
		var q [4]float64
		i := 0
		if i_matMatrix._22 > i_matMatrix._11 { i = 1};
		if i_matMatrix._33 > i_matMatrix( i, i ) { i = 2};
		j := iNext[i]
		k := iNext[j]

		s := math.Sqrt(i_matMatrix( i, i ) - ( i_matMatrix( j, j ) + i_matMatrix( k, k ) ) + 1.0 )
		q[i] = s * 0.5

		if s >= FLT_EPSILON {
			s = 0.5 / s
		}

		q[3] = ( i_matMatrix( k, j ) - i_matMatrix( j, k ) ) * s;
		q[j] = ( i_matMatrix( j, i ) + i_matMatrix( i, j) ) * s;
		q[k] = ( i_matMatrix( k, i ) + i_matMatrix( i, k) ) * s;

		o_qQuatOut.X = q[0];
		o_qQuatOut.Y = q[1];
		o_qQuatOut.Z = q[2];
		o_qQuatOut.W = q[3];
	}

	return o_qQuatOut;
}

func (q Quaternion_t) qQuaternionSLerp( o_qQuatOut *Quaternion_t,i_qQuatA ,i_qQuatB Quaternion_t, i_fLerp float64)*Quaternion_t{
	// http://www.gamasutra.com/features/19980703/quaternions_01.htm

	// calc cosine
	fCosine := i_qQuatA.X * i_qQuatB.X + i_qQuatA.Y * i_qQuatB.Y + i_qQuatA.Z * i_qQuatB.Z + i_qQuatA.W * i_qQuatB.W

	// adjust signs (if necessary)
	var to1 [4]float64
	if( fCosine < 0.0 ){
		fCosine = -fCosine;
		to1[0] = -i_qQuatB.X;
		to1[1] = -i_qQuatB.Y;
		to1[2] = -i_qQuatB.Z;
		to1[3] = -i_qQuatB.W;
	}else	{
		to1[0] = i_qQuatB.X;
		to1[1] = i_qQuatB.Y;
		to1[2] = i_qQuatB.Z;
		to1[3] = i_qQuatB.W;
	}

	fOmega := math.Acos( fCosine );
	fInvSine := 1.0 / math.Sin( fOmega );
	fScale0 := math.Sin( ( 1.0 - i_fLerp ) * fOmega ) * fInvSine;
	fScale1 := math.Sin( i_fLerp * fOmega ) * fInvSine;

	// Calculate final values
	o_qQuatOut.X = fScale0 * i_qQuatB.X + fScale1 * to1[0];
	o_qQuatOut.Y = fScale0 * i_qQuatB.Y + fScale1 * to1[1];
	o_qQuatOut.Z = fScale0 * i_qQuatB.Z + fScale1 * to1[2];
	o_qQuatOut.W = fScale0 * i_qQuatB.W + fScale1 * to1[3];

	return o_qQuatOut;
}

func QuaternionToAxisAngle( i_qQuat *Quaternion_t,o_vAxis Vector_t, o_fAngle *float64) {
	fScale := math.Sqrt( i_qQuat.X * i_qQuat.X + i_qQuat.Y * i_qQuat.Y + i_qQuat.Z * i_qQuat.Z )
	if( fScale >= FLT_EPSILON ) {
		fInvScale := 1.0 / fScale;
		o_vAxis.X = i_qQuat.X * fInvScale;
		o_vAxis.Y = i_qQuat.Y * fInvScale;
		o_vAxis.Z = i_qQuat.Z * fInvScale;
	} else {
		o_vAxis.X = 0.0
		o_vAxis.Y = 1.0
		o_vAxis.Z = 0.0
	}

	*o_fAngle = 2.0 * math.Acos( i_qQuat.W )
}