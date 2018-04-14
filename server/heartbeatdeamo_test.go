package main

import (
	"fmt"
	"testing"
)

func Test_scanAssetdb(t *testing.T) {
	n := scanAssetdb()
	fmt.Println(n)
	for ip, host := range ipHostnameMap {
		fmt.Println(host, ip, ipHeartRecorde[ip])
	}
	heartdeamo()
	//t.Log(err)
}
