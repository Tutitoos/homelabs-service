package shared

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

const (
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
)

func CapturePanic() {
	if r := recover(); r != nil {
		pc := make([]uintptr, 10)
		n := runtime.Callers(3, pc)
		frames := runtime.CallersFrames(pc[:n])

		var (
			file     = "unknown"
			line     = 0
			funcName = "unknown"
		)

		for {
			frame, more := frames.Next()
			if !strings.Contains(frame.Function, "handlePanic") {
				file = frame.File
				line = frame.Line
				funcName = frame.Function
				break
			}
			if !more {
				break
			}
		}

		fmt.Printf("\n%s%s* PANIC DETECTED%s\n", colorBold, colorRed, colorReset)
		fmt.Printf("  %sError:%s     🚨 %v\n", colorYellow, colorReset, r)
		fmt.Printf("  %sFunction:%s  %s\n", colorYellow, colorReset, funcName)
		fmt.Printf("  %sFile:%s      %s:%d\n", colorYellow, colorReset, file, line)
		fmt.Printf("  %sStack Trace:%s\n", colorYellow, colorReset)

		stackLines := strings.Split(string(debug.Stack()), "\n")
		for _, line := range stackLines {
			fmt.Printf("    %s%s%s\n", colorCyan, line, colorReset)
		}
	}

}
