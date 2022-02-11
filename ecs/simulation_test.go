package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleSimulation(t *testing.T) {
	simulation := NewSimpleSimulation()
	simulation.AddEntity(&testComponent{
		a: 5,
		b: 10,
	})
	// simulation.()
}

func TestSimulationEntities(t *testing.T) {
	simulation := NewSimpleSimulation()
	id := simulation.AddEntity(&testComponent{
		a: 5,
		b: 10,
	}, &otherComponent{x: 15})

	var test testComponent
	assert.Equal(t, true, simulation.GetComponent(id, &test), "get component should be ok")
	assert.Equal(t, int32(5), test.a, "read component field 'a' should match")

	// should be a noop
	simulation.AddComponent(id, otherComponent{
		x: 55,
	})

	var other otherComponent
	assert.Equal(t, true, simulation.GetComponent(id, &other), "get component should be ok")
	assert.Equal(t, 15, other.x, "read component field 'x' should match")

	simulation.RemoveComponent(id, otherComponent{})
	assert.Equal(t, false, simulation.GetComponent(id, &other), "get component should fail")

	simulation.AddComponent(id, &otherComponent{
		x: 55,
	})
	assert.Equal(t, true, simulation.GetComponent(id, &other), "get component should be ok ")
	assert.Equal(t, 55, other.x, "read component field 'x' should match")

	simulation.DeleteEntity(id)
	assert.Equal(t, false, simulation.GetComponent(id, &test), "get component should fail")
}
