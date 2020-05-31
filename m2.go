package main

import (
	"fmt"
	"github.com/thinkgos/gomodbus/v2"
	"sync"
	"time"
)

var wg2 sync.WaitGroup

var failCounter = 0

func getData2(num int) {

	fmt.Printf("任务：%d，开始.......\n", num)
	p := modbus.NewTCPClientProvider("127.0.0.1:503",
		modbus.WithEnableLogger())
	client := modbus.NewClient(p)
	err := client.Connect()
	if err != nil {
		fmt.Println("connect failed, ", err)
		wg2.Done()
		failCounter++
		return
	}

	defer client.Close()
	fmt.Println("starting")

	results, err := client.ReadHoldingRegisters(1, 0, 10)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("ReadCoils % x", results)
	}

	fmt.Printf("任务：%d，结束.......\n\n", num)
	wg2.Done()
}

func main() {
	t1 := time.Now() // get current time

	wg2.Add(10000)
	for i := 1; i <= 10000; i++ {
		go getData2(i)
	}

	wg2.Wait()

	elapsed := time.Since(t1)
	fmt.Println("Fail Counter: ", failCounter)
	fmt.Println("App elapsed: ", elapsed)
}
