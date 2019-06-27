package main

import (
	"WsTestTool/kernel"
)

func main() {
	//Need help?
	if kernel.Help() {
		kernel.Usage()
		return
	}

	//Check param format
	err := kernel.CheckParam()
	if err != nil {
		kernel.Usage()
		return
	}

	//Start requests sending
	kernel.Run()

	//Watting result static
	kernel.Wait()
}
