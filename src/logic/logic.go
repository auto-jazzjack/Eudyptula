package logic

type Logic[V any] interface {
	Deserialize([]byte) *V
	DoAction(V) error
	DefaultValue() *V
}
