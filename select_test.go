package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSelectDefaultOnWaitingUnbufferedChannels(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		ch1 <- 1
		ch2 <- 2
		time.Sleep(time.Second * 2)
	}()

	for done := false; !done; {
		select {
		case ch1Val := <-ch1:
			fmt.Println(ch1Val)
		case ch2Val := <-ch2:
			fmt.Println(ch2Val)
		default:
			done = true
		}
	}
}

func TestSelectDefaultOnWaitingUnbufferedChannelsWithSyncWall(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	
	send := func(ch chan<- int, val int, st time.Duration, wg *sync.WaitGroup) {
		ch <- val
		time.Sleep(st)
		wg.Done()
	}

	go send(ch1, 1, time.Second*2, wg)
	go send(ch2, 2, time.Second*2, wg)

	for done := false; !done; {
		select {
		case ch1Val := <-ch1:
			fmt.Println(ch1Val)
		case ch2Val := <-ch2:
			fmt.Println(ch2Val)
		default:
			done = true
		}
	}
	wg.Wait()
}
