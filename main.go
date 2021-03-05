package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	fmt.Println(getCpuUsage())
	fmt.Println(getMemoryUsage())
}

func getCpuUsage() int {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Fatal(err)
	}

	return int(math.Ceil(percent[0]))
}

func getMemoryUsage() int {
	memory, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	return int(math.Ceil(memory.UsedPercent))
}
