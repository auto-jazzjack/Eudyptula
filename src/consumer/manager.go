package consumer

import (
	"go-ka/config"
)

type Manager struct {
	cfg        *config.ProcessorConfigs
	processors *Process
}

type ManagerImpl interface {
	ExecuteAll() map[int32]int
}

func NewManager(configs *config.ProcessorConfigs) *Manager {

	v := NewProcess(configs)
	v.Consume()
	return &Manager{
		cfg:        configs,
		processors: v,
	}
}

// ExecuteAll
/**
Execute all consumer
@Return : map[string]int, Key: Name of processor, Value : executed concurrency
*/
func (m *Manager) ExecuteAll() map[int32]int {
	var retv = make(map[int32]int)
	m.processors.Consume()

	return retv
}
