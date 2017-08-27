package main

import "time"

//message는 단일 메시지
type message struct {
	Name string
	Message string
	When time.Time
}