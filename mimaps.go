package mimaps

type MiMap[K comparable, V any] interface {
	Put(K, V) error
	Get(K) (V, error)
	Delete(K) error
}
