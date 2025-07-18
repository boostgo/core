package orderedmap

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"
)

// Test basic operations
func TestBasicOperations(t *testing.T) {
	om := NewOrderedMap[string, int]()

	// Test Set and Get
	om.Set("one", 1)
	om.Set("two", 2)
	om.Set("three", 3)

	if val, ok := om.Get("two"); !ok || val != 2 {
		t.Errorf("Get failed: expected 2, got %v", val)
	}

	// Test Len
	if om.Len() != 3 {
		t.Errorf("Len failed: expected 3, got %d", om.Len())
	}

	// Test Contains
	if !om.Contains("one") {
		t.Error("Contains failed: 'one' should exist")
	}
	if om.Contains("four") {
		t.Error("Contains failed: 'four' should not exist")
	}

	// Test Delete
	if !om.Delete("two") {
		t.Error("Delete failed: should return true for existing key")
	}
	if om.Delete("four") {
		t.Error("Delete failed: should return false for non-existing key")
	}
	if om.Len() != 2 {
		t.Errorf("Len after delete failed: expected 2, got %d", om.Len())
	}

	// Test Clear
	om.Clear()
	if om.Len() != 0 {
		t.Errorf("Clear failed: expected 0, got %d", om.Len())
	}
}

// Test insertion order is maintained
func TestInsertionOrder(t *testing.T) {
	om := NewOrderedMap[int, string]()

	// Insert in specific order
	insertOrder := []int{5, 2, 8, 1, 9, 3}
	for _, k := range insertOrder {
		om.Set(k, fmt.Sprintf("value-%d", k))
	}

	// Check Keys maintains insertion order
	keys := om.Keys()
	if !reflect.DeepEqual(keys, insertOrder) {
		t.Errorf("Keys order mismatch: expected %v, got %v", insertOrder, keys)
	}

	// Check Values maintains insertion order
	values := om.Values()
	for i, v := range values {
		expected := fmt.Sprintf("value-%d", insertOrder[i])
		if v != expected {
			t.Errorf("Values order mismatch at %d: expected %s, got %s", i, expected, v)
		}
	}
}

// Test update maintains order but updates value
func TestUpdateMaintainsOrder(t *testing.T) {
	om := NewOrderedMap[string, int]()

	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)
	om.Set("b", 20) // Update

	keys := om.Keys()
	expected := []string{"a", "c", "b"} // b moved to end after update
	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Update order failed: expected %v, got %v", expected, keys)
	}

	if val, _ := om.Get("b"); val != 20 {
		t.Errorf("Update value failed: expected 20, got %d", val)
	}
}

// Test thread safety
func TestThreadSafety(t *testing.T) {
	om := NewOrderedMap[int, int](WithThreadSafe[int, int](true))
	var wg sync.WaitGroup
	concurrent := 100
	operations := 1000

	// Concurrent writes
	wg.Add(concurrent)
	for i := 0; i < concurrent; i++ {
		go func(goroutine int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				key := goroutine*operations + j
				om.Set(key, key*2)
			}
		}(i)
	}
	wg.Wait()

	// Verify all writes succeeded
	expected := concurrent * operations
	if om.Len() != expected {
		t.Errorf("Concurrent writes failed: expected %d items, got %d", expected, om.Len())
	}

	// Concurrent reads
	wg.Add(concurrent)
	for i := 0; i < concurrent; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				key := rand.Intn(expected)
				om.Get(key)
			}
		}()
	}
	wg.Wait()

	// Mixed operations
	wg.Add(concurrent * 3)
	for i := 0; i < concurrent; i++ {
		// Writers
		go func(goroutine int) {
			defer wg.Done()
			for j := 0; j < operations/2; j++ {
				key := goroutine*operations + j
				om.Set(key, key*3)
			}
		}(i)

		// Readers
		go func() {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				om.Keys()
			}
		}()

		// Deleters
		go func(goroutine int) {
			defer wg.Done()
			for j := 0; j < operations/4; j++ {
				key := goroutine*operations + j
				om.Delete(key)
			}
		}(i)
	}
	wg.Wait()
}

// Test capacity and eviction
func TestCapacityAndEviction(t *testing.T) {
	capacity := 5
	om := NewOrderedMap[int, string](WithCapacity[int, string](capacity))

	// Add more items than capacity
	for i := 0; i < 10; i++ {
		om.Set(i, fmt.Sprintf("value-%d", i))
	}

	// Check capacity is maintained
	if om.Len() != capacity {
		t.Errorf("Capacity not maintained: expected %d, got %d", capacity, om.Len())
	}

	// Check oldest items were evicted (0-4 should be gone, 5-9 should remain)
	for i := 0; i < 5; i++ {
		if om.Contains(i) {
			t.Errorf("Item %d should have been evicted", i)
		}
	}
	for i := 5; i < 10; i++ {
		if !om.Contains(i) {
			t.Errorf("Item %d should not have been evicted", i)
		}
	}
}

// Test eviction callback
func TestEvictionCallback(t *testing.T) {
	evicted := make(map[int]string)
	var mu sync.Mutex

	config := Config[int, string]{
		Capacity:     3,
		CapacityMode: EvictOldest,
		OnEvict: func(key int, value string) {
			mu.Lock()
			evicted[key] = value
			mu.Unlock()
		},
	}
	om := NewOrderedMapWithConfig(config)

	// Add items
	for i := 0; i < 5; i++ {
		om.Set(i, fmt.Sprintf("value-%d", i))
	}

	// Check eviction callback was called
	mu.Lock()
	defer mu.Unlock()
	if len(evicted) != 2 {
		t.Errorf("Eviction callback count wrong: expected 2, got %d", len(evicted))
	}
	if evicted[0] != "value-0" || evicted[1] != "value-1" {
		t.Errorf("Wrong items evicted: %v", evicted)
	}
}

// Test LRU operations
func TestLRUOperations(t *testing.T) {
	om := NewOrderedMap[string, int]()

	// Setup initial state
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)
	om.Set("d", 4)

	// Test MoveToFront
	if !om.MoveToFront("c") {
		t.Error("MoveToFront should return true for existing key")
	}
	keys := om.Keys()
	if keys[0] != "c" {
		t.Errorf("MoveToFront failed: expected 'c' at front, got %s", keys[0])
	}

	// Test MoveToBack
	if !om.MoveToBack("a") {
		t.Error("MoveToBack should return true for existing key")
	}
	keys = om.Keys()
	if keys[len(keys)-1] != "a" {
		t.Errorf("MoveToBack failed: expected 'a' at back, got %s", keys[len(keys)-1])
	}

	// Test Touch (alias for MoveToBack)
	om.Touch("c")
	keys = om.Keys()
	if keys[len(keys)-1] != "c" {
		t.Errorf("Touch failed: expected 'c' at back, got %s", keys[len(keys)-1])
	}

	// Test GetAndTouch
	val, ok := om.GetAndTouch("b")
	if !ok || val != 2 {
		t.Errorf("GetAndTouch failed: got %v, %v", val, ok)
	}
	keys = om.Keys()
	if keys[len(keys)-1] != "b" {
		t.Errorf("GetAndTouch failed: expected 'b' at back, got %s", keys[len(keys)-1])
	}
}

// Test MoveBefore and MoveAfter
func TestMoveRelative(t *testing.T) {
	om := NewOrderedMap[string, int]()

	// Setup
	for _, k := range []string{"a", "b", "c", "d", "e"} {
		om.Set(k, 1)
	}

	// Test MoveBefore
	om.MoveBefore("e", "b")
	keys := om.Keys()
	expectedBefore := []string{"a", "e", "b", "c", "d"}
	if !reflect.DeepEqual(keys, expectedBefore) {
		t.Errorf("MoveBefore failed: expected %v, got %v", expectedBefore, keys)
	}

	// Test MoveAfter
	om.MoveAfter("a", "c")
	keys = om.Keys()
	expectedAfter := []string{"e", "b", "c", "a", "d"}
	if !reflect.DeepEqual(keys, expectedAfter) {
		t.Errorf("MoveAfter failed: expected %v, got %v", expectedAfter, keys)
	}
}

// Test batch operations
func TestBatchOperations(t *testing.T) {
	om := NewOrderedMap[int, string]()

	// Test SetBatch
	items := []Item[int, string]{
		{Key: 1, Value: "one"},
		{Key: 2, Value: "two"},
		{Key: 3, Value: "three"},
	}
	om.SetBatch(items)

	if om.Len() != 3 {
		t.Errorf("SetBatch failed: expected 3 items, got %d", om.Len())
	}

	// Test GetBatch
	results := om.GetBatch([]int{1, 3, 5})
	if len(results) != 2 {
		t.Errorf("GetBatch failed: expected 2 results, got %d", len(results))
	}
	if results[1] != "one" || results[3] != "three" {
		t.Errorf("GetBatch wrong values: %v", results)
	}

	// Test DeleteBatch
	deleted := om.DeleteBatch([]int{2, 3, 4})
	if deleted != 2 {
		t.Errorf("DeleteBatch failed: expected 2 deleted, got %d", deleted)
	}
	if om.Len() != 1 {
		t.Errorf("DeleteBatch failed: expected 1 item remaining, got %d", om.Len())
	}
}

// Test iterators
func TestIterators(t *testing.T) {
	om := NewOrderedMap[int, string]()
	for i := 1; i <= 5; i++ {
		om.Set(i, fmt.Sprintf("value-%d", i))
	}

	// Test forward iterator
	iter := om.NewIterator()
	var forwardKeys []int
	for iter.HasNext() {
		k, _, ok := iter.Next()
		if !ok {
			t.Error("Iterator Next() returned false")
		}
		forwardKeys = append(forwardKeys, k)
	}
	if !reflect.DeepEqual(forwardKeys, []int{1, 2, 3, 4, 5}) {
		t.Errorf("Forward iterator failed: %v", forwardKeys)
	}

	// Test reverse iterator
	revIter := om.NewReverseIterator()
	var reverseKeys []int
	for revIter.HasNext() {
		k, _, _ := revIter.Next()
		reverseKeys = append(reverseKeys, k)
	}
	if !reflect.DeepEqual(reverseKeys, []int{5, 4, 3, 2, 1}) {
		t.Errorf("Reverse iterator failed: %v", reverseKeys)
	}

	// Test iterator reset
	iter.Reset()
	if !iter.HasNext() {
		t.Error("Iterator reset failed")
	}
}

// Test functional operations
func TestFunctionalOperations(t *testing.T) {
	om := NewOrderedMap[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)
	om.Set("d", 4)

	// Test Filter
	filtered := om.Filter(func(k string, v int) bool {
		return v%2 == 0
	})
	if len(filtered) != 2 {
		t.Errorf("Filter failed: expected 2 items, got %d", len(filtered))
	}

	// Test Map
	mapped := om.Map(func(k string, v int) interface{} {
		return fmt.Sprintf("%s:%d", k, v*2)
	})
	if len(mapped) != 4 {
		t.Errorf("Map failed: expected 4 items, got %d", len(mapped))
	}

	// Test Find
	key, val, found := om.Find(func(k string, v int) bool {
		return v == 3
	})
	if !found || key != "c" || val != 3 {
		t.Errorf("Find failed: got %s, %d, %v", key, val, found)
	}

	// Test Any
	if !om.Any(func(k string, v int) bool { return v > 3 }) {
		t.Error("Any failed: should return true")
	}
	if om.Any(func(k string, v int) bool { return v > 10 }) {
		t.Error("Any failed: should return false")
	}

	// Test All
	if !om.All(func(k string, v int) bool { return v > 0 }) {
		t.Error("All failed: should return true")
	}
	if om.All(func(k string, v int) bool { return v > 2 }) {
		t.Error("All failed: should return false")
	}
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
	om := NewOrderedMap[string, int]()

	// Test operations on empty map
	if val, ok := om.Get("nonexistent"); ok || val != 0 {
		t.Error("Get on empty map should return zero value and false")
	}

	if om.Delete("nonexistent") {
		t.Error("Delete on empty map should return false")
	}

	keys := om.Keys()
	if keys != nil && len(keys) != 0 {
		t.Error("Keys on empty map should return nil or empty slice")
	}

	// Test First/Last on empty map
	if _, _, ok := om.First(); ok {
		t.Error("First on empty map should return false")
	}
	if _, _, ok := om.Last(); ok {
		t.Error("Last on empty map should return false")
	}

	// Test Pop operations on empty map
	if _, ok := om.Pop("key"); ok {
		t.Error("Pop on empty map should return false")
	}
	if _, _, ok := om.PopFirst(); ok {
		t.Error("PopFirst on empty map should return false")
	}
	if _, _, ok := om.PopLast(); ok {
		t.Error("PopLast on empty map should return false")
	}
}

// Test Clone
func TestClone(t *testing.T) {
	om := NewOrderedMap[string, int](WithCapacity[string, int](10))
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	clone := om.Clone()

	// Verify clone has same content
	if !reflect.DeepEqual(om.Keys(), clone.Keys()) {
		t.Error("Clone keys don't match")
	}

	// Verify clone is independent
	clone.Set("d", 4)
	clone.Set("e", 5)
	om.Set("e", 5)

	if clone.Len() == om.Len() {
		t.Error("Clone should be independent")
	}

	// Verify capacity is cloned
	if clone.Capacity() != om.Capacity() {
		t.Error("Clone should have same capacity")
	}
}

// Test capacity modes
func TestCapacityModes(t *testing.T) {
	t.Run("EvictOldest", func(t *testing.T) {
		config := Config[int, string]{
			Capacity:     3,
			CapacityMode: EvictOldest,
		}
		om := NewOrderedMapWithConfig(config)

		for i := 0; i < 5; i++ {
			om.Set(i, fmt.Sprintf("v%d", i))
		}

		// Should keep 2, 3, 4
		for i := 0; i < 2; i++ {
			if om.Contains(i) {
				t.Errorf("EvictOldest: %d should be evicted", i)
			}
		}
		for i := 2; i < 5; i++ {
			if !om.Contains(i) {
				t.Errorf("EvictOldest: %d should be kept", i)
			}
		}
	})

	t.Run("EvictNewest", func(t *testing.T) {
		config := Config[int, string]{
			Capacity:     3,
			CapacityMode: EvictNewest,
		}
		om := NewOrderedMapWithConfig(config)

		for i := 0; i < 5; i++ {
			om.SetWithMode(i, fmt.Sprintf("v%d", i))
		}

		// Should keep early items
		if om.Len() != 3 {
			t.Errorf("EvictNewest: expected 3 items, got %d", om.Len())
		}
	})

	t.Run("RejectNew", func(t *testing.T) {
		config := Config[int, string]{
			Capacity:     3,
			CapacityMode: RejectNew,
		}
		om := NewOrderedMapWithConfig(config)

		// Fill to capacity
		for i := 0; i < 3; i++ {
			if !om.SetWithMode(i, fmt.Sprintf("v%d", i)) {
				t.Errorf("RejectNew: should accept item %d", i)
			}
		}

		// Should reject new items
		if om.SetWithMode(3, "v3") {
			t.Error("RejectNew: should reject item when at capacity")
		}

		// Should still allow updates
		if !om.SetWithMode(1, "updated") {
			t.Error("RejectNew: should allow updates")
		}
	})
}

// Test memory operations
func TestMemoryOperations(t *testing.T) {
	om := NewOrderedMap[int, string]()

	// Test Reserve
	om.Reserve(100)
	// Can't directly test internal capacity, but ensure it doesn't crash

	// Add some items
	for i := 0; i < 50; i++ {
		om.Set(i, fmt.Sprintf("value-%d", i))
	}

	// Test Compact
	om.Compact()
	if om.Len() != 50 {
		t.Errorf("Compact should not change length: got %d", om.Len())
	}

	// Test LoadFactor
	lf := om.LoadFactor()
	if lf <= 0 || lf > 1 {
		t.Errorf("LoadFactor out of range: %f", lf)
	}
}

// Test complex scenarios
func TestComplexScenarios(t *testing.T) {
	t.Run("Deduplication", func(t *testing.T) {
		om := NewOrderedMap[int, string]()

		// Simulate duplicate data from database
		data := []struct {
			ID   int
			Name string
		}{
			{1, "Alice"}, {2, "Bob"}, {1, "Alice2"}, {3, "Charlie"}, {2, "Bob2"},
		}

		for _, row := range data {
			om.Set(row.ID, row.Name)
		}

		// Should have 3 unique items with last values
		if om.Len() != 3 {
			t.Errorf("Deduplication failed: expected 3, got %d", om.Len())
		}

		// Check last value wins
		if val, _ := om.Get(1); val != "Alice2" {
			t.Errorf("Expected Alice2, got %s", val)
		}
		if val, _ := om.Get(2); val != "Bob2" {
			t.Errorf("Expected Bob2, got %s", val)
		}
	})

	t.Run("LRUCache", func(t *testing.T) {
		cache := NewOrderedMap[string, string](WithCapacity[string, string](3))

		// Add items
		cache.Set("a", "1")
		cache.Set("b", "2")
		cache.Set("c", "3")

		// Access 'a' to make it recently used
		cache.GetAndTouch("a")

		// Add new item, should evict 'b'
		cache.Set("d", "4")

		if cache.Contains("b") {
			t.Error("LRU: 'b' should have been evicted")
		}
		if !cache.Contains("a") {
			t.Error("LRU: 'a' should still exist (was touched)")
		}
	})
}

// Benchmarks
func BenchmarkGet(b *testing.B) {
	om := NewOrderedMap[int, int]()
	for i := 0; i < 10000; i++ {
		om.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		om.Get(i % 10000)
	}
}

func BenchmarkSet(b *testing.B) {
	om := NewOrderedMap[int, int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		om.Set(i, i)
	}
}

func BenchmarkDelete(b *testing.B) {
	om := NewOrderedMap[int, int]()
	for i := 0; i < b.N; i++ {
		om.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		om.Delete(i)
	}
}

func BenchmarkIteration(b *testing.B) {
	om := NewOrderedMap[int, int]()
	for i := 0; i < 1000; i++ {
		om.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		om.ForEach(func(k, v int) {
			// Just iterate
		})
	}
}

func BenchmarkConcurrentAccess(b *testing.B) {
	om := NewOrderedMap[int, int](WithThreadSafe[int, int](true))
	for i := 0; i < 1000; i++ {
		om.Set(i, i)
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%3 == 0 {
				om.Set(i, i)
			} else if i%3 == 1 {
				om.Get(i % 1000)
			} else {
				om.Delete(i % 1000)
			}
			i++
		}
	})
}

func BenchmarkBatchOperations(b *testing.B) {
	items := make([]Item[int, int], 1000)
	for i := 0; i < 1000; i++ {
		items[i] = Item[int, int]{Key: i, Value: i}
	}

	b.Run("SetBatch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			om := NewOrderedMap[int, int]()
			om.SetBatch(items)
		}
	})

	b.Run("IndividualSets", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			om := NewOrderedMap[int, int]()
			for _, item := range items {
				om.Set(item.Key, item.Value)
			}
		}
	})
}

// Test for race conditions
func TestRaceConditions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping race test in short mode")
	}

	om := NewOrderedMap[int, string](WithThreadSafe[int, string](true))
	done := make(chan bool)

	// Multiple goroutines performing various operations
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				key := id*100 + j
				om.Set(key, fmt.Sprintf("goroutine-%d-value-%d", id, j))
				om.Get(key)
				om.MoveToBack(key)
				if j%10 == 0 {
					om.Delete(key - 1)
				}
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify map is still consistent
	keys := om.Keys()
	values := om.Values()
	if len(keys) != len(values) {
		t.Error("Keys and values length mismatch after concurrent operations")
	}
}

// Test equals function
func TestEquals(t *testing.T) {
	om1 := NewOrderedMap[string, int]()
	om2 := NewOrderedMap[string, int]()

	// Equal when empty
	if !om1.Equals(om2, func(a, b int) bool { return a == b }) {
		t.Error("Empty maps should be equal")
	}

	// Add same items in same order
	for _, k := range []string{"a", "b", "c"} {
		om1.Set(k, len(k))
		om2.Set(k, len(k))
	}

	if !om1.Equals(om2, func(a, b int) bool { return a == b }) {
		t.Error("Maps with same items in same order should be equal")
	}

	// Different order
	om3 := NewOrderedMap[string, int]()
	om3.Set("b", 1)
	om3.Set("a", 1)
	om3.Set("c", 1)

	if om1.Equals(om3, func(a, b int) bool { return a == b }) {
		t.Error("Maps with different order should not be equal")
	}

	// Different values
	om4 := NewOrderedMap[string, int]()
	om4.Set("a", 2)
	om4.Set("b", 2)
	om4.Set("c", 2)

	if om1.Equals(om4, func(a, b int) bool { return a == b }) {
		t.Error("Maps with different values should not be equal")
	}
}

// Test slice operations
func TestSlice(t *testing.T) {
	om := NewOrderedMap[int, string]()
	for i := 0; i < 10; i++ {
		om.Set(i, fmt.Sprintf("v%d", i))
	}

	// Test normal slice
	slice := om.Slice(2, 5)
	if len(slice) != 3 {
		t.Errorf("Slice(2,5) should return 3 items, got %d", len(slice))
	}
	if slice[0].Key != 2 || slice[2].Key != 4 {
		t.Error("Slice returned wrong items")
	}

	// Test edge cases
	slice = om.Slice(-1, 3)
	if len(slice) != 3 {
		t.Error("Slice with negative start should work")
	}

	slice = om.Slice(8, 20)
	if len(slice) != 2 {
		t.Error("Slice with end > length should work")
	}

	slice = om.Slice(5, 5)
	if slice != nil && len(slice) != 0 {
		t.Error("Empty slice should return nil or empty")
	}
}

// Test misc operations
func TestMiscOperations(t *testing.T) {
	om := NewOrderedMap[string, int]()

	// Test GetOrSet
	val, existed := om.GetOrSet("a", 1)
	if existed || val != 1 {
		t.Error("GetOrSet on new key failed")
	}

	val, existed = om.GetOrSet("a", 2)
	if !existed || val != 1 {
		t.Error("GetOrSet on existing key failed")
	}

	// Test Update
	if !om.Update("a", 10) {
		t.Error("Update on existing key should return true")
	}
	if om.Update("b", 20) {
		t.Error("Update on non-existing key should return false")
	}

	// Test SetIfAbsent
	if !om.SetIfAbsent("b", 2) {
		t.Error("SetIfAbsent on new key should return true")
	}
	if om.SetIfAbsent("b", 3) {
		t.Error("SetIfAbsent on existing key should return false")
	}

	// Test Promote/Demote
	om.Set("c", 3)
	om.Set("d", 4)
	om.Set("e", 5)

	promoted := om.Promote("a", "c", "z") // z doesn't exist
	if promoted != 2 {
		t.Errorf("Promote should return 2, got %d", promoted)
	}

	keys := om.Keys()
	// Should be: b, d, e, a, c (a and c moved to back)
	if keys[len(keys)-2] != "a" || keys[len(keys)-1] != "c" {
		t.Errorf("Promote failed: %v", keys)
	}

	demoted := om.Demote("c", "e")
	if demoted != 2 {
		t.Errorf("Demote should return 2, got %d", demoted)
	}

	// Test SwapPositions
	om.Clear()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)
	om.Set("d", 4)

	if !om.SwapPositions("a", "d") {
		t.Error("SwapPositions should return true for existing keys")
	}

	keys = om.Keys()
	if keys[0] != "d" || keys[3] != "a" {
		t.Errorf("SwapPositions failed: %v", keys)
	}

	// Test swap adjacent
	if !om.SwapPositions("b", "c") {
		t.Error("SwapPositions should handle adjacent keys")
	}
}

// Test position operations
func TestPositionOperations(t *testing.T) {
	om := NewOrderedMap[string, int]()
	for _, k := range []string{"a", "b", "c", "d", "e"} {
		om.Set(k, 1)
	}

	// Test GetPosition
	pos, ok := om.GetPosition("c")
	if !ok || pos != 2 {
		t.Errorf("GetPosition failed: expected 2, got %d", pos)
	}

	_, ok = om.GetPosition("z")
	if ok {
		t.Error("GetPosition should return false for non-existent key")
	}

	// Test MoveToPosition
	if !om.MoveToPosition("e", 1) {
		t.Error("MoveToPosition should return true")
	}

	keys := om.Keys()
	if keys[1] != "e" {
		t.Errorf("MoveToPosition failed: %v", keys)
	}

	// Test edge cases
	om.MoveToPosition("a", -1)  // Should move to front
	om.MoveToPosition("b", 100) // Should move to back

	keys = om.Keys()
	if keys[0] != "a" || keys[len(keys)-1] != "b" {
		t.Errorf("MoveToPosition edge cases failed: %v", keys)
	}
}

// Test eviction operations
func TestEvictionOperations(t *testing.T) {
	om := NewOrderedMap[int, string]()
	for i := 0; i < 10; i++ {
		om.Set(i, fmt.Sprintf("v%d", i))
	}

	// Test TrimToSize
	removed := om.TrimToSize(7)
	if removed != 3 {
		t.Errorf("TrimToSize should remove 3 items, removed %d", removed)
	}
	if om.Len() != 7 {
		t.Errorf("TrimToSize failed: expected 7 items, got %d", om.Len())
	}

	// Test EvictLRU
	evicted := om.EvictLRU(2)
	if len(evicted) != 2 {
		t.Errorf("EvictLRU should return 2 items, got %d", len(evicted))
	}
	if om.Len() != 5 {
		t.Errorf("EvictLRU failed: expected 5 items, got %d", om.Len())
	}

	// Test EvictMRU
	evicted = om.EvictMRU(2)
	if len(evicted) != 2 {
		t.Errorf("EvictMRU should return 2 items, got %d", len(evicted))
	}
	if om.Len() != 3 {
		t.Errorf("EvictMRU failed: expected 3 items, got %d", om.Len())
	}

	// Test edge cases
	evicted = om.EvictLRU(10) // More than available
	if len(evicted) != 3 {
		t.Errorf("EvictLRU should evict all remaining items: got %d", len(evicted))
	}
	if om.Len() != 0 {
		t.Error("Map should be empty after evicting all")
	}
}

// Test reorder operation
func TestReorder(t *testing.T) {
	om := NewOrderedMap[string, int]()
	for _, k := range []string{"a", "b", "c", "d", "e"} {
		om.Set(k, 1)
	}

	// Reorder subset
	reordered := om.Reorder([]string{"d", "b", "a"})
	if reordered != 3 {
		t.Errorf("Reorder should process 3 items, got %d", reordered)
	}

	keys := om.Keys()
	expected := []string{"d", "b", "a", "c", "e"}
	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Reorder failed: expected %v, got %v", expected, keys)
	}
}

// Test merge operation
func TestMerge(t *testing.T) {
	om1 := NewOrderedMap[string, int]()
	om1.Set("a", 1)
	om1.Set("b", 2)

	om2 := NewOrderedMap[string, int]()
	om2.Set("b", 20) // Duplicate key
	om2.Set("c", 3)
	om2.Set("d", 4)

	om1.Merge(om2)

	if om1.Len() != 4 {
		t.Errorf("Merge failed: expected 4 items, got %d", om1.Len())
	}

	// Check merged values
	if val, _ := om1.Get("b"); val != 20 {
		t.Error("Merge should update existing keys")
	}

	// Check order
	keys := om1.Keys()
	expected := []string{"a", "b", "c", "d"}
	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Merge order wrong: expected %v, got %v", expected, keys)
	}

	// Test self-merge
	om1.Merge(om1) // Should be no-op
	if om1.Len() != 4 {
		t.Error("Self-merge should be no-op")
	}
}

// Test First/Last operations
func TestFirstLastOperations(t *testing.T) {
	om := NewOrderedMap[int, string]()

	// Test on empty map
	if _, _, ok := om.First(); ok {
		t.Error("First on empty map should return false")
	}
	if _, _, ok := om.Last(); ok {
		t.Error("Last on empty map should return false")
	}

	// Add items
	om.Set(1, "one")
	om.Set(2, "two")
	om.Set(3, "three")

	// Test First
	if k, v, ok := om.First(); !ok || k != 1 || v != "one" {
		t.Errorf("First failed: got %d, %s, %v", k, v, ok)
	}

	// Test Last
	if k, v, ok := om.Last(); !ok || k != 3 || v != "three" {
		t.Errorf("Last failed: got %d, %s, %v", k, v, ok)
	}
}

// Test pop operations
func TestPopOperations(t *testing.T) {
	om := NewOrderedMap[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	// Test Pop
	if val, ok := om.Pop("b"); !ok || val != 2 {
		t.Errorf("Pop failed: got %d, %v", val, ok)
	}
	if om.Contains("b") {
		t.Error("Pop should remove the key")
	}

	// Test PopFirst
	if k, v, ok := om.PopFirst(); !ok || k != "a" || v != 1 {
		t.Errorf("PopFirst failed: got %s, %d, %v", k, v, ok)
	}

	// Test PopLast
	if k, v, ok := om.PopLast(); !ok || k != "c" || v != 3 {
		t.Errorf("PopLast failed: got %s, %d, %v", k, v, ok)
	}

	// Map should be empty now
	if om.Len() != 0 {
		t.Error("Map should be empty after popping all items")
	}
}

// Test for memory leaks
func TestMemoryLeaks(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory leak test in short mode")
	}

	om := NewOrderedMap[int, []byte]()

	// Add and remove many large items
	for i := 0; i < 10000; i++ {
		// Create 1KB value
		value := make([]byte, 1024)
		om.Set(i, value)

		// Delete older items to maintain size
		if i > 100 {
			om.Delete(i - 100)
		}
	}

	// Clear should release all memory
	om.Clear()

	// Force GC
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	// Map should be completely empty
	if om.Len() != 0 {
		t.Error("Map should be empty after clear")
	}
}

// Test concurrent modifications during iteration
func TestConcurrentModificationDuringIteration(t *testing.T) {
	om := NewOrderedMap[int, string](WithThreadSafe[int, string](true))

	// Add initial items
	for i := 0; i < 100; i++ {
		om.Set(i, fmt.Sprintf("v%d", i))
	}

	// This should not panic
	done := make(chan bool)
	go func() {
		om.ForEach(func(k int, v string) {
			time.Sleep(time.Microsecond) // Simulate work
		})
		done <- true
	}()

	// Modify while iterating
	go func() {
		for i := 100; i < 200; i++ {
			om.Set(i, fmt.Sprintf("v%d", i))
			if i%10 == 0 {
				om.Delete(i - 50)
			}
		}
	}()

	select {
	case <-done:
		// Success - iteration completed without panic
	case <-time.After(5 * time.Second):
		t.Error("Iteration timed out - possible deadlock")
	}
}

// Example usage test
func TestExampleUsage(t *testing.T) {
	// Simulate database deduplication use case
	type Customer struct {
		ID    int
		Name  string
		Email string
	}

	customers := NewOrderedMap[int, Customer]()

	// Simulate rows from database with duplicates
	rows := []Customer{
		{ID: 1, Name: "Alice", Email: "alice@example.com"},
		{ID: 2, Name: "Bob", Email: "bob@example.com"},
		{ID: 1, Name: "Alice Updated", Email: "alice.new@example.com"}, // Duplicate ID
		{ID: 3, Name: "Charlie", Email: "charlie@example.com"},
		{ID: 2, Name: "Bob Updated", Email: "bob.new@example.com"}, // Duplicate ID
	}

	// Process rows - automatically handles deduplication
	for _, customer := range rows {
		customers.Set(customer.ID, customer)
	}

	// Get unique customers
	uniqueCustomers := customers.Values()

	if len(uniqueCustomers) != 3 {
		t.Errorf("Expected 3 unique customers, got %d", len(uniqueCustomers))
	}

	// Verify last update wins
	if c, _ := customers.Get(1); c.Email != "alice.new@example.com" {
		t.Error("Deduplication should keep last value")
	}
}

// Test custom types
func TestCustomTypes(t *testing.T) {
	type CustomKey struct {
		ID   int
		Type string
	}

	om := NewOrderedMap[CustomKey, string]()

	key1 := CustomKey{ID: 1, Type: "A"}
	key2 := CustomKey{ID: 2, Type: "B"}
	key3 := CustomKey{ID: 1, Type: "A"} // Same as key1

	om.Set(key1, "value1")
	om.Set(key2, "value2")
	om.Set(key3, "value3") // Should update key1

	if om.Len() != 2 {
		t.Errorf("Custom key deduplication failed: expected 2, got %d", om.Len())
	}

	if val, _ := om.Get(key1); val != "value3" {
		t.Error("Custom key update failed")
	}
}

// Test stats
func TestStats(t *testing.T) {
	config := Config[string, int]{
		Capacity:     10,
		ThreadSafe:   true,
		CapacityMode: EvictOldest,
	}
	om := NewOrderedMapWithConfig(config)

	stats := om.Stats()
	if stats.Capacity != 10 {
		t.Errorf("Stats capacity wrong: expected 10, got %d", stats.Capacity)
	}
	if !stats.ThreadSafe {
		t.Error("Stats thread safe should be true")
	}
	if stats.CapacityMode != EvictOldest {
		t.Error("Stats capacity mode wrong")
	}

	// Add items and check size
	for i := 0; i < 5; i++ {
		om.Set(fmt.Sprintf("k%d", i), i)
	}

	stats = om.Stats()
	if stats.Size != 5 {
		t.Errorf("Stats size wrong: expected 5, got %d", stats.Size)
	}
}

// Test panic recovery
func TestPanicRecovery(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Operation should not panic: %v", r)
		}
	}()

	om := NewOrderedMap[string, int]()

	// These operations should not panic even with edge cases
	om.Delete("")
	om.MoveToFront("")
	om.GetPosition("")
	om.Pop("")
	om.Update("", 0)

	// Empty key is valid
	om.Set("", 100)
	if val, ok := om.Get(""); !ok || val != 100 {
		t.Error("Empty key should be valid")
	}
}
