package main

import (
	"fmt"
	"time"
)

type Dispatch struct {
	master *Master
	slaves []*Slave
}

func (d *Dispatch) ReadPayload(b []byte, n int) error {
	fmt.Println("read payload:")
	fmt.Println(string(b[:n-1]))

	for _, s := range d.slaves {
		err := s.Do(b[:n])
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return nil
}

func (d *Dispatch) SetMaster(host string, port uint16) {
	if d.master == nil {
		var m Master
		d.master = &m
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

func (d *Dispatch) AddSlave(host string, port uint16) int {
	var s Slave
	s.Host = host
	s.Port = port

	if d.slaves == nil {
		d.slaves = make([]*Slave, 0)
	}

	d.slaves = append(d.slaves, &s)

	return 0
}

func (d *Dispatch) Run() error {
	var err error

	fmt.Println("slave len", len(d.slaves))

	for _, s := range d.slaves {
		err = s.ConnSlave()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	err = d.master.SlaveOf()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	time.Sleep(time.Second * 3000)

	return nil
}

/*
func main() {
	var d Dispatch
	d.SetMaster("127.0.0.1", 6001)
	d.AddSlave("127.0.0.1", 6002)

	d.Run()
}
*/
