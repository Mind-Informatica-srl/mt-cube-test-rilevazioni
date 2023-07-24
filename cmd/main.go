package main

import (
	"fmt"
	mtcubetest "mtcube-test-rilevazioni/internal"
	"os"
)

func main() {
	cmd := mtcubetest.NewCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
