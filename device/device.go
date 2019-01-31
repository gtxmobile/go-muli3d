package device

import (
	"../render"
	"github.com/lxn/win"
	//"github.com/go-gl/gl/v4.6-core/gl"
	//"github.com/go-gl/gl/v4.6-compatibility/gl"
	"golang.org/x/mobile/gl"
	"unsafe"
)

type Swap_chain interface {
	Present()
	Get_surface() 	*render.Surface
}
type Swap_chain_impl struct {

	Renderer_         *render.Renderer
	Surface_          *render.Surface
	Resolved_surface_ *render.Surface
}
func (sc Swap_chain_impl)Get_surface() *render.Surface{
	return sc.Surface_
}
func (sc *Swap_chain_impl)init(	renderer *render.Renderer,render_params *render.Renderer_parameters){
	sc.Renderer_ = renderer
	sc.Surface_ = (*sc.Renderer_).Create_tex2d(
		render_params.Backbuffer_width,
		render_params.Backbuffer_height,
		render_params.Backbuffer_num_samples,
		render_params.Backbuffer_format	).Subresource(0)

	if( render_params.Backbuffer_num_samples > 1 )	{
		sc.Resolved_surface_ = (*sc.Renderer_).Create_tex2d(
			render_params.Backbuffer_width,
			render_params.Backbuffer_height,
			1,
			render_params.Backbuffer_format).Subresource(0)
	} else{
		sc.Resolved_surface_ = sc.Surface_
	}
}


func (sc Swap_chain_impl) Present(){
	//sc.Renderer_.flush()
	if	sc.Resolved_surface_ != sc.Surface_ {
		sc.Surface_.Resolve(sc.Resolved_surface_)
	}

	sc.Present_impl()
}
func (sc *Swap_chain_impl) Present_impl(){

}
func Salviax_create_swap_chain_and_renderer(out_swap_chain *Swap_chain,out_renderer *render.Renderer,
	render_params *render.Renderer_parameters, renderer_type Renderer_types,swap_chain_type Swap_chain_types){

	//out_renderer.reset()
	//out_swap_chain.reset()

	//if	renderer_type == Renderer_sync{
	//	out_renderer = create_sync_renderer()
	//}

	if	renderer_type == Renderer_async{
		out_renderer = render.Create_async_renderer()
	}

	if	out_renderer != nil{
		return
	}

	if	swap_chain_type == Swap_chain_default{
		//#if defined(SALVIAX_D3D11_ENABLED)
		//swap_chain_type = salviax::swap_chain_d3d11
		//#elif defined(SALVIAX_GL_ENABLED)
		swap_chain_type = Swap_chain_gl
	//#else
	//	out_renderer.reset()
	//	return
	//#endif
	}
	if	swap_chain_type == Swap_chain_gl{
	//#if defined(SALVIAX_GL_ENABLED)
		*out_swap_chain = create_gl_swap_chain(out_renderer, render_params)
	//#endif
	} else if(swap_chain_type == Swap_chain_d3d11){
		//#if defined(SALVIAX_D3D11_ENABLED)
		//out_swap_chain = sc.create_d3d11_swap_chain(out_renderer, render_params)
		//#endif
	}
}

type Gl_swap_chain struct {
	Swap_chain_impl
	window_ 	win.HWND
	dc_			win.HDC
	glrc_		win.HGLRC
	tex_		uint32
	width_, height_		uint32
}


func create_gl_swap_chain(renderer *render.Renderer, params *render.Renderer_parameters) Swap_chain{
	if params == nil{
		return Swap_chain_impl{}
	}

	return new_gl_swap_chain(renderer, params)
}

func new_gl_swap_chain(renderer *render.Renderer, params *render.Renderer_parameters) *Gl_swap_chain {
	//
	sci := Swap_chain_impl{}
	sci.init(renderer, params)
	sc := Gl_swap_chain{
		sci,
		0,
		0,
		0,
		0,
		0,
		0}

		sc.window_ = win.HWND(params.Native_window)
		sc.initialize()
		return &sc
	}

func (gsc *Gl_swap_chain)delete_gl_swap_chain(){
		gl.DeleteTextu(1, &gsc.tex_)

}

func (gsc *Gl_swap_chain) initialize(){
		gsc.dc_ = win.GetDC(gsc.window_)

		pfd := win.PIXELFORMATDESCRIPTOR{}
		pfd_size :=unsafe.Sizeof(pfd)
		//memset(&pfd, 0, sizeof(pfd))
		pfd.NSize		= uint16(pfd_size)
		pfd.NVersion	= 1
		pfd.DwFlags		= win.PFD_DRAW_TO_WINDOW | win.PFD_SUPPORT_OPENGL | win.PFD_DOUBLEBUFFER
		pfd.IPixelType	= win.PFD_TYPE_RGBA
		pfd.CColorBits	= 32
		pfd.CDepthBits	= 0
		pfd.CStencilBits = 0
		pfd.ILayerType	= win.PFD_MAIN_PLANE

		pixelFormat := win.ChoosePixelFormat(gsc.dc_, &pfd)
		if pixelFormat == 0 {
			panic("pixelFormat")
		}

		win.SetPixelFormat(gsc.dc_, pixelFormat, &pfd)
		win.DescribePixelFormat(gsc.dc_, pixelFormat, pfd_size, &pfd)

		gsc.glrc_ = win.WglCreateContext(gsc.dc_)
		win.WglMakeCurrent(gsc.dc_, gsc.glrc_)

		{
			ext_str := glGetString(GL_EXTENSIONS)
			if ext_str.find("WGL_EXT_swap_control")	{
				type wglSwapIntervalEXTFUNC func(interval int) bool
				wglSwapIntervalEXT  := wglSwapIntervalEXTFUNC(win.WglGetProcAddress("wglSwapIntervalEXT"))
				wglSwapIntervalEXT(0)
			}
		}

		gl.GlDisable(gl.GL_LIGHTING)
		gl.GlDisable(gl.GL_CULL_FACE)
		gl.GlEnable(GL_TEXTURE_2D)

		gl.glTexEnvi(GL_TEXTURE_ENV, GL_TEXTURE_ENV_MODE, GL_REPLACE)
		gl.glPixelStorei(GL_PACK_ALIGNMENT, 1)
		gl.glPixelStorei(GL_UNPACK_ALIGNMENT, 1)
		gl.glGenTextures(1, gsc.tex_)
	}

	func (gsc *Gl_swap_chain) Present_impl()	{
		surface_width  := gsc.Resolved_surface_.Width()
		surface_height := gsc.Resolved_surface_.Height()
		surf_format := gsc.Resolved_surface_.get_pixel_format()

		glViewport(
			0, 0,uint32(surface_width),
			uint32(surface_height))

		var mapped Mapped_resource
		gsc.Renderer_.map(mapped, gsc.Resolved_surface_, map_read)
		dest := make([]byte,surface_width * surface_height * 4)
		for	y := 0 ;y < surface_height ;y++	{
			dst_line := dest.data() + y * surface_width * 4  //byte*
			src_line := (*byte)(mapped.data) + y * mapped.row_pitch //byte*
			convert_array(Pixel_format_color_rgba8, surf_format,dst_line, src_line, int(surface_width) )
		}
		renderer_.unmap()

		glBindTexture(GL_TEXTURE_2D, gsc.tex_)
		if (gsc.width_ < surface_width) || (gsc.height_ < surface_height){
			gsc.width_	= uint32(surface_width)
			gsc.height_ = uint32(surface_height)
			glTexImage2D(GL_TEXTURE_2D, 0, GL_RGBA8, width_, height_, 0, GL_RGBA, GL_UNSIGNED_BYTE, &dest[0])
			glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST)
			glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST)
		}else
		{
			glTexSubImage2D(GL_TEXTURE_2D, 0, 0, 0,
				GLsizei(surface_width),
				GLsizei(surface_height),
				GL_RGBA, GL_UNSIGNED_BYTE, dest.data())
		}

		fw := float64(surface_width) / gsc.width_
		fh := float64(surface_height) /gsc.height_

		glBegin(GL_TRIANGLE_STRIP)

		glTexCoord2f(0, 0)
		glVertex3f(-1, +1, 0)

		glTexCoord2f(fw, 0)
		glVertex3f(+1, +1, 0)

		glTexCoord2f(0, fh)
		glVertex3f(-1, -1, 0)

		glTexCoord2f(fw, fh)
		glVertex3f(+1, -1, 0)

		glEnd()

		win.SwapBuffers(gsc.dc_)
	}





