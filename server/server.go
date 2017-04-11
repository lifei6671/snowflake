package server

import (
	"fmt"
	"net/http"
	"github.com/lifei6671/snowflake/snowflake"
	"sync/atomic"
	"time"
)


// SnowflakeServer struct .
type SnowflakeServer struct {
	port int
	snowflakes []*tempSnowflake
	channels chan int64
}

type tempSnowflake struct {
	snow *snowflake.Snowflake
	state int32
}

// NewServer 创建一个服务器.
func NewServer(port int, bufferSize int,workerId []int64,timestampBits ,sequenceBits,workerIdBits int,epochSeconds int64) (*SnowflakeServer,error) {

	server := &SnowflakeServer{
		port: port,
		channels : make(chan int64,bufferSize),
	}
	fmt.Println("BufferSize",bufferSize)
	snowflakes := make([]*tempSnowflake,len(workerId))

	for index,worker := range workerId {
		allocator ,err := snowflake.NewBitsAllocator(timestampBits,workerIdBits,sequenceBits)

		if err != nil{
			return server,err
		}
		snow,err :=  snowflake.NewSnowflake(worker,allocator)

		if err != nil {
			return server,err
		}
		snow.SetEpochSeconds(epochSeconds)

		snowflakes[index] = &tempSnowflake{
			snow	: snow,
			state	: 0,
		}

	}
	server.snowflakes = snowflakes

	return server,nil
}

// Run 启动服务器监听.
func (p *SnowflakeServer) Run() error {
	fmt.Printf("listening:%d \r\n",p.port)
	addr := fmt.Sprintf(":%d", p.port)

	go p.produce()

	return http.ListenAndServe(addr, p)
}

func (p *SnowflakeServer) produce()  {

	for _,worker := range p.snowflakes {

		if atomic.LoadInt32(&(worker.state)) == 0 {

			atomic.StoreInt32(&(worker.state),1)
			defer atomic.StoreInt32(&(worker.state) ,0)

			cumulative := 0
			ticker := time.NewTicker(1 * time.Second)

			defer ticker.Stop()

			for{
				if id,err := worker.snow.NextId(); err == nil {

					select {
					case p.channels <- id:
						break
					default:
					}
					select {
					case <- ticker.C:
						fmt.Println("Timeout => WorkerId：",worker.snow.WorkerId())
						return
					default:
					}
				}else{
					cumulative ++
				}
				if cumulative > 1000 {
					return
				}
			}
		}
	}
}

func (p *SnowflakeServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for {
		id := <-p.channels

		if len(p.channels) <= cap(p.channels) / 2 {
			go p.produce()
		}

		fmt.Fprint(w, id)
		return
	}
}

