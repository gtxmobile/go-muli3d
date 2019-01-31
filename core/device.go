package core


//=====================================================================
// 渲染设备
//=====================================================================
type Device_t struct {
	Transform 	Transform_t			// 坐标变换器
	Width		int32				// 窗口宽度
	Height		int32				// 窗口高度
	Framebuffer [][]uint32			// 像素缓存：framebuffer[y] 代表第 y行
	Zbuffer		[][]float64			// 深度缓存：zbuffer[y] 为第 y行指针
	Texture 	[][]uint32			// 纹理：同样是每行索引
	Tex_width	int32				// 纹理宽度
	Tex_height	int32				// 纹理高度
	Max_u		float64				// 纹理最大宽度：tex_width - 1
	Max_v		float64				// 纹理最大高度：tex_height - 1
	Render_state	int32			// 渲染状态
	Background		uint32			// 背景颜色
	Foreground		uint32			// 线框颜色
}

var RENDER_STATE_WIREFRAME int32=1		// 渲染线框
var RENDER_STATE_TEXTURE int32=2		// 渲染纹理
var RENDER_STATE_COLOR   int32=4		// 渲染颜色

// 设备初始化，fb为外部帧缓存，非 NULL 将引用外部帧缓存（每行 4字节对齐）
func device_init(device *Device_t, width int32, height int32, fb *[]uint32) {

	//初始化framebuffer,zbuffer
	var i int32
	if fb != nil {
		for i = 0; i < height; i++ {
			s := width * i
			e := s+ width
			device.Framebuffer = append(device.Framebuffer,(*fb)[s:e])
			ztmp := make([]float64, width)
			device.Zbuffer = append(device.Zbuffer,ztmp)

		}

	}

	device.Tex_width = 2
	device.Tex_height = 2
	device.Max_u = 1.0
	device.Max_v = 1.0
	device.Width = width
	device.Height = height
	device.Background = 0xc0c0c0
	device.Foreground = 0
	transform_init(&device.Transform, width, height)
	device.Render_state = RENDER_STATE_WIREFRAME
}

func device_destroy(device *Device_t){

}
// 设置当前纹理
func device_set_texture(device *Device_t, bits *[][]uint32, pitch int64, w int32,h int32){


}

// 清空 framebuffer 和 zbuffer
func device_clear(device *Device_t,mode int32){
	var x ,y int32
	height := device.Height
	for y = 0; y < height; y++ {
		dst := device.Framebuffer[y]
		cc := (height - 1 - y) * 230 / (height - 1)
		cc = (cc << 16) | (cc << 8) | cc
		if mode == 0 { cc = int32(device.Background)}
		for x = 0;x<device.Width;  x++ {
			dst[x] = uint32(cc)
		}
	}
	for y =0; y < height; y++ {
		dst := device.Zbuffer[y]
		for x = 0;x<device.Width;  x++ {
			dst[x] = 0
		}
	}
}
// 画点
func device_pixel(device *Device_t,x int32,y int32, color uint32){
	if x < device.Width && y < device.Height{
		device.Framebuffer[y][x] = color
	}
}

// 画线
func device_draw_line(device *Device_t,x1 int32,y1 int32,x2 int32,y2 int32,c uint32){
	if x1 == x2 && y1 == y2{
		device_pixel(device,x1,y1,c)
	}else if x1 == x2{
		inc := -1
		if y1 <= y2{ inc =1}
		for y:=y1 ;y!=y2;y+=int32(inc){device_pixel(device,x1,y,c)}
		device_pixel(device,x2,y2,c)
	}else if y1 == y2{
		inc := -1
		if x1 <= x2{ inc =1}
		for x:=x1 ;x!=x2;x+=int32(inc){device_pixel(device,x,y1,c)}
		device_pixel(device,x2,y2,c)
	}else{
		dx := x1 -x2
		if x1 < x2 {dx = -dx}
		dy := y1 - y2
		if y1 < y2 {dy = -dy}
		var rem int32 = 0

		if dx >= dy{
			if x2 < x1 {
				x1,x2 = x2,x1
				y1,y2 = y2,y1
			}
			y := y1
			for x:=x1;x <= x2;x++ {
				device_pixel(device,x,y,c)
				rem += dy
				if rem >= dx{
					rem -= dx
					if y2 >= y1{
						y++
					}else{
						y--
					}
					device_pixel(device,x,y,c)
				}
			}
			device_pixel(device,x2,y2,c)
		}else{
			if y2 < y1{
				x1,x2 = x2,x1
				y1,y2 = y2,y1
			}
			x := x1
			for y:=y1;y<=y2;y++{
				device_pixel(device,x,y,c)
				rem += dx
				if rem >= dy{
					rem -=dy
					if x2 >= x1{
						x++
					}else{
						x--
					}
					device_pixel(device,x,y,c)
				}
			}
			device_pixel(device,x2,y2,c)
		}
	}
}

// 根据坐标读取纹理
func device_texture_read(device *Device_t,u float64,v float64) uint32{
	u = u *device.Max_u
	v = v *device.Max_v
	x := int32(u+0.5)
	y := int32(v+0.5)
	x = CMID(x,0,device.Tex_width -1)
	y = CMID(y,0,device.Tex_height -1)
	return device.Texture[x][y]
}

//=====================================================================
// 渲染实现
//=====================================================================

// 绘制扫描线
func device_draw_scanline(device *Device_t,scanline *Scanline_t){
	framebuffer := device.Framebuffer[scanline.Y]
	zbuffer := device.Zbuffer[scanline.Y]
	x := scanline.X
	w := scanline.W
	width := device.Width
	render_state := device.Render_state
	for ; w>0;w-- {
		if (x >= 0 && x<width){
			rhw := scanline.V.Rhw
			if rhw >= zbuffer[x]{
				iw:= 1/rhw
				zbuffer[x] = rhw
				if (render_state & RENDER_STATE_COLOR) >0{
					r:= scanline.V.Color.R *iw
					g:= scanline.V.Color.G *iw
					b:= scanline.V.Color.B *iw
					R := int32(r * 255)
					G := int32(g * 255)
					B := int32(b * 255)
					R = CMID(R, 0, 255)
					G = CMID(G, 0, 255)
					B = CMID(B, 0, 255)
					framebuffer[x] = uint32((R << 16) | (G << 8) | (B))
				}
				if (render_state & RENDER_STATE_TEXTURE) > 0{
					u := scanline.V.Tc.U * iw
					v := scanline.V.Tc.V * iw
					cc := device_texture_read(device,u,v)
					framebuffer[x] = cc
				}

			}
		}
		vertex_add(&scanline.V,&scanline.Step)
		if x >= width {break}
		x ++
	}

}
// 主渲染函数
func device_render_trap(device *Device_t,trap *Trapezold_t){
	var scanline Scanline_t
	top := int32(trap.Top + 0.5)
	bottom := int32(trap.Bottom + 0.5)
	for j := top;j<bottom; j++ {
		if j >= 0 && j < device.Height {
			trapezoid_edge_interp(trap,float64(j)+0.5)
			trapezoid_init_scan_line(trap,&scanline,j)
			device_draw_scanline(device,&scanline)
		}
		if j >= device.Height {break}
	}
}

// 根据 render_state 绘制原始三角形
func device_draw_primitive(device *Device_t, v1,v2,v3 *Vertex_t){
	var p1, p2, p3, c1, c2, c3 Vector_t
	render_state := device.Render_state
	// 按照 Transform 变化
	transform_apply(&device.Transform, &c1, (*Vector_t)(&v1.Pos))
	transform_apply(&device.Transform, &c2, (*Vector_t)(&v2.Pos))
	transform_apply(&device.Transform, &c3, (*Vector_t)(&v3.Pos))

	// 裁剪，注意此处可以完善为具体判断几个点在 cvv内以及同cvv相交平面的坐标比例
	// 进行进一步精细裁剪，将一个分解为几个完全处在 cvv内的三角形
	if  transform_check_cvv(&c1) != 0 {return}
	if  transform_check_cvv(&c2) != 0 {return}
	if  transform_check_cvv(&c3) != 0 {return}

	// 归一化
	transform_homogenize(&device.Transform, &p1, &c1)
	transform_homogenize(&device.Transform, &p2, &c2)
	transform_homogenize(&device.Transform, &p3, &c3)

	if (render_state & (RENDER_STATE_TEXTURE | RENDER_STATE_COLOR)) > 0{

		t1 := *v1
		t2 := *v2
		t3 := *v3
		t1.Pos = Point_t(p1)
		t2.Pos = Point_t(p2)
		t3.Pos = Point_t(p3)
		t1.Pos.W = c1.W
		t2.Pos.W = c2.W
		t3.Pos.W = c3.W
		vertex_rhw_init(&t1)	// 初始化 w
		vertex_rhw_init(&t2)	// 初始化 w
		vertex_rhw_init(&t3)	// 初始化 w

		var traps [2]Trapezold_t
		// 拆分三角形为0-2个梯形，并且返回可用梯形数量
		n := trapezoid_init_triangle(&traps, &t1, &t2, &t3)
		if n >= 1 {device_render_trap(device, &traps[0])}
		if n >= 2 {device_render_trap(device, &traps[1])}
	}

	if (render_state & RENDER_STATE_WIREFRAME) ==1 {		// 线框绘制
		device_draw_line(device, int32(p1.X), int32(p1.Y), int32(p2.X), int32(p2.Y), device.Foreground)
		device_draw_line(device, int32(p1.X), int32(p1.Y), int32(p3.X), int32(p3.Y), device.Foreground)
		device_draw_line(device, int32(p3.X), int32(p3.Y),int32( p2.X), int32(p2.Y), device.Foreground)
	}
}