package hashmap

type Tuple[K MapKeyConstraint, V comparable] struct {
	key   K
	value V
}
