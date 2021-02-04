package ring_buffer_test

import (
	"fmt"
	"testing"

	"github.com/primalcs/logger/types"

	"github.com/primalcs/logger/ring_buffer"
	"github.com/stretchr/testify/assert"
)

var bufLen = 7

func TestMessageBuffer_AddCell(t *testing.T) {
	buf := ring_buffer.NewMessageBuffer(bufLen)
	_, head, ok := buf.GetCurrent()

	assert.Equal(t, false, ok)
	for i := 0; i < bufLen*2; i++ {
		buf.AddCell(ring_buffer.NewCell(types.LogParams{}, ""))
	}
	_, tail, ok := buf.GetCurrent()
	assert.Equal(t, head, tail)
}

func TestMessageBuffer_EraseCell(t *testing.T) {
	buf := ring_buffer.NewMessageBuffer(bufLen)
	_, head, ok := buf.GetCurrent()

	assert.Equal(t, false, ok)
	for i := 0; i < bufLen*2; i++ {
		buf.AddCell(ring_buffer.NewCell(types.LogParams{}, ""))
	}
	for i := 0; i < bufLen; i++ {
		buf.EraseCell()
	}
	_, tail, ok := buf.GetCurrent()
	assert.Equal(t, head, tail)

	for i := 0; i < bufLen; i++ {
		buf.EraseCell()
	}
	_, tail, ok = buf.GetCurrent()
	assert.Equal(t, head, tail)
}

func TestMessageBuffer_GetNewestCell(t *testing.T) {
	buf := ring_buffer.NewMessageBuffer(bufLen)
	_, head, ok := buf.GetNewestCell()
	assert.Equal(t, false, ok)

	for i := 0; i < bufLen*2; i++ {
		msg := fmt.Sprintf("msg number %d", i)
		buf.AddCell(ring_buffer.NewCell(types.LogParams{}, msg))
		buf.SetCurrent(head)
		cell, _, ok := buf.GetNewestCell()
		assert.Equal(t, true, ok)
		assert.Equal(t, msg, cell.Message)
	}
}

func TestMessageBuffer_GetOldestCell(t *testing.T) {
	buf := ring_buffer.NewMessageBuffer(bufLen)
	_, _, ok := buf.GetOldestCell()
	assert.Equal(t, false, ok)

	for i := 0; i < bufLen*2; i++ {
		msg := fmt.Sprintf("msg number %d", i)
		buf.AddCell(ring_buffer.NewCell(types.LogParams{}, msg))
	}
	cell, _, ok := buf.GetOldestCell()
	assert.Equal(t, true, ok)
	assert.Equal(t, fmt.Sprintf("msg number %d", bufLen), cell.Message)

}
