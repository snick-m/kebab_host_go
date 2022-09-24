package main

import (
	"fmt"
	"io"
	"os"
)

func getPid() int {
	filedata, err := os.ReadFile("_pid")
	if os.IsNotExist(err) {
		return -1
	} else if err != nil {
		panic(err)
	}

	var pid int
	_, err = fmt.Sscanf(string(filedata), "%d", &pid)
	if err == io.EOF {
		return -1
	} else if err != nil {
		panic(err)
	}

	return pid
}
