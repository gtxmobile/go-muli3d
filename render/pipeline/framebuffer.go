package pipeline
import (
	"../fundations"
	"../buffer"
)
type Depth_stencil_op_desc struct {

	stencil_fail_op			Stencil_op
	stencil_depth_fail_op	Stencil_op
	stencil_pass_op			Stencil_op
	stencil_func			fundations.Compare_function
}



type Depth_stencil_desc struct {
	depth_enable				bool
	depth_write_mask				bool
	depth_func				fundations.Compare_function
	stencil_enable				bool
	stencil_read_mask				uint8
	stencil_write_mask				uint8
	front_face				Depth_stencil_op_desc
	back_face				Depth_stencil_op_desc
}
func New_depth_stencil_desc(compare_function_less fundations.Compare_function) Depth_stencil_desc{
 	ret := Depth_stencil_desc{}
 	ret.depth_enable = true
 	ret.depth_write_mask = true
 	ret.depth_func = compare_function_less
 	ret.stencil_enable = false
 	ret.stencil_read_mask = 0xFF
 	ret.stencil_write_mask = 0xFF
	return ret
}
type Mask_stencil_fn   func( stencil,  mask uint32)uint32
type Depth_test_fn		func( ps_depth,   our_depth float64)bool
type Stencil_test_fn   func( ref,      cur_stencil uint32)bool
type Stencil_op_fn     func( ref,      cur_stencil uint32)uint32
type Depth_stencil_state struct {
	Desc_         Depth_stencil_desc
	mask_stencil_ Mask_stencil_fn
	depth_test_   Depth_test_fn
	stencil_test_     [2]Stencil_test_fn
	stencil_op_           [2]Stencil_op_fn
}

depth_stencil_state(const depth_stencil_desc& desc);
const depth_stencil_desc& get_desc() const;

bool        depth_test(float ps_depth, float cur_depth) const;
bool        stencil_test(bool front_face, uint32_t ref, uint32_t cur_stencil) const;
uint32_t    stencil_operation(bool front_face, bool depth_pass, bool stencil_pass, uint32_t ref, uint32_t cur_stencil) const;
uint32_t    mask_stencil(uint32_t stencil, uint32_t stencil_mask) const;


type Framebuffer	struct {
	Color_targets_			[MAX_RENDER_TARGETS]*buffer.Surface
	Ds_target_				*buffer.Surface
	Ds_state_				*Depth_stencil_state
	Stencil_ref_			uint32
	Stencil_read_mask_		uint32
	Stencil_write_mask_		uint32
	Early_z_enabled_		bool
	Sample_count_			uint32
	Px_full_mask_			uint32
	read_depth_stencil_		func(depth float64, stencil, stencil_mask uint32, ds_data interface{});
	write_depth_stencil_	func(ds_data interface{}, depth float64, stencil, stencil_mask uint32);
}

void update_ds_rw_functions(bool ds_format_changed, bool ds_state_changed, bool output_depth_enabled);

public:
void initialize	(render_stages const* stages);
void update		(render_state* state);

framebuffer();
~framebuffer(void);

bool		early_z_enabled() const { return early_z_enabled_; }

void		render_sample(cpp_blend_shader* cpp_bs, size_t x, size_t y, size_t i_sample, const ps_output& ps, float depth, bool front_face);
void		render_sample_quad(cpp_blend_shader* cpp_bs, size_t x, size_t y, uint64_t quad_mask, ps_output const* quad, float const* depth, bool front_face, float const* aa_offset);
uint64_t	early_z_test(size_t x, size_t y, float depth, float const* aa_z_offset);
uint64_t	early_z_test(size_t x, size_t y, uint32_t px_mask, float depth, float const* aa_z_offset);
uint64_t	early_z_test_quad(size_t x, size_t y, float const* depth, float const* aa_z_offset);
uint64_t	early_z_test_quad(size_t x, size_t y, uint64_t quad_mask, float const* depth, float const* aa_z_offset);

static void clear_depth_stencil(surface* tar, uint32_t flag, float depth, uint32_t stencil);
};