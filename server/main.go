package main

import (
	"fmt"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	fmt.Println(judgemap)
	startTCPsrv()

}
