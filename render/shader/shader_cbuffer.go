package shader
import (
	"../buffer"
)
type Shader_cdata_type uint32
const (
	Sdt_none Shader_cdata_type = iota
	Sdt_pod
	Sdt_sampler
)

type Shader_cdata struct {
	offset			uint32
	length			uint32
	array_size		uint32
}
shader_cdata()
: offset(0), length(0), array_size(0)
{
}


type	Variable_table	map[string]Shader_cdata
type 	Sampler_table 	map[string]*buffer.Sampler
type Shader_cbuffer struct {
	variables_		Variable_table
	data_memory_	[]byte
	samplers_		Sampler_table
	textures_		[]*buffer.Texture
}


virtual void set_sampler(eflib::fixed_string const& name, sampler_ptr const& samp);
virtual void set_variable(eflib::fixed_string const& name, void const* data, size_t data_length);

variable_table const&	variables() const
{
return variables_;
}

sampler_table  const&	samplers()  const
{
return samplers_;
}

void const* data_pointer(shader_cdata const& cdata) const
{
if(cdata.length == 0)
{
return nullptr;
}
return data_memory_.data() + cdata.offset;
}

void copy_from(shader_cbuffer const* src)
{
*this = *src;
}
private:
variable_table				variables_;
std::vector<char>			data_memory_;
sampler_table				samplers_;
std::vector<texture_ptr>	textures_;
};

void shader_cbuffer::set_sampler(eflib::fixed_string const& name, sampler_ptr const& samp)
{
auto iter = samplers_.find(name);
if (iter != samplers_.end() )
{
iter->second = samp;
}
else
{
samplers_.emplace(name, samp);
}
}

void shader_cbuffer::set_variable(eflib::fixed_string const& name, void const* data, size_t data_length)
{
shader_cdata* existed_cdata = nullptr;

auto existed_variable_iter = variables_.find(name);
if(existed_variable_iter != variables_.end())
{
existed_cdata = &(existed_variable_iter->second);
}

size_t	offset = 0;
if(existed_cdata != nullptr && existed_cdata->length >= data_length)
{
offset = existed_cdata->offset;
existed_cdata->length = data_length;
}
else
{
offset = data_memory_.size();

shader_cdata cdata;
cdata.array_size= 0;
cdata.length	= data_length;
cdata.offset	= offset;
data_memory_.resize(data_memory_.size() + data_length);

if(existed_cdata == nullptr)
{
variables_.emplace(name, cdata);
}
else
{
*existed_cdata = cdata;
}
}

memcpy(data_memory_.data()+offset, data, data_length);
}
