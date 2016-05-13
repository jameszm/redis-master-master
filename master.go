package main

import (
	"errors"
	"fmt"
	"time"
)

type ReadCallback func(b []byte, n int, priv interface{}) error

type Master struct {
	Host       string
	Port       uint16
	MasterId   string
	BaseOffset int
	ReadCb     ReadCallback
	Priv       interface{}

	addr      string
	redis     *Redis
	offset    int
	replTimer *time.Timer
}

func (m *Master) SlaveOf() error {
	var err error

	m.addr = fmt.Sprintf("%s:%d", m.Host, m.Port)

	/* 连接redis */
	m.redis, err = Dial(m.addr)
	if err != nil {
		return err
	}

	/* PING */
	err = m.redis.PING()
	if err != nil {
		return err
	}

	/* 判断是否为master */
	var role string
	role, err = m.redis.ROLE()
	if err != nil {
		return err
	}
	if role != "master" {
		return errors.New("role not master")
	}

	/* REPLCONF listen-port xxx
	 * 设置listen-port为redis端口+5000 */
	err = m.redis.REPLCONF_listen_port(m.Port + 5000)
	if err != nil {
		return err
	}

	/* REPLCONF capa eof */
	err = m.redis.REPLCONF_capa_eof()
	if err != nil {
		return err
	}

	/* PSYNC */
	_ = m.redis.PSYNC(m.MasterId, m.BaseOffset)
	if err != nil {
		return err
	}

	/* 定时执行REPLCONF ACK xxx */
	m.offset = 1
	m.replTimer = time.NewTimer(time.Millisecond * 500)
	go m.replconfCron()

	/* 循环接收数据 */
	go m.readCron()

	return nil
}

/* 停止监听 */
func (m *Master) SlaveNoOne() error {
	m.redis.Close()
	return nil
}

func (m *Master) replconfCron() {
	for {
		select {
		case <-m.replTimer.C:
			fmt.Println("REPLCONF ACK", m.offset)
			_ = m.redis.REPLCONF_ack(m.offset)
			m.replTimer.Reset(time.Second)
		}
	}
}

func (m *Master) readCron() {
	var err error
	var n int
	buf := make([]byte, 4096)

	for {
		n, err = m.redis.Recv(buf)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		m.offset = m.offset + n
		fmt.Println("offset change to", m.offset)

		if m.ReadCb == nil {
			continue
		}

		err = m.ReadCb(buf, n, m.Priv)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

/*
func main() {
	m, err := SlaveOf("127.0.0.1", 6001, "", -1)
	if err != nil {
		fmt.Println("slave of error", err.Error())
		return
	}

	time.Sleep(time.Second * 30)

	m.SlaveNoOne()
}
*/
