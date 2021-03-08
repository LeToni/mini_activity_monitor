package main

import (
	"github.com/getlantern/systray"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	go func() {
		var result string
		for {
			result = getData()
			systray.SetTitle(result)
		}
	}()

	go func() {
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func getData() string {
	cpuUsage := "CPU: " + strconv.Itoa(getCpuUsage()) + "% "
	memoryUsage := "Mem: " + strconv.Itoa(getMemoryUsage()) + "% "
	return cpuUsage + memoryUsage
}

func onExit() {
	// clean up here
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
