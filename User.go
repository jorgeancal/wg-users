package main

import (
	"time"
)

type User struct {
	name         string
	ip           string
	creation     time.Time
	publicKey    string
	privateKey   string
	presharedKey string
}
