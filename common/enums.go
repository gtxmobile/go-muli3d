package common
type Result uint32
const (
	Ok Result =iota
	Failed
	Outofmemory
	Invalid_parameter
)

type Map_mode uint32
const (
	Map_mode_none Map_mode= iota
	Map_read
	Map_write
	Map_read_write
	Map_write_discard
	Map_write_no_overwrite
)

type Map_result uint32
const (
	Map_succeed Map_result = iota
	Map_failed
	Map_do_not_wait
)

type Resource_usage uint32
const (
	Resource_access_none Resource_usage = iota

	Resource_read		= 0x1	// xxx1
	Resource_write		= 0x2 	// xx1x

	Resource_client		= 0x4	// x1xx
	Resource_device		= 0x8	// 1xxx

	Client_read  = 0x5 // 0101
	Client_write = 0x6 // 0110

	Device_read  = 0x9 // 1001
	Device_write = 0xa // 1010

	Resource_default	= Device_read | Device_write
	Resource_immutable	= Device_read | Client_read
	Resource_dynamic	= Device_read | Client_write
	Resource_staging	= 0xf	// 1111
)

func  RESERVED(i uint32) uint32{
	return 0xFFFF0000 + i
}
type Primitive_topology uint32

var primitive_point_list = Primitive_topology(RESERVED(0))
var primitive_point_sprite = Primitive_topology(RESERVED(1))
const (
	Primitive_line_list Primitive_topology = iota
	Primitive_line_strip
	Primitive_triangle_list
	Primitive_triangle_fan
	Primitive_triangle_strip
	Primivite_topology_count
)

type  Cull_mode uint32
const (
	Cull_none Cull_mode = iota
	Cull_front
	Cull_back
	Cull_mode_count
)

type Fill_mode uint32
const (
	Fill_wireframe Fill_mode = iota
	Fill_solid
	Fill_mode_count
)

type Prim_type uint32
const (
	Pt_none Prim_type = iota
	Pt_point
	Pt_line
	Pt_solid_tri
	Pt_wireframe_tri
)

type Texture_type uint32
const (
	Texture_type_1d Texture_type = iota
	Texture_type_2d
	Texture_type_cube
	Texture_type_count
)

type Address_mode uint32
const (
	Address_wrap Address_mode = iota
	Address_mirror
	Address_clamp
	Address_border
	Address_mode_count
)

type Filter_type uint32
const (
	Filter_point Filter_type = iota
	Filter_linear
	Filter_anisotropic
	Filter_type_count
)

type Sampler_state uint32
const (
	Sampler_state_min Sampler_state = iota
	Sampler_state_mag = 1
	Sampler_state_mip = 2
	Sampler_state_count = 3
)

type Sampler_axis uint32
const (
	Sampler_axis_u Sampler_axis = iota
	Sampler_axis_v
	Sampler_axis_w
	Sampler_axis_count
)

type Cubemap_faces uint32
const (
	Cubemap_face_positive_x Cubemap_faces = iota
	Cubemap_face_negative_x
	Cubemap_face_positive_y
	Cubemap_face_negative_y
	Cubemap_face_positive_z
	Cubemap_face_negative_z
	Cubemap_faces_count
)

// Usage describes the default component value will be filled to unfilled component if source data isn't a 4-components vector.
// Position means fill to (0 0 0 1)
// Attrib means fill to (0000)
type Input_register_usage_decl uint32
const (
	Input_register_usage_position Input_register_usage_decl = iota
	Input_register_usage_attribute
	Input_register_usage_decl_count
)

type Render_target uint32
const (
	Render_target_none = 0
	Render_target_color = 1
	Render_target_depth_stencil = 2
	Render_target_count = 3
)

type Compare_function uint32
const (
	Compare_function_never Compare_function = iota
	Compare_function_less = 1
	Compare_function_equal = 2
	Compare_function_less_equal = 3
	Compare_function_greater = 4
	Compare_function_not_equal = 5
	Compare_function_greater_equal = 6
	Compare_function_always = 7
)

type Stencil_op uint32
const (
	Stencil_op_keep Stencil_op = iota
	Stencil_op_zero = 2
	Stencil_op_replace = 3
	Stencil_op_incr_sat = 4
	Stencil_op_decr_sat = 5
	Stencil_op_invert = 6
	Stencil_op_incr_wrap = 7
	Stencil_op_decr_wrap = 8
)

type Clear_flag uint32
const (
	Clear_depth Clear_flag = 0x1
	Clear_stencil Clear_flag = 0x2
)

type Async_object_ids uint32
const (
	None Async_object_ids = iota
	Event
	Occlusion
	Pipeline_statistics
	Occlusion_predicate
	Internal_statistics
	Pipeline_profiles
	Aoi_count
)

type  Async_status  uint32
const (
	Error Async_status = iota
	Timeout
	Ready
)