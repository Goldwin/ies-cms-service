package data

const (
	ModeLocal = "local"
	ModeRedis = "redis"
	ModeMongo = "mongo"
)

type WorkerConfig struct {
	Mode           string `env:"WORKER_MODE" yaml:"mode"`
	DB             string `env:"WORKER_DB" yaml:"db"`
	UseTransaction bool   `env:"WORKER_USE_TRANSACTION" yaml:"useTransaction"`
}

type DataLayerConfig struct {
	CommandConfig *WorkerConfig `yaml:"command"`
	QueryConfig   *WorkerConfig `yaml:"query"`
}
