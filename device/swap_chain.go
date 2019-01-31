package device

type Swap_chain_types int
const(
	Swap_chain_none Swap_chain_types= iota
	Swap_chain_default
	Swap_chain_d3d11
	Swap_chain_gl
)
type Renderer_types int
const(
	Renderer_none	Renderer_types = iota
	Renderer_async
	Renderer_sync
)