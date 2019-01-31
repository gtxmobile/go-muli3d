package render
import (
	"../common"
	"../smath"
	"math"
)

type  Sampler_desc struct{
	min_filter 		common.Filter_type
	mag_filter 		common.Filter_type
	mip_filter 		common.Filter_type
 	addr_mode_u 	common.Address_mode
 	addr_mode_v 	common.Address_mode
 	addr_mode_w 	common.Address_mode
	mip_lod_bias	float64
	max_anisotropy	uint32
 	comparison_func	common.Compare_function
	border_color	Color_rgba32f
	min_lod			float64
	max_lod			float64
}
func New_sampler_desc(){
	sampler_desc := Sampler_desc{}
	sampler_desc.min_filter = common.Filter_point
	sampler_desc.mag_filter = common.Filter_point
	sampler_desc.mip_filter = common.Filter_point
	sampler_desc.addr_mode_u= common.Address_wrap
	sampler_desc.addr_mode_v= common.Address_wrap
	sampler_desc.addr_mode_w= common.Address_wrap
	sampler_desc.mip_lod_bias = 0
	sampler_desc.max_anisotropy = 0
	sampler_desc.comparison_func = common.Compare_function_always
	sampler_desc.border_color = Color_rgba32f{0.0, 0.0, 0.0, 0.0}
	sampler_desc.min_lod = -1e20
	sampler_desc.max_lod = 1e20

}


type Filter_op_type func(surf *Surface, x, y float64, sample uint32, border_color *Color_rgba32f)Color_rgba32f
type Sampler struct {
	Desc_ 	Sampler_desc
	Tex_	*Texture
	Filters_  [common.Sampler_state_count]Filter_op_type
}


func New_sampler(desc *Sampler_desc, tex *Texture){}
func Sample( coordx,  coordy,  miplevel float64) Color_rgba32f{return Color_rgba32f{}}
func Sample_2d_lod( proj_coord smath.Vec_2, lod float64 ) Color_rgba32f{return Color_rgba32f{}}
func Sample_2d_grad(proj_coord, ddx, ddy smath.Vec_2,  lod_bias float64) Color_rgba32f{return Color_rgba32f{}}
func Sample_2d_proj(proj_coord,ddx,  ddy smath.Vector_t) Color_rgba32f{return Color_rgba32f{}}
func Sample_cube( coordx,  coordy,  coordz, miplevel float64) Color_rgba32f{return Color_rgba32f{}}


func Calc_lod( size,  ddx, ddy smath.Vector_t,  bias float64 ) float64{return 0}
func Calc_lod_2d(ddx ,ddy smath.Vec_2) float64{return 0}
func Calc_anisotropic_lod(size, ddx,  ddy smath.Vector_t,
	bias float64,out_lod, out_ratio *float64, out_long_axis smath.Vector_t) float64{return 0}

func sample_surface( surf *Surface, x, y float64, sample uint32,ss common.Sampler_state) Color_rgba32f{return Color_rgba32f{}}
//template <bool IsCubeTexture>
func sample_impl(face int32,  coordx, coordy float64, sample uint32,  miplevel, ratio float64,
	long_axis smath.Vector_t) Color_rgba32f{return Color_rgba32f{}}

//addresser

type Wrap struct {

}
func (Wrap)Do_coordf(coord float64, size int32) float64{
	return (coord - math.Floor(coord)) * float64(size) - 0.5
}

func (Wrap)Do_coordi_point_1d(coord, size int) int{
	return (size * 8192 + coord) % size
}

func (Wrap)Do_coordi_point_2d(coord smath.Vector_t, size [4]int) [4]int{

	sub_v := smath.V_sub(coord,smath.Vector_t{math.Floor(coord.X), math.Floor(coord.Y), math.Floor(coord.Z), math.Floor(coord.W)})
	cross_b := smath.Vector_t{float64(size[0]), float64(size[1]), float64(size[2]), float64(size[3])}
	o_coord := smath.Vector_crossproduct(sub_v,cross_b)
	coord_ipart := [4]int{ int(math.Floor(o_coord.X)), int(math.Floor(o_coord.Y)), 0, 0}

	return [4]int{(size[0] * 8192 + coord_ipart[0]) % size[0], (size[1] * 8192 + coord_ipart[1]) % size[1],0, 0}
}

func (Wrap)Do_coordi_linear_2d( low, up [4]int,frac smath.Vector_t, coord smath.Vector_t,size [4]int){

	sub_v := smath.V_sub(coord , smath.Vector_t{math.Floor(coord.X), math.Floor(coord.Y), math.Floor(coord.Z), math.Floor(coord.W)})
	cross_b := smath.Vector_t{float64(size[0]), float64(size[1]), float64(size[2]), float64(size[3])}

	o_coord := smath.Vector_crossproduct(sub_v,cross_b)
	//o_coord - 0.5f
	o_coord.Sub_d(0.5)
	coord_ipart := [4]int{ int(math.Floor(o_coord.X)), int(math.Floor(o_coord.Y)), 0, 0}
	frac = smath.V_sub(o_coord , smath.Vector_t{float64(coord_ipart[0]), float64(coord_ipart[1]), 0, 0})

	low = [4]int{(size[0] * 8192 + coord_ipart[0]) % size[0], (size[1] * 8192 + coord_ipart[1]) % size[1], 0, 0}
	up = [4]int{(size[0] * 8192 + coord_ipart[0] + 1) % size[0],(size[1] * 8192 + coord_ipart[1] + 1) % size[1],0, 0}

}

type Mirror struct {
}
func (Mirror)Do_coordf(coord float64, size int32) float64{
	var ret float64
	selection_coord := math.Floor(coord)
	if int(selection_coord) & 1 >0{
		ret	= 1 + selection_coord - coord
	}else{
		ret = coord - selection_coord
	}
	return ret * float64(size) - 0.5
}

func  (Mirror)Do_coordi_point_1d( coord,  size int)int32{
	return int32(common.Clamp(float64(coord), 0, float64(size - 1)))
}
func  (Mirror)Do_coordi_point_2d(coord smath.Vector_t, size [4]int) [4]int{
	 selection_coord_x := math.Floor(coord.X)
	 selection_coord_y := math.Floor(coord.Y)
	 var vx,vy float64
	 if int(selection_coord_x) & 1 > 0{
		 vx = 1 + selection_coord_x - coord.X
	 }else{
		 vx = coord.X - selection_coord_x
	 }
	if int(selection_coord_y) & 1 > 0{
		vy = 1 + selection_coord_y - coord.Y
	}else{
		vy = coord.Y - selection_coord_y
	}
	 o_coord := smath.Vector_t{vx*float64(size[0]),vy* float64(size[1]),0, 0}


	coord_ipart := [4]int{int(math.Floor(o_coord.X)), int(math.Floor(o_coord.Y)), 0, 0}

	return [4]int{int(common.Clamp(float64(coord_ipart[0]), 0, float64(size[0] - 1))),
		int(common.Clamp(float64(coord_ipart[1]), 0, float64(size[1] - 1))),	0, 0}
}

func (Mirror)Do_coordi_linear_2d(low,up *[4]int,frac,coord smath.Vector_t, size [4]int){

	selection_coord_x := math.Floor(coord.X)
	selection_coord_y := math.Floor(coord.Y)
	var vx,vy float64
	if int(selection_coord_x) & 1 > 0{
		vx = 1 + selection_coord_x - coord.X
	}else{
		vx = coord.X - selection_coord_x
	}
	if int(selection_coord_y) & 1 > 0{
		vy = 1 + selection_coord_y - coord.Y
	}else{
		vy = coord.Y - selection_coord_y
	}
	o_coord := smath.Vector_t{vx*float64(size[0])-0.5,vy* float64(size[1])-0.5,0, 0}

	coord_ipart := [4]int{int(math.Floor(o_coord.X)), int(math.Floor(o_coord.Y)), 0, 0}

	frac = smath.V_sub(o_coord, smath.Vector_t{float64(coord_ipart[0]),float64(coord_ipart[1]), 0, 0})

	*low = [4]int{int(common.Clamp(float64(coord_ipart[0]), 0, float64(size[0] - 1))),
		int(common.Clamp(float64(coord_ipart[1]), 0, float64(size[1] - 1))),	0, 0}

	*up = [4]int{int(common.Clamp(float64(coord_ipart[0]+1), 0, float64(size[0] - 1))),
		int(common.Clamp(float64(coord_ipart[1]+1), 0, float64(size[1] - 1))),	0, 0}
}

type Clamp struct {

}
func (Clamp)Do_coordf(coord float64, size float64) float64{
	return common.Clamp(coord * size, 0.5, size - 0.5) - 0.5
}

func (Clamp)Do_coordi_point_1d(coord, size float64) float64{
	return common.Clamp(coord, 0, size - 1)
}
func (Clamp)do_coordi_point_2d( coord *smath.Vector_t,  size *[4]int) [4]int{

 	o_coord := smath.Vector_t{common.Clamp(coord.X * float64(size[0]), 0.5, float64(size[0]) - 0.5),
		common.Clamp(coord.Y * float64(size[1]), 0.5, float64(size[1]) - 0.5),0, 0}
 	coord_ipart := [4]int{int(math.Floor(o_coord.X)), int(math.Floor(o_coord.Y)), 0, 0}

	return [4]int{int(common.Clamp(float64(coord_ipart[0]), 0, float64(size[0] - 1))),
		int(common.Clamp(float64(coord_ipart[1]), 0, float64(size[1] - 1))),0, 0}

}

func (Clamp)Do_coordi_linear_2d( low, up *[4]int, frac, coord smath.Vector_t,  size [4]int){

	selection_coord_x := math.Floor(coord.X)
	selection_coord_y := math.Floor(coord.Y)
	var vx,vy float64
	if int(selection_coord_x) & 1 > 0{
		vx = 1 + selection_coord_x - coord.X
	}else{
		vx = coord.X - selection_coord_x
	}
	if int(selection_coord_y) & 1 > 0{
		vy = 1 + selection_coord_y - coord.Y
	}else{
		vy = coord.Y - selection_coord_y
	}
	o_coord := smath.Vector_t{vx*float64(size[0])-0.5,vy* float64(size[1])-0.5,0, 0}

	coord_ipart := [4]int{int(math.Floor(o_coord.X)), int(math.Floor(o_coord.Y)), 0, 0}

	frac = smath.V_sub(o_coord, smath.Vector_t{float64(coord_ipart[0]),float64(coord_ipart[1]), 0, 0})

	*low = [4]int{int(common.Clamp(float64(coord_ipart[0]), 0, float64(size[0] - 1))),
		int(common.Clamp(float64(coord_ipart[1]), 0, float64(size[1] - 1))),	0, 0}

	*up = [4]int{int(common.Clamp(float64(coord_ipart[0]+1), 0, float64(size[0] - 1))),
		int(common.Clamp(float64(coord_ipart[1]+1), 0, float64(size[1] - 1))),	0, 0}

}

type Border struct {

}

func (Border) Do_coordf(coord float64, size int) float64{
	return common.Clamp(coord * float64(size), -0.5, float64(size) + 0.5) - 0.5
}

func (Border) Do_coordi_point_1d(coord, size int) int{
	if coord >= size {
		return -1
	}
	return  coord
}
func (Border) Do_coordi_point_2d(coord smath.Vector_t,  size [4]int) [4]int{




	o_coord := smath.Vector_t{common.Clamp(coord.X * float64(size[0]), 0.5, float64(size[0]) - 0.5),
		common.Clamp(coord.Y * float64(size[1]), 0.5, float64(size[1]) - 0.5),0, 0}
	coord_ipart := [4]int{int(math.Floor(o_coord.X)), int(math.Floor(o_coord.Y)), 0, 0}

	return [4]int{int(common.Clamp(float64(coord_ipart[0]), 0, float64(size[0] - 1))),
		int(common.Clamp(float64(coord_ipart[1]), 0, float64(size[1] - 1))),0, 0}
}

func (Border) Do_coordi_linear_2d(low, up *[4]int, frac, coord smath.Vector_t,  size [4]int){

	selection_coord_x := math.Floor(coord.X)
	selection_coord_y := math.Floor(coord.Y)
	var vx,vy float64
	if int(selection_coord_x) & 1 > 0{
		vx = 1 + selection_coord_x - coord.X
	}else{
		vx = coord.X - selection_coord_x
	}
	if int(selection_coord_y) & 1 > 0{
		vy = 1 + selection_coord_y - coord.Y
	}else{
		vy = coord.Y - selection_coord_y
	}
	o_coord := smath.Vector_t{vx*float64(size[0])-0.5,vy* float64(size[1])-0.5,0, 0}

	coord_ipart := [4]int{int(math.Floor(o_coord.X)), int(math.Floor(o_coord.Y)), 0, 0}

	frac = smath.V_sub(o_coord, smath.Vector_t{float64(coord_ipart[0]),float64(coord_ipart[1]), 0, 0})

	*low = [4]int{int(common.Clamp(float64(coord_ipart[0]), 0, float64(size[0] - 1))),
		int(common.Clamp(float64(coord_ipart[1]), 0, float64(size[1] - 1))),	0, 0}

	*up = [4]int{int(common.Clamp(float64(coord_ipart[0]+1), 0, float64(size[0] - 1))),
		int(common.Clamp(float64(coord_ipart[1]+1), 0, float64(size[1] - 1))),	0, 0}
}