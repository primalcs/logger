package ring_buffer

import (
	"container/ring"
	"time"

	"github.com/primalcs/logger/types"
)

// MessageBuffer wraps ring.Ring for creating ring buffer
type MessageBuffer struct {
	ringBuf *ring.Ring
}

// NewMessageBuffer creates new instance of MessageBuffer with n buffer length
func NewMessageBuffer(n int) *MessageBuffer {
	return &MessageBuffer{ringBuf: ring.New(n)}
}

// RingCell is a type of value stored in MessageBuffer
type RingCell struct {
	LogParams types.LogParams
	Message   string
	timestamp int64
}

// GetOldestCell returns newest (by timestamp) ring value, pointer to current ring item
// and true if value was found; returns false otherwise
func (mb *MessageBuffer) GetOldestCell() (*RingCell, *ring.Ring, bool) {
	return mb.getCell(func(a, b *RingCell) bool {
		return a.timestamp <= b.timestamp
	}, mb.ringBuf.Prev)
}

// GetNewestCell returns newest (by timestamp) ring value, pointer to current ring item
// and true if value was found; returns false otherwise
func (mb *MessageBuffer) GetNewestCell() (*RingCell, *ring.Ring, bool) {
	return mb.getCell(func(a, b *RingCell) bool {
		return a.timestamp >= b.timestamp
	}, mb.ringBuf.Next)
}

func (mb *MessageBuffer) getCell(compare func(a, b *RingCell) bool, dir func() *ring.Ring) (*RingCell, *ring.Ring, bool) {
	head := mb.ringBuf
	for i := 0; i < mb.ringBuf.Len(); i++ {
		if mb.ringBuf.Value != nil {
			break
		}
		mb.ringBuf = dir()
	}
	if mb.ringBuf.Value == nil {
		return nil, head, false
	}
	val, ok := mb.ringBuf.Value.(*RingCell)
	if !ok {
		return nil, head, false
	}
	current := mb.ringBuf
	res := val
	for i := 0; i < mb.ringBuf.Len(); i++ {
		mb.ringBuf = dir()
		if mb.ringBuf.Value == nil {
			continue
		}
		val, ok := mb.ringBuf.Value.(*RingCell)
		if !ok {
			continue
		}
		if compare(val, res) {
			res = val
			current = mb.ringBuf
		} else {
			break
		}
	}

	return res, current, true
}

// EraseCell deletes value from pointed ring item
func (mb *MessageBuffer) EraseCell(r *ring.Ring) {
	mb.ringBuf = r
	r.Value = nil
	mb.ringBuf = mb.ringBuf.Prev()
}

// AddCell writes value to the ring item after Newest
// it may rewrite existing ring item value
func (mb *MessageBuffer) AddCell(r *RingCell) {
	_, cur, _ := mb.GetNewestCell()
	cur = cur.Next()
	cur.Value = r
}

// NewCell creates new value for ring buffer item
func NewCell(lp types.LogParams, m string) *RingCell {
	return &RingCell{
		LogParams: lp,
		Message:   m,
		timestamp: time.Now().UTC().UnixNano(),
	}
}
