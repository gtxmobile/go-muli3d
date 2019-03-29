package pipeline

import (
	"bytes"
	"fmt"
	"github.com/lxn/win"
	"../../common"

)


type Viewport struct {
	X,Y,W,H,MinZ,MaxX float64
}

type Renderer_parameters struct {

	Backbuffer_width uint32
	Backbuffer_height uint32
	Backbuffer_num_samples uint32
	Backbuffer_format	Pixel_format
	Native_window		win.HWND
}
type Mapped_resource struct{
	data interface{};
	row_pitch, depth_pitch uint32
};
type Renderer interface{

	Create_buffer(size int) *bytes.Buffer
	Create_tex2d( width,  height,  num_samples uint32,fmt Pixel_format) *Texture
	Create_texcube(width, height,  num_samples uint32,fmt Pixel_format) *Texture
	Create_sampler( desc Sampler_desc, tex *Texture) *Sampler
	Create_query(id common.Async_object_ids) *Async_object
	Create_input_layout(elem_descs *Input_element_desc, elems_count uint32,code *Shader_object) *Input_layout
	//create_input_layout(elem_descs *Input_element_desc, elems_count uint32, vs *Cpp_vertex_shader) *Input_layout

	Mmap( m *Mapped_resource, buf *bytes.Buffer, mm common.Map_mode) fundations.Result
	//map(mapped_resource&, surface_ptr const& buf, map_mode mm) 
	Unmap() fundations.Result

	// State set
	Set_vertex_buffers (starts_slot,buffers_count uint32, buffers *bytes.Buffer,strides, offsets *uint32) fundations.Result
	Set_index_buffer(hbuf *bytes.Buffer, index_fmt Format)
	Set_input_layout(layout *Input_layout)
	Set_vertex_shader(hvs *Cpp_vertex_shader)
	Set_primitive_topology( primtopo common.Primitive_topology)
	Set_vertex_shader_code(*Shader_object )
	Set_vs_variable_value( name string, pvariable interface{}, sz uint32)
	Set_vs_variable_pointer( name string , pvariable interface{}, sz uint32 )
	Set_vs_sampler( name string, samp *Sampler)
	Set_rasterizer_state(rs *Raster_state)
	Set_ps_variable( name string, data interface{}, sz uint32)
	Set_ps_sampler( name string, samp *Sampler)
	Set_blend_shader(hbs *Cpp_blend_shader)
	Set_pixel_shader(hps *Cpp_pixel_shader)
	Set_pixel_shader_code( so *Shader_object )
	Set_depth_stencil_state( dss *Depth_stencil_state, Stencil_ref int32)
	Set_render_targets(color_target_count uint32, color_targets *Surface, ds_target *Surface)
	Set_viewport(vp *Viewport)

	//template <typename T>
	//result set_vs_variable( std::string const& name, T const* data )
	//{
	Set_vs_variable( name string, data interface{})
	//}
	//template <typename T>
	//result set_ps_variable( std::string const& name, T const* data )
	//{
	//Set_ps_variable( name, static_cast<void const*>(data), sizeof(T) );
	//}
	Get_surface() Surface
	// State get
	Get_index_buffer() *bytes.Buffer
	Get_index_format() *Format
	Get_primitive_topology() Primitive_topology
	Get_vertex_shader() *Cpp_vertex_shader
	Get_vertex_shader_code() *Shader_object
	Get_rasterizer_state() *Raster_state
	Get_pixel_shader() *Cpp_pixel_shader
	Get_pixel_shader_code() *Shader_object
	Get_blend_shader() *Cpp_blend_shader
	Get_viewport() Viewport

	//render operations
	Begin(async_obj Async_object) fundations.Result
	End(async_obj Async_object) fundations.Result
	Get_data(async_obj Async_object, data interface{}, do_not_wait bool) Async_status

	Draw(startpos, primcnt uint32) fundations.Result
	Draw_index(startpos,primcnt uint32, basevert int32) fundations.Result

	Clear_color(color_target *Surface, c Color_rgba32f) fundations.Result
	Clear_depth_stencil(depth_stencil_target *Surface, f uint32, d float64, s uint32) fundations.Result
	Flush() fundations.Result
}




type Async_object struct {

}

type Internal_statistics struct {

}
type Pipeline_profiles struct {

}
type Pipeline_statistics struct {

}


const USE_ASYNC_RENDERER =1
func create_software_renderer() *Renderer {
	if USE_ASYNC_RENDERER > 0 {
		return create_async_renderer()
	} else {
		return create_sync_renderer()
	}
}

func create_benchmark_renderer() *Renderer{
	return create_sync_renderer()
}

func compile_with_log(code string,profile *shader_profile, logs *shader_log)*Shader_object{

	var external_funcs []external_function_desc
	external_funcs=external_funcs.append( external_function_desc(tex2Dlod,		"sasl.vs.tex2d.lod",	true) )
	external_funcs=external_funcs.append( external_function_desc(texCUBElod,	"sasl.vs.texCUBE.lod",	true) )
	external_funcs=external_funcs.append( external_function_desc(tex2Dlod_ps,	"sasl.ps.tex2d.lod" ,	true) )
	external_funcs=external_funcs.append( external_function_desc(tex2Dgrad_ps,	"sasl.ps.tex2d.grad",	true) )
	external_funcs=external_funcs.append( external_function_desc(tex2Dbias_ps,	"sasl.ps.tex2d.bias",	true) )
	external_funcs=external_funcs.append( external_function_desc(tex2Dproj_ps,	"sasl.ps.tex2d.proj",	true) )
	var ret *Shader_object
	modules.host.compile(ret, logs, code, profile, external_funcs)

	return ret
}

func compile_prof(code string, profile Shader_profile) *Shader_object{
	var log *Shader_log
	ret := compile_with_log(code, profile, log)

	if ret!= nil {
		fmt.Println("Shader was compiled failed!")
		for i := 0 ;i < log.Count();i+=1 {
			fmt.Println(log.Log_string(i))
		}
	}
	return ret
}

func compile(code string, languages lang)	*Shader_object{
	var prof =Shader_profile{}
	prof.language = lang
	return compile_prof(code, prof)
}

func compile_from_file_with_log(file_name string,profile *Shader_profile, logs Shader_log) *Shader_object{
	var external_funcs []external_function_desc

	external_funcs = external_funcs.append( external_function_desc(tex2Dlod,		"sasl.vs.tex2d.lod",	true) )
	external_funcs = external_funcs.append( external_function_desc(texCUBElod,		"sasl.vs.texCUBE.lod",	true) )
	external_funcs = external_funcs.append( external_function_desc(tex2Dlod_ps,		"sasl.ps.tex2d.lod" ,	true) )
	external_funcs = external_funcs.append( external_function_desc(tex2Dgrad_ps,	"sasl.ps.tex2d.grad",	true) )
	external_funcs = external_funcs.append( external_function_desc(tex2Dbias_ps,	"sasl.ps.tex2d.bias",	true) )
	external_funcs = external_funcs.append( external_function_desc(tex2Dproj_ps,	"sasl.ps.tex2d.proj",	true) )

	var ret *Shader_object
	modules.host.compile_from_file(ret, logs, file_name, profile, external_funcs)
	return ret
}

func compile_from_file_with_prof(file_name string, profile shader) *Shader_object{
	var log *Shader_log
	ret := compile_from_file(file_name, profile, log)

	if	ret!=nil {
		fmt.Println("Shader was compiled failed!")
		for i:=0 ;i<log.count() ;i++  {
			fmt.Println(log.log_string(i))
		}
	}
	return ret
}

func compile_from_file(file_name string, languages lang) *Shader_object{
	var prof shader_profile
	prof.language = lang
	return compile_from_file(file_name, prof)
}

type Renderer_impl struct {
	Resource_pool_ *Resource_manager
}

func Create_buffer(size int) *bytes.Buffer{
	return nil
}
func (render *Renderer_impl)Create_tex2d(width, height, num_samples uint32,fmt Pixel_format) *Texture{
	return render.Resource_pool_.Create_texture_2d(width, height, num_samples, fmt).Texture
}


func (render *Renderer_impl)Create_texcube(width, height,  num_samples uint32,fmt Pixel_format) *Texture{
	return render.Resource_pool_.Create_texture_cube(width, height, num_samples, fmt).Texture
}
func Create_sampler( desc Texture_desc, tex *Texture) *Sampler
func Create_query(id Async_object_ids) *Async_object
func Create_input_layout(elem_descs *Input_element_desc, elems_count uint32,code *Shader_object) *Input_layout
//create_input_layout(elem_descs *Input_element_desc, elems_count uint32, vs *Cpp_vertex_shader) *Input_layout

func (render *Renderer_impl)Mmap( m *Mapped_resource, buf *Buffer, mm Map_mode) fundations.Result{

}
//map(mapped_resource&, surface_ptr const& buf, map_mode mm)
func (render *Renderer_impl)Unmap() fundations.Result{

}

// State set
func (render *Renderer_impl)Set_vertex_buffers (starts_slot,buffers_count uint32, buffers *bytes.Buffer,strides, offsets *uint32) fundations.Result{

}
func (render *Renderer_impl)Set_index_buffer(hbuf *bytes.Buffer, index_fmt Format){

}
func (render *Renderer_impl)Set_input_layout(layout *Input_layout){

}
func (render *Renderer_impl)Set_vertex_shader(hvs *Cpp_vertex_shader){

}
func (render *Renderer_impl)Set_primitive_topology( primtopo Primitive_topology){

}
func (render *Renderer_impl)Set_vertex_shader_code(*Shader_object ){

}
func (render *Renderer_impl)Set_vs_variable_value( name string, pvariable interface{}, sz uint32){}
func (render *Renderer_impl)Set_vs_variable_pointer( name string , pvariable interface{}, sz uint32 ){}
func (render *Renderer_impl)Set_vs_sampler( name string, samp *Sampler){}
func (render *Renderer_impl)Set_rasterizer_state(rs *Raster_state){}
func (render *Renderer_impl)Set_ps_variable( name string, data interface{}, sz uint32){}
func (render *Renderer_impl)Set_ps_sampler( name string, samp *Sampler){}
func (render *Renderer_impl)Set_blend_shader(hbs *Cpp_blend_shader){}
func (render *Renderer_impl)Set_pixel_shader(hps *Cpp_pixel_shader){}
func (render *Renderer_impl)Set_pixel_shader_code( so *Shader_object ){}
func (render *Renderer_impl)Set_depth_stencil_state( dss *Depth_stencil_state, Stencil_ref int32){}
func (render *Renderer_impl)Set_render_targets(color_target_count uint32, color_targets *Surface, ds_target *Surface){}
func (render *Renderer_impl)Set_viewport(vp *Viewport){}

//template <typename T>
//result set_vs_variable( std::string const& name, T const* data )
//{
func (render *Renderer_impl)Set_vs_variable( name string, data interface{}){}
//}
//template <typename T>
//result set_ps_variable( std::string const& name, T const* data )
//{
//Set_ps_variable( name, static_cast<void const*>(data), sizeof(T) );
//}

// State get
func (render *Renderer_impl)Get_index_buffer() *bytes.Buffer{}
func (render *Renderer_impl)Get_index_format() *Format{}
func (render *Renderer_impl)Get_primitive_topology() Primitive_topology{}
func (render *Renderer_impl)Get_vertex_shader() *Cpp_vertex_shader{}
func (render *Renderer_impl)Get_vertex_shader_code() *Shader_object{}
func (render *Renderer_impl)Get_rasterizer_state() *Raster_state{}
func (render *Renderer_impl)Get_pixel_shader() *Cpp_pixel_shader{}
func (render *Renderer_impl)Get_pixel_shader_code() *Shader_object{}
func (render *Renderer_impl)Get_blend_shader() *Cpp_blend_shader{}
func (render *Renderer_impl)Get_viewport() Viewport{}

//render operations
func (render *Renderer_impl)Begin(async_obj Async_object) fundations.Result{}
func (render *Renderer_impl)End(async_obj Async_object) fundations.Result{}
func (render *Renderer_impl)Get_data(async_obj Async_object, data interface{}, do_not_wait bool) Async_status{}
func (render *Renderer_impl)Draw(startpos, primcnt uint32) fundations.Result{}
func (render *Renderer_impl)Draw_index(startpos,primcnt uint32, basevert int32) fundations.Result{}
func (render *Renderer_impl)Clear_color(color_target *Surface, c Color_rgba32f) fundations.Result{}
func (render *Renderer_impl)Clear_depth_stencil(depth_stencil_target *Surface, f uint32, d float64, s uint32) fundations.Result{}
func (render *Renderer_impl)Flush() fundations.Result{}