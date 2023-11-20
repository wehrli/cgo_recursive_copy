package main

import (
	"fmt"
	"os"

	"technical_test/cgo_recursive/core"
	"technical_test/cgo_recursive/logger"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: copy_recursive_bin <source directory> <destination directory>")
		return
	}

	srcDir := os.Args[1]
	destDir := os.Args[2]

	logger := logger.NewLoggerWrapper("log.txt")
	defer logger.Close()
	dirProcessor := core.NewProcessor(logger)

	err := dirProcessor.CopyDirectory(srcDir, destDir)
	if err != nil {
		fmt.Println(err)
		return
	}
}
