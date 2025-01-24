package hashmap

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"hash/maphash"
	"math"

	"golang.org/x/exp/constraints"
)

type number interface {
	constraints.Float | constraints.Integer
}

type MapKeyConstraint interface {
	~string | number | ~struct{}
}

type HashMap[K MapKeyConstraint, V comparable] struct {
	buckets    Buckets[K, V]
	size       uint
	seed       maphash.Seed
	capacity   uint
	loadFactor float64
}

func New[K MapKeyConstraint, V comparable](capacity uint) HashMap[K, V] {
	return HashMap[K, V]{
		buckets:    make([]*Node[K, V], capacity),
		size:       0,
		seed:       maphash.MakeSeed(),
		capacity:   capacity,
		loadFactor: 0.6,
	}
}

func NewFromSlices[K MapKeyConstraint, V comparable](keys []K, values []V) (HashMap[K, V], error) {
	keySize := len(keys)
	valueSize := len(values)
	if keySize != valueSize {
		return HashMap[K, V]{}, errors.New("cannot create map from arrays of different length")
	}

	items := make([]*Node[K, V], keySize*2)
	m := HashMap[K, V]{
		buckets:    items,
		size:       uint(keySize) * 3,
		seed:       maphash.MakeSeed(),
		capacity:   uint(keySize) * 3,
		loadFactor: 0.6,
	}

	for index, key := range keys {
		m.Set(key, values[index])
	}

	return m, nil
}

func newKey[K MapKeyConstraint](keyLiteral K, seed maphash.Seed) *Key[K] {
	return &Key[K]{
		literal: keyLiteral,
		hash:    hashKey(keyLiteral, seed),
	}
}

func hashKey[K MapKeyConstraint](keyLiteral K, seed maphash.Seed) uint {
	var hashFunc maphash.Hash
	keyBinary := encodeKey(keyLiteral)
	hashFunc.SetSeed(seed)
	hashFunc.Write(keyBinary)
	return uint(hashFunc.Sum64())
}

type Key[K MapKeyConstraint] struct {
	literal K
	hash    uint
}

func (k *Key[K]) getHash() uint {
	return k.hash
}

func (m *HashMap[K, V]) calculateBucket(keyLiteral K) uint {
	hash := hashKey(keyLiteral, m.seed)
	return m.calculateBucketFromHash(hash)
}

func (m *HashMap[K, V]) calculateBucketFromHash(hash uint) uint {
	return hash % m.capacity
}

func (m *HashMap[K, V]) findParentNode(bucketIndex uint, keyLiteral K) *Node[K, V] {
	node := m.buckets[bucketIndex]
	if node == nil {
		return nil
	}
	prevNode := node
	for keyLiteral != node.key.literal && node.hasNext() {
		prevNode = node
		node = node.child
	}
	if keyLiteral == node.key.literal {
		return prevNode
	}
	return nil
}

func (m *HashMap[K, V]) grow() {
	newBuckets := NewBuckets[K, V](m.capacity * 2)
	var newSize uint = 0
	m.capacity = m.capacity * 2
	for _, currentBucket := range m.buckets {
		if currentBucket == nil {
			continue
		}
		currentBucket.forEach(func(node *Node[K, V]) {
			newBucketIndex := m.calculateBucketFromHash(currentBucket.key.hash)
			newBuckets.createNodeInBucket(newBucketIndex, currentBucket.key, currentBucket.value)
			newSize++
		})
	}
	m.buckets = newBuckets
	m.size = newSize
}

func (m *HashMap[K, V]) reorderItems() {}

func (m *HashMap[K, V]) findNode(bucketIndex uint, keyLiteral K) *Node[K, V] {
	node := m.buckets[bucketIndex]
	if node == nil {
		return nil
	}
	for keyLiteral != node.key.literal && node.hasNext() {
		node = node.child
	}
	if keyLiteral == node.key.literal {
		return node
	}
	return nil
}

func (m *HashMap[K, V]) appendNodeToBucket(bucket uint, key *Key[K], value V) {
	m.buckets.createNodeInBucket(bucket, key, value)
	m.size++
}

func (m *HashMap[K, V]) Set(keyLiteral K, value V) {
	if float64(m.size+1)/float64(m.capacity) >= m.loadFactor {
		m.grow()
	}
	bucket := m.calculateBucket(keyLiteral)

	key := newKey(keyLiteral, m.seed)
	m.buckets.createNodeInBucket(bucket, key, value)
	m.size++
}

func (m *HashMap[K, V]) Get(keyLiteral K) V {
	bucket := m.calculateBucket(keyLiteral)
	node := m.findNode(bucket, keyLiteral)
	if node == nil {
		var nullValue V
		return nullValue
	}
	return node.value
}

func (m *HashMap[K, V]) Remove(keyLiteral K) {
	bucket := m.calculateBucket(keyLiteral)
	if m.buckets[bucket].child == nil {
		m.buckets[bucket] = nil
		return
	}
	parentNode := m.findParentNode(bucket, keyLiteral)
	node := parentNode.child
	if node.child != nil {
		parentNode.child = node.child
	}
	parentNode.child = nil
}

func encodeKey[K MapKeyConstraint](keyLiteral K) []byte {
	var keyBinary []byte
	switch t := any(keyLiteral).(type) {
	case string:
		keyBinary = []byte(t)
	case float32:
		keyBinary = make([]byte, 4)
		binary.LittleEndian.PutUint32(keyBinary, math.Float32bits(t))
	case float64:
		keyBinary = make([]byte, 8)
		binary.LittleEndian.PutUint64(keyBinary, math.Float64bits(t))
	case int:
		keyBinary = make([]byte, 8)
		binary.LittleEndian.PutUint64(keyBinary, uint64(t))
	case int8:
		keyBinary = make([]byte, 2)
		binary.LittleEndian.PutUint16(keyBinary, uint16(t))
	case int16:
		keyBinary = make([]byte, 2)
		binary.LittleEndian.PutUint16(keyBinary, uint16(t))
	case int32:
		keyBinary = make([]byte, 4)
		binary.LittleEndian.PutUint32(keyBinary, uint32(t))
	case int64:
		keyBinary = make([]byte, 8)
		binary.LittleEndian.PutUint64(keyBinary, uint64(t))
	case struct{}:
		var buf bytes.Buffer
		encoder := gob.NewEncoder(&buf)
		encoder.Encode(t)
		keyBinary = buf.Bytes()
	default:
		panic("invalid key type")
	}
	return keyBinary
}
