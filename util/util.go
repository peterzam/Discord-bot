package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"syscall"
	"time"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

type DiskUsage struct {
	stat *syscall.Statfs_t
}

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return

}

func GetCPUinfo() string {
	idle0, total0 := getCPUSample()
	time.Sleep(3 * time.Second)
	idle1, total1 := getCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	busyTicks := totalTicks - idleTicks

	cpuUsageStr := fmt.Sprintf("%f", (100 * (totalTicks - idleTicks) / totalTicks))
	totalTicksStr := fmt.Sprintf("%f", totalTicks)
	busyStr := fmt.Sprintf("%f", busyTicks)

	message :=
		"\nCPU usage- " +
			"\nRunning     : " + cpuUsageStr + "%" +
			"\nTotal Ticks : " + totalTicksStr +
			"\nBusy Ticks  : " + busyStr

	return (message)
}

func GetMeminfo() string {

	stat, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Fatal("stat read fail")
	}

	free := strconv.FormatUint(stat.MemFree/1024, 10)
	total := strconv.FormatUint(stat.MemTotal/1024, 10)
	ava := strconv.FormatUint(stat.MemAvailable/1024, 10)
	cached := strconv.FormatUint(stat.Cached/1024, 10)
	cachedPercent := fmt.Sprintf("%f", (float32(stat.Cached)/float32(stat.MemTotal))*100)

	message := "---------------" +
		"\nMemory Usage- " +
		"\nFree         : " + free + " MB" +
		"\nTotal        : " + total + " MB" +
		"\nAvailable : " + ava + " MB" +
		"\nCached     : " + cached + " MB(" + cachedPercent + "%)"

	return message
}

func GetDiskinfo() string {

	var stat syscall.Statfs_t
	syscall.Statfs("/root", &stat)
	free := strconv.FormatUint(DiskUsage{&stat}.stat.Bfree*(uint64(DiskUsage{&stat}.stat.Bsize))/(1024*1024), 10)
	total := strconv.FormatUint(DiskUsage{&stat}.stat.Blocks*(uint64(DiskUsage{&stat}.stat.Bsize))/(1024*1024), 10)
	available := strconv.FormatUint(DiskUsage{&stat}.stat.Bavail*(uint64(DiskUsage{&stat}.stat.Bsize))/(1024*1024), 10)
	used := strconv.FormatUint((DiskUsage{&stat}.stat.Blocks*(uint64(DiskUsage{&stat}.stat.Bsize))-DiskUsage{&stat}.stat.Bfree*(uint64(DiskUsage{&stat}.stat.Bsize)))/(1024*1024), 10)
	usedPercent := fmt.Sprintf("%f", float32((DiskUsage{&stat}.stat.Blocks*(uint64(DiskUsage{&stat}.stat.Bsize))-DiskUsage{&stat}.stat.Bfree*(uint64(DiskUsage{&stat}.stat.Bsize))))/float32(DiskUsage{&stat}.stat.Blocks*(uint64(DiskUsage{&stat}.stat.Bsize)))*100)

	message := "---------------" +
		"\nDisk Usage- " +
		"\nFree        : " + free + " MB" +
		"\nTotal       : " + total + " MB" +
		"\nAvailable   : " + available + " MB" +
		"\nUsed        : " + used + " MB(" + usedPercent + "%)"

	return message
}

func NewDiskUsage(volumePath string) *DiskUsage {

	var stat syscall.Statfs_t
	syscall.Statfs(volumePath, &stat)
	return &DiskUsage{&stat}
}

func (this *DiskUsage) Free() uint64 {
	return this.stat.Bfree * uint64(this.stat.Bsize)
}

func (this *DiskUsage) Available() uint64 {
	return this.stat.Bavail * uint64(this.stat.Bsize)
}

func (this *DiskUsage) Size() uint64 {
	return this.stat.Blocks * uint64(this.stat.Bsize)
}

func (this *DiskUsage) Used() uint64 {
	return this.Size() - this.Free()
}

func (this *DiskUsage) Usage() float32 {
	return float32(this.Used()) / float32(this.Size())
}
