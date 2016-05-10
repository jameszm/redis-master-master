package main

type Dispatch struct {
	master *Master
	slaves []*Slave
}

func (*Dispatch) SetMaster(addr string) {
}

func (*Dispatch) AddSlave(addr string) int {
}

func (*Dispatch) Run() error {
}
