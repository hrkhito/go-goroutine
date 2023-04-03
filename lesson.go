package main

import (
	"fmt"
	"sync"
	"time"
)

// goroutineとsync.WaitGroup (goroutineの処理が終わってなくてもgoの処理は終了する)

func goroutine(s string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		fmt.Println(s)
	}
}

func normal(s string) {
	for i := 0; i < 5; i++ {
		fmt.Println(s)
	}
}

// channel

func goroutine1(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}

// channelのrangeとclose

func goroutine2(ss []int, cc chan int) {
	sum := 0
	for _, vv := range ss {
		sum += vv
		cc <- sum
	}
	close(cc)
}

// producerとconsumer

// producer
func producer(ch1 chan int, i int) {
	// Something
	ch1 <- i * 2
}

// consumer
func consumer(ch1 chan int, wg1 *sync.WaitGroup) {
	for i := range ch1 {
		func() {
			defer wg1.Done()
			fmt.Println("process", i*1000)
		}()
	}
	fmt.Println("#################")
}

// fan-out fan-in

func fanproducer(first chan int) {
	defer close(first)
	for i := 0; i < 10; i++ {
		first <- i
	}
}

func multi2(first <-chan int, second chan<- int) {
	defer close(second)
	for i := range first {
		second <- i * 2
	}
}

func multi4(second chan int, third chan int) {
	defer close(third)
	for i := range second {
		third <- i * 4
	}
}

// channelとselect

func selectgoroutine1(ch chan string) {
	for {
		ch <- "packet from 1"
		time.Sleep(3 * time.Second)
	}
}

func selectgoroutine2(ch chan int) {
	for {
		ch <- 100
		time.Sleep(1 * time.Second)
	}
}

// sync.Mutex

type Counter struct {
	v   map[string]int
	mux sync.Mutex
}

func (c *Counter) Inc(key string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.v[key]++
}

func (c *Counter) Value(key string) int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.v[key]
}

func main() {
	// goroutineとsync.WaitGroup (goroutineの処理が終わってなくてもgoの処理は終了する)

	var wg sync.WaitGroup
	wg.Add(1)
	go goroutine("world", &wg)
	normal("hello")
	wg.Wait()

	// channel

	s := []int{1, 2, 3, 4, 5}
	c := make(chan int) // 15, 15
	go goroutine1(s, c)
	go goroutine1(s, c)
	x := <-c
	fmt.Println(x)
	y := <-c
	fmt.Println(y)

	// Buffered Channels

	ch := make(chan int, 2)
	ch <- 100
	fmt.Println(len(ch))
	ch <- 200
	fmt.Println(len(ch))

	close(ch)
	for c := range ch {
		fmt.Println(c)
	}

	// channelのrangeとclose
	ss := []int{1, 2, 3, 4, 5}
	cc := make(chan int, len(ss))
	go goroutine2(ss, cc)
	for i := range cc {
		fmt.Println(i)
	}

	// producerとconsumer

	var wg1 sync.WaitGroup
	ch1 := make(chan int)

	// Producer
	for i := 0; i < 10; i++ {
		wg1.Add(1)
		go producer(ch1, i)
	}

	// Consumer
	go consumer(ch1, &wg1)
	wg1.Wait()
	close(ch1)
	time.Sleep(2 * time.Second)
	fmt.Println("Done")

	// fan-out fan-in

	first := make(chan int)
	second := make(chan int)
	third := make(chan int)

	go fanproducer(first)
	go multi2(first, second)
	go multi4(second, third)
	for result := range third {
		fmt.Println(result)
	}

	// channelとselect

	// selectc1 := make(chan string)
	// selectc2 := make(chan int)
	// go selectgoroutine1(selectc1)
	// go selectgoroutine2(selectc2)

	// for {
	// 	select {
	// 	case msg1 := <-selectc1:
	// 		fmt.Println(msg1)
	// 	case msg2 := <-selectc2:
	// 		fmt.Println(msg2)
	// 	}
	// }

	// Default Selectionとfor break

	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
OuterLoop:
	for {
		select {
		case t := <-tick:
			fmt.Println("tick.", t)
		case <-boom:
			fmt.Println("BOOM!")
			break OuterLoop
		default:
			fmt.Println("   .")
			time.Sleep(50 * time.Millisecond)
		}
	}
	fmt.Println("##########")

	// sync.Mutex

	ccc := Counter{v: make(map[string]int)}
	go func() {
		for i := 0; i < 10; i++ {
			ccc.Inc("Key")
		}
	}()
	go func() {
		for i := 0; i < 10; i++ {
			ccc.Inc("Key")
		}
	}()
	time.Sleep(1 * time.Second)
	fmt.Println(&ccc, ccc.Value("Key"))

}
