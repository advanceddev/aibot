package utils

import (
	"fmt"
	"runtime"
)

// PrintAlloc - выводит инофрмацию о затраченных ресурсах
func PrintAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d KB\n", m.Alloc/1024)
}
