package shared

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v3/log"
)

func Collector() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	go func() {
		for range time.Tick(10 * time.Minute) {
			runtime.ReadMemStats(&m)

			heapAllocMB := float64(m.HeapAlloc) / 1024 / 1024
			if heapAllocMB <= 1024 {
				continue
			}

			text0 := "starting garbage collector..."
			log.Debug(text0)
			fmt.Printf("%s\n", text0)
			Logger.Debug(text0)

			text1 := fmt.Sprintf("before GC: HeapAlloc = %.2f MB, TotalAlloc = %.2f MB, NumGC = %v",
				float64(m.HeapAlloc)/1024/1024, float64(m.TotalAlloc)/1024/1024, m.NumGC)
			log.Debug(text1)
			fmt.Printf("%s\n", text1)
			Logger.Debug(text1)

			runtime.GC()

			runtime.ReadMemStats(&m)

			text2 := fmt.Sprintf("after GC: HeapAlloc = %.2f MB, TotalAlloc = %.2f MB, NumGC = %v",
				float64(m.HeapAlloc)/1024/1024, float64(m.TotalAlloc)/1024/1024, m.NumGC)

			log.Debug(text2)
			fmt.Printf("%s\n", text2)
			Logger.Debug(text2)

			text3 := "garbage collector executed!"
			log.Debug(text3)
			fmt.Printf("%s\n", text3)
			Logger.Debug(text3)
		}
	}()
}
