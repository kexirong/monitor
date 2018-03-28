package main

import (
	"fmt"
	"testing"
)

func Test_scanAssetdb(t *testing.T) {
	n := scanAssetdb()
	fmt.Println(n)
	for host, ip := range hostIPMap {
		fmt.Println(host, ip, hostHeartRecorde[host])
	}
	heartdeamo()
	//t.Log(err)
}
