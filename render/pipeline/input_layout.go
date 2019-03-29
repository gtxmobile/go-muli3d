package pipeline

import (
	"../../common"
	"math"
)
type Input_classifications uint32
const Input_per_vertex Input_classifications = iota

// TODO:
// input_per_instance
//};

type Input_element_desc struct{
	Semantic_name 			string
	Semantic_index 			uint32
	Data_format 			Format
	Input_slot 				uint32
	Aligned_byte_offset 	uint32
	Slot_class 				Input_classifications
	Instance_data_step_rate 	uint32
};

//input_element_desc()
//: semantic_index(0)
//, data_format(format_unknown)
//, input_slot(0), aligned_byte_offset(0xFFFFFFFF)
//, slot_class(input_per_vertex), instance_data_step_rate(0)
//{}

//input_element_desc(
//const char* semantic_name,
//uint32_t semantic_index,
//format data_format,
//uint32_t input_slot,
//uint32_t aligned_byte_offset,
//input_classifications slot_class,
//uint32_t instance_data_step_rate
//)
//: semantic_name( semantic_name )
//, semantic_index( semantic_index )
//, data_format( data_format )
//, input_slot( input_slot )
//, aligned_byte_offset( aligned_byte_offset )
//, slot_class( slot_class )
//, instance_data_step_rate( instance_data_step_rate )
//{}


//EFLIB_DECLARE_CLASS_SHARED_PTR(input_layout);
var input_layout_ptr *Input_layout

type Input_layout struct{
	Descs []Input_element_desc;
}

func (il Input_layout)slot_range(  min_slot,  max_slot *uint32){
	*min_slot = math.MaxUint32
	*max_slot = 0;

	for _,elem_desc := range il.Descs {
		*min_slot = uint32(math.Min( float64(*min_slot), float64(elem_desc.Input_slot )))
		*max_slot = uint32(math.Max( float64(*max_slot), float64(elem_desc.Input_slot )))
	}

	if( *max_slot < *min_slot ){
		*max_slot =0
		*min_slot = 0
	}
}

func (il Input_layout)Desc_begin() Input_element_desc {
	return il.Descs[0];
}

func (il Input_layout)Desc_end() Input_element_desc{
	return il.Descs[len(il.Descs)-1];
}

semantic_value input_layout::get_semantic( iterator it ) const{
return semantic_value( it->semantic_name, it->semantic_index );
}

size_t input_layout::find_slot( semantic_value const& v ) const{
input_element_desc const* elem_desc = find_desc( v );
if( elem_desc ){
return elem_desc->input_slot;
}
return std::numeric_limits<size_t>::max();
}

input_element_desc const* input_layout::find_desc( size_t slot ) const{
for( size_t i_desc = 0; i_desc < descs.size(); ++i_desc ){
if( slot == descs[i_desc].input_slot ){
return &( descs[i_desc] );
}
}
return NULL;
}

input_element_desc const* input_layout::find_desc( semantic_value const& v ) const{
for( size_t i_desc = 0; i_desc < descs.size(); ++i_desc ){
if( semantic_value( descs[i_desc].semantic_name,  descs[i_desc].semantic_index ) == v ){
return &( descs[i_desc] );
}
}
return NULL;
}

input_layout_ptr input_layout::create( input_element_desc const* pdesc, size_t desc_count, shader_object_ptr const& /*vs*/ ){
input_layout_ptr ret = make_shared<input_layout>();
ret->descs.assign( pdesc, pdesc + desc_count );

// Check shader code.
// Caculate member offset.

return ret;
}

input_layout_ptr input_layout::create( input_element_desc const* pdesc, size_t desc_count, cpp_vertex_shader_ptr const& /*vs*/ ){
input_layout_ptr ret = make_shared<input_layout>();
ret->descs.assign( pdesc, pdesc + desc_count );

// Check vertex shader.
// Caculate member offset.

return ret;
}

size_t hash_value(input_element_desc const& v)
{
size_t seed = 0;

boost::hash_combine(seed, v.aligned_byte_offset);
boost::hash_combine(seed, static_cast<size_t>(v.data_format) );
boost::hash_combine(seed, v.input_slot);
boost::hash_combine(seed, v.instance_data_step_rate);
boost::hash_combine(seed, v.semantic_index);
boost::hash_combine(seed, v.semantic_name);
boost::hash_combine(seed, static_cast<size_t>(v.slot_class) );

return seed;
}

size_t hash_value(input_layout const& v)
{
return boost::hash_range( v.desc_begin(), v.desc_end() );
}