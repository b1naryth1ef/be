package ecs

type SystemSetup interface {
	Setup(*Simulation) error
}

type System interface {
	Update(*SimulationFrame)
	Render(*SimulationFrame)
}

type SystemExecutor interface {
	System

	Setup(*Simulation) error
	All() []System
}

/// SequentialSystemExecutor executes systems sequentially in the order they where
///  added.
type SequentialSystemExecutor struct {
	systems []System
}

func NewSequentialSystemExecutor() *SequentialSystemExecutor {
	return &SequentialSystemExecutor{systems: []System{}}
}

func (s *SequentialSystemExecutor) Setup(sim *Simulation) error {
	for _, system := range s.systems {
		if setupSystem, ok := system.(SystemSetup); ok {
			err := setupSystem.Setup(sim)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *SequentialSystemExecutor) Update(frame *SimulationFrame) {
	for _, system := range s.systems {
		system.Update(frame)
	}
}

func (s *SequentialSystemExecutor) Render(frame *SimulationFrame) {
	for _, system := range s.systems {
		system.Render(frame)
	}
}

func (s *SequentialSystemExecutor) Add(systems ...System) {
	s.systems = append(s.systems, systems...)
}

func (s *SequentialSystemExecutor) All() []System {
	return s.systems
}
