package logger

import (
	"container/ring"
	"time"
)

type MessageBuffer struct {
	ringBuf *ring.Ring
}

func NewMessageBuffer(n int) *MessageBuffer {
	return &MessageBuffer{ringBuf: ring.New(n)}
}

type RingCell struct {
	LogParams LogParams
	Message   string
	timestamp int64
}

func (mb *MessageBuffer) GetOldestCell() (*RingCell, *ring.Ring, bool) {
	return mb.GetCell(func(a, b *RingCell) bool {
		return a.timestamp <= b.timestamp
	}, mb.ringBuf.Prev)
}
func (mb *MessageBuffer) GetNewestCell() (*RingCell, *ring.Ring, bool) {
	return mb.GetCell(func(a, b *RingCell) bool {
		return a.timestamp >= b.timestamp
	}, mb.ringBuf.Next)
}

func (mb *MessageBuffer) GetCell(compare func(a, b *RingCell) bool, dir func() *ring.Ring) (*RingCell, *ring.Ring, bool) {
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

func (mb *MessageBuffer) EraseCell(r *ring.Ring) {
	mb.ringBuf = r
	r.Value = nil
	mb.ringBuf = mb.ringBuf.Prev()
}

func (mb *MessageBuffer) AddCell(r *RingCell) {
	_, cur, _ := mb.GetNewestCell()
	cur = cur.Next()
	cur.Value = r
}

func NewCell(lp LogParams, m string) *RingCell {
	return &RingCell{
		LogParams: lp,
		Message:   m,
		timestamp: time.Now().UTC().UnixNano(),
	}
}
