package main

import (
	"fmt"
	"time"
)

type Dispatch struct {
	master *Master
	slave  *Slave
}

func (d *Dispatch) ReadPayload(b []byte, n int) error {
	fmt.Println("read payload:")
	fmt.Println(string(b[:n]))

	err := d.slave.Sync(b[:n])
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (d *Dispatch) SetMaster(host string, port uint16) {
	if d.master == nil {
		d.master = new(Master)
	}

	d.master.MasterId = "?"
	d.master.BaseOffset = -1
	d.master.Port = port
	d.master.Host = host

	d.master.Priv = d
	d.master.ReadCb = func(b []byte, n int, priv interface{}) error {
		d, ok := priv.(*Dispatch)
		if !ok {
			return nil
		}

		return d.ReadPayload(b, n)
	}
}

func (d *Dispatch) SetSlave(host string, port uint16) {
	if d.slave == nil {
		d.slave = new(Slave)
	}

	d.slave.Host = host
	d.slave.Port = port
}

func (d *Dispatch) Start() error {
	var err error

	err = d.slave.ConnSlave()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = d.master.SlaveOf()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (d *Dispatch) Stop() {
}

func main() {
	var d Dispatch
	/*
		d.SetMaster("127.0.0.1", 6001)
		d.SetSlave("127.0.0.1", 6002)
	*/
	d.SetMaster("127.0.0.1", 6002)
	d.SetSlave("127.0.0.1", 6001)

	d.Start()

	time.Sleep(time.Second * 3000)
}
