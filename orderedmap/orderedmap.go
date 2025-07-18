package orderedmap

import (
	"sync"
)

// CapacityMode defines eviction behavior when capacity is reached
type CapacityMode int

const (
	// EvictOldest removes least recently used items (default)
	EvictOldest CapacityMode = iota
	// EvictNewest removes most recently used items
	EvictNewest
	// RejectNew prevents new items from being added
	RejectNew
)

// Config extended configuration for OrderedMap
type Config[K comparable, V any] struct {
	Capacity     int
	ThreadSafe   bool
	InitialSize  int
	CapacityMode CapacityMode
	OnEvict      func(key K, value V) // Callback when item is evicted
}

// Node represents a node in the doubly-linked list
type Node[K comparable, V any] struct {
	key   K
	value V
	prev  *Node[K, V]
	next  *Node[K, V]
}

// Item represents a key-value pair
type Item[K comparable, V any] struct {
	Key   K
	Value V
}

// OrderedMap maintains insertion order while providing O(1) key access
type OrderedMap[K comparable, V any] struct {
	// Core data structures
	data map[K]*Node[K, V]
	head *Node[K, V] // sentinel head
	tail *Node[K, V] // sentinel tail
	size int

	// Configuration
	capacity     int
	threadSafe   bool // Add this field
	capacityMode CapacityMode
	onEvict      func(key K, value V)

	// Synchronization
	mu *sync.RWMutex

	// Object pool for nodes (performance optimization)
	nodePool *sync.Pool
}

// Option is a functional option for configuring OrderedMap
type Option[K comparable, V any] func(*OrderedMap[K, V])

// WithCapacity sets the maximum capacity with LRU eviction
func WithCapacity[K comparable, V any](capacity int) Option[K, V] {
	return func(om *OrderedMap[K, V]) {
		if capacity > 0 {
			om.capacity = capacity
		}
	}
}

// WithThreadSafe enables thread-safe operations
func WithThreadSafe[K comparable, V any](threadSafe bool) Option[K, V] {
	return func(om *OrderedMap[K, V]) {
		om.threadSafe = threadSafe
		if threadSafe {
			om.mu = &sync.RWMutex{}
		}
	}
}

// WithInitialCapacity pre-allocates the internal map
func WithInitialCapacity[K comparable, V any](initialCap int) Option[K, V] {
	return func(om *OrderedMap[K, V]) {
		if initialCap > 0 {
			om.data = make(map[K]*Node[K, V], initialCap)
		}
	}
}

// NewOrderedMap creates a new ordered map with options
func NewOrderedMap[K comparable, V any](opts ...Option[K, V]) *OrderedMap[K, V] {
	om := &OrderedMap[K, V]{
		data:     make(map[K]*Node[K, V]),
		head:     &Node[K, V]{}, // sentinel
		tail:     &Node[K, V]{}, // sentinel
		capacity: 0,             // no limit by default
		nodePool: &sync.Pool{
			New: func() interface{} {
				return &Node[K, V]{}
			},
		},
	}

	// Initialize sentinels
	om.head.next = om.tail
	om.tail.prev = om.head

	// Apply options
	for _, opt := range opts {
		opt(om)
	}

	return om
}

// NewOrderedMapWithConfig creates an ordered map with detailed configuration
func NewOrderedMapWithConfig[K comparable, V any](config Config[K, V]) *OrderedMap[K, V] {
	om := &OrderedMap[K, V]{
		data:     make(map[K]*Node[K, V], config.InitialSize),
		head:     &Node[K, V]{},
		tail:     &Node[K, V]{},
		capacity: config.Capacity,
		nodePool: &sync.Pool{
			New: func() interface{} {
				return &Node[K, V]{}
			},
		},
	}

	// Initialize sentinels
	om.head.next = om.tail
	om.tail.prev = om.head

	// Set thread safety
	if config.ThreadSafe {
		om.mu = &sync.RWMutex{}
	}

	// Store config for advanced features
	om.capacityMode = config.CapacityMode
	om.onEvict = config.OnEvict

	return om
}

// Get retrieves a value by key - O(1)
func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	om.rlock()
	defer om.runlock()

	node, exists := om.data[key]
	if !exists {
		var zero V
		return zero, false
	}

	return node.value, true
}

// Set adds or updates a key-value pair - O(1) average
func (om *OrderedMap[K, V]) Set(key K, value V) {
	om.lock()
	defer om.unlock()

	// Check if key already exists
	if node, exists := om.data[key]; exists {
		// Update existing value
		node.value = value
		// Move to back to maintain LRU order (most recently used)
		om.moveToBack(node)
		return
	}

	// Create new node
	node := om.newNode(key, value)
	om.data[key] = node
	om.addToBack(node)
	om.size++

	// Check capacity and evict if necessary
	if om.capacity > 0 && om.size > om.capacity {
		om.evictOldest() // This will now call the callback
	}
}

// Delete removes a key-value pair - O(1)
func (om *OrderedMap[K, V]) Delete(key K) bool {
	om.lock()
	defer om.unlock()

	node, exists := om.data[key]
	if !exists {
		return false
	}

	om.removeNode(node)
	delete(om.data, key)
	om.releaseNode(node)
	om.size--

	return true
}

// Clear removes all items
func (om *OrderedMap[K, V]) Clear() {
	om.lock()
	defer om.unlock()

	// Release all nodes to pool
	current := om.head.next
	for current != om.tail {
		next := current.next
		om.releaseNode(current)
		current = next
	}

	// Reset data structures
	om.data = make(map[K]*Node[K, V])
	om.head.next = om.tail
	om.tail.prev = om.head
	om.size = 0
}

// Contains checks if a key exists - O(1)
func (om *OrderedMap[K, V]) Contains(key K) bool {
	om.rlock()
	defer om.runlock()

	_, exists := om.data[key]
	return exists
}

// GetOrSet returns existing value or sets and returns new value - O(1)
func (om *OrderedMap[K, V]) GetOrSet(key K, value V) (V, bool) {
	om.lock()
	defer om.unlock()

	if node, exists := om.data[key]; exists {
		// Move to back for LRU
		om.moveToBack(node)
		return node.value, true
	}

	// Create new node
	node := om.newNode(key, value)
	om.data[key] = node
	om.addToBack(node)
	om.size++

	// Check capacity and evict if necessary
	if om.capacity > 0 && om.size > om.capacity {
		om.evictOldest()
	}

	return value, false
}

// Update updates value only if key exists - O(1)
func (om *OrderedMap[K, V]) Update(key K, value V) bool {
	om.lock()
	defer om.unlock()

	node, exists := om.data[key]
	if !exists {
		return false
	}

	node.value = value
	// Move to back for LRU
	om.moveToBack(node)
	return true
}

// Pop removes and returns the value for a key - O(1)
func (om *OrderedMap[K, V]) Pop(key K) (V, bool) {
	om.lock()
	defer om.unlock()

	node, exists := om.data[key]
	if !exists {
		var zero V
		return zero, false
	}

	value := node.value
	om.removeNode(node)
	delete(om.data, key)
	om.releaseNode(node)
	om.size--

	return value, true
}

// PopFirst removes and returns the first (oldest) item
func (om *OrderedMap[K, V]) PopFirst() (K, V, bool) {
	om.lock()
	defer om.unlock()

	if om.size == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	first := om.head.next
	key := first.key
	value := first.value

	om.removeNode(first)
	delete(om.data, key)
	om.releaseNode(first)
	om.size--

	return key, value, true
}

// PopLast removes and returns the last (newest) item
func (om *OrderedMap[K, V]) PopLast() (K, V, bool) {
	om.lock()
	defer om.unlock()

	if om.size == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	last := om.tail.prev
	key := last.key
	value := last.value

	om.removeNode(last)
	delete(om.data, key)
	om.releaseNode(last)
	om.size--

	return key, value, true
}

// First returns the first (oldest) item without removing it
func (om *OrderedMap[K, V]) First() (K, V, bool) {
	om.rlock()
	defer om.runlock()

	if om.size == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	first := om.head.next
	return first.key, first.value, true
}

// Last returns the last (newest) item without removing it
func (om *OrderedMap[K, V]) Last() (K, V, bool) {
	om.rlock()
	defer om.runlock()

	if om.size == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	last := om.tail.prev
	return last.key, last.value, true
}

// SetIfAbsent sets value only if key doesn't exist - O(1)
func (om *OrderedMap[K, V]) SetIfAbsent(key K, value V) bool {
	om.lock()
	defer om.unlock()

	if _, exists := om.data[key]; exists {
		return false
	}

	// Create new node
	node := om.newNode(key, value)
	om.data[key] = node
	om.addToBack(node)
	om.size++

	// Check capacity and evict if necessary
	if om.capacity > 0 && om.size > om.capacity {
		om.evictOldest()
	}

	return true
}

// Merge adds all items from another ordered map
func (om *OrderedMap[K, V]) Merge(other *OrderedMap[K, V]) {
	if om == other {
		return // avoid self-merge
	}

	other.rlock()
	defer other.runlock()

	// Iterate through other map in order
	current := other.head.next
	for current != other.tail {
		om.Set(current.key, current.value)
		current = current.next
	}
}

// Keys returns all keys in order - O(n)
func (om *OrderedMap[K, V]) Keys() []K {
	om.rlock()
	defer om.runlock()

	if om.size == 0 {
		return nil
	}

	keys := make([]K, 0, om.size)
	current := om.head.next
	for current != om.tail {
		keys = append(keys, current.key)
		current = current.next
	}

	return keys
}

// Values returns all values in order - O(n)
func (om *OrderedMap[K, V]) Values() []V {
	om.rlock()
	defer om.runlock()

	if om.size == 0 {
		return nil
	}

	values := make([]V, 0, om.size)
	current := om.head.next
	for current != om.tail {
		values = append(values, current.value)
		current = current.next
	}

	return values
}

// Items returns all key-value pairs in order - O(n)
func (om *OrderedMap[K, V]) Items() []Item[K, V] {
	om.rlock()
	defer om.runlock()

	if om.size == 0 {
		return nil
	}

	items := make([]Item[K, V], 0, om.size)
	current := om.head.next
	for current != om.tail {
		items = append(items, Item[K, V]{
			Key:   current.key,
			Value: current.value,
		})
		current = current.next
	}

	return items
}

// ForEach iterates over all items in order
func (om *OrderedMap[K, V]) ForEach(fn func(key K, value V)) {
	om.rlock()
	defer om.runlock()

	current := om.head.next
	for current != om.tail {
		fn(current.key, current.value)
		current = current.next
	}
}

// ForEachReverse iterates over all items in reverse order
func (om *OrderedMap[K, V]) ForEachReverse(fn func(key K, value V)) {
	om.rlock()
	defer om.runlock()

	current := om.tail.prev
	for current != om.head {
		fn(current.key, current.value)
		current = current.prev
	}
}

// Map applies a function to all values and returns a new slice
func (om *OrderedMap[K, V]) Map(fn func(key K, value V) interface{}) []interface{} {
	om.rlock()
	defer om.runlock()

	if om.size == 0 {
		return nil
	}

	result := make([]interface{}, 0, om.size)
	current := om.head.next
	for current != om.tail {
		result = append(result, fn(current.key, current.value))
		current = current.next
	}

	return result
}

// Filter returns items that match the predicate
func (om *OrderedMap[K, V]) Filter(predicate func(key K, value V) bool) []Item[K, V] {
	om.rlock()
	defer om.runlock()

	items := make([]Item[K, V], 0)
	current := om.head.next
	for current != om.tail {
		if predicate(current.key, current.value) {
			items = append(items, Item[K, V]{
				Key:   current.key,
				Value: current.value,
			})
		}
		current = current.next
	}

	return items
}

// Iterator provides an iterator interface
type Iterator[K comparable, V any] struct {
	om      *OrderedMap[K, V]
	current *Node[K, V]
	reverse bool
}

// NewIterator creates a forward iterator
func (om *OrderedMap[K, V]) NewIterator() *Iterator[K, V] {
	om.rlock()
	defer om.runlock()

	return &Iterator[K, V]{
		om:      om,
		current: om.head,
		reverse: false,
	}
}

// NewReverseIterator creates a reverse iterator
func (om *OrderedMap[K, V]) NewReverseIterator() *Iterator[K, V] {
	om.rlock()
	defer om.runlock()

	return &Iterator[K, V]{
		om:      om,
		current: om.tail,
		reverse: true,
	}
}

// HasNext checks if there are more items
func (it *Iterator[K, V]) HasNext() bool {
	it.om.rlock()
	defer it.om.runlock()

	if it.reverse {
		return it.current.prev != it.om.head
	}
	return it.current.next != it.om.tail
}

// Next returns the next item
func (it *Iterator[K, V]) Next() (K, V, bool) {
	it.om.rlock()
	defer it.om.runlock()

	if it.reverse {
		if it.current.prev == it.om.head {
			var zeroK K
			var zeroV V
			return zeroK, zeroV, false
		}
		it.current = it.current.prev
	} else {
		if it.current.next == it.om.tail {
			var zeroK K
			var zeroV V
			return zeroK, zeroV, false
		}
		it.current = it.current.next
	}

	return it.current.key, it.current.value, true
}

// Reset resets the iterator to the beginning
func (it *Iterator[K, V]) Reset() {
	it.om.rlock()
	defer it.om.runlock()

	if it.reverse {
		it.current = it.om.tail
	} else {
		it.current = it.om.head
	}
}

// Range provides a Go 1.23+ style iterator (if needed)
// This allows: for k, v := range om.Range() { ... }
func (om *OrderedMap[K, V]) Range() func(func(K, V) bool) {
	return func(yield func(K, V) bool) {
		om.rlock()
		defer om.runlock()

		current := om.head.next
		for current != om.tail {
			if !yield(current.key, current.value) {
				break
			}
			current = current.next
		}
	}
}

// Slice returns a slice of values for keys in the given range [start:end)
func (om *OrderedMap[K, V]) Slice(start, end int) []Item[K, V] {
	om.rlock()
	defer om.runlock()

	if start < 0 {
		start = 0
	}
	if end > om.size {
		end = om.size
	}
	if start >= end {
		return nil
	}

	items := make([]Item[K, V], 0, end-start)
	current := om.head.next

	// Skip to start
	for i := 0; i < start && current != om.tail; i++ {
		current = current.next
	}

	// Collect items
	for i := start; i < end && current != om.tail; i++ {
		items = append(items, Item[K, V]{
			Key:   current.key,
			Value: current.value,
		})
		current = current.next
	}

	return items
}

// Find returns the first item matching the predicate
func (om *OrderedMap[K, V]) Find(predicate func(key K, value V) bool) (K, V, bool) {
	om.rlock()
	defer om.runlock()

	current := om.head.next
	for current != om.tail {
		if predicate(current.key, current.value) {
			return current.key, current.value, true
		}
		current = current.next
	}

	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// Any returns true if any item matches the predicate
func (om *OrderedMap[K, V]) Any(predicate func(key K, value V) bool) bool {
	om.rlock()
	defer om.runlock()

	current := om.head.next
	for current != om.tail {
		if predicate(current.key, current.value) {
			return true
		}
		current = current.next
	}

	return false
}

// All returns true if all items match the predicate
func (om *OrderedMap[K, V]) All(predicate func(key K, value V) bool) bool {
	om.rlock()
	defer om.runlock()

	if om.size == 0 {
		return true
	}

	current := om.head.next
	for current != om.tail {
		if !predicate(current.key, current.value) {
			return false
		}
		current = current.next
	}

	return true
}

// MoveToFront moves an existing key to the front (oldest position) - O(1)
func (om *OrderedMap[K, V]) MoveToFront(key K) bool {
	om.lock()
	defer om.unlock()

	node, exists := om.data[key]
	if !exists {
		return false
	}

	om.moveToFront(node)
	return true
}

// MoveToBack moves an existing key to the back (newest position) - O(1)
func (om *OrderedMap[K, V]) MoveToBack(key K) bool {
	om.lock()
	defer om.unlock()

	node, exists := om.data[key]
	if !exists {
		return false
	}

	om.moveToBack(node)
	return true
}

// MoveBefore moves key before mark - O(1)
func (om *OrderedMap[K, V]) MoveBefore(key, mark K) bool {
	om.lock()
	defer om.unlock()

	if key == mark {
		return true // no-op
	}

	keyNode, keyExists := om.data[key]
	markNode, markExists := om.data[mark]

	if !keyExists || !markExists {
		return false
	}

	// Remove key node from its current position
	om.removeNode(keyNode)

	// Insert before mark
	keyNode.prev = markNode.prev
	keyNode.next = markNode
	markNode.prev.next = keyNode
	markNode.prev = keyNode

	return true
}

// MoveAfter moves key after mark - O(1)
func (om *OrderedMap[K, V]) MoveAfter(key, mark K) bool {
	om.lock()
	defer om.unlock()

	if key == mark {
		return true // no-op
	}

	keyNode, keyExists := om.data[key]
	markNode, markExists := om.data[mark]

	if !keyExists || !markExists {
		return false
	}

	// Remove key node from its current position
	om.removeNode(keyNode)

	// Insert after mark
	keyNode.prev = markNode
	keyNode.next = markNode.next
	markNode.next.prev = keyNode
	markNode.next = keyNode

	return true
}

// Touch moves an existing key to the back (marks as recently used) - O(1)
// Alias for MoveToBack for LRU semantics
func (om *OrderedMap[K, V]) Touch(key K) bool {
	return om.MoveToBack(key)
}

// GetAndTouch gets value and moves to back (LRU update) - O(1)
func (om *OrderedMap[K, V]) GetAndTouch(key K) (V, bool) {
	om.lock()
	defer om.unlock()

	node, exists := om.data[key]
	if !exists {
		var zero V
		return zero, false
	}

	om.moveToBack(node)
	return node.value, true
}

// Promote moves multiple keys to back in order - O(k) where k is number of keys
func (om *OrderedMap[K, V]) Promote(keys ...K) int {
	om.lock()
	defer om.unlock()

	promoted := 0
	for _, key := range keys {
		if node, exists := om.data[key]; exists {
			om.moveToBack(node)
			promoted++
		}
	}

	return promoted
}

// Demote moves multiple keys to front in order - O(k) where k is number of keys
func (om *OrderedMap[K, V]) Demote(keys ...K) int {
	om.lock()
	defer om.unlock()

	demoted := 0
	// Process in reverse order to maintain relative ordering
	for i := len(keys) - 1; i >= 0; i-- {
		if node, exists := om.data[keys[i]]; exists {
			om.moveToFront(node)
			demoted++
		}
	}

	return demoted
}

// SwapPositions swaps the positions of two keys - O(1)
func (om *OrderedMap[K, V]) SwapPositions(key1, key2 K) bool {
	om.lock()
	defer om.unlock()

	if key1 == key2 {
		return true // no-op
	}

	node1, exists1 := om.data[key1]
	node2, exists2 := om.data[key2]

	if !exists1 || !exists2 {
		return false
	}

	// Special case: adjacent nodes
	if node1.next == node2 {
		// node1 is directly before node2
		om.removeNode(node2)
		node2.prev = node1.prev
		node2.next = node1
		node1.prev.next = node2
		node1.prev = node2
	} else if node2.next == node1 {
		// node2 is directly before node1
		om.removeNode(node1)
		node1.prev = node2.prev
		node1.next = node2
		node2.prev.next = node1
		node2.prev = node1
	} else {
		// Non-adjacent nodes - save positions
		prev1, next1 := node1.prev, node1.next
		prev2, next2 := node2.prev, node2.next

		// Remove both nodes
		om.removeNode(node1)
		om.removeNode(node2)

		// Insert node1 at node2's position
		node1.prev = prev2
		node1.next = next2
		prev2.next = node1
		next2.prev = node1

		// Insert node2 at node1's position
		node2.prev = prev1
		node2.next = next1
		prev1.next = node2
		next1.prev = node2
	}

	return true
}

// GetPosition returns the position of a key (0-based) - O(n)
func (om *OrderedMap[K, V]) GetPosition(key K) (int, bool) {
	om.rlock()
	defer om.runlock()

	if _, exists := om.data[key]; !exists {
		return -1, false
	}

	position := 0
	current := om.head.next
	for current != om.tail {
		if current.key == key {
			return position, true
		}
		position++
		current = current.next
	}

	return -1, false
}

// MoveToPosition moves a key to a specific position - O(n)
func (om *OrderedMap[K, V]) MoveToPosition(key K, position int) bool {
	om.lock()
	defer om.unlock()

	node, exists := om.data[key]
	if !exists {
		return false
	}

	if position < 0 {
		position = 0
	}
	if position >= om.size {
		om.moveToBack(node)
		return true
	}

	// Find target position
	current := om.head.next
	for i := 0; i < position && current != om.tail; i++ {
		current = current.next
	}

	if current == node {
		return true // already at position
	}

	// Remove from current position
	om.removeNode(node)

	// Insert at target position
	node.prev = current.prev
	node.next = current
	current.prev.next = node
	current.prev = node

	return true
}

// TrimToSize removes oldest items to reach target size
func (om *OrderedMap[K, V]) TrimToSize(targetSize int) int {
	om.lock()
	defer om.unlock()

	if targetSize < 0 {
		targetSize = 0
	}

	removed := 0
	for om.size > targetSize {
		om.evictOldest()
		removed++
	}

	return removed
}

// EvictLRU removes the least recently used (oldest) items
func (om *OrderedMap[K, V]) EvictLRU(count int) []Item[K, V] {
	om.lock()
	defer om.unlock()

	if count <= 0 || om.size == 0 {
		return nil
	}

	if count > om.size {
		count = om.size
	}

	evicted := make([]Item[K, V], 0, count)

	for i := 0; i < count && om.size > 0; i++ {
		oldest := om.head.next
		evicted = append(evicted, Item[K, V]{
			Key:   oldest.key,
			Value: oldest.value,
		})

		om.removeNode(oldest)
		delete(om.data, oldest.key)
		om.releaseNode(oldest)
		om.size--
	}

	return evicted
}

// EvictMRU removes the most recently used (newest) items
func (om *OrderedMap[K, V]) EvictMRU(count int) []Item[K, V] {
	om.lock()
	defer om.unlock()

	if count <= 0 || om.size == 0 {
		return nil
	}

	if count > om.size {
		count = om.size
	}

	evicted := make([]Item[K, V], 0, count)

	for i := 0; i < count && om.size > 0; i++ {
		newest := om.tail.prev
		evicted = append(evicted, Item[K, V]{
			Key:   newest.key,
			Value: newest.value,
		})

		om.removeNode(newest)
		delete(om.data, newest.key)
		om.releaseNode(newest)
		om.size--
	}

	return evicted
}

// Reorder rearranges items to match the given key order
func (om *OrderedMap[K, V]) Reorder(keys []K) int {
	om.lock()
	defer om.unlock()

	reordered := 0

	// Process keys in reverse order to maintain desired order
	for i := len(keys) - 1; i >= 0; i-- {
		if node, exists := om.data[keys[i]]; exists {
			om.moveToFront(node)
			reordered++
		}
	}

	return reordered
}

// SetWithMode adds or updates a key-value pair with capacity mode handling
func (om *OrderedMap[K, V]) SetWithMode(key K, value V) bool {
	om.lock()
	defer om.unlock()

	// Check if key already exists
	if node, exists := om.data[key]; exists {
		node.value = value
		om.moveToBack(node)
		return true
	}

	// Check capacity before adding
	if om.capacity > 0 && om.size >= om.capacity {
		switch om.capacityMode {
		case RejectNew:
			return false
		case EvictNewest:
			om.evictNewest()
		default: // EvictOldest
			om.evictOldest()
		}
	}

	// Create new node
	node := om.newNode(key, value)
	om.data[key] = node
	om.addToBack(node)
	om.size++

	return true
}

// evictNewest removes the most recently used item
func (om *OrderedMap[K, V]) evictNewest() {
	if om.size == 0 {
		return
	}

	newest := om.tail.prev

	// Call eviction callback if set
	if om.onEvict != nil {
		om.onEvict(newest.key, newest.value)
	}

	om.removeNode(newest)
	delete(om.data, newest.key)
	om.releaseNode(newest)
	om.size--
}

// Override evictOldest to support callback
func (om *OrderedMap[K, V]) evictOldestWithCallback() {
	om.evictOldest() // Since evictOldest now handles the callback
}

// Stats provides runtime statistics
type Stats struct {
	Size         int
	Capacity     int
	CapacityMode CapacityMode
	ThreadSafe   bool
}

// Stats returns current statistics
func (om *OrderedMap[K, V]) Stats() Stats {
	om.rlock()
	defer om.runlock()

	return Stats{
		Size:         om.size,
		Capacity:     om.capacity,
		CapacityMode: om.capacityMode,
		ThreadSafe:   om.mu != nil,
	}
}

// Reserve pre-allocates space for n items
func (om *OrderedMap[K, V]) Reserve(n int) {
	om.lock()
	defer om.unlock()

	if n <= len(om.data) {
		return
	}

	// Create new map with larger capacity
	newData := make(map[K]*Node[K, V], n)

	// Copy existing entries
	for k, v := range om.data {
		newData[k] = v
	}

	om.data = newData
}

// Compact removes excess capacity
func (om *OrderedMap[K, V]) Compact() {
	om.lock()
	defer om.unlock()

	if len(om.data) == om.size {
		return
	}

	// Create new map with exact size
	newData := make(map[K]*Node[K, V], om.size)

	// Copy existing entries
	for k, v := range om.data {
		newData[k] = v
	}

	om.data = newData
}

// Clone creates a deep copy of the ordered map
func (om *OrderedMap[K, V]) Clone() *OrderedMap[K, V] {
	om.rlock()
	defer om.runlock()

	// Create new map with same configuration
	clone := &OrderedMap[K, V]{
		data:         make(map[K]*Node[K, V], len(om.data)),
		head:         &Node[K, V]{},
		tail:         &Node[K, V]{},
		size:         0,
		capacity:     om.capacity,
		capacityMode: om.capacityMode,
		threadSafe:   om.threadSafe, // Use the actual field
		onEvict:      om.onEvict,    // Copy the callback too
		nodePool: &sync.Pool{
			New: func() interface{} {
				return &Node[K, V]{}
			},
		},
	}

	// Initialize sentinels
	clone.head.next = clone.tail
	clone.tail.prev = clone.head

	// Set thread safety
	if clone.threadSafe {
		clone.mu = &sync.RWMutex{}
	}

	// Copy all items in order
	current := om.head.next
	for current != om.tail {
		node := clone.newNode(current.key, current.value)
		clone.data[current.key] = node
		clone.addToBack(node)
		clone.size++
		current = current.next
	}

	return clone
}

// Equals checks if two ordered maps are equal (same items in same order)
func (om *OrderedMap[K, V]) Equals(other *OrderedMap[K, V], valueEqual func(V, V) bool) bool {
	if om == other {
		return true
	}

	om.rlock()
	other.rlock()
	defer om.runlock()
	defer other.runlock()

	if om.size != other.size {
		return false
	}

	current1 := om.head.next
	current2 := other.head.next

	for current1 != om.tail && current2 != other.tail {
		if current1.key != current2.key {
			return false
		}

		if valueEqual != nil {
			if !valueEqual(current1.value, current2.value) {
				return false
			}
		}

		current1 = current1.next
		current2 = current2.next
	}

	return true
}

// LoadFactor returns the load factor of the internal map
func (om *OrderedMap[K, V]) LoadFactor() float64 {
	om.rlock()
	defer om.runlock()

	if len(om.data) == 0 {
		return 0
	}

	return float64(om.size) / float64(len(om.data))
}

// Optimized batch operations

// SetBatch adds multiple items efficiently
func (om *OrderedMap[K, V]) SetBatch(items []Item[K, V]) {
	om.lock()
	defer om.unlock()

	// Pre-allocate if needed
	expectedSize := om.size + len(items)
	if expectedSize > len(om.data) {
		newData := make(map[K]*Node[K, V], expectedSize)
		for k, v := range om.data {
			newData[k] = v
		}
		om.data = newData
	}

	// Add all items
	for _, item := range items {
		if node, exists := om.data[item.Key]; exists {
			node.value = item.Value
			om.moveToBack(node)
		} else {
			node := om.newNode(item.Key, item.Value)
			om.data[item.Key] = node
			om.addToBack(node)
			om.size++

			// Handle capacity
			if om.capacity > 0 && om.size > om.capacity {
				if om.capacityMode == RejectNew {
					// Remove the just-added item
					om.removeNode(node)
					delete(om.data, item.Key)
					om.releaseNode(node)
					om.size--
					break
				} else if om.capacityMode == EvictNewest {
					om.evictNewest()
				} else {
					om.evictOldest() // This now properly calls the callback
				}
			}
		}
	}
}

// DeleteBatch removes multiple items efficiently
func (om *OrderedMap[K, V]) DeleteBatch(keys []K) int {
	om.lock()
	defer om.unlock()

	deleted := 0
	for _, key := range keys {
		if node, exists := om.data[key]; exists {
			om.removeNode(node)
			delete(om.data, key)
			om.releaseNode(node)
			om.size--
			deleted++
		}
	}

	return deleted
}

// GetBatch retrieves multiple values efficiently
func (om *OrderedMap[K, V]) GetBatch(keys []K) map[K]V {
	om.rlock()
	defer om.runlock()

	result := make(map[K]V, len(keys))
	for _, key := range keys {
		if node, exists := om.data[key]; exists {
			result[key] = node.value
		}
	}

	return result
}

// Len returns the number of items
func (om *OrderedMap[K, V]) Len() int {
	om.rlock()
	defer om.runlock()
	return om.size
}

// IsEmpty returns true if the map is empty
func (om *OrderedMap[K, V]) IsEmpty() bool {
	om.rlock()
	defer om.runlock()
	return om.size == 0
}

// Capacity returns the capacity limit (0 means unlimited)
func (om *OrderedMap[K, V]) Capacity() int {
	om.rlock()
	defer om.runlock()
	return om.capacity
}

// SetCapacity updates the capacity limit
func (om *OrderedMap[K, V]) SetCapacity(capacity int) {
	om.lock()
	defer om.unlock()

	om.capacity = capacity

	// Evict items if necessary
	if om.capacity > 0 {
		for om.size > om.capacity {
			om.evictOldest()
		}
	}
}

// Internal helper methods for thread-safety
func (om *OrderedMap[K, V]) lock() {
	if om.mu != nil {
		om.mu.Lock()
	}
}

func (om *OrderedMap[K, V]) unlock() {
	if om.mu != nil {
		om.mu.Unlock()
	}
}

func (om *OrderedMap[K, V]) rlock() {
	if om.mu != nil {
		om.mu.RLock()
	}
}

func (om *OrderedMap[K, V]) runlock() {
	if om.mu != nil {
		om.mu.RUnlock()
	}
}

// Internal helper methods for node management
func (om *OrderedMap[K, V]) newNode(key K, value V) *Node[K, V] {
	node := om.nodePool.Get().(*Node[K, V])
	node.key = key
	node.value = value
	node.prev = nil
	node.next = nil
	return node
}

func (om *OrderedMap[K, V]) releaseNode(node *Node[K, V]) {
	// Clear references to help GC
	var zeroK K
	var zeroV V
	node.key = zeroK
	node.value = zeroV
	node.prev = nil
	node.next = nil
	om.nodePool.Put(node)
}

// Internal linked list operations
func (om *OrderedMap[K, V]) addToBack(node *Node[K, V]) {
	node.prev = om.tail.prev
	node.next = om.tail
	om.tail.prev.next = node
	om.tail.prev = node
}

func (om *OrderedMap[K, V]) removeNode(node *Node[K, V]) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (om *OrderedMap[K, V]) moveToBack(node *Node[K, V]) {
	om.removeNode(node)
	om.addToBack(node)
}

func (om *OrderedMap[K, V]) moveToFront(node *Node[K, V]) {
	om.removeNode(node)
	node.prev = om.head
	node.next = om.head.next
	om.head.next.prev = node
	om.head.next = node
}

// evictOldest removes the least recently used item (from front)
func (om *OrderedMap[K, V]) evictOldest() {
	if om.size == 0 {
		return
	}

	oldest := om.head.next

	// Call eviction callback if set
	if om.onEvict != nil {
		om.onEvict(oldest.key, oldest.value)
	}

	om.removeNode(oldest)
	delete(om.data, oldest.key)
	om.releaseNode(oldest)
	om.size--
}
