package core

import (
	"fmt"
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

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
var screen_fb []uint32

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
	wc.HInstance = win.GetModuleHandle(nil)
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

	var lpBits unsafe.Pointer
	hDC := win.GetDC(screen_handle)
	screen_dc = win.CreateCompatibleDC(hDC)
	win.ReleaseDC(screen_handle, hDC)

	screen_hb = win.CreateDIBSection(screen_dc, &bi.BmiHeader, win.DIB_RGB_COLORS, &lpBits, 0, 0)
	switch screen_hb {
	case 0, win.ERROR_INVALID_PARAMETER:
		fmt.Println("CreateDIBSection failed")
		return -3
	}

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
	// Fill the bit map image
	screen_fb = (*[1<<23]uint32)(lpBits)[:]
	//frambuffer := (*[w*h*4]uint8)(lpBits)
	//*frambuffer = [w*h*4]uint8{}
	for i :=0;i<int(w*h);i++{
		screen_fb[i] = 0
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
