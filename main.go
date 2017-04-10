package main

import (
	"fmt"
	"github.com/lifei6671/snowflake/snowflake"
	"flag"
	"time"
	"github.com/lifei6671/snowflake/ring"
)

func main() {
	var port int
	var workerId int64
	var timeDate string
	var sequence int64

	flag.IntVar(&port,"port",8181,"Snowflake Server listen port.")
	flag.Int64Var(&workerId,"worker",23,"Snowflake Server ID.")
	flag.StringVar(&timeDate,"time","2017-04-10","Snowflake start time.")
	flag.Int64Var(&sequence,"sequence",0,"Snowflake sequence number.")
	flag.Parse()

	t,err := time.Parse("",timeDate)
	if err != nil {
		t = time.Now()
	}

	epochSeconds := t.Unix()

	allocator ,_ := snowflake.NewBitsAllocator(31,23,9)

	snow,_ := snowflake.NewSnowflake(10,allocator)

	snow.SetEpochSeconds(epochSeconds)

	fmt.Println(snow.NextId())

	r,_ := ring.NewRingBuffer(100,1)

	go func() {

	}()

	for i:= 10; i<= 1000;i ++ {
		r.Put(int64(i))
	}

	for i:= 10; i<= 1000;i ++ {
		v,err := r.Take();
		if err != nil{
			fmt.Println(err)
		}
		fmt.Println(v)
	}

	fmt.Printf("%+v %v",r,  10 & 50)
}
