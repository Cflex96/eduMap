package hashmap

import (
	"hash/maphash"
	"testing"
)

func TestCreateNodeInBucket(t *testing.T) {
	buckets := NewBuckets[string, int](10)
	seed := maphash.MakeSeed()
	testingKey := newKey("testing", seed)
	testingKey2 := newKey("testing2", seed)
	buckets.createNodeInBucket(1, testingKey, 1)
	buckets.createNodeInBucket(1, testingKey2, 2)

	if buckets[1].value != 1 {
		t.Log(buckets[1])
		t.Fatalf("expected value 1 in bucket 1 but recieved %v", buckets[1].value)
	}
	if buckets[1].key != testingKey {
		t.Log(buckets[1])
		t.Fatalf("expected %v but recieved %v", testingKey, buckets[1].key)
	}
	if buckets[1].child.value != 2 {
		t.Log(buckets[1].child)
		t.Fatalf("expected value 2 in bucket 1 but recieved %v", buckets[1].value)
	}
	if buckets[1].child.key != testingKey2 {
		t.Log(buckets[1].child)
		t.Fatalf("expected %v but recieved %v", testingKey2, buckets[1].child.key)
	}
}
