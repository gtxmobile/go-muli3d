package main

import (
	"fmt"
	"math"

	"github.com/lxn/win"
	"unsafe"
	"syscall"
	"time"
	"./smath"
)


type Point_t smath.Vector_t
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




//=====================================================================
// 坐标变换
//=====================================================================
type Transform_t struct {
	World      smath.Matrix_t
	View       smath.Matrix_t
	Projection smath.Matrix_t
	Transform  smath.Matrix_t
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
func transform_init(ts *Transform_t, width ,height int32){
	aspect := float64(width)/float64(height)
	ts.World.Set_identity()
	matrix_set_identiry(&ts.View)
	matrix_set_perspective(&ts.Projection,math.Pi *0.5,aspect,1.0,500.0)
	ts.W = float64(width)
	ts.H = float64(height)
	transform_update(ts)

}
// 将矢量 x 进行 project
func transform_apply(ts *Transform_t,y , x *smath.Vector_t){
	matrix_apply(y,x,&ts.Transform)
}
// 检查齐次坐标同 cvv 的边界用于视锥裁剪
func transform_check_cvv(v *smath.Vector_t)(int){
	w := v.W
	check := 0
	if v.Z < 0 {check |= 1}
	if v.Z > w {check |= 2}
	if v.X < -w {check |= 4}
	if v.X > w {check |= 8}
	if v.Y < -w {check |= 16}
	if v.Y > w {check |= 32}
	return check
}
// 归一化，得到屏幕坐标
func transform_homogenize(ts *Transform_t, y , x *smath.Vector_t){
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
	Left Edge_t
	Right Edge_t
}

type Scanline_t struct {
	V Vertex_t
	Step Vertex_t
	X int32
	Y int32
	W int32
}

func vertex_rhw_init(v *Vertex_t){
	rhw := 1/ v.Pos.W
	v.Rhw *= rhw
	v.Tc.U *= rhw
	v.Tc.V *= rhw
	v.Color.R *= rhw
	v.Color.G *= rhw
	v.Color.B *= rhw
}

func vertex_interp(y,x1,x2 *Vertex_t,t float64){
	vector_interp((*Vector_t)(&y.Pos),(*Vector_t)(&x1.Pos),(*Vector_t)(&x2.Pos),t )
	y.Tc.U = interp(x1.Tc.U,x2.Tc.U,t)
	y.Tc.V = interp(x1.Tc.V,x2.Tc.V,t)
	y.Color.R = interp(x1.Color.R,x2.Color.R,t)
	y.Color.G = interp(x1.Color.G,x2.Color.G,t)
	y.Color.B = interp(x1.Color.B,x2.Color.B,t)
	y.Rhw = interp(x1.Rhw,x2.Rhw,t)
}

func vertex_division(y,x1,x2 *Vertex_t,w float64){
	inv := 1/w
	y.Pos.X = (x2.Pos.X - x1.Pos.X) * inv
	y.Pos.Y = (x2.Pos.Y - x1.Pos.Y) * inv
	y.Pos.Z = (x2.Pos.Z - x1.Pos.Z) * inv
	y.Pos.W = (x2.Pos.W - x1.Pos.W) * inv
	y.Tc.U = (x2.Tc.U - x1.Tc.U) * inv
	y.Tc.V = (x2.Tc.V - x1.Tc.V) * inv
	y.Color.R = (x2.Color.R-x1.Color.R) * inv
	y.Color.G = (x2.Color.G-x1.Color.G) * inv
	y.Color.B = (x2.Color.B-x1.Color.B) * inv
	y.Rhw = (x2.Rhw-x1.Rhw) * inv
}

func vertex_add(y,x *Vertex_t){
	y.Pos.X += x.Pos.X
	y.Pos.Y += x.Pos.Y
	y.Pos.Z += x.Pos.Z
	y.Pos.W += x.Pos.W
	y.Rhw += x.Rhw
	y.Tc.U += x.Tc.U
	y.Tc.V += x.Tc.V
	y.Color.R += x.Color.R
	y.Color.G += x.Color.G
	y.Color.B += x.Color.B
}
// 根据三角形生成 0-2 个梯形，并且返回合法梯形的数量
func trapezoid_init_triangle(trap *[2]Trapezold_t,p1,p2,p3 *Vertex_t)(int){
	if p1.Pos.Y > p2.Pos.Y { p1,p2 = p2,p1 }
	if p1.Pos.Y > p3.Pos.Y { p1,p3 = p3,p1 }
	if p2.Pos.Y > p3.Pos.Y { p2,p3 = p3,p2 }
	//直线的情况
	if p1.Pos.Y == p2.Pos.Y && p1.Pos.Y == p3.Pos.Y {return 0}
	if p1.Pos.X == p2.Pos.X && p1.Pos.X == p3.Pos.X {return 0}
	//

	if p1.Pos.Y == p2.Pos.Y {
		if p1.Pos.X > p2.Pos.X {p1,p2 = p2,p1}
		(*trap)[0].Top = p1.Pos.Y
		(*trap)[0].Bottom = p3.Pos.Y
		(*trap)[0].Left.V1 = *p1
		(*trap)[0].Left.V2 = *p3
		(*trap)[0].Right.V1 = *p2
		(*trap)[0].Right.V2 = *p3
		if (*trap)[0].Top < (*trap)[0].Bottom {
			return 1
		}else{
			return 0
		}
	}

	if p2.Pos.Y == p3.Pos.Y {
		if p2.Pos.X > p3.Pos.X {p3,p2 = p2,p3}
		(*trap)[0].Top = p1.Pos.Y
		(*trap)[0].Bottom = p3.Pos.Y
		(*trap)[0].Left.V1 = *p1
		(*trap)[0].Left.V2 = *p2
		(*trap)[0].Right.V1 = *p1
		(*trap)[0].Right.V2 = *p3
		if (*trap)[0].Top < (*trap)[0].Bottom {
			return 1
		}else{
			return 0
		}
	}

	(*trap)[0].Top = p1.Pos.Y
	(*trap)[0].Bottom = p2.Pos.Y
	(*trap)[1].Top = p2.Pos.Y
	(*trap)[1].Bottom = p3.Pos.Y

	k := (p3.Pos.Y - p1.Pos.Y) / (p2.Pos.Y  - p1.Pos.Y)
	x := p1.Pos.X + (p2.Pos.X - p1.Pos.X) * k

	if x <= p3.Pos.X {
		(*trap)[0].Left.V1 = *p1
		(*trap)[0].Left.V2 = *p2
		(*trap)[0].Right.V1 = *p1
		(*trap)[0].Right.V2 = *p3
		(*trap)[1].Left.V1 = *p2
		(*trap)[1].Left.V2 = *p3
		(*trap)[1].Right.V1 = *p1
		(*trap)[1].Right.V2 = *p3
	} else {
		(*trap)[0].Left.V1 = *p1
		(*trap)[0].Left.V2 = *p3
		(*trap)[0].Right.V1 = *p1
		(*trap)[0].Right.V2 = *p2
		(*trap)[1].Left.V1 = *p1
		(*trap)[1].Left.V2 = *p3
		(*trap)[1].Right.V1 = *p2
		(*trap)[1].Right.V2 = *p3
	}
	return 2
}

// 按照 Y 坐标计算出左右两条边纵坐标等于 Y 的顶点
func trapezoid_edge_interp(trap *Trapezold_t, y float64){
	s1 := trap.Left.V2.Pos.Y  - trap.Left.V1.Pos.Y
	s2 := trap.Right.V2.Pos.Y  - trap.Right.V1.Pos.Y
	t1 := (y - trap.Left.V1.Pos.Y) / s1
	t2 := (y - trap.Right.V1.Pos.Y) / s2
	vertex_interp(&trap.Left.V,&trap.Left.V1,&trap.Left.V2,t1)
	vertex_interp(&trap.Right.V,&trap.Right.V1,&trap.Right.V2,t2)
}

// 根据左右两边的端点，初始化计算出扫描线的起点和步长
func trapezoid_init_scan_line(trap *Trapezold_t,scanline *Scanline_t,y int32){
	width := trap.Right.V.Pos.X - trap.Left.V.Pos.X
	scanline.X = int32(trap.Left.V.Pos.X + 0.5)
	scanline.W = int32(trap.Right.V.Pos.X + 0.5) - scanline.X
	scanline.Y = y
	scanline.V = trap.Left.V
	if trap.Left.V.Pos.X >= trap.Right.V.Pos.X {
		scanline.W = 0
	}
	vertex_division(&scanline.Step,&trap.Left.V,&trap.Right.V,width)
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
	device_draw_primitive(device, &p3, &p4, &p1)
}

func draw_box(device *Device_t,theta float64){
	var m Matrix_t
	matrix_set_rotate(&m,-1,-0.5,1,theta)
	device.Transform.World = m
	transform_update(&device.Transform)
	draw_plane(device, 0, 1, 2, 3)
	draw_plane(device, 4, 5, 6, 7)
	draw_plane(device, 0, 4, 5, 1)
	draw_plane(device, 1, 5, 6, 2)
	draw_plane(device, 2, 6, 7, 3)
	draw_plane(device, 3, 7, 4, 0)
}
func camera_at_zero(device *Device_t,x,y,z float64){
	eye := Vector_t{x,y,z,1}
	at := Vector_t{0,0,0,1}
	up := Vector_t{0,0,1,1}
	matrix_set_lookat(&device.Transform.View,&eye,&at,&up)
	transform_update(&device.Transform)

}

func init_texture(device *Device_t){
	w :=256
	h :=256
	device.Texture = make([][]uint32,h)
	for j:=0;j<256;j++{
		device.Texture[j] = make([]uint32,w)
		for i:=0;i<256;i++{
			x := i/32
			y := j/32
			if ((x+y)&1) ==1 {
				device.Texture[j][i] = 0xffffff
			}else{
				device.Texture[j][i] = 0x3fbcef
			}
		}
	}
	if w > 1024 || h > 1024 {
		panic("must less than 1024*1024")
	}
	// 重新计算每行纹理的指针

	device.Tex_height = int32(h)
	device.Tex_width = int32(w)
	device.Max_u = float64(w-1)
	device.Max_v = float64(h-1)
}


func init(){
	fmt.Println("init")
}

func main(){
	fmt.Println("main")
	var device Device_t

	pos := 3.5
	var alpha float64 = 1
	kbhit := 0
	indicator := 0
	states := []int32{ RENDER_STATE_TEXTURE, RENDER_STATE_COLOR, RENDER_STATE_WIREFRAME }
	if screen_init(800, 600, "go mini3d") !=0 {
		return
	}

	device_init(&device,800,600,&screen_fb)
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
			//fmt.Println("Render State: ",device.Render_state)

		}else {
			kbhit = 0
		}

		draw_box(&device, alpha)
		screen_update()
		time.Sleep(1)
	}
}
