package consumer

type Executor[V interface{}] interface {
	DeSerialize(data []byte) V
	DoAction(data V)
}
