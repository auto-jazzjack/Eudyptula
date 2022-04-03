package logic

type Logic[V any] interface {
	Deserialize([]byte) V
	DoAction(V)
	DefaultValue() V
}
