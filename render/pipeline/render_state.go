package pipeline
import (
	"../../common"
	"../fundations"
	"../shader"
	"../buffer"
	"bytes"
)
//render status
type Command_id uint32
const(
	Draw Command_id=iota
	Draw_index
	Clear_depth_stencil
	Clear_color
	Async_begin
	Async_end
)

type Render_state struct {

	Cmd							Command_id

	Index_buffer				*bytes.Buffer
	Index_format				Format
	Prim_topo					fundations.Primitive_topology
	Base_vertex					int32
	Start_index					uint32
	Prim_count					uint32

	Str_state					Stream_state
	Layout						*Input_layout

	Vp							Viewport
	Ras_state					*Raster_state
	Stencil_ref					int32
	Ds_state					*Depth_stencil_state
	Cpp_vs						*shader.Cpp_vertex_shader
	Cpp_ps						*shader.Cpp_pixel_shader
	Cpp_bs						*shader.Cpp_blend_shader
	Vx_shader					*shader.Shader_object
	Px_shader					*shader.Shader_object
	Vx_cbuffer					shader.Shader_cbuffer
	Px_cbuffer					shader.Shader_cbuffer
	Vsi_ops						*shader.Vs_input_op
	Ps_proto					*shader.Pixel_shader_unit
	Color_targets				[]*buffer.Surface
	Depth_stencil_target		buffer.Surface
	Asyncs						[fundations.Aoi_count]*Async_object
	Current_async				*Async_object
	Target_vp					Viewport
	Target_sample_count			uint32
	Clear_color_target			*buffer.Surface
	Clear_ds_target				*buffer.Surface
	Clear_f						uint32
	Clear_z						float64
	Clear_stencil				uint32
	Clear_color					fundations.Color_rgba32f
}
func copy_using_state(dest ,src *Render_state){
	switch(src.Cmd)	{
		case Draw:
		case Draw_index:
			*dest = *src
			if(src.Cpp_vs != nil) {
				dest.Cpp_vs = src.Cpp_vs.Clone()
			}
			if	src.Cpp_ps!= nil{
				dest.Cpp_ps = src.Cpp_ps.Clone()
			}
			if	src.Cpp_bs != nil{
				dest.Cpp_bs = src.Cpp_bs.Clone()
			}
			break;
		case Clear_color:
		case Clear_depth_stencil:
			dest.Cmd                = src.Cmd
			dest.Clear_color_target = src.Clear_color_target
			dest.Clear_ds_target    = src.Clear_ds_target

			dest.Clear_f			= src.Clear_f
			dest.Clear_z            = src.Clear_z
			dest.Clear_stencil      = src.Clear_stencil
			dest.Clear_color        = src.Clear_color
			break;
		case Async_begin:
		case Async_end:
			dest.Cmd                = src.Cmd
			dest.Current_async      = src.Current_async
		break;
}
}