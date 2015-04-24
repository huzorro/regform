package sexredis

import (
	"errors"
	"github.com/gosexy/redis"
)

//  Example use
//	q := sexredis.New()
//	q.Create("myqueue", "localhost", uint(6379))
//	for i := 0; i < 5; i++ {
//		q.Put(mjson)
//	}

//	q.Worker(uint(5), func(msg string, err error) {
//		if err != nil {
//			log.Printf(err.Error())
//		} else {
//			log.Printf("... %s\n", msg)
//		}

//	})
type Queue struct {
	name    string
	rc      *redis.Client
	msgchan chan Msg
}

type Msg struct {
	Content interface{}
	Err     error
}

type Processor interface {
	//串行处理
	SProcess(msg *Msg)
}

func New() *Queue {
	q := new(Queue)
	return q
}

func (self *Queue) SetRClient(name string, rc *redis.Client) {
	self.name = name
	self.rc = rc
}

func (self *Queue) Create(name string, host string, port uint) error {
	self.name = name
	self.rc = redis.New()
	return self.rc.Connect(host, port)
}

func (self *Queue) Active() (response string, err error) {
	return self.rc.Ping()
}

func (self *Queue) LLen() (l int64, err error) {
	return self.rc.LLen(self.name)
}

func (self *Queue) Keys() (name string) {
	return self.name
}

func (self *Queue) Clear() (recode int64, err error) {
	return self.rc.Del(self.name)
}

func (self *Queue) Close() (string, error) {
	return self.rc.Quit()
}

func (self *Queue) Get(block bool, timeout uint64) (msg Msg) {
	if block {
		if rs, err := self.rc.BLPop(timeout, self.name); err != nil {
			return Msg{"", err}
		} else {
			if len(rs) > 1 {
				return Msg{rs[1], nil}
			} else {
				return Msg{"", errors.New("queue is empty")}
			}
		}
	} else {
		m, err := self.rc.LPop(self.name)
		return Msg{m, err}
	}
}

func (self *Queue) Put(msg interface{}) (recode int64, err error) {
	return self.rc.RPush(self.name, msg)
}

/*
use channel implement like python yield
*/
func (self *Queue) consume() {
	self.msgchan = make(chan Msg)
	for {
		//		if msg, err := self.Get(true, uint64(0)); err != nil {
		//			return err
		//		} else {
		//			self.msgchan <- msg
		//		}

		self.msgchan <- self.Get(true, uint64(0))
	}
}

func (self *Queue) yield() (msg Msg) {
	return <-self.msgchan
}

//func (self *Queue) Worker(pnum uint, fn func(msg string, err error)) {
//	var err error
//	control := make(chan string, pnum)
//	go func() {
//		self.consume()
//	}()

//	go func() {
//		for {
//			if err != nil {
//				fn("", err)
//				return
//			}
//			msg := self.yield()
//			control <- msg
//			go func() {
//				fn(msg, nil)
//				<-control
//			}()
//		}
//	}()
//}

func (self *Queue) Worker(pnum uint, serial bool, ps ...Processor) {
	control := make(chan Msg, pnum)
	go func() {
		self.consume()
	}()

	go func() {
		for {
			msg := self.yield()
			control <- msg
			go func() {
				if serial {
					for _, sp := range ps {
						sp.SProcess(&msg)
						if msg.Err != nil {
							break
						}
					}
				}
				<-control
			}()
		}
	}()
}
