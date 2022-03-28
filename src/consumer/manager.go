package consumer

import (
	"go-ka/config"
	reflect "reflect"
)

type Manager struct {
	cfg        *config.ProcessorConfigs
	processors map[string]*Process[any]
}

type ManagerImpl interface {
	ExecuteAll() map[string]int
}

func NewManager(configs *config.ProcessorConfigs) *Manager {

	results := make(map[string]*Process[any])

	for k, v := range configs.Processors {
		a := (interface{}(reflect.New(*v.Clazz))).(Executor[interface{}])

		results[k] = NewProcess(v, a)
	}

	return &Manager{
		cfg:        configs,
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
