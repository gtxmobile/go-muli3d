package pipeline

import (
	"debug/macho"
	"sync"
	"sync/atomic"
)

type Async_status uint32

const(
	error Async_status = iota
	timeout
	ready
)
const ASYNC_RENDER_QUEUE_SIZE uint32 = 32
type  Async_renderer struct {
	Renderer_impl
	State_queue_ [ASYNC_RENDER_QUEUE_SIZE]Render_state
	State_pool_ [ASYNC_RENDER_QUEUE_SIZE]Render_state
	Waiting_exit_ bool
	state_pool_mutex_ sync.Mutex
	state_queue_ []*Render_state
	state_pool_ []*Render_state

	rendering_thread_  macho.Thread
	waiting_exit_ atomic.Value
}

func (ar *Async_renderer)init_async_render(){
	//ar.State_queue_  = make([]Render_state,ASYNC_RENDER_QUEUE_SIZE)
	//ar.State_pool_	=  make([]Render_state,ASYNC_RENDER_QUEUE_SIZE)
	ar.Waiting_exit_ = false
	for _,state := range ar.State_pool_{
		state.Reset(Render_state{})
	}
}

func (ar *Async_renderer)flush()fundations.Result{
	for object_count_in_pool() != MAX_COMMAND_QUEUE {
		//yield()
	}
	return common.Ok
}

func (ar *Async_renderer) run(){
	//rendering_thread_ = boost::thread(&async_renderer::do_rendering, this)
}

func (ar *Async_renderer) alloc_render_state() render_state_ptr{
	for{
		pool_lock := sync.Mutex{}

		if	len(ar.State_pool_)!= 0 {
			ret := ar.State_pool_.back()
			ar.State_pool_.pop_back()
			return ret
		}
		//boost::thread::yield();
	}
}

func (ar *Async_renderer)free_render_state(state *Render_state){
	pool_lock := sync.Mutex{}
	//state_pool_mutex_);
	ar.State_pool_.push_back(state);
}

func (ar *Async_renderer)object_count_in_pool() uint32{
	pool_lock := sync.Mutex{}
	//state_pool_mutex_);
	return ar.State_pool_.size()
}

func (ar *Async_renderer)commit_state_and_command() fundations.Result{
	dest_state := alloc_render_state()
	copy_using_state(dest_state.get(), state_.get())
	ar.State_queue_.push_front(dest_state)

	return common.Ok
}

func (ar *Async_renderer)release() fundations.Result{
	if  ar.Rendering_thread_.joinable(){
		ar.Waiting_exit_ = true
		ar.State_queue_.push_front(render_state_ptr())
		ar.Rendering_thread_.join()
	}
	return common.Ok
}

func (ar *Async_renderer)do_rendering() {
	for ar.Waiting_exit_ {
		var rendering_state *Render_state
		ar.State_queue_.pop_back(&rendering_state)

		if ar.rendering_state {
			core_.update(rendering_state)
			core_.execute()
			free_render_state(rendering_state)
		}
	}
}


func Create_async_renderer() *Async_renderer{
	ret := Async_renderer{}
	ret.run()
	return &ret
}