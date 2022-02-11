package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testComponent struct {
	a int32
	b int32
}

func fillStorage(storage EntityStorage, count int) {
	for n := 0; n < count; n++ {
		storage.Add(EntityId(n), &testComponent{a: int32(n), b: int32(n + 1)})
	}
}

func testStorageAdd(t *testing.T, storage EntityStorage) {
	idSet := map[EntityId]struct{}{}
	for n := 0; n < 100000; n++ {
		storage.Add(EntityId(n), &testComponent{a: int32(n), b: int32(n)})
		if _, exists := idSet[EntityId(n)]; exists {
			t.Fatalf("duplicate id produced on adding to storage: %v", n)
		}
		idSet[EntityId(n)] = struct{}{}
	}
}

func testStorageEmptyQuery(t *testing.T, storage EntityStorage) {
	fillStorage(storage, 10000)
	count := 0

	query := NewQuery[AllEntities]()
	assert.NotNil(t, query, "query should not be nil")

	iter := query.ExecuteStorage(storage)
	for iter.Next() {
		count += 1
	}
	assert.Equal(t, count, 10000, "empty query should return all entities")

	iter = query.ExecuteStorage(storage)
	assert.Len(t, iter.ToList(), 10000, "safe query list for iterator should produce all entities")
}

type otherComponent struct {
	x int
}

func testStorageSimpleQuery(t *testing.T, storage EntityStorage) {
	fillStorage(storage, 10000)
	count := 0

	storage.Add(10001, &otherComponent{x: 5})

	query := NewQuery[struct {
		Id   EntityId
		Test *testComponent
	}]()

	iter := query.ExecuteStorage(storage)
	for iter.Next() {
		assert.Equal(t, iter.Item.Test.a+1, iter.Item.Test.b, "test component fields should be set for query result")
		count += 1
	}

	iter = query.ExecuteStorage(storage)
	for _, item := range iter.ToList() {
		assert.Equal(t, item.Test.a+1, item.Test.b, "test component fields should be set for query result")
	}

	assert.Equal(t, 10000, count, "simple query should return all matching entities")
}
