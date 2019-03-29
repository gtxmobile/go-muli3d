package buffer
import (
	"../../common"
	"../../smath"
	"math"
	"reflect"
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


func Calc_lod( size,  ddx, ddy [4]uint32,  bias float64 ) float64{return 0}
func Calc_lod_2d(ddx ,ddy smath.Vec_2) float64{return 0}
func Calc_anisotropic_lod(size, ddx,  ddy [4]uint32,
	bias float64,out_lod, out_ratio *float64, out_long_axis *smath.Vector_t) float64{return 0}

func sample_surface( surf *Surface, x, y float64, sample uint32,ss common.Sampler_state) Color_rgba32f{return Color_rgba32f{}}
//template <bool IsCubeTexture>
func sample_impl(IsCubeTexture bool,face int32,  coordx, coordy float64, sample uint32,  miplevel, ratio float64,
	long_axis smath.Vector_t) Color_rgba32f{return Color_rgba32f{}}

//addresser
//namespace coord_calculator
//{
//template <typename Addresser_type>
type Addresser_type interface {
	Do_coordf(coord float64, size int) float64
	Do_coordi_point_1d(coord, size int) int
	Do_coordi_point_2d(coord smath.Vector_t,  size [4]int) [4]int
	Do_coordi_linear_2d(low, up *[4]int, frac, coord smath.Vector_t,  size [4]int)
}
type Wrap struct {

}
func (Wrap)Do_coordf(coord float64, size int) float64{
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

func (Wrap)Do_coordi_linear_2d( low, up *[4]int,frac smath.Vector_t, coord smath.Vector_t,size [4]int){

	sub_v := smath.V_sub(coord , smath.Vector_t{math.Floor(coord.X), math.Floor(coord.Y), math.Floor(coord.Z), math.Floor(coord.W)})
	cross_b := smath.Vector_t{float64(size[0]), float64(size[1]), float64(size[2]), float64(size[3])}

	o_coord := smath.Vector_crossproduct(sub_v,cross_b)
	//o_coord - 0.5f
	o_coord.Sub_d(0.5)
	coord_ipart := [4]int{ int(math.Floor(o_coord.X)), int(math.Floor(o_coord.Y)), 0, 0}
	frac = smath.V_sub(o_coord , smath.Vector_t{float64(coord_ipart[0]), float64(coord_ipart[1]), 0, 0})

	*low = [4]int{(size[0] * 8192 + coord_ipart[0]) % size[0], (size[1] * 8192 + coord_ipart[1]) % size[1], 0, 0}
	*up = [4]int{(size[0] * 8192 + coord_ipart[0] + 1) % size[0],(size[1] * 8192 + coord_ipart[1] + 1) % size[1],0, 0}

}

type Mirror struct {
}
func (Mirror)Do_coordf(coord float64, size int) float64{
	var ret float64
	selection_coord := math.Floor(coord)
	if int(selection_coord) & 1 >0{
		ret	= 1 + selection_coord - coord
	}else{
		ret = coord - selection_coord
	}
	return ret * float64(size) - 0.5
}

func  (Mirror)Do_coordi_point_1d( coord,  size int)int{
	return int(common.Clamp(float64(coord), 0, float64(size - 1)))
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
func (Clamp)Do_coordf(coord float64, size int) float64{
	fsize := float64(size)
	return common.Clamp(coord * fsize, 0.5, fsize - 0.5) - 0.5
}

func (Clamp)Do_coordi_point_1d(coord, size int) int{
	return int(common.Clamp(float64(coord), 0, float64(size - 1)))
}
func (Clamp)Do_coordi_point_2d( coord smath.Vector_t,  size [4]int) [4]int{

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


func point_cc(addresser_type Addresser_type,coord float64, size int) int{
	o_coord := addresser_type.Do_coordf(coord, size)
	coord_ipart := int(math.Floor(o_coord + 0.5))
	return addresser_type.Do_coordi_point_1d(coord_ipart, size)
}

//template <typename Addresser_type>
func linear_cc(addresser_type Addresser_type,low, up *int, frac *float64,coord float64, size int){
	o_coord := addresser_type.Do_coordf(coord, size);
	coord_ipart := int(math.Floor(o_coord))
	*low = addresser_type.Do_coordi_point_1d(coord_ipart, size)
	*up = addresser_type.Do_coordi_point_1d(coord_ipart + 1, size)
	*frac = o_coord - float64(coord_ipart)
}

func point_cc_2(addresser_type Addresser_type, coord smath.Vector_t, size [4]int) [4]int{
	return addresser_type.Do_coordi_point_2d(coord, size)
}

//template <typename Addresser_type>
func linear_cc_2(addresser_type Addresser_type ,low,up *[4]int, frac smath.Vector_t, coord smath.Vector_t,  size [4]int){
	addresser_type.Do_coordi_linear_2d(low, up, frac, coord, size)
}

//namespace surface_sampler
//{
//template <typename addresser_type_u, typename addresser_type_v>
type Point struct {
	addresser_type_u, addresser_type_v Addresser_type
}
func (p Point)op( surf *Surface,  x, y float64, sample uint32,  border_color *Color_rgba32f) Color_rgba32f{
 	ix :=	point_cc(p.addresser_type_u,x, int(surf.Width()))
	iy :=   point_cc(p.addresser_type_v,y, int(surf.Height()))
	if(ix < 0 || iy < 0) {
		return *border_color
	}
	return surf.Get_texel(uint32(ix),uint32(iy),sample)
}

//template <typename addresser_type_u, typename addresser_type_v>
type Linear struct {
	addresser_type_u, addresser_type_v Addresser_type
}
func (l Linear)op(surf *Surface,  x, y float64, sample uint32,  border_color *Color_rgba32f) Color_rgba32f{

	var xpos0, ypos0, xpos1, ypos1 int
	var tx, ty *float64
	linear_cc(l.addresser_type_u,&xpos0, &xpos1, tx, x, int(surf.Width()))
	linear_cc(l.addresser_type_v,&ypos0, &ypos1, ty, y, int(surf.Height()));

	return surf.Get_texel_6(uint32(xpos0), uint32(ypos0), uint32(xpos1), uint32(ypos1), *tx, *ty, sample)
}


//template <typename addresser_type_uv>
func (p Point) op_(surf *Surface,  x, y , sample uint32,  border_color *Color_rgba32f)Color_rgba32f{
	region_size := [4]int{int(surf.Width()),int(surf.Height()), 0, 0}
	ixy := point_cc_2(p.addresser_type_u,smath.Vector_t{float64(x), float64(y), 0, 0}, region_size)
	if 0 <= ixy[0] && ixy[0] < region_size[0] && 0 <= ixy[1] && ixy[1] < region_size[1] {
		return surf.Get_texel(uint32(ixy[0]), uint32(ixy[1]), sample)
	}
	return *border_color;
}

//template <typename addresser_type_uv>
//struct linear<addresser_type_uv, addresser_type_uv>

func (l Linear)op_(surf *Surface, x, y float64, sample uint32, border_color *Color_rgba32f )Color_rgba32f{
	var pos0, pos1 [4]int
	var t smath.Vector_t
	linear_cc_2(l.addresser_type_u,&pos0, &pos1, t, smath.Vector_t{x, y, 0, 0},
	[4]int{int(surf.Width()), int(surf.Height()), 0, 0})

// printf("%d, %d\n", pos0[0], pos0[1]);

	return surf.Get_texel_6(uint32(pos0[0]), uint32(pos0[1]), uint32(pos1[0]), uint32(pos1[1]), t.X, t.Y, sample);
}


var filter_table  = [common.Filter_type_count][common.Address_mode_count][common.Address_mode_count]Filter_op_type{
	{
		{
			Point{Wrap{}, Wrap{}}.op,
			Point{Wrap{}, Mirror{}}.op,
			Point{Wrap{}, Clamp{}}.op,
			Point{Wrap{}, Border{}}.op},
		{
			Point{Mirror{}, Wrap{}}.op,
			Point{Mirror{}, Mirror{}}.op,
			Point{Mirror{}, Clamp{}}.op,
			Point{Mirror{}, Border{}}.op},
		{
			Point{Clamp{}, Wrap{}}.op,
			Point{Clamp{}, Mirror{}}.op,
			Point{Clamp{}, Clamp{}}.op,
			Point{Clamp{}, Border{}}.op	},
		{
			Point{Border{}, Wrap{}}.op,
			Point{Border{}, Mirror{}}.op,
			Point{Border{}, Clamp{}}.op,
			Point{Border{},  Border{}}.op}},
	{
		{
			Linear{Wrap{},Wrap{}}.op,
			Linear{Wrap{},Mirror{}}.op,
			Linear{Wrap{},Clamp{}}.op,
			Linear{Wrap{}, Border{}}.op	},
		{
			Linear{Mirror{}, Wrap{}}.op,
			Linear{Mirror{}, Mirror{}}.op,
			Linear{Mirror{}, Clamp{}}.op,
			Linear{Mirror{},  Border{}}.op},
		{
			Linear{Clamp{},Wrap{}}.op,
			Linear{Clamp{},Mirror{}}.op,
			Linear{Clamp{},Clamp{}}.op,
			Linear{Clamp{}, Border{}}.op},
		{
			Linear{Border{},Wrap{}}.op,
			Linear{Border{},Mirror{}}.op,
			Linear{Border{},Clamp{}}.op,
			Linear{Border{},Border{}}.op}}}

 func (s Sampler)calc_lod( size *[4]uint32, ddx, ddy *smath.Vector_t, bias float64) float64{

	var rho, lambda float64

	size_vec4 := smath.Vector_t{float64(size[0]), float64(size[1]), float64(size[2]), 0 }

	ddx_ts := ddx.Multiply(size_vec4)
	ddy_ts := ddy.Multiply(size_vec4)

	ddx_rho := ddx_ts.Length();
	ddy_rho := ddy_ts.Length();

	rho = math.Max(ddx_rho, ddy_rho)

	if(rho == 0.0){ rho = 0.000001}
	lambda = math.Log2(rho)
	return lambda + bias

}

func (s Sampler)sample_surface(surf *Surface,x, y float64, sample uint32, ss common.Sampler_state) Color_rgba32f {
	return s.Filters_[ss](
		surf,
		x, y, sample,
		&s.Desc_.border_color)
}

func NewSampler(desc Sampler_desc,  tex *Texture) Sampler{
 	s  := Sampler{Desc_:desc,Tex_:tex}
	s.Filters_[common.Sampler_state_min] = filter_table[s.Desc_.min_filter][s.Desc_.addr_mode_u][s.Desc_.addr_mode_v];
	s.Filters_[common.Sampler_state_mag] = filter_table[s.Desc_.mag_filter][s.Desc_.addr_mode_u][s.Desc_.addr_mode_v];
	s.Filters_[common.Sampler_state_mip] = filter_table[s.Desc_.mip_filter][s.Desc_.addr_mode_u][s.Desc_.addr_mode_v];
	return s
}

func compute_cube_subresource(btype bool, face, lod_level uint32) uint32{
	if btype {
		return lod_level * 6 + face
	} else{
		return lod_level
	}
}

//func compute_cube_subresource( false_type bool,face, lod_level uint32)uint32{
//	return lod_level;
//}
const true_type = true
const false_type = false
//template <bool IsCubeTexture>
func (s Sampler)Sample_impl(IsCubeTexture bool, face int, coordx, coordy float64,  sample uint32, miplevel, ratio float64, long_axis smath.Vector_t) Color_rgba32f{
	//std::integral_constant<bool, IsCubeTexture> dummy;
	face_sz := uint32(face);
	var is_mag bool
	if s.Desc_.mip_filter == common.Filter_point {
		is_mag = miplevel < 0.5
	}else{
		is_mag = miplevel < 0
	}
	if(is_mag){
		subres_index := compute_cube_subresource(IsCubeTexture, face_sz, s.Tex_.Max_lod_);
		return sample_surface(s.Tex_.Subresource(subres_index), coordx, coordy, sample, common.Sampler_state_mag);
	}

	if(s.Desc_.mip_filter == common.Filter_point)	{
		ml := math.Floor(miplevel + 0.5);

		ml = common.Clamp(ml, float64(s.Tex_.Max_lod_),float64(s.Tex_.Min_lod_))

		subres_index := compute_cube_subresource(IsCubeTexture, face_sz, uint32(ml))
		return sample_surface(s.Tex_.Subresource(subres_index), coordx, coordy, sample, common.Sampler_state_min);
	}

	if(s.Desc_.mip_filter == common.Filter_linear){
		lo := math.Floor(miplevel);
		hi := lo + 1;

		frac := miplevel - lo;

		lo_sz := common.Clamp(lo, float64(s.Tex_.Max_lod_), float64(s.Tex_.Min_lod_));
		hi_sz := common.Clamp(hi, float64(s.Tex_.Max_lod_), float64(s.Tex_.Min_lod_));

		subres_index_lo := compute_cube_subresource(IsCubeTexture, face_sz, uint32(lo_sz));
		subres_index_hi := compute_cube_subresource(IsCubeTexture, face_sz, uint32(hi_sz));

	 	c0 := sample_surface(s.Tex_.Subresource(subres_index_lo), coordx, coordy, sample, common.Sampler_state_min);
	 	c1 := sample_surface(s.Tex_.Subresource(subres_index_hi), coordx, coordy, sample, common.Sampler_state_min);

		return Lerp(c0, c1, frac);
	}

	if(s.Desc_.mip_filter == common.Filter_anisotropic){
		start_relative_distance := - 0.5 * (ratio - 1.0);

		sample_coord_x := coordx + long_axis.X * start_relative_distance;
		sample_coord_y := coordy + long_axis.Y * start_relative_distance;

		lo := math.Floor(miplevel);
		hi := lo + 1;

		frac := miplevel - lo;
		lo = math.Max(lo, 0);
		hi = math.Max(hi, 0);

		lo_sz := common.Clamp(lo, float64(s.Tex_.Max_lod_), float64(s.Tex_.Min_lod_));
		hi_sz := common.Clamp(hi, float64(s.Tex_.Max_lod_), float64(s.Tex_.Min_lod_));

		color := smath.Vector_t{0.0, 0.0, 0.0, 0.0};
		for i_sample := 0; i_sample < int(ratio); i_sample++{
			subres_index_lo := compute_cube_subresource(IsCubeTexture, face_sz, uint32(lo_sz));
			subres_index_hi := compute_cube_subresource(IsCubeTexture, face_sz, uint32(hi_sz));
			c0 := sample_surface(s.Tex_.Subresource(subres_index_lo), sample_coord_x, sample_coord_y, sample, common.Sampler_state_min);
			c1 := sample_surface(s.Tex_.Subresource(subres_index_hi), sample_coord_x, sample_coord_y, sample, common.Sampler_state_min);

			color.Add_apply(Lerp(c0, c1, frac).get_vec4())

			sample_coord_x += long_axis.X;
			sample_coord_y += long_axis.Y;
		}

		color.Divide(ratio)

		return Color_rgba32f(color)
	}

	//assert(false, "Mip filters is error.");
	return s.Desc_.border_color;
}

func sample( coordx,  coordy, miplevel float64) Color_rgba32f{
	return sample_impl(false,0, coordx, coordy, 0, miplevel, 1.0, smath.Vector_t{0.0, 0, 0.0, 0.0})
}

func  (samp Sampler)Sample_cube( coordx,  coordy,  coordz, miplevel float64) Color_rgba32f{
	var major_dir common.Cubemap_faces;
	var s, t, m float64;

	x := coordx;
	y := coordy;
	z := coordz;

	ax := math.Abs(x);
	ay := math.Abs(y);
	az := math.Abs(z);

	if(ax > ay && ax > az) {
		// x max
		m = ax;
		if(x > 0){
			//+x
			s = 0.5 * (z / m + 1.0);
			t = 0.5 * (y / m + 1.0);
			major_dir = common.Cubemap_face_positive_x;
		} else {
			s = 0.5 * (-z / m + 1.0);
			t = 0.5 * (y / m + 1.0);
			major_dir = common.Cubemap_face_negative_x;
		}
	} else {
	if(ay > ax && ay > az){
		m = ay;
		if(y > 0){
			//+y
			s =0.5 * (x / m + 1.0);
			t = 0.5 * (z / m + 1.0);
			major_dir = common.Cubemap_face_positive_y;
		} else {
			s = 0.5 * (x / m + 1.0);
			t = 0.5 * (-z / m + 1.0);
			major_dir = common.Cubemap_face_negative_y;
		}
	} else {
		m = az;
		if(z > 0){
		//+z
		s = 0.5 * (-x / m + 1.0);
		t = 0.5 * (y / m + 1.0);
		major_dir = common.Cubemap_face_positive_z;
		} else {
		s = 0.5 * (x / m + 1.0);
		t = 0.5 * (y / m + 1.0);
		major_dir = common.Cubemap_face_negative_z;
		}
	}
}

	//EFLIB_ASSERT(tex_->get_texture_type() != texture_type_cube , "texture is not a cube texture.");

	return samp.Sample_impl(true,int(major_dir), s, t, 0, miplevel, 1.0, smath.Vector_t{0.0, 0.0, 0.0, 0.0});
}

func (s Sampler)Calc_lod_2d( ddx,  ddy smath.Vec_2) float64{
 	size := s.Tex_.Size_;

	ddx_vec4 := [4]uint32{uint32(ddx.X), uint32(ddx.Y), 0.0, 0.0};
	ddy_vec4 := [4]uint32{uint32(ddy.X), uint32(ddy.Y), 0.0, 0.0};

	var lod, ratio float64;
	var  long_axis smath.Vector_t;
	if( s.Desc_.mip_filter == common.Filter_anisotropic && s.Desc_.max_anisotropy > 1 ) {
		Calc_anisotropic_lod(size, ddx_vec4, ddy_vec4, 0.0, &lod, &ratio, &long_axis);
	} else{
		lod = Calc_lod(size, ddx_vec4, ddy_vec4, 0.0);
		ratio = 1.0;
	}

	return lod;
}

func (s Sampler) sample_2d_lod( proj_coord smath.Vec_2, lod float64) Color_rgba32f{
	return sample( proj_coord.X, proj_coord.Y, lod )
}

func (s Sampler) sample_2d_grad( proj_coord, ddx,  ddy smath.Vec_2,  lod_bias float64) Color_rgba32f{
	size := s.Tex_.Size_

	ddx_vec4 := [4]uint32{uint32(ddx.X), uint32(ddx.Y), 0.0, 0.0}
	ddy_vec4 := [4]uint32{uint32(ddy.X), uint32(ddy.Y), 0.0, 0.0}

	var lod, ratio float64;
	var long_axis smath.Vector_t;
	if( s.Desc_.mip_filter == common.Filter_anisotropic && s.Desc_.max_anisotropy > 1 )	{
		Calc_anisotropic_lod(size, ddx_vec4, ddy_vec4, lod_bias, &lod, &ratio, &long_axis);
	} else {
		lod = Calc_lod(size, ddx_vec4, ddy_vec4, lod_bias);
		ratio = 1.0;
	}

	return sample_impl(false,0, proj_coord.X, proj_coord.Y, 0, lod, ratio, long_axis);
}

func (s Sampler)calc_anisotropic_lod( size [4]uint32, ddx, ddy smath.Vector_t,  bias float64,
	out_lod, out_ratio *float64, out_long_axis  smath.Vector_t){
	size_vec4 := smath.Vector_t{ float64(size[0]), float64(size[1]), float64(size[2]), 0 }

	ddx_ts := ddx.Multiply(size_vec4)
	ddy_ts := ddy.Multiply(size_vec4)

	ddx_len := ddx_ts.Length();
	ddy_len := ddy_ts.Length();

	diag0_len := ddx_ts.Sub_apply(*ddy_ts).Length();
	diag1_len := (ddx_ts.Sub_apply(*ddy_ts)).Length();

	minor_axis_len := math.Min( math.Min(diag0_len, diag1_len), math.Min(ddx_len, ddy_len) );

	if	(minor_axis_len == 0.0){
		minor_axis_len = 0.000001
	};

	var long_axis *smath.Vector_t
	var long_axis_len float64;
	if(ddx_len > ddy_len){
		long_axis_len = ddx_len;
		long_axis = ddx_ts;
	}else{
		long_axis_len = ddy_len;
		long_axis = ddy_ts;
	}

	probe_count := long_axis_len / minor_axis_len;

	if probe_count > float64(s.Desc_.max_anisotropy){
		probe_count = float64(s.Desc_.max_anisotropy);
	}else{
		probe_count = math.Floor(probe_count+0.5);
	}

	if (probe_count == 1.0){
		*out_lod = math.Log2(long_axis_len) + bias;
		*out_ratio = 1.0
		return;
	} else{
		step := long_axis_len / probe_count;
		minor_axis_len = 2 * long_axis_len / (probe_count + 1);
		*out_lod = math.Log2(minor_axis_len) + bias;
		&out_long_axis = long_axis.Scale(step).Divide(long_axis_len)
		//out_long_axis = 1.0 / size_vec4;
		*out_ratio = probe_count;
	}
}