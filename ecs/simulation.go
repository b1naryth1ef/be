package ecs

import (
	"reflect"
	"time"
)

type SimulationFrame struct {
	Frame         uint64
	Tick          uint64
	Sim           *Simulation
	Delta         float64
	LastFrameTime uint32
	Data          map[string]interface{}
}

func WithFrameData[T any](frame *SimulationFrame, name string) T {
	return frame.Data[name].(T)
}

func WithSingleton[T any](sim *Simulation) T {
	var t T
	res := sim.Singleton(t)
	if res == nil {
		return t
	}
	return res.(T)
}

func (s *SimulationFrame) Set(key string, value interface{}) {
	s.Data[key] = value
}

type Simulation struct {
	Storage  EntityStorage
	Executor SystemExecutor
	Frame    *SimulationFrame

	id EntityId
}

func NewSimulation(storage EntityStorage, executor SystemExecutor) *Simulation {
	sim := &Simulation{
		Storage:  storage,
		Executor: executor,
		id:       1,
	}
	sim.Frame = &SimulationFrame{
		Frame:         0,
		Sim:           sim,
		Delta:         0,
		LastFrameTime: 0,
		Data:          map[string]interface{}{},
	}
	return sim
}

// NewSimpleSimulation creates a new simulation with simple defaults
func NewSimpleSimulation() *Simulation {
	return &Simulation{
		Storage:  NewEntitySimpleStorage(),
		Executor: NewSequentialSystemExecutor(),
	}
}

func (s *Simulation) AddEntity(components ...interface{}) EntityId {
	id := s.id
	s.id += 1
	s.Storage.Add(id, components...)
	return id
}

func (s *Simulation) DeleteEntity(id EntityId) {
	s.Storage.Delete(id)
}

func (s *Simulation) GetComponent(id EntityId, component interface{}) bool {
	result := s.Storage.GetComponent(id, reflect.TypeOf(component))
	if result == nil {
		return false
	}

	reflect.ValueOf(component).Elem().Set(reflect.ValueOf(result).Elem())
	return true
}

func (s *Simulation) GetComponentRef(id EntityId, component interface{}) bool {
	result := s.Storage.GetComponent(id, reflect.TypeOf(component).Elem())
	if result == nil {
		return false
	}

	reflect.ValueOf(component).Elem().Set(reflect.ValueOf(result))
	return true

}

func (s *Simulation) RemoveComponent(id EntityId, component interface{}) {
	componentType := reflect.TypeOf(component)
	if componentType.Kind() == reflect.Struct {
		componentType = reflect.PointerTo(componentType)
	} else if componentType.Kind() != reflect.Ptr || componentType.Elem().Kind() != reflect.Struct {
		return
	}

	s.Storage.RemoveComponent(id, componentType)
}

func (s *Simulation) AddComponent(id EntityId, component interface{}) {
	s.Storage.AddComponent(id, component)
}

func (s *Simulation) AddComponents(id EntityId, components ...interface{}) {
	for _, component := range components {
		s.Storage.AddComponent(id, component)
	}
}

func (s *Simulation) Singleton(component interface{}) interface{} {
	res, _ := s.Storage.FindOne(reflect.TypeOf(component))
	return res
}

func (s *Simulation) Setup() error {
	return s.Executor.Setup(s)
}

func (s *Simulation) Update() {
	s.Executor.Update(s.Frame)
	s.Frame.Tick += 1
}

func (s *Simulation) Render() {
	start := time.Now().UnixMicro()
	s.Executor.Render(s.Frame)
	s.Frame.LastFrameTime = uint32((time.Now().UnixMicro() - start))
	s.Frame.Frame += 1
}
