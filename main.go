package main

import (
	"fmt"
	"flag"
	"time"
	"github.com/lifei6671/snowflake/server"
	"strings"
	"strconv"
)

func main() {
	var port int
	var workerIdBits int
	var timeDate string
	var sequenceBits int
	var timestampBits int
	var size int
	var worker string



	flag.IntVar(&port,"port",8181,"Snowflake Server listen port.")
	flag.IntVar(&workerIdBits,"workerSize",10,"WorkerBits number.")
	flag.StringVar(&worker,"worker","1","Worder id list.Please use , split.")
	flag.StringVar(&timeDate,"time","2017-04-10","Snowflake start time.")
	flag.IntVar(&sequenceBits,"sequenceSize",12,"SequenceBits sequence number.")
	flag.IntVar(&timestampBits,"timestampSize",41,"TimestampBits number.")
	flag.IntVar(&size,"size",1000,"Buffer Size.")

	flag.Parse()

	ws := strings.Split(worker,",")

	workers := make([]int64,len(ws))

	for index, ws := range ws {
		if r,err := strconv.ParseInt(ws,10,64);err == nil {
			workers[index] = r
		}
	}

	t,err := time.Parse("",timeDate)
	if err != nil {
		t = time.Now()
	}

	epochSeconds := t.Unix()

	s,err := server.NewServer(port,size,workers,timestampBits,sequenceBits,workerIdBits,epochSeconds)

	if err != nil {
		fmt.Println(err)
	}else{
		s.Run()
	}

}
