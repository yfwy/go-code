package main

import "github.com/garyburd/redigo/redis"

var pool* redis.Pool

func init()  {
	pool=&redis.Pool{
		MaxIdle:8,
		MaxActive:0,
		IdleTimeout:100,
	}
}
