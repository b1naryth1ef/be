# ECS

This package provides an Entity Component System abstraction.

## Example

### Basic Setup

```go
// storage defines how the underlying components get stored / accessed
storage := ecs.NewEntitySimpleStorage()

// system executors manage how executors are registered and ran
executor := ecs.NewSequentialSystemExecutor()

// create our actual simulation
sim := ecs.NewSimulation(storage, executor)

for {
    // update the simulation and its systems (designed to be ran once per game update / tick)
    sim.Update()
   
    // render the simulation and its systems (can be safely called many times per frame)
    sim.Render()
}
```

### Manipulating Entities and Components

```go
type Position struct {
    X float32
    Y float32
}

type Health struct {
    Current int
    Total int
}

type PositionAndHealth struct {
    Position *Position
    Health *Health
}

// add an entity, returning its id
id := sim.AddEntity(&Position{
    X: 15.0,
    y: 15.0,
}, &Health{
    Current: 5,
    Total: 5,
})

// the most efficient way to randomly query data for entities is via an EntityView
view := ecs.NewEntityView[PositionAndHealth](sim)

// data is of type PositionAndHealth
data := view.Get(id)

// components can be removed
sim.RemoveComponent(id, Health{})

// or added
sim.AddComponent(id, &SomeOtherComponent{})

// and entities can be deleted
sim.DeleteEntity(id)
```

### Querying Entities

```go
var query = ecs.NewQuery[struct {
    Id EntityId
    *Health
    *Animation
}]()

// execute an iterate over all matching entities
result := query.Execute(sim)
for result.Next() {
    log.Printf("Health: %v / Animation: %v", result.Item.Health, result.Item.Animation)
}

// turn into a list of results
result = query.Execute(sim)
resultList := result.ToList()
```