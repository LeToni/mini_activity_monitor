package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"log"
	"math"
	"strconv"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type DiskStatus struct {
	All   uint64 `json:"all"`
	Used  uint64 `json:"used"`
	Avail uint64 `json:"avail"`
}

func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)

	if err != nil {
		log.Fatal(err)
	}

	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Avail = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Avail
	return
}

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

	model, numberOfCores := getCpuInfo()
	systray.AddMenuItem(fmt.Sprintf("CPU: %s", model), "CPU Model")
	systray.AddMenuItem(fmt.Sprintf("Cores: %s", strconv.Itoa(int(numberOfCores))), "Number of Cores")

	disk := DiskUsage("/")
	systray.AddMenuItem(fmt.Sprintf("Total Disk Space: %.2f GB", float64(disk.All)/float64(GB)), "Total disk space")
	systray.AddMenuItem(fmt.Sprintf("Avail Disk Space: %.2f GB", float64(disk.Avail)/float64(GB)), "Available disk space")
	systray.AddMenuItem(fmt.Sprintf("Used Disk Space: %.2f GB", float64(disk.Used)/float64(GB)), "Used disk space")
	systray.AddSeparator()

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

func onExit() {
	// clean up here
}

func getCpuInfo() (string, int32) {
	info, err := cpu.Info()
	if err != nil {
		log.Fatal(err)
	}

	return info[0].ModelName, info[0].Cores
}
