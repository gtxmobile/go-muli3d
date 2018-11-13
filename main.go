package main

import (
	"fmt"
	"math"

	"github.com/lxn/win"
	"unsafe"
	"syscall"
	"time"
	"github.com/lxn/walk/declarative"
)

type Matrix_t struct{
	M 	[4][4]float64
}
type Vector_t struct {
	X	float64
	Y 	float64
	Z 	float64
	W 	float64
}

type Point_t Vector_t
func CMID(x int32,min int32, max int32) int32 {
	if x<min{
		return min
	} else{
		if x>max{
			return max
		}else {
			return x
		}
	}
}

func interp(x1 float64,x2 float64,t float64) float64 {
	return x1+(x2-x1)*t
}

func vector_length(v *Vector_t) float64{
	sq := v.X * v.X + v.Y*v.Y + v.Z*v.Z
	return math.Sqrt((sq))
}



func vector_add(z *Vector_t,x *Vector_t, y *Vector_t){
	z.X = x.X + y.X
	z.Y = x.Y + y.Y
	z.Z = x.Z + y.Z
	z.W = 1.0
}

func vector_sub(z *Vector_t,x *Vector_t, y *Vector_t){
	z.X = x.X - y.X
	z.Y = x.Y - y.Y
	z.Z = x.Z - y.Z
	z.W = 1.0
}

func vector_dotproduct(x *Vector_t,y *Vector_t) float64 {
	return x.X * y.X + x.Y * y.Y + x.Z *y.Z
}

func vector_crossproduct(z *Vector_t,x *Vector_t,y*Vector_t){
	z.X = x.Y * y.Z - x.Z*y.Y
	z.Y = x.Z * y.X - x.X*y.Z
	z.Z = x.X * y.Y - x.Y*y.X
	z.W = 1.0
}

// 矢量插值，t取值 [0, 1]
func vector_interp(z *Vector_t, x1 *Vector_t,x2 *Vector_t, t float64){
	z.X = interp(x1.X,x2.X,t)
	z.Y = interp(x1.Y,x2.Y,t)
	z.Z = interp(x1.Z,x2.Z,t)
	z.W = 1.0
}

func vector_normalize(v *Vector_t){
	length := vector_length(v)
	if length != 0.0 {
		inv := 1.0/length
		v.X *= inv
		v.Y *= inv
		v.Z *= inv
	}
}

func matrix_add(c *Matrix_t,a *Matrix_t, b *Matrix_t){
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			c.M[i][j] = a.M[i][j] + b.M[i][j]
		}
	}
}

func matrix_sub(c *Matrix_t,a *Matrix_t, b *Matrix_t){
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			c.M[i][j] = a.M[i][j] - b.M[i][j]
		}
	}
}

func matrix_mul(c *Matrix_t,a *Matrix_t, b *Matrix_t){
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			c.M[j][i] = a.M[j][0] * b.M[0][i] +
				a.M[j][1] * b.M[1][i] +
				a.M[j][2] * b.M[2][i] +
				a.M[j][3] * b.M[3][i]
		}
	}
}

func matrix_scale(c *Matrix_t,a *Matrix_t, f float64){
	for i := 0; i<4 ; i++{
		for j := 0; j<4; j++{
			c.M[i][j] = a.M[i][j] * f
		}
	}
}

func matrix_apply(y *Vector_t, x *Vector_t,m *Matrix_t){
	X,Y,Z,W := x.X,x.Y,x.Z,x.W
	y.X = X * m.M[0][0] + Y * m.M[1][0]+Z * m.M[2][0] + W * m.M[3][0]
	y.Y = X * m.M[0][1] + Y * m.M[1][1]+Z * m.M[2][1] + W * m.M[3][1]
	y.Z = X * m.M[0][2] + Y * m.M[1][2]+Z * m.M[2][2] + W * m.M[3][2]
	y.W = X * m.M[0][3] + Y * m.M[1][3]+Z * m.M[2][3] + W * m.M[3][3]

}

func matrix_set_identiry(m *Matrix_t){
	m.M[0][0] ,m.M[1][1] , m.M[2][2] , m.M[3][3] = 1.0,1.0,1.0,1.0
	m.M[0][1] , m.M[0][2] , m.M[0][3] = 0,0,0
	m.M[1][0] , m.M[1][2] , m.M[1][3] = 0,0,0
	m.M[2][0] ,m.M[2][1] , m.M[2][3] = 0,0,0
	m.M[3][0] ,m.M[3][1] , m.M[3][2]  = 0,0,0
}

func matrix_set_zero(m *Matrix_t){
	m.M[0][0] ,m.M[1][1] , m.M[2][2] , m.M[3][3] = 0,0,0,0
	m.M[0][1] , m.M[0][2] , m.M[0][3] = 0,0,0
	m.M[1][0] , m.M[1][2] , m.M[1][3] = 0,0,0
	m.M[2][0] ,m.M[2][1] , m.M[2][3] = 0,0,0
	m.M[3][0] ,m.M[3][1] , m.M[3][2]  = 0,0,0
}

//平移
func matrix_set_translate(m *Matrix_t, x float64,y float64, z float64){
	matrix_set_identiry(m)
	m.M[3][0] = x
	m.M[3][1] = y
	m.M[3][2] = z

}
//平移
func matrix_set_scale(m *Matrix_t, x float64,y float64, z float64){
	matrix_set_identiry(m)
	m.M[0][0] = x
	m.M[1][1] = y
	m.M[2][2] = z

}
//从四元数构造旋转矩阵
func matrix_set_rotate(m *Matrix_t, x float64,y float64, z float64, theta float64){
	qsin := math.Sin(theta * 0.5)
	qcos := math.Cos(theta * 0.5)
	vec := Vector_t{x,y,z,1.0	}
	w := qcos
	vector_normalize(&vec)
	x = vec.X * qsin
	y = vec.Y * qsin
	z = vec.Z * qsin


	x2 := x * x
	y2 := y * y
	z2 := z * z

	xy := x * y
	xz := x * z
	yz := y * z

	wx := w * x
	wy := w * y
	wz := w * z
	m.M[0][0] = 1 - 2*y2 - 2*z2
	m.M[1][0] = 2*xy - 2*wz
	m.M[2][0] = 2*xz + 2*wy

	m.M[0][1] = 2*xy + 2*wz
	m.M[1][1] = 1-2*x2 -2*z2
	m.M[2][1] = 2*yz - 2*wx

	m.M[0][2] = 2*xz - 2*wy
	m.M[1][2] = 2*yz + 2*wx
	m.M[2][2] = 1-2*x2 - 2*y2

	m.M[0][3] , m.M[1][3],m.M[2][3] = 0,0,0
	m.M[3][0] , m.M[3][1],m.M[3][2] = 0,0,0
	m.M[3][3] = 1

	}

func matrix_set_lookat(m *Matrix_t,eye *Vector_t,at *Vector_t,up *Vector_t) {
	var xaxis,yaxis,zaxis Vector_t
	vector_sub(&zaxis,at,eye)
	vector_normalize(&zaxis)
	vector_crossproduct(&xaxis,up,&zaxis)
	vector_normalize(&xaxis)
	vector_crossproduct(&yaxis,&zaxis,&xaxis)

	m.M[0][0] = xaxis.X
	m.M[1][0] = xaxis.Y
	m.M[2][0] = xaxis.Z
	m.M[3][0] = -vector_dotproduct(&xaxis,eye)

	m.M[0][1] = yaxis.X
	m.M[1][1] = yaxis.Y
	m.M[2][1] = yaxis.Z
	m.M[3][1] = -vector_dotproduct(&yaxis,eye)

	m.M[0][2] = zaxis.X
	m.M[1][2] = zaxis.Y
	m.M[2][2] = zaxis.Z
	m.M[3][2] = -vector_dotproduct(&zaxis,eye)

	m.M[0][3],m.M[1][3],m.M[2][3] = 0,0,0
	m.M[3][3] = 1

}

func matrix_set_perspective(m *Matrix_t,fovy float64,aspect float64,zn float64,zf float64){
	fax := 1/math.Tan(fovy * 0.5)
	matrix_set_zero(m)
	m.M[0][0] = fax / aspect
	m.M[1][1] = fax
	m.M[2][2] = zf/(zf - zn)
	m.M[3][2] = - zn * zf /(zf-zn)
	m.M[2][3] = 1

}

//=====================================================================
// 坐标变换
//=====================================================================
type Transform_t struct {
	World      Matrix_t
	View       Matrix_t
	Projection Matrix_t
	Transform  Matrix_t
	W          float64
	H          float64
}
//矩阵更新

func transform_update(ts *Transform_t){
	var m Matrix_t
	matrix_mul(&m, &ts.World,&ts.View)
	matrix_mul(&ts.Transform,&m,&ts.Projection)

}



//初始化、设置屏幕长宽
func transform_init(ts *Transform_t, width int32,height int32){
	aspect := float64(width)/float64(height)
	matrix_set_identiry(&ts.World)
	matrix_set_identiry(&ts.View)
	matrix_set_perspective(&ts.Projection,math.Pi *0.5,aspect,1.0,500.0)
}
// 将矢量 x 进行 project
func transform_apply(ts *Transform_t,y *Vector_t, x *Vector_t){
	matrix_apply(y,x,&ts.Transform)
}
// 检查齐次坐标同 cvv 的边界用于视锥裁剪
func transform_check_cvv(v *Vector_t)(int){
	w := v.W
	check := 0
	if v.Z < 0 {check |= 1}
	if v.Z > w {check |= 2}
	if v.X < -w {check |= 4}
	if v.X > w {check |= 8}
	if v.Y < -w {check |= 16}
	if v.Y < w {check |= 32}
	return check
}
// 归一化，得到屏幕坐标
func transform_homogenize(ts *Transform_t, y *Vector_t, x *Vector_t){
	rhw := 1/x.W
	y.X = (x.X * rhw +1)*ts.W*0.5
	y.Y = (1-x.Y*rhw) * ts.H * 0.5
	y.Z = x.Z * rhw
	y.W = 1
}


//=====================================================================
// 几何计算：顶点、扫描线、边缘、矩形、步长计算
//=====================================================================
type Color_t struct {
	R float64
	G float64
	B float64
}

type Texcoord_t struct {
	U float64
	V float64
}

type Vertex_t struct {
	Pos Point_t
	Tc Texcoord_t
	Color Color_t
	Rhw float64
}

type Edge_t struct {
	V Vertex_t
	V1 Vertex_t
	V2 Vertex_t
}

type Trapezold_t struct{
	Top float64
	Bottom float64
	left Edge_t
	right Edge_t
}

type Scanline_t struct {
	V Vertex_t
	Step Vertex_t
	X int
	Y int
	Z int
}

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
func device_init(device *Device_t, width int32, height int32, fb *[][]uint32) {
	//need := sizeof(void*) * (height * 2 + 1024) + width * height * 8
	//char *ptr = (char*)malloc(need + 64);
	//char *framebuf, *zbuf;

	//assert(ptr);
	//device.Framebuffer = (IUINT32**)ptr;
	//device.zbuffer = (float**)(ptr + sizeof(void*) * height);
	//ptr += sizeof(void*) * height * 2;
	//device.texture = (IUINT32**) ptr;
	//ptr += sizeof(void*) * 1024;
	//framebuf = (char*)ptr;
	//zbuf = (char*)ptr + width * height * 4;
	//ptr += width * height * 8;
	//if (fb != NULL) framebuf = (char*)fb;
	//for j = 0; j < height; j++ {
	//	device->framebuffer[j] = (IUINT32*)(framebuf + width * 4 * j);
	//	device->zbuffer[j] = (float*)(zbuf + width * 4 * j);
	//}
	//初始化framebuffer,zbuffer
	var i int32
	for i = 0; i < height; i++ {
		ftmp := make([]uint32, width)
		device.Framebuffer = append(device.Framebuffer,ftmp)
		ztmp := make([]float64, width)
		device.Zbuffer = append(device.Zbuffer,ztmp)
	}

	device.Texture[0] = make([]uint32,64)
	device.Texture[1] = make([]uint32,64)
	//memset(device->texture[0], 0, 64);
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
func device_set_texture(device *Device_t, bits [][]uint32, pitch int64, w int32,h int32){
	if w > 1024 || h > 1024 {
		panic("must less than 1024*1024")
	}
	// 重新计算每行纹理的指针
	var j int32
	for j=0 ;j<h; j++{
		device.Texture[j] = bits[j]
	}

	device.Tex_height = h
	device.Tex_width = w
	device.Max_u = float64(w-1)
	device.Max_v = float64(h-1)

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
		for x = device.Width; x > 0;  x-- {
			dst[x] = uint32(cc)
		}
	}
	for y =0; y < height; y++ {
		dst := device.Zbuffer[y]
		for x = device.Width; x>0; x-- {
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
			x1,x2 = x2,x1
			y1,y2 = y2,y1
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

}
// 主渲染函数
func device_render_trap(device *Device_t,trap *Trapezold_t){

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
	if transform_check_cvv(&c1) != 0 {return}
	if  transform_check_cvv(&c2) != 0 {return}
	if  transform_check_cvv(&c3) != 0 {return}

	// 归一化
	transform_homogenize(&device.Transform, &p1, &c1)
	transform_homogenize(&device.Transform, &p2, &c2)
	transform_homogenize(&device.Transform, &p3, &c3)


	if render_state & RENDER_STATE_WIREFRAME ==1 {		// 线框绘制
		device_draw_line(device, int32(p1.X), int32(p1.Y), int32(p2.X), int32(p2.Y), device.Foreground)
		device_draw_line(device, int32(p1.X), int32(p1.Y), int32(p3.X), int32(p3.Y), device.Foreground)
		device_draw_line(device, int32(p3.X), int32(p3.Y),int32( p2.X), int32(p2.Y), device.Foreground)
	}
}

//=====================================================================
// Win32 窗口及图形绘制：为 device 提供一个 DibSection 的 FB
//=====================================================================
var screen_ob win.HBITMAP
var screen_w ,screen_h int32
var screen_exit int32 = 0
var screen_pitch int64
var screen_keys [512]int32
var screen_dc win.HDC
var screen_hb win.HBITMAP
var screen_handle win.HWND
var screen_fb *uint8

func screen_init(w int32,h int32,title string)(int){
	//hInst := win.GetModuleHandle(nil)
	//hIcon := win.LoadIcon(0, MAKEINTRESOURCE(IDI_APPLICATION))
	//hCursor := LoadCursor(0, MAKEINTRESOURCE(IDC_ARROW))
	var wc = win.WNDCLASSEX{uint32(unsafe.Sizeof(win.WNDCLASSEX{})),
		win.CS_BYTEALIGNCLIENT,
		syscall.NewCallback(screen_events),
		0,
		0,
		0,
		0,
		0,
		0,
		nil,
		syscall.StringToUTF16Ptr("SCREEN3.1415926"),
		0}
	var bi = win.BITMAPINFO {
		win.BITMAPINFOHEADER{
			uint32(unsafe.Sizeof(win.BITMAPINFOHEADER{})), w, -h, 1, 32, win.BI_RGB,
			uint32(w * h * 4), 0, 0, 0, 0 },nil}

	var rect = win.RECT { 0, 0, w, h }
	screen_close()
	wc.HbrBackground = win.HBRUSH(win.GetStockObject(win.BLACK_BRUSH))
	wc.HInstance = win.GetModuleHandle(nil);
	wc.HCursor = win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_ARROW))
	if win.RegisterClassEx(&wc) == 0{
		return -1
	}

	screen_handle = win.CreateWindowEx(
		0,
		syscall.StringToUTF16Ptr("SCREEN3.1415926"),
		syscall.StringToUTF16Ptr(title),
		win.WS_OVERLAPPED | win.WS_CAPTION | win.WS_SYSMENU | win.WS_MINIMIZEBOX,
		0, 0, 0, 0, 0, 0, wc.HInstance, nil)
	if screen_handle == 0 {
		return -2
	}

	var ptr unsafe.Pointer
	//screen_exit := 0
	hDC := win.GetDC(screen_handle)
	screen_dc = win.CreateCompatibleDC(hDC)
	win.ReleaseDC(screen_handle, hDC)

	screen_hb = win.CreateDIBSection(screen_dc, &bi.BmiHeader, win.DIB_RGB_COLORS, &ptr, 0, 0)
	switch screen_hb {
	case 0, win.ERROR_INVALID_PARAMETER:
		fmt.Println("CreateDIBSection failed")
		return -3
	}
	//if (screen_hb == 0){
	//	return -3
	//}
	screen_ob = win.HBITMAP(win.SelectObject(screen_dc, win.HGDIOBJ(screen_hb)))
	screen_w = w
	screen_h = h
	screen_pitch = int64(w * 4)

	win.AdjustWindowRect(&rect, uint32(win.GetWindowLong(screen_handle, win.GWL_STYLE)), false)
	wx := rect.Right - rect.Left
	wy := rect.Bottom - rect.Top
	sx := (win.GetSystemMetrics(win.SM_CXSCREEN) - wx) / 2
	sy := (win.GetSystemMetrics(win.SM_CYSCREEN) - wy) / 2
	if sy < 0 {
		sy = 0
	}
	win.SetWindowPos(screen_handle, 0, sx, sy, wx, wy, (win.SWP_NOCOPYBITS | win.SWP_NOZORDER | win.SWP_SHOWWINDOW))
	win.SetForegroundWindow(screen_handle)

	win.ShowWindow(screen_handle, win.SW_NORMAL)
	screen_dispatch()

	screen_fb = (*uint8)(ptr)
	for i :=0;i<int(w*h*4);i++{
		*screen_fb[i] = 0
	}
	// Fill the image
	bitmap_array := (*[1 << 30]byte)(unsafe.Pointer(lpBits))
	i := 0
	for y := im.Bounds().Min.Y; y != im.Bounds().Max.Y; y++ {
		for x := im.Bounds().Min.X; x != im.Bounds().Max.X; x++ {
			r, g, b, a := im.At(x, y).RGBA()
			bitmap_array[i+3] = byte(a >> 8)
			bitmap_array[i+2] = byte(r >> 8)
			bitmap_array[i+1] = byte(g >> 8)
			bitmap_array[i+0] = byte(b >> 8)
			i += 4
		}
	}
	return 0
}

func screen_close()(int){

	if screen_dc != 0 {
		if screen_ob != 0 {
			win.SelectObject(screen_dc, win.HGDIOBJ(screen_ob))
			screen_ob = 0
		}
		win.DeleteDC(screen_dc)
		screen_dc = 0
	}
	if screen_hb != 0 {
		win.DeleteObject(win.HGDIOBJ(screen_hb))
		screen_hb = 0
	}
	if screen_handle != 0 {
		win.CloseHandle(win.HANDLE(screen_handle))
		screen_handle = 0
	}
	return 0
}
func screen_events(hWnd win.HWND,msg uint32,wParam ,lParam uintptr)(uintptr){

	switch (msg) {
	case win.WM_CLOSE: screen_exit = 1; break
	case win.WM_KEYDOWN: screen_keys[wParam & 511] = 1; break
	case win.WM_KEYUP: screen_keys[wParam & 511] = 0; break
	default: return win.DefWindowProc(hWnd, msg, wParam, lParam)
	}
	return 0
}
func screen_dispatch(){
	var msg win.MSG
	for {
		if !win.PeekMessage(&msg, 0, 0, 0, win.PM_NOREMOVE){ break}
		if win.GetMessage(&msg, 0, 0, 0) == 0 {break}
		win.DispatchMessage(&msg)
	}
}

func screen_update() {
	hDC := win.GetDC(screen_handle)
	win.BitBlt(hDC, 0, 0, screen_w, screen_h, screen_dc, 0, 0, win.SRCCOPY)
	win.ReleaseDC(screen_handle, hDC)
	screen_dispatch()
}

//=====================================================================
// 主程序
//=====================================================================
var mesh = [8]Vertex_t{
{ Point_t{  1, -1,  1, 1 }, Texcoord_t{ 0, 0 }, Color_t{ 1.0, 0.2, 0.2 }, 1 },
{ Point_t{ -1, -1,  1, 1 }, Texcoord_t{ 0, 1 }, Color_t{ 0.2, 1.0, 0.2 }, 1 },
{ Point_t{ -1,  1,  1, 1 }, Texcoord_t{ 1, 1 }, Color_t{ 0.2, 0.2, 1.0 }, 1 },
{ Point_t{  1,  1,  1, 1 }, Texcoord_t{ 1, 0 }, Color_t{ 1.0, 0.2, 1.0 }, 1 },
{ Point_t{  1, -1, -1, 1 }, Texcoord_t{ 0, 0 }, Color_t{ 1.0, 1.0, 0.2 }, 1 },
{ Point_t{ -1, -1, -1, 1 }, Texcoord_t{ 0, 1 }, Color_t{ 0.2, 1.0, 1.0 }, 1 },
{ Point_t{ -1,  1, -1, 1 }, Texcoord_t{ 1, 1 }, Color_t{ 1.0, 0.3, 0.3 }, 1 },
{ Point_t{  1,  1, -1, 1 }, Texcoord_t{ 1, 0 }, Color_t{ 0.2, 1.0, 0.3 }, 1 },
}

func draw_plane(device *Device_t ,a,b,c,d int){
	p1 := mesh[a]
	p2 := mesh[b]
	p3 := mesh[c]
	p4 := mesh[d]
	p1.Tc.U = 0
	p1.Tc.V = 0
	p2.Tc.U = 0
	p2.Tc.V = 1
	p3.Tc.U = 1
	p3.Tc.V = 1
	p4.Tc.U = 1
	p4.Tc.V = 0
	device_draw_primitive(device, &p1, &p2, &p3)
}

func draw_box(device *Device_t,theta float64){
	var m Matrix_t
	matrix_set_rotate(&m,-1,-0.5,1,theta)
	device.Transform.World = m
	transform_update(&device.Transform)
	draw_plane(device, 0, 1, 2, 3)
}
func camera_at_zero(device *Device_t,x,y,z float64){
	eye := Vector_t{x,y,z,1}
	at := Vector_t{0,0,0,1}
	up := Vector_t{0,0,1,1}
	matrix_set_lookat(&device.Transform.View,&eye,&at,&up)
	transform_update(&device.Transform)

}

func init_texture(device *Device_t){
	var textrue [][]uint32
	for j:=0;j<256;j++{
		for i:=0;i<256;i++{
			x := i/32
			y := j/32
			if ((x+y)&1) ==1 {
				textrue[j][i] = 0xffffff
			}else{
				textrue[j][i] = 0x3fbcef
			}
		}
	}
	device_set_texture(device,textrue,256*4,256,256)
}


func init(){
	fmt.Println("init")
}

func main(){
	fmt.Println("main")
	var device Device_t
	//var inTe, outTE *walk.TextEdit

	//MainWindow{
	//	Title:	"SCREAMO",
	//	MinSize: Size{600,800},
	//	Layout: VBox{},
	//}.Run()
	pos := 3.5
	var alpha float64 = 1
	kbhit := 0
	indicator := 0
	states := []int32{ RENDER_STATE_TEXTURE, RENDER_STATE_COLOR, RENDER_STATE_WIREFRAME }
	screen_init(800, 600, "go mini3d")

	device_init(&device,800,600,screen_fb)
	camera_at_zero(&device,pos,0,0)

	init_texture(&device)
	device.Render_state = RENDER_STATE_TEXTURE
	for screen_exit == 0 && screen_keys[win.VK_ESCAPE] ==0 {
		screen_dispatch()
		device_clear(&device, 1)
		camera_at_zero(&device, pos, 0, 0)

		if screen_keys[win.VK_UP] != 0{ pos -= 0.01}
		if screen_keys[win.VK_DOWN] != 0{pos += 0.01}
		if screen_keys[win.VK_LEFT] != 0 {alpha += 0.01}
		if screen_keys[win.VK_RIGHT] != 0 {alpha -= 0.01}

		if screen_keys[win.VK_SPACE] !=0 {
			if kbhit == 0 {
				kbhit = 1
				indicator +=1
				if indicator >= 3 {indicator = 0}
				device.Render_state = states[indicator]
			}
		}else {
			kbhit = 0
		}

		draw_box(&device, alpha)
		screen_update()
		time.Sleep(1)
	}
}
