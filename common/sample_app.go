package common

import (
	"../device"
	"../render"
	"fmt"
	"github.com/lxn/win"
	"sort"
	"time"
)
type App_modes int


type Quit_conditions int

const (
	unknown App_modes = iota
	benchmark		// Run as benchmark. It will generate some benchmark results.
	test			// Run as regression test. It will generate final frames as image file for test.
	interactive	// Interactive mode.
	replay			// Play mode.
	count
)
const (
	user_defined	Quit_conditions = iota
	frame_limits
	time_out
)

type Frame_data struct {
	pipeline_stat render.Pipeline_statistics
	internal_stat render.Internal_statistics
	pipeline_prof render.Pipeline_profiles
}
type sample_app_data struct{
	Benchmark_name	string
	//mode
	Screen_width	uint32
	Screen_height	uint32
	Screen_aspect_ratio	uint32
	Screen_vp		render.Viewport

	is_sync_renderer	bool
	sc_type				device.Swap_chain_types
	mode				App_modes
	gui					*device.Gui
	swap_chain 			*device.Swap_chain
	renderer			*render.Renderer
	color_target		*render.Surface
	resolved_color_target	*render.Surface
	ds_target			*render.Surface

	total_elapsed_sec	float64
	elapsed_sec			float64
	frame_count			uint32

	//eflib::profiler				prof

	pipeline_stat_obj	*render.Async_object
	internal_stat_obj	*render.Async_object
	pipeline_prof_obj	*render.Async_object

	frame_profs			[]Frame_data
	quit_cond			Quit_conditions
	quit_cond_data		uint32
	runnable			bool
	quiting				bool
	//frame_timer			time.Timer
	frame_timer			MyTimer
	//second_timer		time.Timer
	second_timer		MyTimer
	frames_in_second	int
}

type sample_app struct {
	data_ *sample_app_data
	app_name string
}

func (app *sample_app) create_app(app_name string){
	app.data_.Benchmark_name = app_name
	app.data_.mode = unknown
	app.data_.quiting = false
	app.data_.runnable = true
	app.data_.frame_count = 0
	app.data_.total_elapsed_sec = 0.0
	app.data_.elapsed_sec = 0.0
	app.data_.gui = nil
	//app.data_.is_sync_renderer = boost::indeterminate
	app.data_.frames_in_second = 0
	app.data_.quit_cond = user_defined
	app.data_.quit_cond_data = 0
}

func (app *sample_app) init_params(){

}
func (app *sample_app) on_init(){

}
func (app *sample_app) on_frame(){

}

func (app *sample_app) init(){
	app.init_params()
	if	app.data_.runnable {
		if	app.data_.mode == benchmark{


			//process_handle := win.GetCurrentProcess()
			//SetPriorityClass(process_handle, HIGH_PRIORITY_CLASS)
			//app.data_.prof.start(data_->benchmark_name, 0)
		}

		app.on_init()
	}
}

func (app *sample_app) create_devices_and_targets(width, height, sample_count uint32,
	color_fmt, ds_format render.Pixel_format){
	if	app.data_.mode == unknown {
		return
	}

	switch app.data_.mode {
	case interactive:
	case replay:
	//if defined(EFLIB_WINDOWS)
		fmt.Println("Create GUI ...")
		app.data_.gui = device.Create_win_gui()
		break
	//#else
	//return
	//#endif
	}

	var wnd_handle win.HWND
	if app.data_.gui != nil {
		app.data_.gui.Create_window( int32(width), int32(height) )
		wnd_handle = app.data_.gui.Main_wnd_.Hwnd_
		if	wnd_handle != 0{
			fmt.Println("Error: window creation is failed.")
			app.data_.runnable = false
			return
		}
		app.data_.gui.Main_wnd_.Set_title(app.data_.Benchmark_name)
	}

	fmt.Println("Create devices and targets ..." )

	rparams := render.Renderer_parameters{width, height,
		sample_count, color_fmt, wnd_handle}

	rtype := device.Renderer_none
	var sc_type device.Swap_chain_types
	if app.data_.gui != nil {
		sc_type = device.Swap_chain_default
	}else{
		sc_type = device.Swap_chain_none
	}

	if app.data_.is_sync_renderer{
		rtype = device.Renderer_sync
	} else if !app.data_.is_sync_renderer{
		rtype = device.Renderer_async
	} else {
		switch	app.data_.mode {
			case interactive:
			case replay:
			case test:
				rtype = device.Renderer_async
				break
			case benchmark:
				rtype = device.Renderer_sync
				break
		}
	}
	device.Salviax_create_swap_chain_and_renderer(app.data_.swap_chain, app.data_.renderer, &rparams, rtype, sc_type)

	if app.data_.renderer != nil{
		return
	}
	if app.data_.swap_chain != nil{
		app.data_.color_target = (*app.data_.swap_chain).Get_surface()
	} else {
		app.data_.color_target = (*app.data_.renderer).Create_tex2d(width, height, sample_count, color_fmt).Subresource(0)
	}

	if sample_count > 1 {
		app.data_.resolved_color_target = (*app.data_.renderer).Create_tex2d(width, height, 1, color_fmt).Subresource(0)
	} else{
		app.data_.resolved_color_target = app.data_.color_target
	}

	if	ds_format != render.Pixel_format_invalid{
		app.data_.ds_target = (*app.data_.renderer).Create_tex2d(width, height, sample_count, ds_format).Subresource(0)
	}

	(*app.data_.renderer).Set_render_targets(1, app.data_.color_target, app.data_.ds_target)

	if	app.data_.gui != nil{
		app.data_.gui.Main_wnd_.Set_idle_handler( app.on_gui_idle )
		app.data_.gui.Main_wnd_.Set_draw_handler( app.on_gui_draw )
	}

	if app.data_.mode == benchmark{
		app.data_.pipeline_stat_obj = (*app.data_.renderer).Create_query(Pipeline_statistics)
		app.data_.internal_stat_obj = (*app.data_.renderer).Create_query(Internal_statistics)
		app.data_.pipeline_prof_obj = (*app.data_.renderer).Create_query(Pipeline_profiles)
	}
}

func (app *sample_app)draw_frame(){
	app.data_.elapsed_sec = app.data_.frame_timer.Elapsed()
	app.data_.total_elapsed_sec += app.data_.elapsed_sec
	app.data_.frame_timer.Restart()

	switch app.data_.quit_cond{
		case frame_limits:
			if	app.data_.frame_count >= app.data_.quit_cond_data{
				app.quit()
			}
			break
		case time_out:
			if app.data_.total_elapsed_sec > float64(app.data_.quit_cond_data / 1000.0) {
				app.quit()
			}
			break
	}

	if	app.data_.quiting{
		return
	}
	if	app.data_.mode == benchmark{
		(*app.data_.renderer).Begin(*app.data_.pipeline_stat_obj)
		(*app.data_.renderer).Begin(*app.data_.internal_stat_obj)
		(*app.data_.renderer).Begin(*app.data_.pipeline_prof_obj)
	}

	app.on_frame()

	if	app.data_.mode == benchmark {
		(*app.data_.renderer).End(*app.data_.pipeline_stat_obj)
		(*app.data_.renderer).End(*app.data_.internal_stat_obj)
		(*app.data_.renderer).End(*app.data_.pipeline_prof_obj)
	})
	if	app.data_.quiting {
		return
	}

	if	app.data_.swap_chain != nil{
		(*app.data_.swap_chain).Present()
	}

	if	app.data_.mode == benchmark {
		var frame_prof Frame_data
		(*app.data_.renderer).Get_data(*app.data_.pipeline_stat_obj, &frame_prof.pipeline_stat, false)
		(*app.data_.renderer).Get_data(*app.data_.internal_stat_obj, &frame_prof.internal_stat, false)
		(*app.data_.renderer).Get_data(*app.data_.pipeline_prof_obj, &frame_prof.pipeline_prof, false)
		app.data_.frame_profs = append(app.data_.frame_profs,frame_prof)
	}

	app.data_.frame_count +=1
	app.data_.frames_in_second +=1

	current_time := app.data_.second_timer.Elapsed()
	frame_elapsed := app.data_.frame_timer.Elapsed()
	if( (current_time >= 1.0) || (1.0 - current_time < frame_elapsed) ) {

		fmt.Printf("Frame: #%d  FPS: %f",app.data_.frame_count,float64(app.data_.frames_in_second) / current_time )

		app.data_.second_timer.Restart()
		app.data_.frames_in_second = 0
	}

	switch	app.data_.mode {
		case test:
			(*app.data_.renderer).Flush()
			if (app.data_.color_target != app.data_.resolved_color_target) {
				app.data_.color_target.Resolve(app.data_.resolved_color_target)
			}
			app.Save_frame(app.data_.resolved_color_target)
			break
	}
}

func (app *sample_app)on_gui_idle(){
	app.draw_frame()
}
func (app *sample_app)on_gui_draw(){
}

func (app *sample_app)run(){
	if (app.data_.runnable) {
		return
	}
	fmt.Println("Start running ...")
	app.data_.frame_count = 0
	if app.data_.gui != nil{
		app.data_.gui.Run()
	} else{
		for !app.data_.quiting{
			app.draw_frame()
		}
	}

	//if	app.data_.mode == benchmark{
	//	app.data_.prof.end(app.data_.benchmark_name)
	//	app.save_profiling_result()
	//	app.print_profiler(&app.data_.prof, 3)
	//}

	fmt.Println("Running done.")
}

func (app *sample_app)quit_at_frame(frame_cnt uint32){
	app.data_.quit_cond = frame_limits
	app.data_.quit_cond_data = frame_cnt
}

func (app *sample_app)quit_if_time_out(milli_sec uint32){
	app.data_.quit_cond = time_out
	app.data_.quit_cond_data = milli_sec
}

func (app *sample_app)quit(){
	fmt.Println("Exiting ...")
	app.data_.quiting = true
}

// Utilities
//func (app *sample_app)profiling(stage_name string, fn interface{}){
//	if(app.data_.mode == benchmark){
//		app.data_.prof.start(stage_name, 0)
//		fn()
//		app.data_.prof.end(stage_name)
//	} else{
//		fn()
//	}
//}

func (app *sample_app) Save_frame(surf *render.Surface){
	ss := fmt.Sprintf("%s_%d.png",app.data_.Benchmark_name,app.data_.frame_count - 1)
	//pixel_format_color_bgra8 由宏生成，这里要处理一下
	device.Save_surface((*app.data_.renderer).Get(), surf, ss, app.Pixel_format_color_bgra8)
}

//func (app *sample_app)min_max(in_out_min, in_out_max, v interface{}){
//	if(v < in_out_min){
//		in_out_min = v
//	} else if (v > in_out_max){
//		in_out_max = v
//	}
//}

//template <typename ValueT, typename IterT, typename TransformT>
//func reduce_and_output(IterT beg, IterT end, TransformT trans, ptree& parent, path string){
//	ValueT minv
//	ValueT maxv
//	ValueT total = 0
//	ValueT avg = 0
//
//	if(beg == end) {
//		minv = maxv = total = avg = 0
//	} else {
//		it = beg
//		minv = maxv = total = trans(*it)
//		it+=1
//
//		count := 1
//		for ;it != end;it+=1{
//			v := trans(*it)
//			min_max(minv, maxv, v)
//			total += v
//			count +=1
//		}
//		avg = total / count
//	}
//
//	ptree vnode
//	vnode.put("min", minv)
//	vnode.put("max", maxv)
//	vnode.put("total", total)
//	vnode.put("avg", avg)
//	parent.put_child(path, vnode)
//}

const OUTPUT_PROFILER_LEVEL = 3
func (app *sample_app)save_profiling_result(){
	//app.data_.prof.merge_items()
	//ss := fmt.Sprintf("%s_Profiling.json",app.data_.benchmark_name)
	//root := make_ptree(&app.data_.prof, OUTPUT_PROFILER_LEVEL)
	//
	//// Statistic
	//root.put("frames", data_->frame_count)

	//reduce_and_output<uint64_t>(data_->frame_profs.begin(), data_->frame_profs.end(), [](frame_data const& v) { return v.pipeline_stat.cinvocations			}, root, "async.pipeline_stat.cinvocations")
	//reduce_and_output<uint64_t>(data_->frame_profs.begin(), data_->frame_profs.end(), [](frame_data const& v) { return v.pipeline_stat.cprimitives			}, root, "async.pipeline_stat.cprimitives")
	//reduce_and_output<uint64_t>(data_->frame_profs.begin(), data_->frame_profs.end(), [](frame_data const& v) { return v.pipeline_stat.ia_primitives		}, root, "async.pipeline_stat.ia_primitives")
	//reduce_and_output<uint64_t>(data_->frame_profs.begin(), data_->frame_profs.end(), [](frame_data const& v) { return v.pipeline_stat.ia_vertices			}, root, "async.pipeline_stat.ia_vertices")
	//reduce_and_output<uint64_t>(data_->frame_profs.begin(), data_->frame_profs.end(), [](frame_data const& v) { return v.pipeline_stat.vs_invocations		}, root, "async.pipeline_stat.vs_invocations")
	//reduce_and_output<uint64_t>(data_->frame_profs.begin(), data_->frame_profs.end(), [](frame_data const& v) { return v.pipeline_stat.ps_invocations		}, root, "async.pipeline_stat.ps_invocations")
	//reduce_and_output<uint64_t>(data_->frame_profs.begin(), data_->frame_profs.end(), [](frame_data const& v) { return v.internal_stat.backend_input_pixels	}, root, "async.internal_stat.backend_input_pixels")
	//
	//reduce_and_output<uint64_t>(data_->frame_profs.begin(), data_->frame_profs.end(), [](frame_data const& v) { return v.pipeline_prof.gather_vtx			}, root, "async.pipeline_prof.gather_vtx")
	//reduce_and_output<uint64_t>(data_->frame_profs.begin(), data_->frame_profs.end(), [](frame_data const& v) { return v.pipeline_prof.vtx_proc				}, root, "async.pipeline_prof.vtx_proc")
	//reduce_and_output<uint64_t>(data_->frame_profs.begin(), data_->frame_profs.end(), [](frame_data const& v) { return v.pipeline_prof.clipping				}, root, "async.pipeline_prof.clipping")
	//reduce_and_output<uint64_t>(data_->frame_profs.begin(), data_->frame_profs.end(), [](frame_data const& v) { return v.pipeline_prof.compact_clip			}, root, "async.pipeline_prof.compact_clip")
	//reduce_and_output<uint64_t>(data_->frame_profs.begin(), data_->frame_profs.end(), [](frame_data const& v) { return v.pipeline_prof.vp_trans				}, root, "async.pipeline_prof.vp_trans")
	//reduce_and_output<uint64_t>(data_->frame_profs.begin(), data_->frame_profs.end(), [](frame_data const& v) { return v.pipeline_prof.tri_dispatch			}, root, "async.pipeline_prof.tri_dispatch")
	//reduce_and_output<uint64_t>(data_->frame_profs.begin(), data_->frame_profs.end(), [](frame_data const& v) { return v.pipeline_prof.ras					}, root, "async.pipeline_prof.ras")

	//write_json(ss.str(), root)
}

