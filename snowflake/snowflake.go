// Package snowflake .
package snowflake

import (
	"time"
	"errors"
	"strconv"
)



//  Snowflake struct .
type Snowflake struct {
	allocator *BitsAllocator
	epochSeconds int64
	lastTimestamp int64
	sequence int64
	workerId int64
}

func NewSnowflake(workerId int64, allocator *BitsAllocator) (*Snowflake,error) {

	if workerId > allocator.MaxWorkerId() {
		return nil,errors.New("Worker id " +  strconv.FormatInt(workerId,10) + " exceeds the max " + strconv.FormatInt(allocator.MaxWorkerId(),10));
	}
	p := &Snowflake{
		allocator : allocator ,
		epochSeconds : 1491794590000,
		lastTimestamp: -1,
		sequence : 0,
		workerId : workerId,
	}


	return p,nil
}

func (p *Snowflake) getCurrentSecond() (int64,error) {
	currentSecond := time.Now().UnixNano()  / 1000000

	if currentSecond < p.epochSeconds {
		return 0,errors.New("epochSeconds error.")
	}
	return currentSecond,nil
}

func (p *Snowflake) SetEpochSeconds(epochSeconds int64)  {
	p.epochSeconds = epochSeconds
}

func (p *Snowflake) NextId() (int64,error) {
	currentSecond,err := p.getCurrentSecond()
	if err != nil {
		return 0,err
	}

	if currentSecond < p.lastTimestamp {
		refusedSeconds := strconv.FormatInt(p.lastTimestamp - currentSecond,10);

		return 0,errors.New("Clock moved backwards. Refusing for "+ refusedSeconds + " seconds")
	}

	if currentSecond == p.lastTimestamp {
		p.sequence = (p.sequence + 1) & p.allocator.SequenceMask();
		if p.sequence == 0 {
			timestamp,err := p.getNextSecond();
			if  err != nil {
				return 0,err
			}
			currentSecond = timestamp
		}
	}else {
		p.sequence = 0
	}
	p.lastTimestamp = currentSecond

	return p.allocator.Allocate(currentSecond - p.epochSeconds, p.workerId, p.sequence),nil;
}

func (p *Snowflake) getNextSecond() (int64,error) {
	timestamp,err := p.getCurrentSecond();
	if err != nil {
		return 0,err
	}

	for timestamp <= p.lastTimestamp {
		timestamp,err = p.getCurrentSecond();
		if err != nil {
			return 0,err
		}
	}

	return timestamp,nil;
}

func (p *Snowflake) WorkerId() int64  {
	return p.workerId
}