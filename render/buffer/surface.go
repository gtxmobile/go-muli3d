package buffer

import (
	"../../common"
	"fmt"
	"github.com/stretchr/testify/assert"
)
type Surface struct {
	Elem_size_ uint32
	Sample_count_ uint32
	Size_		[4]uint32
	Format_		Pixel_format
	Datas_		[16]byte
}
func New_surface( w,  h,  samp_count uint32, fmt Pixel_format) *Surface{
	surf := &Surface{}
	surf.Format_ = fmt
 	surf.Size_ = [4]uint32{w,h, 1, 0}
 	surf.Sample_count_ = samp_count
 	surf.Elem_size_ = color_infos[fmt].size

	//#if SALVIA_TILED_SURFACE
	//tile_size_[0] = (Width + TILE_SIZE - 1) >> TILE_BITS
	//tile_size_[1] = (Height + TILE_SIZE - 1) >> TILE_BITS
	//
	//tile_mode_ = true
	//if ((TILE_SIZE > size_[0]) || (TILE_SIZE > size_[1])){
	//tile_mode_ = false
	//}
	//
	//if (tile_mode_){
	//datas_.resize(tile_size_[0] * tile_size_[1] * TILE_SIZE * TILE_SIZE * sample_count_ * elem_size_)
	//}
	//else{
	//datas_.resize(size_[0] * size_[1] * sample_count_ * elem_size_)
	//}
	//#else
	//surf.Datas_.resize( surf.Pitch() * h )
	//#endif

	to_rgba32_func_         := get_convertor_func(pixel_format_color_rgba32f, format_)
	from_rgba32_func_       := get_convertor_func(format_, pixel_format_color_rgba32f)
	to_rgba32_array_func_   := get_array_convertor_func(pixel_format_color_rgba32f, format_)
	from_rgba32_array_func_ := get_array_convertor_func(format_, pixel_format_color_rgba32f)
	lerp_1d_func_           := get_lerp_1d_func(format_)
	lerp_2d_func_           := get_lerp_2d_func(format_)
}

func (s *Surface) Width() uint32{
	return s.Size_[0]
}

func (s *Surface) Height() uint32{
	return s.Size_[1]
}
func (s *Surface) Pitch() uint32{
	return s.Width() * s.Sample_count_ * s.Elem_size_
}

 func (surf *Surface) Make_mip_surface(filter common.Filter_type) *Surface {
 	mip_w := ( surf.Width()  + 1 ) / 2;
 	mip_h := ( surf.Height() + 1 ) / 2;

	ret := New_surface(mip_w, mip_h, surf.Sample_count_, surf.Format_);
	var x,y,s uint32
	switch (filter){
	case common.Filter_point:
		for y = 0; y < mip_h; y++{
			for x = 0; x < mip_w; x++{
				for s = 0; s < surf.Sample_count_; s++{
					c := surf.Get_texel(x*2, y*2, s);
					ret.Set_texel(x, y, s, c);
				}
			}
		}
		break;

	case common.Filter_linear:
		for y = 0; y < mip_h; y++{
			for x = 0; x < mip_w; x++{
				for s = 0; s < surf.Sample_count_; s++{
					 c := [4]Color_rgba32f{
						surf.Get_texel(x*2+0, y*2+0, s),
						surf.Get_texel(x*2+1, y*2+0, s),
						surf.Get_texel(x*2+0, y*2+1, s),
						surf.Get_texel(x*2+1, y*2+1, s)
					}
					Color_rgba32f dst_color((c[0].get_vec4() + c[1].get_vec4() + c[2].get_vec4() + c[3].get_vec4()) * 0.25);
					ret.set_texel(x, y, s, dst_color);
				}
			}
		}
		break;
	}

	return ret;
}

func (surf *Surface) Smap( mapped *Internal_mapped_resource, mm map_mode)fundations.Result{
	//#if SALVIA_TILED_SURFACE
	//	//// Unimplemented
	//	//this->untile(mapped_data_);
	//	//#else
	switch mm{
		case common.Map_read:
			mapped.data = mapped.reallocator(datas_.size())
			copy( mapped.data, datas_.data())
			break
		case common.Map_read_write:
		case common.Map_write_no_overwrite:
		case common.Map_write_discard:
		case common.Map_write:
			mapped.data = datas_.data()
		break
	}

	mapped.row_pitch = surf.Pitch()
	mapped.depth_pitch = mapped.row_pitch * size_[1]

	return common.Ok
	//#endif
}

func (surf *Surface) Unmap(*internal_mapped_resource, map_mode) fundations.Result{
	// No intermediate buffer needed in linear mode.
	return common.Ok
}

func (surf *Surface) Resolve(target *Surface){
	if target.Sample_count_ != 1{
		panic("Resolve's target can't be a multi-sample surface")
	}

	var clr,tmp color_rgba32f

	for y := 0; y < size_[1]; y++ {
		for x := 0; x < size_[0]; x++ {
			clr = color_rgba32f(0, 0, 0, 0)
			for s := 0; s < surf.Sample_count_; s++	{
				surf.To_rgba32_func_( &tmp, texel_address(x, y, s) )
				clr.get_vec4() += tmp.get_vec4()
			}
			clr.get_vec4() /= float64(surf.Sample_count_)

			target.set_texel(x, y, 0, clr)
		}
	}
}

func (surf *Surface)Get_texel(x,y, sample uint32)Color_rgba32f{
	var color Color_rgba32f;
	to_rgba32_func_(&color, texel_address(x, y, sample))
	return color;
}

func (surf *Surface)Get_texel_void( color interface{},  x,  y,  sample uint32){
	//memcpy(color, Texel_address(x, y, sample), elem_size_);
}

func (surf *Surface)Get_texel_6( x0,  y0,  x1,  y1 uint32,  tx,  ty float64, sample uint32) Color_rgba32f{
	void const* addrs[] ={
	texel_address(x0, y0, sample),
	texel_address(x1, y0, sample),
	texel_address(x0, y1, sample),
	texel_address(x1, y1, sample)
	};

	return lerp_2d_func_(addrs[0], addrs[1], addrs[2], addrs[3], tx, ty)
}

func (surf *Surface)Set_texel( x,  y,  sample uint32, color *Color_rgba32f){
	from_rgba32_func_(surf.Texel_address(x, y, sample), &color)
}

func (surf *Surface)set_texel( x,  y,  sample uint32,  color interface{}){
	memcpy(texel_address(x, y, sample), color, surf.Elem_size_)
}

func (surf *Surface)fill_texels( sx,  sy,  width,  height uint32, color *Color_rgba32f){
	 var pix_clr [4 * 4 * 8]byte
	 from_rgba32_func_(pix_clr, color)

	 //if SALVIA_TILED_SURFACE
	 if tile_mode_{
		if ((0 == sx) && (0 == sy) && (width == surf.Size_[0]) && (height == surf.Size_[1])){
			for x := 0; x < TILE_SIZE; x++ {
				for s := 0; s < surf.Sample_count_;s++{
					memcpy(surf.Datas_[surf.texel_offset(x, 0, s)], pix_clr, surf.Elem_size_)
				}
			}
			for  y := 1; y < TILE_SIZE; y++{
				memcpy(suef.Datas_[surf.texel_offset(0, y, 0)], surf.Datas_[surf.texel_offset(0, 0, 0)], TILE_SIZE * surf.Sample_count_ * surf.Elem_size_)
			}
			for tx := 1; tx < tile_size_[0]; tx++{
				memcpy(surf.Datas_[surf.texel_offset(tx << TILE_BITS, 0, 0)], surf.Datas_[surf.texel_offset(0, 0, 0)],
					TILE_SIZE * TILE_SIZE * surf.Sample_count_ * surf.Elem_size_);
			}
			for ty := 1; ty < tile_size_[1]; ty++{
				memcpy(surf.Datas_[surf.texel_offset(0, ty << TILE_BITS, 0)], surf.Datas_[surf.texel_offset(0, 0, 0)], TILE_SIZE * TILE_SIZE * tile_size_[0] * surf.Sample_count_ * surf.Elem_size_)
			}
		} else{
			begin_tile_x := sx >> TILE_BITS
			begin_x_in_tile := sx & TILE_MASK
			end_tile_x := (sx + width - 1) >> TILE_BITS
			end_x_in_tile := (sx + width - 1) & TILE_MASK
			{
				for x := begin_x_in_tile; x < TILE_SIZE; x++{
					for s := 0; s < surf.Sample_count_;  s++{
						memcpy(surf.Datas_[surf.texel_offset((begin_tile_x << TILE_BITS) + x, sy, s)], pix_clr, surf.Elem_size_)
					}
				}
				for y := sy + 1; y < sy + height; y++{
					memcpy(surf.Datas_[surf.texel_offset(sx, y, 0)],
					surf.Datas_[surf.texel_offset(sx, sy, 0)], (TILE_SIZE - begin_x_in_tile) * surf.Sample_count_ * surf.Elem_size_)
				}
			}

			for tx := begin_tile_x + 1; tx < end_tile_x; tx++{
				for x := 0; x < TILE_SIZE; x++{
					for s := 0; s < surf.Sample_count_; s++{
						memcpy(surf.Datas_[surf.texel_offset((tx << TILE_BITS) + x, sy, s)],pix_clr, surf.Elem_size_)
					}
				}
				for y := sy + 1; y < sy + height; y++{
					memcpy(surf.Datas_[surf.texel_offset(tx << TILE_BITS, y, 0)],
						surf.Datas_[surf.texel_offset(tx << TILE_BITS, sy, 0)], TILE_SIZE * surf.Sample_count_ * surf.Elem_size_)
				}
			}

			{
				for x := 0; x < end_x_in_tile; x++{
					for s := 0; s < surf.Sample_count_; s++{
						memcpy(surf.Datas_[surf.texel_offset((end_tile_x << TILE_BITS) + x, sy, s)],
						pix_clr, surf.Elem_size_)
					}
				}
				for  y := sy + 1; y < sy + height; y++{
					memcpy(surf.Datas_[surf.texel_offset(end_tile_x << TILE_BITS, y, 0)],
						surf.Datas_[surf.texel_offset(end_tile_x << TILE_BITS, sy, 0)], end_x_in_tile * surf.Sample_count_ * surf.Elem_size_);
				}
			}
		}
	} else{
		for x := sx; x < sx + width; x++ {
			for s := 0; s < surf.Sample_count_; s++{
				memcpy(surf.Datas_[((surf.Size_[0] * sy + x) * surf.Sample_count_ + s) * surf.Elem_size_], pix_clr, surf.Elem_size_)
			}
		}
		for y := sy + 1; y < sy + height; y++{
			memcpy(surf.Datas_[(surf.Size_[0] * y + sx) * surf.Sample_count_ * surf.Elem_size_],
				surf.Datas_[(surf.Size_[0] * sy + sx) * surf.Sample_count_ * surf.Elem_size_],
				surf.Sample_count_ * surf.Elem_size_ * width)
		}
	}
	//#else

	for x := sx; x < sx + width; x++{
		for s := 0; s < surf.Sample_count_; s++{
			memcpy(surf.Datas_[((surf.Size_[0] * sy + x) * surf.Sample_count_ + s) * surf.Elem_size_], pix_clr, surf.Elem_size_)
		}
	}

	for  y := sy + 1; y < sy + height; y++ {
		memcpy(surf.Datas_[(surf.Size_[0] * y + sx) * surf.Sample_count_ * surf.Elem_size_], surf.Datas_[(surf.Size_[0] * sy + sx) * surf.Sample_count_ * surf.Elem_size_], surf.Sample_count_ * surf.Elem_size_ * width);
	}
	//#endif
}

func (surf *Surface)Fill_texels( color *Color_rgba32f){
	surf.fill_texels(0, 0, surf.Size_[0], surf.Size_[1], color)
}

func (surf *Surface)texel_offset( x,  y,  sample uint32) {
	//#if SALVIA_TILED_SURFACE
	//if (tile_mode_)
	//{
	//size_t const tile_x		= x >> TILE_BITS;
	//size_t const tile_y		= y >> TILE_BITS;
	//size_t const x_in_tile	= x & TILE_MASK;
	//size_t const y_in_tile	= y & TILE_MASK;
	//return (((tile_y * tile_size_[0] + tile_x) * TILE_SIZE * TILE_SIZE + (y_in_tile * TILE_SIZE + x_in_tile)) * sample_count_ + sample) * elem_size_;
	//}
	//else
	//{
	//return ((y * size_[0] + x) * sample_count_ + sample) * elem_size_;
	//}
	//#else
	return ((y * surf.Size_[0] + x) * surf.Sample_count_ + sample) * surf.Elem_size_;
	//#endif
}

//#if SALVIA_TILED_SURFACE
//void surface::tile(internal_mapped_resource const& mapped)
//{
//if (tile_mode_)
//{
//for (size_t ty = 0; ty < tile_size_[1]; ++ ty)
//{
//const size_t by = ty << TILE_BITS;
//const size_t rest_height = std::min(size_[1] - by, TILE_SIZE);
//for (size_t tx = 0; tx < tile_size_[0]; ++ tx)
//{
//const size_t bx = tx << TILE_BITS;
//const size_t tile_id = ty * tile_size_[0] + tx;
//const size_t rest_width = std::min(size_[0] - bx, TILE_SIZE);
//for (size_t y = 0; y < rest_height; ++ y)
//{
//memcpy(&datas_[((tile_id * TILE_SIZE + y) * TILE_SIZE) * sample_count_ * elem_size_],
//&tile_data[((by + y) * size_[0] + bx) * sample_count_ * elem_size_],
//rest_width * sample_count_ * elem_size_);
//}
//}
//}
//}
//else
//{
//datas_ = tile_data;
//}
//}
//
//void surface::untile(std::vector<byte>& untile_data)
//{
//if (tile_mode_)
//{
//for (size_t ty = 0; ty < tile_size_[1]; ++ ty)
//{
//const size_t by = ty << TILE_BITS;
//const size_t rest_height = std::min(size_[1] - by, TILE_SIZE);
//for (size_t tx = 0; tx < tile_size_[0]; ++ tx)
//{
//const size_t bx = tx << TILE_BITS;
//const size_t tile_id = ty * tile_size_[0] + tx;
//const size_t rest_width = std::min(size_[0] - bx, TILE_SIZE);
//for (size_t y = 0; y < rest_height; ++ y)
//{
//memcpy(&untile_data[((by + y) * size_[0] + bx) * sample_count_ * elem_size_],
//&datas_[((tile_id * TILE_SIZE + y) * TILE_SIZE) * sample_count_ * elem_size_],
//rest_width * sample_count_ * elem_size_);
//}
//}
//}
//}
//else
//{
//untile_data = datas_;
//}
//}
//
//#endif

func (surf *Surface) Texel_address( x,  y,  sample uint32) uintptr{
	ret := surf.Datas_ + surf.texel_offset(x, y, sample)
	return uintptr(&ret)
}

func (surf *Surface)texel_address( x,  y,  sample uint32) uintptr{
	ret := surf.Datas_ + surf.texel_offset(x, y, sample)
	return uintptr(&ret)
}
