package main

import (
	"github.com/garyburd/redigo/redis"
)

type Master struct {
	Addr string
	conn redis.Conn
}

func SlaveOf(addr string) (*Master, error) {
}

func (m *Master) SlaveNoOne() error {
}

func (m *Master) replconfCron() {
}

func (m *Master) replconfCron() {
}

func (m *Master) readSyncPayload() {
}
