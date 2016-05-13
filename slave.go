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

func IsSyncCommand(c []byte) bool {
	for i, k := range c {
		if i%2 == 0 && k <= 'z' && k >= 'a' {
			continue
		}
		if i%2 == 1 && k <= 'Z' && k >= 'A' {
			continue
		}

		return false
	}

	return true
}

func SetSyncCommand(c []byte) {
	for i, k := range c {
		if i%2 == 1 && k <= 'z' && k >= 'a' {
			c[i] = c[i] - 0x20
			continue
		}
		if i%2 == 0 && k <= 'Z' && k >= 'A' {
			c[i] = c[i] + 0x20
			continue
		}
	}
}

func CanSendToSlave(b []byte) bool {
	cmd, ok := GetRedisCommand(b)
	if !ok {
		return false
	}

	fmt.Println("cmd is", string(cmd), len(cmd))

	/* 不同步PING */
	if strings.EqualFold(string(cmd), "PING") {
		return false
	}

	if IsSyncCommand(cmd) {
		return false
	}

	SetSyncCommand(cmd)
	return true
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

	if !CanSendToSlave(b) {
		return nil
	}

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
