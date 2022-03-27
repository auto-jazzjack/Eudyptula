package config

type ProcessorConfigs struct {
	Processors map[string]ProcessorConfig
}

type ProcessorConfig struct {
	BoostrapServer string
	GroupId        string
	Offset         string
}
