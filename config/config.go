// Runtime configuration for networking, persistence, and eviction.
package config

const (
	RedisHost      string  = "localhost"
	RedisPort      string  = "8080"
	MaxStoreSize   int     = 10
	AofPath        string  = ""
	KeyLimit       uint8   = 100
	EvicationRatio float64 = 0.4
	EcitType       string  = "simple-first"
)
