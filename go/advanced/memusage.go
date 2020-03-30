package advanced

import (
	"fmt"
	"runtime"
)

// PrintMemoryUsage print memory stats
func PrintMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	printMem := func(info string, b uint64) {
		fmt.Printf("%s %d bytes %d mb\n", info, b, b/1024/1024)
	}
	printMem("Alloc", m.Alloc)
	printMem("Total ALloc", m.TotalAlloc)
	printMem("Sys Alloc", m.Sys)
	printMem("HeapAlloc", m.HeapAlloc)
	fmt.Println("NumGC=", m.NumGC)
}

func toMb(b uint64) uint64 {
	return b / 1024 / 1024
}
