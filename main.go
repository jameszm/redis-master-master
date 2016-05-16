package main

import (
	//"encoding/json"
	//"fmt"
	"time"
)

func main() {
	/*
		var d1, d2 Dispatch

		d1.SetMaster("127.0.0.1", 6001, "", 0)
		d1.SetSlave("127.0.0.1", 6002)

		d2.SetMaster("127.0.0.1", 6002, "", 0)
		d2.SetSlave("127.0.0.1", 6001)

		d1.Start()
		d2.Start()
	*/

	var d1 Dispatch

	d1.SetMaster("127.0.0.1", 6001, "dd9c508557638d393f95a56d405dd4344886a2e6", 1400)
	d1.SetSlave("127.0.0.1", 6002)

	d1.Start()

	time.Sleep(time.Second * 30000)
}
