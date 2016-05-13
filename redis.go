package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type Redis struct {
	addr string
	conn net.Conn
}

func Dial(addr string) (*Redis, error) {
	var r Redis
	var err error
	r.conn, err = net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	r.addr = addr

	return &r, nil
}

func (r *Redis) Close() {
	r.conn.Close()
}

func (r *Redis) PING() error {
	var err error
	_, err = r.Send([]byte("PING\r\n"))
	if err != nil {
		return err
	}

	b := make([]byte, 4096)
	_, err = r.Recv(b)
	if err != nil {
		return err
	}

	return nil
}

/* 获取所属角色
 * 返回 slave|master */
func (r *Redis) ROLE() (string, error) {
	/* TODO */
	return "master", nil
}

func (r *Redis) REPLCONF_ack(offset int) error {
	cmd := fmt.Sprintf("*3\r\n$8\r\nREPLCONF\r\n$3\r\nACK\r\n$%d\r\n%d\r\n", len(strconv.Itoa(offset)), offset)
	_, err := r.Send([]byte(cmd))
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) REPLCONF_capa_eof() error {
	var err error

	_, err = r.Send([]byte("REPLCONF capa eof\r\n"))
	if err != nil {
		return err
	}

	b := make([]byte, 4096)
	_, err = r.Recv(b)
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) REPLCONF_listen_port(port uint16) error {
	var err error

	cmd := fmt.Sprintf("REPLCONF listening-port %d\r\n", port)
	_, err = r.Send([]byte(cmd))
	if err != nil {
		return err
	}

	b := make([]byte, 4096)
	_, err = r.Recv(b)
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) PSYNC(masterId string, offset int) error {
	var err error

	cmd := fmt.Sprintf("PSYNC %s %d\r\n", masterId, offset)
	_, err = r.Send([]byte(cmd))
	if err != nil {
		return err
	}

	time.Sleep(time.Millisecond * 100)

	b := make([]byte, 4096)
	_, err = r.Recv(b)
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Recv(b []byte) (int, error) {
	n, err := r.conn.Read(b)
	return n, err
}

func (r *Redis) Send(b []byte) (int, error) {
	n, err := r.conn.Write(b)
	return n, err
}
