package consumer

import (
	"go-ka/config"
)

type Manager[V any] struct {
	cfg        *config.ProcessorConfigs
	processors *Process[V]
}

type ManagerImpl interface {
	ExecuteAll() map[int32]int
}

func NewManager[V any](configs *config.ProcessorConfigs) *Manager[V] {

	v := NewProcess[V](configs)
	v.Consume()
	return &Manager[V]{
		cfg:        configs,
		processors: v,
	}
}

// ExecuteAll
/**
Execute all consumer
@Return : map[string]int, Key: Name of processor, Value : executed concurrency
*/
func (m *Manager[V]) ExecuteAll() map[int32]int {
	var retv = make(map[int32]int)
	m.processors.Consume()

	return retv
}
