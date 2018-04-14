package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func Test_alarmdo(t *testing.T) {

	alarmdo()
	path, _ := os.Getwd()
	fmt.Println(path, time.Now().Format("2006-01-02 15:04:05"))
	wait := time.Tick(1000 * time.Millisecond)

	for range wait {

		fmt.Println("wait..", time.Now().Format("2006-01-02 15:04:05"))
	}

}
