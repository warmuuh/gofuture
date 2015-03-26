package main

import (
	. "../../common"
	"errors"
	"log"
	"time"
)

func main() {
	var s Service = Defer(func(r Req) Res {
		log.Println("service1: ", r)
		time.Sleep(2 * time.Second)
		log.Println("service1 finished")
		return "res from s1"
	})

	var s2 Service = Defer(func(r Req) Res {
		log.Println("service2: ", r)
		time.Sleep(2 * time.Second)
		log.Println("service2 finished")
		return "res from s2"
	})

	var identity1 Filter = func(r Req, s Service) Future {
		log.Println("identity Filter 1")
		res := s(r)
		log.Println("identity Filter 1 finished")
		return res
	}

	var identity2 Filter = func(r Req, s Service) Future {
		log.Println("identity Filter 2")
		res := s(r)
		log.Println("identity Filter 2 finished")
		return res
	}

	var timeoutFilter Filter = func(r Req, s Service) Future {
		log.Println("timeout filter")
		notify := make(chan bool)
		timeout := time.NewTicker(1 * time.Second).C

		var future Future
		go func() {
			future = s(r)
			future.Get()
			notify <- true
		}()

		select {
		case <-notify:
			log.Println("no timeout!!")
			return future
		case <-timeout:
			log.Println("timeout!!")
			return FailedFuture(errors.New("timeout"))
		}
	}

	middleware := identity1.AndThen(identity2)

	log.Println("executing service")
	f := middleware.And(s)("req")

	log.Println("adding then")
	f = f.Success(middleware.AndThen(timeoutFilter).And(s2))
	log.Println("waiting for it")

	res, err := f.Get()

	log.Println(res, err)
	log.Println(f.Get())
}
