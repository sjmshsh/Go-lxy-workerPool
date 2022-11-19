package main

import (
	"github.com/sjmshsh/model"
	"time"
)

func main() {
	p := model.New(5)
	for i := 0; i < 10; i++ {
		err := p.Schedule(func() {
			time.Sleep(time.Second * 3)
		})
		if err != nil {
			println("task: ", i, "err:", err)
		}
	}
	p.Free()
}
