package hashmap

type Node[K MapKeyConstraint, V comparable] struct {
	value V
	key   *Key[K]
	child *Node[K, V]
}

func (n *Node[K, V]) hasNext() bool {
	return n.child != nil
}

func (n *Node[K, V]) forEach(fn func(*Node[K, V])) {
	current := n
	for ok := true; ok; ok = current.hasNext() {
		fn(n)
		if n.child != nil {
			current = n.child
		}
	}
}
