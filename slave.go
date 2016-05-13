package main

import (
	"fmt"
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

func (s *Slave) Do(b []byte) error {
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

/*
func main() {
	var s Slave
	s.Host = "127.0.0.1"
	s.Port = 6002

	err := s.ConnSlave()
	if err != nil {
		fmt.Println(err.Error())
	}

	s.Do([]byte("*1\r\n$4\r\nPING\r\n"))
	//s.Do([]byte("PING"))
}
*/
