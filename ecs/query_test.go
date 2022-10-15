package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type componentA struct {
	A float64
}

type componentB struct {
	B int64
}

type componentC struct {
	C bool
}

func TestQueryToList(t *testing.T) {
	sim := NewSimpleSimulation()
	for n := 0; n < 1000; n++ {
		sim.AddEntity(&componentA{
			A: float64(n),
		})
	}

	query := NewQuery[struct {
		Id EntityId
		A  *componentA
	}]()
	iter := query.Execute(sim)

	result := iter.ToList()
	assert.Len(t, result, 1000)
	for _, item := range result {
		assert.Equal(t, item.Id, EntityId(item.A.A))
	}
}
func TestQueryFirst(t *testing.T) {
	sim := NewSimpleSimulation()
	sim.AddEntity(&componentA{
		A: float64(1),
	})
	sim.AddEntity(&componentA{
		A: float64(2),
	})

	query := NewQuery[struct {
		Id EntityId
		A  *componentA
	}]()
	iter := query.Execute(sim)

	item, ok := iter.First()
	assert.True(t, ok)
	assert.Equal(t, item.A.A, float64(1))
}

func TestQueryOptionalField(t *testing.T) {
	sim := NewSimpleSimulation()
	for n := 0; n < 1000; n++ {
		sim.AddEntity(&componentA{
			A: float64(n),
		}, &componentB{
			B: int64(n),
		})
	}

	queryNone := NewQuery[struct {
		A *componentA
		C *componentC
	}]()
	iterNone := queryNone.Execute(sim)
	assert.Len(t, iterNone.ToList(), 0)

	queryAll := NewQuery[struct {
		A *componentA
		C *componentC `ecs:"optional"`
	}]()
	iterAll := queryAll.Execute(sim)
	assert.Len(t, iterAll.ToList(), 1000)
}
