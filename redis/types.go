package redis

import "github.com/redis/go-redis/v9"

type XAddArgs redis.XAddArgs

func (a XAddArgs) Value() *redis.XAddArgs {
	return (*redis.XAddArgs)(&a)
}

type XReadGroupArgs redis.XReadGroupArgs

func (a XReadGroupArgs) Value() *redis.XReadGroupArgs {
	return (*redis.XReadGroupArgs)(&a)
}

type XStream redis.XStream
