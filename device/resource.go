package device
import (
	"../render"
	"../common"
	"github.com/lxn/win"
)

//template<typename FIColorT>
//func copy_image_to_surface_impl(surf *render.Surface, image *win.FIBITMAP,typename FIUC<FIColorT>::CompT default_alpha = (typename FIUC<FIColorT>::CompT)(0) ) bool{
func Copy_image_to_surface_impl(surf *render.Surface, image *win.BITMAP,default_alpha interface{} ) bool{

	if image == nil{
		return false
	}

	image_pitch := FreeImage_GetPitch(image);
	image_bpp := (FreeImage_GetBPP(image) >> 3);
	surface_format := surf.Format_;
	inter_format = salvia_rgba_color_type<FIColorT>::fmt;
	source_line := FreeImage_GetBits(image);
	var x,y uint32
	for y = 0; y < surf.Height(); y++ {
		src_pixel = source_line
		for  x = 0 ;x < surf.Width() ; x++{
			//FIUC<FIColorT> uc((typename FIUC<FIColorT>::CompT*)src_pixel, default_alpha)
			//typename salvia_rgba_color_type<FIColorT>::type c(uc.r, uc.g, uc.b, uc.a)
			render.Convert(surface_format, inter_format, surf.Texel_address(x, y, 0), &c)
			src_pixel += image_bpp
		}
		source_line += image_pitch
	}

	return true
}

// Copy region of image to dest region of surface.
// If the size of source and destination are different, it will be stretch copy with bi-linear interpolation.
func Copy_image_to_surface(surf *render.Surface, img *win.BITMAP)bool{
	image_type := FreeImage_GetImageType( img )

	if(image_type == FIT_RGBAF) {
		return copy_image_to_surface_impl<FIRGBAF>(surf, img)
	}
	if(image_type == FIT_BITMAP) {
		if(FreeImage_GetColorType(img) == FIC_RGBALPHA) {
			return copy_image_to_surface_impl<RGBQUAD>(surf, img)
		}else	{
			return copy_image_to_surface_impl<RGBTRIPLE>(surf, img)
		}
	}
	return false
}

// Load image file to new texture
func Load_texture( rend *render.Renderer, filename *string, tex_format render.Pixel_format) *render.Texture{
	img := load_image(filename)
	var ret *render.Texture
	src_w := FreeImage_GetWidth(img)
	src_h := FreeImage_GetHeight(img)

	ret = (*rend).Create_tex2d(src_w, src_h, 1, tex_format)

	if( !copy_image_to_surface(ret.Subresource(0), img) ) {
		//ret.reset()
	}

	FreeImage_Unload(img)

	return ret
}


// Create cube texture by six images.
// Size of first texture is the size of cube face.
// If other textures are not same size as first, just stretch it.
func Load_cube(rend *render.Renderer, filenames []string, tex_format render.Pixel_format) *render.Texture{
	ret :=render.Texture{}

	image_deleter  := [](FIBITMAP* bmp){ FreeImage_Unload(bmp) }
 	tex_width  := 0
 	tex_height := 0
	var i_cubeface uint32= 0
	for ; i_cubeface < 6; i_cubeface++ {
		//std::unique_ptr<FIBITMAP, decltype(image_deleter)> cube_img(load_image(filenames[i_cubeface]), image_deleter);

		if cube_img.get() == nullptr{
			//ret.reset();
			return &ret;
		}

 		img_w := FreeImage_GetWidth (cube_img.get());
 		img_h := FreeImage_GetHeight(cube_img.get());

		// Create texture cube while first face is created.
		//if ret != nil {
		//	tex_width  = img_w;
		//	tex_height = img_h;
		//	ret = *(*rend).Create_texcube(img_w, img_h, 1, tex_format);
		//} else {
		//	if (tex_width != img_w || tex_height != img_h) {
		//		//ret.reset();
		//		return ret;
		//	}
		//}

		cube_tex := ret
		face_surface := cube_tex.Subresource(i_cubeface, 0);
		copy_image_to_surface( face_surface, cube_img.get() );
	}

	return ret;
}



// Save surface as PNG or HRD formatted file.
func Save_surface(rend *render.Renderer , surf *render.Surface, filename string, image_format render.Pixel_format){
	fit := FIT_UNKNOWN
	fif := FIF_UNKNOWN

	var image *FIBITMAP;
	surface_width  := surf.Width()
	surface_height := surf.Height()

	switch(image_format){
		case pixel_format_color_bgra8:
			fit = FIT_BITMAP;
			fif = FIF_PNG;
			image = FreeImage_AllocateT(fit, surface_width, surface_height, 32, 0x0000FF, 0x00FF00, 0xFF0000);
			break;
		case pixel_format_color_rgb32f:
			fit := FIT_RGBF;
			fif = FIF_HDR;
			image = FreeImage_AllocateT(fit, surface_width, surface_height, 96);
			break;
		default:
			EFLIB_ASSERT(false, "Unsupport format was used.");
		return;
	}

	mapped := Mapped_resource{}
	rend.Smap(mapped, surf, common.Map_read)

	byte* 		 surf_data = reinterpret_cast<byte*>(mapped.data);
	byte*		 img_data = FreeImage_GetBits(image);
	pixel_format surf_format = surf->get_pixel_format();
	size_t		 height = surf->height();
	size_t		 width = surf->width();

	for(size_t y = 0; y < height; ++y){
	pixel_format_convertor::convert_array(
	image_format, surf_format,
	img_data, surf_data, int(width)
	);
	surf_data +=  mapped.row_pitch;
	img_data += FreeImage_GetPitch(image);
	}

	rend->unmap();

	FreeImage_Save(fif, image, to_ansi_string(filename).c_str());
	FreeImage_Unload(image);
}