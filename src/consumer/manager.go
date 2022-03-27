package consumer

import (
	"go-ka/config"
)

type Manager struct {
	processors map[string]*Process
}

type ManagerImpl interface {
	ExecuteAll() map[string]int
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

// ExecuteAll
/**
Execute all consumer
@Return : map[string]int, Key: Name of processor, Value : executed concurrency
*/
func (m *Manager) ExecuteAll() map[string]int {
	var retv = make(map[string]int)
	for k, v := range m.processors {
		retv[k] = v.Consume()
	}
	return retv
}
