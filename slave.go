package main

type Slave struct {
	Addr string
	Conn redis.Conn
}
