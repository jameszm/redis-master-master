package main

import (
	"fmt"
	"strings"
)

type Slave struct {
	Host string
	Port uint16

	addr  string
	redis *Redis
}

func (s *Slave) ConnSlave() error {
	var err error

	s.addr = fmt.Sprintf("%s:%d", s.Host, s.Port)

	/* 连接redis */
	s.redis, err = Dial(s.addr)
	if err != nil {
		return err
	}

	/* PING */
	err = s.redis.PING()
	if err != nil {
		return err
	}

	fmt.Println("slave connect success")

	return nil
}

func (s *Slave) Sync(b []byte) error {
	var err error
	var n int

	fmt.Println("slave do")
	fmt.Println(string(b))

	_, err = s.redis.Send(b)
	if err != nil {
		return err
	}

	r := make([]byte, 4096)
	n, err = s.redis.Recv(r)
	if err != nil {
		return err
	}

	fmt.Println(string([]byte(r[:n-1])))

	return nil
}
