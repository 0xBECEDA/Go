package server

import "time"

type Config struct {
	Address        string        `yaml:"address"`
	MaxConn        int           `yaml:"max_conn"`
	TCPAlivePeriod time.Duration `yaml:"tcp_alive_period"`
	Timeout        time.Duration `yaml:"timeout"`
}
