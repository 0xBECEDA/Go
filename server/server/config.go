package server

type Config struct {
	Address string `yaml:"address"`
	MaxConn int    `yaml:"max_conn"`
}
