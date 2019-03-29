package shader
import (
	"../buffer"
	"../fundations"
	"../pipeline"
	"../../smath"
)
type Shader_object interface {
	get_reflection()
	native_function()
}
type Languages uint32
const (
	Lang_none	Languages = iota
	Lang_general
	Lang_vertex_shader
	Lang_pixel_shader
	Lang_blending_shader
	Lang_count
)

type System_values uint32
const (
	Sv_none System_values = iota
	
	Sv_position
	Sv_texcoord
	Sv_normal
	
	Sv_blend_indices
	Sv_blend_weights
	Sv_psize
	
	Sv_target
	Sv_depth
	
	Sv_customized
)
type register_array  [pipeline.MAX_VS_OUTPUT_ATTRS+1]smath.Vector_t

type Vs_input struct {
	registers_ register_array ;
	//vs_output(const vs_output& rhs);
	//vs_output& operator = (const vs_output& rhs);
}
public:
vs_input()
{}

eflib::vec4& attribute(size_t index)
{
return attributes_[index];
}

eflib::vec4 const& attribute(size_t index) const
{
return attributes_[index];
}

private:
typedef boost::array<
eflib::vec4, MAX_VS_INPUT_ATTRS > attribute_array;
attribute_array attributes_;

vs_input(const vs_input& rhs);
vs_input& operator=(const vs_input& rhs);
};

#include <eflib/include/platform/disable_warnings.h>
class EFLIB_ALIGN(16) vs_output
{
public:
/*
BUG FIX:
	Type A is aligned by A bytes, new A is A bytes aligned too, but address of new A[] is not aligned.
	It is known bug but not fixed yet.
	operator new/delete will ensure the address is aligned.
*/
static void* operator new[] (size_t size)
{
if(size == 0) { size = 1; }
return eflib::aligned_malloc(size, 16);
}

static void operator delete[](void* p)
{
if(p == nullptr) return;
return eflib::aligned_free(p);
}

enum attrib_modifier_type
{
am_linear = 1UL << 0,
am_centroid = 1UL << 1,
am_nointerpolation = 1UL << 2,
am_noperspective = 1UL << 3,
am_sample = 1UL << 4
};

public:
eflib::vec4& position()
{
return registers_[0];
}

eflib::vec4 const& position() const
{
return registers_[0];
}

eflib::vec4* attribute_data()
{
return registers_.data() + 1;
}

eflib::vec4 const* attribute_data() const
{
return registers_.data() + 1;
}

eflib::vec4* raw_data()
{
return registers_.data();
}

eflib::vec4 const* raw_data() const
{
return registers_.data();
}

eflib::vec4 const& attribute(size_t index) const
{
return attribute_data()[index];
}

eflib::vec4& attribute(size_t index)
{
return attribute_data()[index];
}

vs_output()
{}

private:

};
#include <eflib/include/platform/enable_warnings.h>

#if defined(EFLIB_MSVC)
#pragma warning(push)
#pragma warning(disable: 4324)	// warning C4324: Structure was padded due to __declspec(align())
#endif

struct triangle_info
{
vs_output const*			v0;
bool						front_face;
EFLIB_ALIGN(16)	eflib::vec4	bounding_box;
EFLIB_ALIGN(16)	eflib::vec4	edge_factors[3];
vs_output					ddx;
vs_output					ddy;

triangle_info() {}
triangle_info(triangle_info const& /*rhs*/)
{
}
triangle_info& operator = (triangle_info const& /*rhs*/)
{
return *this;
}
};

#if defined(EFLIB_MSVC)
#pragma warning(pop)
#endif

//vs_output compute_derivate
struct ps_output
{
boost::array<eflib::vec4, MAX_RENDER_TARGETS> color;
};

struct pixel_accessor
{
pixel_accessor(surface** const& color_buffers, surface* ds_buffer)
{
color_buffers_ = color_buffers;
ds_buffer_ = ds_buffer;
}

void set_pos(size_t x, size_t y)
{
x_ = x;
y_ = y;
}

color_rgba32f color(size_t target_index, size_t sample_index) const
{
if(color_buffers_[target_index] == nullptr)
{
return color_rgba32f(0.0f, 0.0f, 0.0f, 0.0f);
}
return color_buffers_[target_index]->get_texel(x_, y_, sample_index);
}

void color(size_t target_index, size_t sample, const color_rgba32f& clr)
{
if(color_buffers_[target_index] != nullptr)
{
color_buffers_[target_index]->set_texel(x_, y_, sample, clr);
}
}

void* depth_stencil_address(size_t sample) const
{
return ds_buffer_->texel_address(x_, y_, sample);
}

private:
pixel_accessor(const pixel_accessor& rhs);
pixel_accessor& operator = (const pixel_accessor& rhs);

surface**   color_buffers_;
surface*    ds_buffer_;
size_t      x_, y_;
};

type vs_input_construct func(out Vs_input,attrs *smath.Vector_t) Vs_input;
type vs_input_copy func(out Vs_input,  in Vs_input)Vs_input;
type  Vs_input_op struct{
	construct vs_input_construct
	copy	vs_input_copy
};
type aligned_allocator struct {

}
type aligned_vector [32]aligned_allocator ;

type Pixel_shader_unit struct {
	code *Shader_object
	used_samplers 		[]buffer.Sampler      // For take ownership
	stream_data			aligned_vector
	buffer_data			aligned_vector
	stream_odata		aligned_vector
	buffer_odata		aligned_vector
}
public:
pixel_shader_unit();
~pixel_shader_unit();

pixel_shader_unit( pixel_shader_unit const& );
pixel_shader_unit& operator = ( pixel_shader_unit const& );

boost::shared_ptr<pixel_shader_unit> clone() const;

void initialize( shader_object const* );
void reset_pointers();

void set_variable( std::string const&, void const* data );
void set_sampler( std::string const&, sampler_ptr const& samp );

void update( vs_output* inputs, shader_reflection const* vs_abi );
void execute(ps_output* outs, float* depths);

public:
shader_object const* code;

std::vector<sampler_ptr>									used_samplers;	// For take ownership

typedef std::vector<char, eflib::aligned_allocator<char, 32> > aligned_vector;

aligned_vector stream_data;
aligned_vector buffer_data;

aligned_vector stream_odata;
aligned_vector buffer_odata;
};

EFLIB_DECLARE_CLASS_SHARED_PTR(vx_shader_unit);
class vx_shader_unit
{
public:
virtual uint32_t output_attributes_count() const = 0;
virtual uint32_t output_attribute_modifiers(size_t index) const = 0;

virtual void execute(size_t ivert, void* out_data) = 0;
virtual void execute(size_t ivert, vs_output& out) = 0;

virtual ~vx_shader_unit(){}
};

namespace vs_output_functions
{
using eflib::vec4;

typedef vs_output& (*construct)		(vs_output& out, vec4 const& position, vec4 const* attrs);
typedef vs_output& (*copy)			(vs_output& out, const vs_output& in);

typedef vs_output& (*project)		(vs_output& out, const vs_output& in);
typedef vs_output& (*unproject)		(vs_output& out, const vs_output& in);

typedef vs_output& (*add)			(vs_output& out, const vs_output& vso0, const vs_output& vso1);
typedef vs_output& (*sub)			(vs_output& out, const vs_output& vso0, const vs_output& vso1);
typedef vs_output& (*mul)			(vs_output& out, const vs_output& vso0, float f);
typedef vs_output& (*div)			(vs_output& out, const vs_output& vso0, float f);

typedef void (*compute_derivative)	(vs_output& ddx, vs_output& ddy, vs_output const& e01, vs_output const& e02, float inv_area);

typedef vs_output& (*lerp)			(vs_output& out, const vs_output& start, const vs_output& end, float step);
typedef vs_output& (*step_2d_unproj)(
vs_output& out, vs_output const& start,
float step0, vs_output const& derivation0,
float step1, vs_output const& derivation1);
typedef vs_output& (*step_2d_unproj_quad)(
vs_output* out, vs_output const& start,
float step0, vs_output const& derivation0,
float step1, vs_output const& derivation1);
}

struct vs_output_op
{
vs_output_functions::construct		construct;
vs_output_functions::copy			copy;

vs_output_functions::project		project;
vs_output_functions::unproject		unproject;

vs_output_functions::add			add;
vs_output_functions::sub			sub;
vs_output_functions::mul			mul;
vs_output_functions::div			div;

vs_output_functions::lerp			lerp;
vs_output_functions::step_2d_unproj	step_2d_unproj_pos;
vs_output_functions::step_2d_unproj	step_2d_unproj_attr;
vs_output_functions::step_2d_unproj_quad
step_2d_unproj_pos_quad;
vs_output_functions::step_2d_unproj_quad
step_2d_unproj_attr_quad;

vs_output_functions::compute_derivative
compute_derivative;

typedef boost::array<uint32_t, MAX_VS_OUTPUT_ATTRS> interpolation_modifier_array;
interpolation_modifier_array		attribute_modifiers;
};

vs_input_op& get_vs_input_op(uint32_t n);
vs_output_op& get_vs_output_op(uint32_t n);
float compute_area(const vs_output& v0, const vs_output& v1, const vs_output& v2);
void viewport_transform(eflib::vec4& position, viewport const& vp);

type Semantic_value struct {
	name			string
	sv				System_values
	index			uint32
}



public:
static std::string lower_copy( std::string const& name )
{
std::string ret(name);
for( size_t i = 0; i < ret.size(); ++i )
{
if( 'A' <= ret[i] && ret[i] <= 'Z' )
{
ret[i] = ret[i] - ('A' - 'a');
}
}
return ret;
}

semantic_value(): sv(sv_none), index(0){}

explicit semantic_value( std::string const& name, uint32_t index = 0 )
{
assert( !name.empty() );

std::string lower_name = lower_copy(name);

if( lower_name == "position" || lower_name == "sv_position" ){
sv = sv_position;
} else if ( lower_name == "normal" ){
sv = sv_normal;
} else if ( lower_name == "texcoord" ){
sv = sv_texcoord;
} else if ( lower_name == "color" || lower_name == "sv_target" ){
sv = sv_target;
} else if ( lower_name == "depth" || lower_name == "sv_depth" ) {
sv = sv_depth;
} else if ( lower_name == "blend_indices" ){
sv = sv_blend_indices;
} else if ( lower_name == "blend_weights" ){
sv = sv_blend_weights;
} else if ( lower_name == "psize" ){
sv = sv_psize;
} else {
sv = sv_customized;
this->name = lower_name;
}
this->index = index;
}

semantic_value( system_values sv, uint32_t index = 0 ){
assert( sv_none <= sv && sv < sv_customized );
this->sv = sv;
this->index = index;
}

std::string const& get_name() const{
return name;
}

system_values const& get_system_value() const{
return sv;
}

uint32_t get_index() const{
return index;
}

bool operator < ( semantic_value const& rhs ) const{
return sv < rhs.sv || name < rhs.name || index < rhs.index;
}

bool operator == ( semantic_value const& rhs ) const{
return is_same_sv(rhs) && index == rhs.index;
}

bool operator == ( system_values rhs ) const{
return sv == rhs && index == 0;
}

template <typename T>
bool operator != ( T const& v ) const{
return !( *this == v );
}

semantic_value advance_index(size_t i) const
{
semantic_value ret;
ret.name	= name;
ret.sv		= sv;
ret.index	= static_cast<uint32_t>(index + i);
return ret;
}

bool valid() const
{
return sv != sv_none || !name.empty();
}

private:


bool is_same_sv( semantic_value const& rhs ) const{
if( sv != rhs.sv ) return false;
if( sv == sv_customized ) return rhs.name == name;
return true;
}
};

inline size_t hash_value( semantic_value const& v ){
size_t seed = v.get_index();
if(v.get_system_value() != sv_customized )
{
boost::hash_combine( seed, static_cast<size_t>( v.get_system_value() ) );
}
else
{
boost::hash_combine( seed, v.get_name() );
}
return seed;
}

type Shader_profile struct {
	language	Languages
}

type Cpp_shader interface {

	set_sampler(varname string, samp *buffer.Sampler ) fundations.Result
	set_constant(varname string, pval interface{},moreval ...interface{}) fundations.Result
	//set_constant(varname string, pval interface{}, index uint32)	fundations.Result
	find_register( sv Semantic_value, index uint32) fundations.Result
	//unordered_map<semantic_value, size_t> const& get_register_map() = 0;
	clone() *Cpp_shader
}

template <typename T>
boost::shared_ptr<T> clone()
{
auto ret = boost::dynamic_pointer_cast<T>( clone() );
assert(ret);
return ret;
}
type Cpp_shader_impl struct {
	varmap_			variable_map
	contmap_		container_variable_map	
	regmap_			register_map
	sampmap_		sampler_map				
}
public:
result set_sampler(std::_tstring const& samp_name, sampler_ptr const& samp)
{
auto samp_it = sampmap_.find(samp_name);
if( samp_it == sampmap_.end() )
{
return result::failed;
}
*(samp_it->second) = samp;
return result::ok;
}

result set_constant(const std::_tstring& varname, shader_constant::const_voidptr pval){
variable_map::iterator var_it = varmap_.find(varname);
if( var_it == varmap_.end() ){
return result::failed;
}
if(shader_constant::assign(var_it->second, pval)){
return result::ok;
}
return result::failed;
}

result set_constant(const std::_tstring& varname, shader_constant::const_voidptr pval, size_t index)
{
container_variable_map::iterator cont_it = contmap_.find(varname);
if( cont_it == contmap_.end() ){
return result::failed;
}
cont_it->second->set(pval, index);
return result::ok;
}

result find_register( semantic_value const& sv, size_t& index );
boost::unordered_map<semantic_value, size_t> const& get_register_map();
void bind_semantic( char const* name, size_t semantic_index, size_t register_index );
void bind_semantic( semantic_value const& s, size_t register_index );

template< class T >
result declare_constant(const std::_tstring& varname, T& var)
{
varmap_[varname] = shader_constant::voidptr(&var);
return result::ok;
}

result declare_sampler(const std::_tstring& varname, sampler_ptr& var)
{
sampmap_[varname] = &var;
return result::ok;
}

template<class T>
result declare_container_constant(const std::_tstring& varname, T& var)
{
return declare_container_constant_impl(varname, var, var[0]);
}

private:
typedef std::map<std::_tstring, shader_constant::voidptr>
variable_map;
typedef std::map<std::_tstring, sampler_ptr*>
sampler_map;
typedef std::map<std::_tstring, boost::shared_ptr<detail::container> >
container_variable_map;
typedef boost::unordered_map<semantic_value, size_t>
register_map;



template<class T, class ElemType>
result declare_container_constant_impl(const std::_tstring& varname, T& var, const ElemType&)
{
varmap_[varname] = shader_constant::voidptr(&var);
contmap_[varname] = boost::shared_ptr<detail::container>(new detail::container_impl<T, ElemType>(var));
return result::ok;
}
};

type Cpp_vertex_shader struct{
public:
void execute(const vs_input& in, vs_output& out);
virtual void shader_prog(const vs_input& in, vs_output& out) = 0;
virtual uint32_t num_output_attributes() const = 0;
virtual uint32_t output_attribute_modifiers(uint32_t index) const = 0;
};

type Cpp_pixel_shader struct {
	bool front_face_;
	vs_output const *        px_;
	vs_output const *        quad_;
	uint64_t                lod_flag_;
	float                    lod_[MAX_VS_OUTPUT_ATTRS];
}
protected:
bool front_face() const	{ return front_face_; }

eflib::vec4	  ddx(size_t iReg) const;
eflib::vec4   ddy(size_t iReg) const;

color_rgba32f tex2d(const sampler& s, size_t iReg);
color_rgba32f tex2dlod(const sampler& s, size_t iReg);
color_rgba32f tex2dlod(sampler const& s, eflib::vec4 const& coord_with_lod);
color_rgba32f tex2dproj(const sampler& s, size_t iReg);

color_rgba32f texcube(const sampler& s, const eflib::vec4& coord, const eflib::vec4& ddx, const eflib::vec4& ddy, float bias = 0);
color_rgba32f texcube(const sampler&s, size_t iReg);
color_rgba32f texcubelod(const sampler& s, size_t iReg);
color_rgba32f texcubeproj(const sampler& s, size_t iReg);
color_rgba32f texcubeproj(const sampler&s, const eflib::vec4& v, const eflib::vec4& ddx, const eflib::vec4& ddy);

public:
void update_front_face(bool v)
{
front_face_ = v;
}

uint64_t execute(vs_output const* quad_in, ps_output* px_out, float* depth);

virtual bool shader_prog(vs_output const& in, ps_output& out) = 0;
virtual bool output_depth() const;
};

//it is called when render a shaded pixel into framebuffer
type Cpp_blend_shader struct {

}
public:
void execute(size_t sample, pixel_accessor& inout, const ps_output& in);
virtual bool shader_prog(size_t sample, pixel_accessor& inout, const ps_output& in) = 0;
};