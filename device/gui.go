package device

import (
	"fmt"
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)
type Gui struct {
	Hinst_ win.HINSTANCE
	Main_wnd_ Win_windows
}
var g_gui *Gui
type Window interface {

	Set_idle_handler( handler interface{})
	Set_draw_handler( handler interface{})
	Set_create_handler( handler interface{})
	Set_title( string )
	View_handle_as_void()
	View_handle()
	Refresh()
	Show()

}
type Win_windows struct {
	Wnd_class_      win.ATOM
	App_            *Gui
	Wnd_class_name_ string
	Hwnd_           win.HWND
	on_idle         chan struct{}
	on_paint        chan struct{}
	on_create       chan struct{}
}
/** Event handlers
@{ */
func (*Win_windows)Set_idle_handler( handler interface{}) {}
func (*Win_windows)Set_draw_handler( handler interface{}) {}
func (*Win_windows)Set_create_handler( handler interface{}) {}
/** @} */

/** Properties @{ */
func (*Win_windows) View_handle_as_void() {}
func (*Win_windows) View_handle()         {}
func (*Win_windows)Set_title( string ){
	//( hwnd_
}
/** @} */

func (*Win_windows) Refresh() {}

func (w *Win_windows) Register_window_class(hinst win.HINSTANCE)win.ATOM{
	var wcex win.WNDCLASSEX

	wcex.CbSize = uint32(unsafe.Sizeof(win.WNDCLASSEX{}))

	wcex.Style			= win.CS_HREDRAW | win.CS_VREDRAW;
	wcex.LpfnWndProc	= syscall.NewCallback(w.Win_proc)
	wcex.CbClsExtra		= 0
	wcex.CbWndExtra		= 0
	wcex.HInstance		= win.GetModuleHandle(nil)
	wcex.HIcon			= win.LoadIcon(hinst, win.MAKEINTRESOURCE(win.IDI_APPLICATION))
	wcex.HCursor		= win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_ARROW))
	wcex.HbrBackground	= win.HBRUSH(win.GetStockObject(win.COLOR_WINDOW+1))
	wcex.LpszMenuName	= nil;
	wcex.LpszClassName	= syscall.StringToUTF16Ptr(w.Wnd_class_name_)
	wcex.HIconSm		= win.LoadIcon(wcex.HInstance, win.MAKEINTRESOURCE(win.IDI_APPLICATION))

	return win.RegisterClassEx(&wcex)
}
func (w *Win_windows) Process_message(message uint32,wParam ,lParam uintptr) uintptr {
	var ps win.PAINTSTRUCT
	//var hdc win.HDC

	switch message{
		case win.WM_CREATE:
			<-w.on_create
			break
		case win.WM_PAINT:
			win.BeginPaint(w.Hwnd_, &ps)
			<-w.on_paint
			win.EndPaint(w.Hwnd_, &ps)
			break;
		case win.WM_DESTROY:
			win.PostQuitMessage(0)
			break;
		default:
			return win.DefWindowProc(w.Hwnd_, message, wParam, lParam)
	}
	return 0
}


func ( w *Win_windows)Create(width, height int32) bool{
	if w.Wnd_class_ == 0{
		w.Wnd_class_ = w.Register_window_class( w.App_.Hinst_)
	}
	style := uint32(win.WS_OVERLAPPEDWINDOW)
	var rc =win.RECT{0, 0, width, height};
	win.AdjustWindowRect(&rc, style, false)
	//rc.Right - rc.Left, rc.Bottom - rc.Top
	win.CreateWindowEx(
		0,
		syscall.StringToUTF16Ptr(w.Wnd_class_name_),
		syscall.StringToUTF16Ptr(""),
		win.WS_OVERLAPPED | win.WS_CAPTION | win.WS_SYSMENU | win.WS_MINIMIZEBOX,
		0, 0, 0, 0, 0, 0, w.App_.Hinst_, nil)
	if w.Hwnd_ != 0{
		//err := win.GetLastError()
		//
		//LPTSTR lpMsgBuf;
		//FormatMessage(
		//FORMAT_MESSAGE_ALLOCATE_BUFFER |
		//FORMAT_MESSAGE_FROM_SYSTEM |
		//FORMAT_MESSAGE_IGNORE_INSERTS,
		//NULL,
		//err,
		//MAKELANGID(LANG_NEUTRAL, SUBLANG_DEFAULT),
		//(LPTSTR) &lpMsgBuf,
		//0, NULL );
		//// Display the error message and exit the process
		//std::wstringstream ss;
		//ss << L"Window created failed. Error: " << std::hex << eflib::to_wide_string( std::_tstring(lpMsgBuf) ) << ".";
		//OutputDebugStringW( ss.str().c_str() );
		//LocalFree(lpMsgBuf);
		fmt.Println("error msg")
		return false
	}

	win.ShowWindow(w.Hwnd_, win.SW_SHOW);
	win.UpdateWindow(w.Hwnd_);
	return true;
}
func (w *Win_windows) Show(){
	win.ShowWindow(w.Hwnd_, win.SW_SHOW);
}

func (w *Win_windows) Win_proc(hWnd win.HWND,msg uint32,wParam ,lParam uintptr)(uintptr){
	wnd := g_gui.Main_wnd_
	wnd.Hwnd_ = w.Hwnd_;
	return wnd.Process_message(msg, wParam, lParam)
}

func New_win_gui() *Gui{
	win.DefWindowProc(0, 0, 0, 0)
	var gui = Gui{Hinst_: win.GetModuleHandle(nil),
		Main_wnd_:Win_windows{}}
	return &gui
}

func Create_win_gui() *Gui{
	//if g_gui{
	//assert(false);
	//exit(1);
	//}
	g_gui := New_win_gui()
	return g_gui;
}

func (gui *Gui)Create_window( width,height int32) int{

	if !gui.Main_wnd_.Create(width, height) {
		//OutputDebugString( _EFLIB_T("Main window creation failed!\n") );
		return 0
	}
	return 1
}

func (gui *Gui) Run() int{
	// Message loop.
	gui.Main_wnd_.Show()

	var msg win.MSG
	for {
		if win.PeekMessage(&msg, 0, 0, 0, win.PM_REMOVE){
			if win.WM_QUIT == msg.Message {
				break; // WM_QUIT, exit message loop
			}

			win.TranslateMessage(&msg)
			win.DispatchMessage(&msg)
		} else{
			<-gui.Main_wnd_.on_idle
		}
	}

	return int(msg.WParam)
}
