package render
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
