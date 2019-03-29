package pipeline

import "bytes"

type Stream_buffer_desc struct{

	slot,stride,offset uint32
	buf *bytes.Buffer
};

type Stream_state struct {

	buffer_descs [MAX_INPUT_SLOTS]Stream_buffer_desc
}

func (ss Stream_state)update(starts_slot,buffers_count uint32,  bufs []*bytes.Buffer, strides ,offsets  []uint32){
	if( bufs == nil || strides == nil || offsets == nil ) {
		return;
	}

	if( starts_slot + buffers_count >= MAX_INPUT_SLOTS ) {
		return;
	}
	var i uint32
	for i = 0; i < buffers_count; i++{
		ss.buffer_descs[starts_slot+i].buf		= bufs[i];
		ss.buffer_descs[starts_slot+i].offset	= offsets[i];
		ss.buffer_descs[starts_slot+i].stride	= strides[i];
		ss.buffer_descs[starts_slot+i].slot	= starts_slot+i;
	}
}