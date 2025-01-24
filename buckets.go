package hashmap

type Buckets[K MapKeyConstraint, V comparable] []*Node[K, V]

func NewBuckets[K MapKeyConstraint, V comparable](capacity uint) Buckets[K, V] {
	return make([]*Node[K, V], capacity)
}

func (b Buckets[K, V]) createNodeInBucket(bucketIndex uint, key *Key[K], value V) {
	nextNode := &Node[K, V]{
		value: value,
		key:   key,
		child: nil,
	}

	if b[bucketIndex] == nil {
		b[bucketIndex] = nextNode
	} else {
		node := b[bucketIndex]
		for node.hasNext() {
			node = node.child
		}
		node.child = nextNode
	}
}
