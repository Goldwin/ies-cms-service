package data

const (
	ModeLocal = "local"
	ModeRedis = "redis"
	ModeMongo = "mongo"
)

type WorkerConfig struct {
	Mode           string `yaml:"mode"`
	UseTransaction bool   `yaml:"useTransaction"`
}

type DataLayerConfig struct {
	CommandConfig *WorkerConfig `yaml:"command"`
	QueryConfig   *WorkerConfig `yaml:"query"`
}
