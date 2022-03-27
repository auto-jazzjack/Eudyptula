package consumer

import (
	"config"
)

type Manager struct {
	processors map[string]*Process
}

func NewManager(configs config.ProcessorConfigs) *Manager {

	results := make(map[string]*Process)

	for k, v := range configs.Processors {
		results[k] = NewProcess(v)
	}

	return &Manager{
		processors: results,
	}
}
