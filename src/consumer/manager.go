package consumer

import (
	"fmt"
	"go-ka/config"
	"net/http"
	"time"
)

type Manager[V any] struct {
	cfg        *config.ProcessorConfigs[V]
	processors *Process[V]
}

type ManagerImpl interface {
	ExecuteAll() map[int32]int
	Rewind(date string) (map[int32]int, error)
}

func NewManager[V any](configs *config.ProcessorConfigs[V]) *Manager[V] {

	v := NewProcess(configs)
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
func (m *Manager[V]) ExecuteAll() map[string]int32 {
	return m.processors.Consume()
}

// Rewind
/**
date should be yyyy-mm-dd-hh-mm
*/
func (m *Manager[V]) Rewind(date string) (map[string][]string, error) {
	yourDate, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		return nil, fmt.Errorf("%d occurs check date format", http.StatusBadRequest)
	}

	return m.processors.Rewind(yourDate), nil
}
