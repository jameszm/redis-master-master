package main

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

type Master struct {
	host     string
	port     uint16
	addr     string
	redis    redis.Conn
	masterId string
	offset   uint64
}

func SlaveOf(host string, port uint16, masterId string, offset uint64) (*Master, error) {
	var m Master
	var err error

	m.host = host
	m.port = port
	m.addr = host + ":" + strconv.Itoa(port)

	if masterId == "" {
		m.masterId = "?"
	} else {
		m.masterId = masterId
	}

	/* 连接redis */
	m.redis, err = redis.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	/* PING */
	_, err = m.redis.Do("PING")
	if err != nil {
		return nil, err
	}

	/* 判断是否为master */
	var role string
	role, err = m.role()
	if err != nil {
		return nil, err
	}
	if role != "master" {
		return nil, errors.New("role not master")
	}

	/* REPLCONF listen-port xxx
	 * 设置listen-port为redis端口+5000 */
	_, err = m.redis.Do("REPLCONF listen-port " + strconv.Itoa(m.port+5000))
	if err != nil {
		return err
	}

	/* PSYNC */
}

/* 获取所属角色
 * 返回 slave|master */
func (m *Master) role() (string, error) {
}

func (m *Master) SlaveNoOne() error {
}

func (m *Master) replconfCron() {
}

func (m *Master) replconfCron() {
}

func (m *Master) readSyncPayload() {
}
