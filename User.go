package main

import (
	"net"
	"time"
)

type User struct {
	name     string
	ip       net.IP
	creation time.Time
}
