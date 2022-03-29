package main

import (
	"fmt"
	"testing"
)

func TestSth(t *testing.T) {
	fmt.Println("test PASS")
	fmt.Println(getConfig("F:\\gitrepo\\hook_deploy\\config,json"))
	// Listen(12356)
}
