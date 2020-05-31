package main

import (
	"fmt"
	"github.com/goburrow/modbus"
	"log"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup
var failCounter = 0

func main() {
	t1 := time.Now() // get current time

	wg.Add(10000)
	for i := 1; i <= 10000; i++ {
		go getData(i)
	}

	wg.Wait()

	elapsed := time.Since(t1)
	fmt.Println("App elapsed: ", elapsed)
}

func getData(num int) {
	fmt.Printf("任务：%d，开始.......\n", num)

	handler := modbus.NewTCPClientHandler("localhost:503")
	handler.Timeout = 5 * time.Second
	handler.SlaveId = 1
	handler.Logger = log.New(os.Stdout, "Modbus TCP: ", log.LstdFlags)

	//handler := modbus.NewRTUClientHandler("COM2")
	//handler.BaudRate = 9600
	//handler.DataBits = 8
	//handler.Parity = "N"
	//handler.StopBits = 1
	//handler.SlaveId = 1
	//handler.Timeout = 5 * time.Second

	// Connect manually so that multiple requests are handled in one session
	err := handler.Connect()
	defer handler.Close()

	client := modbus.NewClient(handler)

	//_, err = client.WriteMultipleRegisters(0, 4, []byte{0, 10, 0, 255, 1, 5, 0, 3})
	if err != nil {
		fmt.Printf("%v\n", err)
		wg.Done()
		failCounter++
		return
	}

	results, err := client.ReadHoldingRegisters(0, 10)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("results %v\n", results)
	fmt.Printf("任务：%d，结束.......\n\n", num)
	fmt.Println("Fail Counter: ", failCounter)
	wg.Done()
}
