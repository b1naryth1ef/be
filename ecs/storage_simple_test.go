package ecs

import "testing"

func TestSimpleStorageAdd(t *testing.T) {
	testStorageAdd(t, NewEntitySimpleStorage())
}

func TestSimpleStorageEmptyQuery(t *testing.T) {
	testStorageEmptyQuery(t, NewEntitySimpleStorage())
}

func TestSimpleStorageSimpleQuery(t *testing.T) {
	testStorageSimpleQuery(t, NewEntitySimpleStorage())
}
