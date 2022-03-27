package consumer

import (
	"config"
)

type Manager struct {
	processors map[string]*Process
}

type ManagerImpl interface {
	ExecuteAll()
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

// ExecuteAll /**Execute all consumer*/
func (m *Manager) ExecuteAll() {
	for _, v := range m.processors {
		v.Consume()
	}
}
