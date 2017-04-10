package ring

import (
	"errors"
	"sync/atomic"
)

var(
	ErrDisposed = errors.New(`queue: disposed`)
	ErrEmptyQueue = errors.New(`queue: empty queue`)
)

const (
	CAN_PUT_FLAG = 0
	CAN_TAKE_FLAG = 1
	START_POINT = -1
)

type RingBuffer struct {
	_padding0 [8]uint64
	bufferSize uint64
	_padding1 [8]uint64
	indexMask uint64
	_padding2 [8]uint64
	slots []int64
	_padding3 [8]uint64
	flags []int64
	_padding4 [8]uint64
	// tail 写入者最后所在位置 .
	tail int64
	_padding5 [8]uint64
	// cursor 消费者最后所在位置.
	cursor int64
	_padding6 [8]uint64
	paddingThreshold int


}

func NewRingBuffer(size uint64,padding int)  (*RingBuffer,error){

	this := &RingBuffer{
		slots : make([]int64,size),
		bufferSize: size,
		indexMask: size -1,
		paddingThreshold: int(size)* padding/100,
		flags : make([]int64,size),
		tail : START_POINT,
		cursor : START_POINT,
	}

	for i := 0 ;i < int(size); i ++ {
		this.flags[i] = CAN_PUT_FLAG
	}

	return this,nil
}

func (p *RingBuffer) Put(v int64) error {

	currentTail := atomic.LoadInt64(&(p.tail))
	currentCursor := atomic.LoadInt64(&(p.cursor))

	if currentCursor == START_POINT {
		currentCursor = 0
	}


	distance := currentTail - currentCursor

	if distance == int64(p.bufferSize - 1 ){
		return errors.New("cache full.")
	}

	nextTailIndex := uint64(currentTail + 1) % p.bufferSize

	p.slots[nextTailIndex] = v
	p.flags[nextTailIndex] = CAN_TAKE_FLAG
	atomic.AddInt64(&(p.tail),1)

	return  nil
}

func (p *RingBuffer) Take() (int64,error)  {
	currentTail := atomic.LoadInt64(&(p.tail));
	currentCursor := atomic.LoadInt64(&(p.cursor))

	if currentCursor + 1 >= currentTail {
		return 0,errors.New("Curosr can't move back")
	}

	nextCursorIndex := uint64(currentCursor + 1) % p.bufferSize

	v := p.slots[nextCursorIndex]
	atomic.AddInt64(&(p.cursor),1)
	p.flags[nextCursorIndex] = CAN_PUT_FLAG
	atomic.AddInt64(&(p.tail),1)
	return v,nil
}