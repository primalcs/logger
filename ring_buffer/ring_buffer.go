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
// and true if value was found; returns false otherwise; mostly for debug purposes
func (mb *MessageBuffer) GetOldestCell() (*RingCell, *ring.Ring, bool) {
	return mb.getCell(func(a, b *RingCell) bool {
		return a.timestamp < b.timestamp
	}, false)
}

// GetNewestCell returns newest (by timestamp) ring value, pointer to current ring item
// and true if (non-nil) value was found; returns false otherwise; mostly for debug purposes
func (mb *MessageBuffer) GetNewestCell() (*RingCell, *ring.Ring, bool) {
	return mb.getCell(func(a, b *RingCell) bool {
		return a.timestamp > b.timestamp
	}, true)
}

func (mb *MessageBuffer) getCell(compare func(a, b *RingCell) bool, dirForward bool) (*RingCell, *ring.Ring, bool) {
	head := mb.ringBuf
	defer func() {
		mb.ringBuf = head
	}()
	for i := 0; i < mb.ringBuf.Len(); i++ {
		if mb.ringBuf.Value != nil {
			break
		}
		if dirForward {
			mb.ringBuf = mb.ringBuf.Next()
		} else {
			mb.ringBuf = mb.ringBuf.Prev()
		}
	}
	if mb.ringBuf.Value == nil {
		return nil, head, false
	}
	current := mb.ringBuf
	val, ok := current.Value.(*RingCell)
	if !ok {
		return nil, head, false
	}
	res := val
	for i := 0; i < mb.ringBuf.Len(); i++ {
		if dirForward {
			mb.ringBuf = mb.ringBuf.Next()
		} else {
			mb.ringBuf = mb.ringBuf.Prev()
		}
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

// EraseCell deletes value from current ring item and sets Current to previous
func (mb *MessageBuffer) EraseCell() {
	mb.ringBuf.Value = nil
	mb.ringBuf = mb.ringBuf.Prev()
}

// AddCell writes value to the ring item after Current
// it may rewrite existing ring item value
func (mb *MessageBuffer) AddCell(r *RingCell) {
	mb.ringBuf = mb.ringBuf.Next()
	mb.ringBuf.Value = r
}

// GetCurrent gets current ring cell Value, pointer to current cell and bool that indicates if Value was not nil
func (mb *MessageBuffer) GetCurrent() (*RingCell, *ring.Ring, bool) {
	if mb.ringBuf.Value == nil {
		return nil, mb.ringBuf, false
	}
	val, ok := mb.ringBuf.Value.(*RingCell)
	if !ok {
		return nil, mb.ringBuf, false
	}
	return val, mb.ringBuf, true
}

//  SetCurrent sets ring current position
func (mb *MessageBuffer) SetCurrent(r *ring.Ring) {
	mb.ringBuf = r
}

// NewCell creates new value for ring buffer item
func NewCell(lp types.LogParams, m string) *RingCell {
	return &RingCell{
		LogParams: lp,
		Message:   m,
		timestamp: time.Now().UTC().UnixNano(),
	}
}
