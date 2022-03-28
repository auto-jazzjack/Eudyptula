package consumer

type Executor[V any] interface {
	DeSerialize(data []byte, target V) V
	DefaultValue() V
	DoAction(data V)
}
