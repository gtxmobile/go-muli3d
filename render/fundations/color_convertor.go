package fundations

import (
	"../../smath"
	"../../common"
)
/** R32BG32B32A32 unormalized float type.
*/
type Pixel_format int


type component_type float64
const component  = 4
//Color_rgba32f(const float* color):r(color[0]), g(color[1]), b(color[2]), a(color[3]){}
func vec4_to_color_rgba32f( v smath.Vector_t) Color_rgba32f{
	return Color_rgba32f{v.X,v.Y,v.Z,v.W}
}

type  Color_rgba32f struct{
	r, g, b, a float64
}
func ( this *Color_rgba32f) get_pointer() *smath.Vector_t{
	return &smath.Vector_t{this.r,this.g,this.b,this.a}
}



func ( this *Color_rgba32f) get_vec4() *smath.Vector_t{
	return &smath.Vector_t{this.r,this.g,this.b,this.a}
}


//color_rgb32f(const comp_t* color):r(color[0]), g(color[1]), b(color[2]){}

type Color_rgb32f struct {
	r, g, b float64
}


func (this *Color_rgb32f) to_rgba32f() Color_rgba32f{
	return Color_rgba32f{this.r, this.g, this.b, 1.0}
}

func (this *Color_rgb32f) assign(rhs *Color_rgba32f) *Color_rgb32f{
	this.r = rhs.r
	this.g = rhs.g
	this.b = rhs.b
	return this
}

type Color_r32f struct {
	r float64
}

func (this Color_r32f) to_rgba32f() Color_rgba32f {
	return Color_rgba32f{this.r, 0.0, 0.0, 0.0}
}

func (this *Color_r32f) assign( rhs *Color_rgba32f ) *Color_r32f {

	this.r = rhs.r
	return this

}

type Color_rg32f struct{
	r, g float64
}


func (this *Color_rg32f) color_rg32f( color []float64){
	this.r = color[0]
	this.g = color[1]
}
//color_rg32f(comp_t r, comp_t g):r(r), g(g){}




func (this *Color_rg32f) to_rgba32f() Color_rgba32f{
	return Color_rgba32f{this.r, this.g, 0.0, 0.0}
}

func (this *Color_rg32f) assign(rhs *Color_rgba32f) *Color_rg32f{
	this.r = rhs.r
	this.g = rhs.g
	return this
}

type Color_rgba8 struct {
	r, g, b, a uint8
}
func (this *Color_rgba8) color_rgba8( color []uint8){
	this.r=color[0]
	this.g=color[1]
	this.b=color[2]
	this.a=color[3]
}

func (this *Color_rgba8)to_rgba32f() Color_rgba32f{
	const inv_255 = 1.0 / 255
	return Color_rgba32f{float64(this.r) * inv_255, float64(this.g) * inv_255, float64(this.b) * inv_255, float64(this.a) * inv_255}
}

func (this *Color_rgba8) assign(rhs *Color_rgba32f)Color_rgba8{
	//#ifndef EFLIB_NO_SIMD
	//const __m128 f255 = _mm_set_ps1(255.0f)
	//__m128 m4 = _mm_loadu_ps(&rhs.r)
	//m4 = _mm_mul_ps(m4, f255)
	//m4 = _mm_max_ps(m4, _mm_setzero_ps())
	//m4 = _mm_min_ps(m4, f255)
	//__m128i mi4 = _mm_cvtps_epi32(m4)
	//mi4 = _mm_or_si128(mi4, _mm_srli_si128(mi4, 3))
	//mi4 = _mm_or_si128(mi4, _mm_srli_si128(mi4, 6))
	//*reinterpret_cast<int*>(&r) = _mm_cvtsi128_si32(mi4)
	//#else
	this.r = uint8( common.Clamp(rhs.r * 255.0, 0.0, 255.0) )
	this.g = uint8( common.Clamp(rhs.g * 255.0, 0.0, 255.0) )
	this.b = uint8( common.Clamp(rhs.b * 255.0, 0.0, 255.0) )
	this.a = uint8( common.Clamp(rhs.a * 255.0, 0.0, 255.0) )
	//#endif

	return *this
}

type Color_bgra8 struct {
	b, g, r, a uint8
}

func (this *Color_bgra8) color_bgra8(color []uint8 ) {
	this.b=color[0]
	this.g=color[1]
	this.r=color[2]
	this.a=color[3]
}



func (this *Color_bgra8)  to_rgba32f() Color_rgba32f{
	const  inv_255 = 1.0 / 255
	return Color_rgba32f{float64(this.r) * inv_255, float64(this.g) * inv_255, float64(this.b) * inv_255, float64(this.a) * inv_255}
}

func (this *Color_bgra8) assign(rhs *Color_rgba32f)Color_bgra8{
	//#ifndef EFLIB_NO_SIMD
	//const __m128 f255 = _mm_set_ps1(255.0f)
	//__m128 m4 = _mm_loadu_ps(&rhs.r)
	//m4 = _mm_shuffle_ps(m4, m4, _MM_SHUFFLE(3, 0, 1, 2))
	//m4 = _mm_mul_ps(m4, f255)
	//m4 = _mm_max_ps(m4, _mm_setzero_ps())
	//m4 = _mm_min_ps(m4, f255)
	//__m128i mi4 = _mm_cvtps_epi32(m4)
	//mi4 = _mm_or_si128(mi4, _mm_srli_si128(mi4, 3))
	//mi4 = _mm_or_si128(mi4, _mm_srli_si128(mi4, 6))
	//*reinterpret_cast<int*>(&b) = _mm_cvtsi128_si32(mi4)
	//#else
	this.r = uint8( common.Clamp(rhs.r * 255.0 + 0.5, 0.0, 255.0))
	this.g = uint8( common.Clamp(rhs.g * 255.0 + 0.5, 0.0, 255.0))
	this.b = uint8( common.Clamp(rhs.b * 255.0 + 0.5, 0.0, 255.0))
	this.a = uint8( common.Clamp(rhs.a * 255.0 + 0.5, 0.0, 255.0))

	return *this
}
type  comp_t32 int32
type Color_r32i struct {
	r int32
}

func (this *Color_r32i) to_rgba32f() Color_rgba32f{
	return Color_rgba32f{float64(this.r), 0.0, 0.0, 0.0}
}

func (this *Color_r32i) assign(rhs *Color_rgba32f)*Color_r32i{
	this.r = int32( rhs.r + 0.5 )
	return this
}

func  lerp_rgba32f(c0 Color_rgba32f,c1 Color_rgba32f, t float64) Color_rgba32f{
	//#ifndef EFLIB_NO_SIMD
	//__m128 mc0 = _mm_loadu_ps(&c0.r)
	//__m128 mc1 = _mm_loadu_ps(&c1.r)
	//__m128 mret = _mm_add_ps(mc0, _mm_mul_ps(_mm_sub_ps(mc1, mc0), _mm_set1_ps(t)))
	//Color_rgba32f ret
	//_mm_storeu_ps(&ret.r, mret)
	//return ret
	//#else
	return vec4_to_color_rgba32f(smath.V_add(c0.get_vec4(),smath.V_sub(c1.get_vec4(),c0.get_vec4()).Scale(t)))
	//#endif
}
func  lerp_rgb32(c0 Color_rgb32f,c1 Color_rgb32f, t float64) Color_rgba32f{
	return Color_rgb32f{c0.r + (c1.r - c0.r) * t, c0.r + (c1.g - c0.g) * t, c0.r + (c1.b - c0.b) * t}.to_rgba32f()
}
func  lerp_2bgra8(c0 Color_bgra8, c1 Color_bgra8, t float64) Color_rgba32f{
	ret := lerp_rgba32f(Color_rgba32f{float64(c0.r), float64(c0.g), float64(c0.b), float64(c0.a)},
		Color_rgba32f{float64(c1.r), float64(c1.g), float64(c1.b), float64(c1.a)}, t)
	ret.get_vec4().Divide(255.0)
	return ret
}

func  lerp_r(c0 *Color_r32f, c1 *Color_r32f, t float64) Color_rgba32f{
	return Color_r32f{c0.r + (c1.r - c0.r) * t}.to_rgba32f()
}
func  lerp_rg(c0 *Color_rg32f, c1 *Color_rg32f, t float64)Color_rgba32f{
	return Color_rg32f{c0.r + (c1.r - c0.r) * t, c0.r + (c1.g - c0.g) * t}.to_rgba32f()
}
func  lerp_ri(c0 *Color_r32i, c1 *Color_r32i, t float64)Color_rgba32f{
	return Color_r32i{c0.r + int32(float64(c1.r - c0.r) * t)}.to_rgba32f()
}

func  lerp_vector(c0 , c1 , c2 , c3 *Color_rgba32f, tx, ty float64)Color_rgba32f{
	//#ifndef EFLIB_NO_SIMD
	//__m128 mc0 = _mm_loadu_ps(&c0.r)
	//__m128 mc1 = _mm_loadu_ps(&c1.r)
	//__m128 mc2 = _mm_loadu_ps(&c2.r)
	//__m128 mc3 = _mm_loadu_ps(&c3.r)
	//__m128 mc01 = _mm_add_ps(mc0, _mm_mul_ps(_mm_sub_ps(mc1, mc0), _mm_set1_ps(tx)))
	//__m128 mc23 = _mm_add_ps(mc2, _mm_mul_ps(_mm_sub_ps(mc3, mc2), _mm_set1_ps(tx)))
	//__m128 mret = _mm_add_ps(mc01, _mm_mul_ps(_mm_sub_ps(mc23, mc01), _mm_set1_ps(ty)))
	//Color_rgba32f ret
	//_mm_storeu_ps(&ret.r, mret)
	//return ret
	//#else
	c01 := vec4_to_color_rgba32f(smath.V_add(c0.get_vec4() , smath.V_sub(c1.get_vec4(), c0.get_vec4()).Scale(tx)))
	c23 := vec4_to_color_rgba32f(smath.V_add(c2.get_vec4() , smath.V_sub(c3.get_vec4(), c2.get_vec4()).Scale(tx)))
	return vec4_to_color_rgba32f(smath.V_add(c01.get_vec4(), smath.V_sub(c23.get_vec4(),c01.get_vec4()).Scale(ty)))
	//#endif
}
func  lerp_rgb32f(c0 ,  c1,  c2,  c3 *Color_rgb32f, tx, ty float64)Color_rgba32f{
	c01 := Color_rgb32f{c0.r + (c1.r - c0.r) * tx, c0.r + (c1.g - c0.g) * tx, c0.r + (c1.b - c0.b) * tx}
	c23 := Color_rgb32f{c2.r + (c3.r - c2.r) * tx, c2.r + (c3.g - c2.g) * tx, c2.r + (c3.b - c2.b) * tx}
	return Color_rgb32f{c01.r + (c23.r - c01.r) * ty, c01.r + (c23.g - c01.g) * ty, c01.r + (c23.b - c01.b) * ty}.to_rgba32f()
}
func  lerp_4bgra8( c0,c1,  c2,  c3 *Color_bgra8, tx, ty float64)Color_rgba32f{
	//#ifndef EFLIB_NO_SIMD
	//__m128i mzero = _mm_setzero_si128()
	//__m128i mci = _mm_cvtsi32_si128(*reinterpret_cast<const int*>(&c0.r))
	//mci = _mm_unpacklo_epi8(mci, mzero)
	//mci = _mm_unpacklo_epi16(mci, mzero)
	//__m128 mc0 = _mm_cvtepi32_ps(_mm_shuffle_epi32(mci, _MM_SHUFFLE(3, 0, 1, 2)))
	//mci = _mm_cvtsi32_si128(*reinterpret_cast<const int*>(&c1.r))
	//mci = _mm_unpacklo_epi8(mci, mzero)
	//mci = _mm_unpacklo_epi16(mci, mzero)
	//__m128 mc1 = _mm_cvtepi32_ps(_mm_shuffle_epi32(mci, _MM_SHUFFLE(3, 0, 1, 2)))
	//mci = _mm_cvtsi32_si128(*reinterpret_cast<const int*>(&c2.r))
	//mci = _mm_unpacklo_epi8(mci, mzero)
	//mci = _mm_unpacklo_epi16(mci, mzero)
	//__m128 mc2 = _mm_cvtepi32_ps(_mm_shuffle_epi32(mci, _MM_SHUFFLE(3, 0, 1, 2)))
	//mci = _mm_cvtsi32_si128(*reinterpret_cast<const int*>(&c3.r))
	//mci = _mm_unpacklo_epi8(mci, mzero)
	//mci = _mm_unpacklo_epi16(mci, mzero)
	//__m128 mc3 = _mm_cvtepi32_ps(_mm_shuffle_epi32(mci, _MM_SHUFFLE(3, 0, 1, 2)))
	//
	//__m128 mc01 = _mm_add_ps(mc0, _mm_mul_ps(_mm_sub_ps(mc1, mc0), _mm_set1_ps(tx)))
	//__m128 mc23 = _mm_add_ps(mc2, _mm_mul_ps(_mm_sub_ps(mc3, mc2), _mm_set1_ps(tx)))
	//__m128 mret = _mm_add_ps(mc01, _mm_mul_ps(_mm_sub_ps(mc23, mc01), _mm_set1_ps(ty)))
	//mret = _mm_mul_ps(mret, _mm_set1_ps(1.0f / 255))
	//Color_rgba32f ret
	//_mm_storeu_ps(&ret.r, mret)
	//return ret
	//#else
	ret := lerp_vector(&Color_rgba32f{float64(c0.r), float64(c0.g), float64(c0.b), float64(c0.a)},
		&Color_rgba32f{float64(c1.r), float64(c1.g), float64(c1.b), float64(c1.a)},
		&Color_rgba32f{float64(c2.r), float64(c2.g), float64(c2.b), float64(c2.a)},
		&Color_rgba32f{float64(c3.r), float64(c3.g), float64(c3.b), float64(c3.a)}, tx, ty)
	ret.get_vec4().Divide(255.0)
	return ret
	//#endif
}
//func  lerp( c0, c1, c2, c3 *Color_rgba8, tx, ty float64)Color_rgba32f{
//	//#ifndef EFLIB_NO_SIMD
//	//__m128i mzero = _mm_setzero_si128()
//	//__m128i mci = _mm_cvtsi32_si128(*reinterpret_cast<const int*>(&c0.r))
//	//mci = _mm_unpacklo_epi8(mci, mzero)
//	//__m128 mc0 = _mm_cvtepi32_ps(_mm_unpacklo_epi16(mci, mzero))
//	//mci = _mm_cvtsi32_si128(*reinterpret_cast<const int*>(&c1.r))
//	//mci = _mm_unpacklo_epi8(mci, mzero)
//	//__m128 mc1 = _mm_cvtepi32_ps(_mm_unpacklo_epi16(mci, mzero))
//	//mci = _mm_cvtsi32_si128(*reinterpret_cast<const int*>(&c2.r))
//	//mci = _mm_unpacklo_epi8(mci, mzero)
//	//__m128 mc2 = _mm_cvtepi32_ps(_mm_unpacklo_epi16(mci, mzero))
//	//mci = _mm_cvtsi32_si128(*reinterpret_cast<const int*>(&c3.r))
//	//mci = _mm_unpacklo_epi8(mci, mzero)
//	//__m128 mc3 = _mm_cvtepi32_ps(_mm_unpacklo_epi16(mci, mzero))
//	//
//	//__m128 mc01 = _mm_add_ps(mc0, _mm_mul_ps(_mm_sub_ps(mc1, mc0), _mm_set1_ps(tx)))
//	//__m128 mc23 = _mm_add_ps(mc2, _mm_mul_ps(_mm_sub_ps(mc3, mc2), _mm_set1_ps(tx)))
//	//__m128 mret = _mm_add_ps(mc01, _mm_mul_ps(_mm_sub_ps(mc23, mc01), _mm_set1_ps(ty)))
//	//mret = _mm_mul_ps(mret, _mm_set1_ps(1.0f / 255))
//	//Color_rgba32f ret
//	//_mm_storeu_ps(&ret.r, mret)
//	//return ret
//	//#else
//	ret := lerp(Color_rgba32f{c0.r, c0.g, c0.b, c0.a}, Color_rgba32f{c1.r, c1.g, c1.b, c1.a},
//	Color_rgba32f{c2.r, c2.g, c2.b, c2.a}, Color_rgba32f{c3.r, c3.g, c3.b, c3.a}, tx, ty)
//	ret.get_vec4() /= 255.0
//	return ret
//	//#endif
//}
func lerp_r32f(c0,c1,c2, c3 *Color_r32f, tx, ty float64)Color_rgba32f{
	c01 := Color_r32f{c0.r + (c1.r - c0.r) * tx}
	c23 := Color_r32f{c2.r + (c3.r - c2.r) * tx}
	return Color_r32f{c01.r + (c23.r - c01.r) * ty}.to_rgba32f()
}
func  lerp_rg32f(c0,c1, c2, c3 *Color_rg32f, tx, ty float64)Color_rgba32f{
	c01 := Color_rg32f{c0.r + (c1.r - c0.r) * tx, c0.r + (c1.g - c0.g) * tx}
	c23 := Color_rg32f{c2.r + (c3.r - c2.r) * tx, c2.r + (c3.g - c2.g) * tx}
	return Color_rg32f{c01.r + (c23.r - c01.r) * ty, c01.r + (c23.g - c01.g) * ty}.to_rgba32f()
}
func  lerp_r32i(c0,c1,c2, c3 *Color_r32i, tx, ty float64) Color_rgba32f{
	c01 := Color_r32f{float64(c0.r + (c1.r - c0.r)) * tx}
	c23 := Color_r32f{float64(c2.r + (c3.r - c2.r)) * tx}
	return Color_r32f{float64(c01.r + (c23.r - c01.r)) * ty}.to_rgba32f()
}




type  Color_max struct {

}

type Pixel_type_to_fmt struct {

}

type Pixel_fmt_to_type struct {

}

//func decl_type_fmt_pair(color_type, fmt_code){

//const pixel_format pixel_format_##color_type = fmt_code\
//template<>\
//struct pixel_type_to_fmt< color_type >\
//{ static const pixel_format fmt = fmt_code }\
//\
//template<>\
//struct pixel_fmt_to_type< fmt_code >\
//{ typedef color_type type }

type pixel_format int

// ------------------------------------------
// enumeration in compiling time and translation between different pixel format

type Pixel_information struct {
 	Size int
	Describe [32]byte
}

func decl_color_info(pixel_type interface{}) Pixel_information{
	return Pixel_information{
		len(pixel_type),pixel_type
	}
}

//decl_type_fmt_pair(color_rgba32f, 0)
//decl_type_fmt_pair(color_rgb32f, 1)
//decl_type_fmt_pair(color_bgra8, 2)
//decl_type_fmt_pair(color_rgba8, 3)
//decl_type_fmt_pair(color_r32f, 4)
//decl_type_fmt_pair(color_rg32f, 5)
//decl_type_fmt_pair(color_r32i, 6)
//decl_type_fmt_pair(color_max, 7)

const pixel_format_color_ub = pixel_format_color_max - 1
const Pixel_format_invalid = -1

// Pixel format informations
var a = []int{1,2}

var color_infos  = [8]*Pixel_information{
	decl_color_info(Color_rgba32f),
	decl_color_info(Color_rgb32f),
	decl_color_info(Color_bgra8),
	decl_color_info(Color_rgba8),
	decl_color_info(Color_r32f),
	decl_color_info(Color_rg32f),
	decl_color_info(Color_r32i)
}

func get_color_info( pf pixel_format )*Pixel_information{
	return color_infos[pf]
}

type  Pixel_format_convertor struct {
	//convertors[outfmt][infmt]
	convertors [pixel_format_color_max][pixel_format_color_max]Pixel_convertor
	array_convertors [pixel_format_color_max][pixel_format_color_max]Pixel_array_convertor
 	lerpers_1d [pixel_format_color_max]pixel_lerp_1d
	lerpers_2d [pixel_format_color_max]pixel_lerp_2d
}
	//template <int outColor, int inColor> friend struct color_convertor_initializer
func (pfc *Pixel_format_convertor)Convert(outfmt , infmt Pixel_format , outpixel , inpixel interface{}){
	(pfc.convertors[outfmt][infmt])(outpixel, inpixel)
}

func (pfc *Pixel_format_convertor)convert_array( outfmt,  infmt Pixel_format, outpixel, inpixel interface{},
		count,  outstride ,  instride int){
 	if outstride == 0 {
		outstride = color_infos[outfmt].size
	}
	if  instride == 0 {
		instride = color_infos[infmt].size
	}
	pfc.array_convertors[outfmt][infmt](outpixel, inpixel, count, outstride, instride)
}

type pixel_convertor  func(outcolor,incolor interface{})
type pixel_array_convertor  func(outcolor, incolor interface{}, count, outstride, instride int)
type pixel_lerp_1d  func(incolor0, incolor1 interface{}, t float64)Color_rgba32f
type pixel_lerp_2d  func( incolor0, incolor1, incolor2, incolor3 interface{}, tx, ty float64)Color_rgba32f

func (pfc *Pixel_format_convertor)get_convertor_func(outfmt,infmt Pixel_format)pixel_convertor {
	return pfc.convertors[outfmt][infmt]
}
func (pfc *Pixel_format_convertor) get_array_convertor_func( outfmt, infmt Pixel_format)pixel_array_convertor {
	return pfc.array_convertors[outfmt][infmt]
}
func (pfc *Pixel_format_convertor)get_lerp_1d_func(infmt Pixel_format)pixel_lerp_1d{
	return pfc.lerpers_1d[infmt]
}
func (pfc *Pixel_format_convertor)  get_lerp_2d_func(infmt Pixel_format)pixel_lerp_2d{
	return pfc.lerpers_2d[infmt]
}

func (pfc *Pixel_format_convertor) pixel_format_convertor(){
	//color_convertor_initializer<pixel_format_color_max - 1, pixel_format_color_max - 1>
		init_(pfc.convertors[0][0], pfc.array_convertors[0][0], pfc.lerpers_1d[0], pfc.lerpers_2d[0])
}
