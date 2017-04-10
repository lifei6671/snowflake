package snowflake

import (
	"errors"
)

const TOTAL_BITS = 1 << 6

type BitsAllocator struct {
	// signBits 首位标识位 .
	signBits int
	// timestampBits 时间戳的位数 .
	timestampBits int
	//workerIdBits 工作ID的位数 .
	workerIdBits int
	// sequenceBits 顺序位数.
	sequenceBits int

	// sequenceMask 生成序列的掩码 .
	sequenceMask int64

	maxDeltaSeconds int64
	maxWorkerId int64
	maxSequence int64

	timestampShift uint64
	workerIdShift uint64
}


func NewBitsAllocator(timestampBits ,workerIdBits ,sequenceBits int) (*BitsAllocator,error) {

	if (1 + timestampBits + workerIdBits + sequenceBits) > TOTAL_BITS {
		return nil,errors.New("Allocate not enough 64 bits")
	}
	p := &BitsAllocator{
		signBits: 1,
		timestampBits: timestampBits,
		workerIdBits: workerIdBits,
		sequenceBits: sequenceBits,
		sequenceMask: ^(-1 << uint64(sequenceBits)),
		maxDeltaSeconds: ^(-1 << uint64(timestampBits)),
		maxSequence: ^(-1 << uint64(sequenceBits)),
		maxWorkerId:^(-1 << uint64(workerIdBits)),
		timestampShift: uint64(workerIdBits + sequenceBits),
		workerIdShift: uint64(sequenceBits),
	}

	return p,nil
}

func (p *BitsAllocator) Allocate(deltaSeconds,workerId,sequence int64) (int64) {

	return (deltaSeconds << p.timestampShift) | (workerId << p.workerIdShift) | sequence;
}

func (p *BitsAllocator) TimestampBits() int  {
	return  p.timestampBits
}

func (p *BitsAllocator) WorkerIdBits() int  {
	return  p.workerIdBits
}
func (p *BitsAllocator) SequenceBits() int  {
	return  p.sequenceBits
}

func (p *BitsAllocator) MaxDeltaSeconds() int64  {
	return  p.maxDeltaSeconds
}
func (p *BitsAllocator) MaxSequence() int64  {
	return  p.maxSequence
}
func (p *BitsAllocator) MaxWorkerId() int64  {
	return  p.maxWorkerId
}

func (p *BitsAllocator) TimestampShift() uint64  {
	return  p.timestampShift
}
func (p *BitsAllocator) WorkerIdShift() uint64  {
	return  p.workerIdShift
}
func (p *BitsAllocator) SequenceMask() int64 {
	return p.sequenceMask
}