package pipeline
import (
	"../fundations"
)

type Raster_desc struct {
	fm                       	fundations.Fill_mode;
	cm                       	fundations.Cull_mode;
	front_ccw                	bool;
	depth_bias               	int32;
	depth_bias_clamp         	float64;
	slope_scaled_depth_bias  	float64;
	depth_clip_enable        	bool;
	scissor_enable           	bool;
	multisample_enable       	bool;
	anti_aliased_line_enable 	bool;
}

func New_raster_desc()Raster_desc{
	return Raster_desc{fundations.Fill_solid, fundations.Cull_back,
	false,0, 0, 0,
	true, false,
	true, false}
}

var Clipper_ptr *Clipper
var Cpp_pixel_shader_ptr *Cpp_pixel_shader;
var Pixel_shader_unit_ptr *Pixel_shader_unit

//struct clip_context;
//struct clip_results;
//class  rasterizer;
//class  raster_multi_prim_context;
//class  vs_output;
//struct viewport;

type  Raster_state struct {
	Desc_		Raster_desc
	Cull_		Cull_func;
	Prim_		Prim_type;
}


typedef bool (*cull_func)				(float area);

raster_state(raster_desc const& desc);

inline raster_desc const& get_desc() const
{
return desc_;
}

inline cull_func get_cull_func() const
{
return cull_;
}

inline bool cull(float area) const
{
return cull_(area);
}
};

bool cull_mode_none(float /*area*/)
{
return false;
}

bool cull_mode_ccw(float area)
{
return area <= 0;
}

bool cull_mode_cw(float area)
{
return area >= 0;
}

raster_state::raster_state(const raster_desc& desc)
: desc_(desc)
{
switch (desc.cm)
{
case cull_none:
cull_ = cull_mode_none;
break;

case cull_front:
cull_ = desc.front_ccw ? cull_mode_ccw : cull_mode_cw;
break;

case cull_back:
cull_ = desc.front_ccw ? cull_mode_cw : cull_mode_ccw;
break;

default:
EFLIB_ASSERT_UNEXPECTED();
break;
}
}