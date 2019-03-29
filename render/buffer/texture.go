package buffer

import (
	"math"
	"../../common"
)

type Resource_manager struct {

}

type Texture_inter interface {
	et_texture_type() common.Texture_type
	Gen_mipmap(filter common.Filter_type, auto_gen bool)
}

type Texture struct {
	Fmt_          Pixel_format
	Sample_count_ uint32
	Min_lod_      uint32
	Max_lod_      uint32
	Size_         [4]uint32 //uint4
	Surfs_        []*Surface
}

func max(a,b uint32) uint32{
	return uint32(math.Max(float64(a),float64(b)))
}
func Calc_lod_limit(sz [4]uint32) uint32{
	rv := 0
	max_sz := max(sz[0], max(sz[1], sz[2]) )

	for max_sz >0{
		max_sz >>= 1
		rv++
	}
	return uint32(rv)
}
func (t *Texture) Subresource(index uint32) *Surface{
	if t.Max_lod_ <= index && index <= index{
		return t.Surfs_[index]
	}
	return &Surface{}
}
func (t *Texture)size(subresource_index uint32) [4]uint32{
	subres := t.Subresource(subresource_index)
	if subres != nil{
		return [4]uint32{0, 0, 0, 0}
	}
	return subres.Size_
}

func (t *Texture) max_lod(miplevel uint32) {
	t.Max_lod_ = miplevel
}

func (t *Texture) min_lod(miplevel uint32){
	t.Min_lod_ = miplevel
}

type Texture_2d struct {
	*Texture
}

func (*Resource_manager) Create_texture_2d(width, height, num_samples uint32,fmt Pixel_format) *Texture_2d{
	return new_texture_2d(width, height, num_samples,fmt)
}
func (*Resource_manager) Create_texture_cube(width, height, num_samples uint32,fmt Pixel_format) *Texture_cube{
	return new_texture_cube(width, height, num_samples,fmt)
}

func new_texture_2d( width,  height,  num_samples uint32, format Pixel_format) *Texture_2d{
	var t Texture_2d
	t.Fmt_ = format
	t.Sample_count_ = num_samples
	t.Size_ = [4]uint32{width,height, 1, 0}
	t.Surfs_ = append(t.Surfs_,New_surface(width, height, num_samples, format))
	return &t
}
func (t *Texture_2d)Get_texture_type() common.Texture_type{
	return common.Texture_type_2d
}
func (t *Texture_2d)Gen_mipmap(filter common.Filter_type, auto_gen bool){
	if auto_gen{
		t.Max_lod_ = 0
		t.Min_lod_ = Calc_lod_limit(t.Size_) - 1
	}

	//t.Surfs_.reserve(min_lod_ + 1)

	for lod_level := t.Max_lod_;lod_level < t.Min_lod_;lod_level++{
		t.Surfs_ = append( t.Surfs_,t.Surfs_[len(t.Surfs_) - 1].Make_mip_surface(filter) )
	}
}

type Texture_cube struct {
	*Texture
}


func new_texture_cube( width,  height,  num_samples uint32,  format Pixel_format) *Texture_cube {
	var t *Texture_cube
	for i := 0; i < 6; i++{
		t.Surfs_ = append(t.Surfs_,New_surface( width, height, num_samples, format) )
	}
	return t
}

func (t *Texture_cube) Get_texture_type() common.Texture_type{
	return common.Texture_type_cube
}

func (t *Texture_cube)  Subresource(face, lod uint32) *Surface{
	return t.Texture.Subresource(lod * 6 + face)
}

func (t *Texture_cube) Gen_mipmap(filter common.Filter_type, auto_gen bool){
	if auto_gen	{
		t.Max_lod_ = 0;
		t.Min_lod_ = Calc_lod_limit(t.Size_) - 1;
	}

	//t.Surfs_.reserve( (t.Min_lod_ + 1) * 6 );

	for lod_level := t.Max_lod_; lod_level < t.Min_lod_; lod_level++	{
		for i_face := 0; i_face < 6; i_face++ {
			t.Surfs_= append(t.Surfs_, t.Subresource(uint32(i_face), lod_level).Make_mip_surface(filter) )
		}
	}
}