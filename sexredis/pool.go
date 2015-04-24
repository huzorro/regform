package sexredis

import (
	"github.com/gosexy/redis"
)

type RedisPool struct {
	Connections chan *redis.Client
	ConnFn      func() (*redis.Client, error) // function to create new connection.
}

func (this *RedisPool) Get() (*redis.Client, error) {
	var conn *redis.Client
	select {
	case conn = <-this.Connections:
	default:
		conn, err := this.ConnFn()
		if err != nil {
			return nil, err
		}

		return conn, nil
	}

	if err := this.pingConn(conn); err != nil {
		return this.Get() // if connection is bad, get the next one in line until base case is hit, then create new client
	}

	return conn, nil
}

func (this *RedisPool) Close(conn *redis.Client) {
	select {
	case this.Connections <- conn:
		return
	default:
		conn.Quit()
	}
}

func (this *RedisPool) pingConn(conn *redis.Client) error {
	if _, err := conn.Ping(); err != nil {
		conn.Quit()
		return err
	}

	return nil
}
