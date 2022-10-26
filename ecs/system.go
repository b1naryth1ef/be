package ecs

type System interface {
	Update(*SimulationFrame)
}

type SystemSetup interface {
	Setup(*Simulation) error
}

type SystemRender interface {
	Render(*SimulationFrame)
}

type SystemExecutor interface {
	System
	SystemSetup
	SystemRender
}

type SystemStage struct {
	Label   string
	Enabled bool

	executors []SystemExecutor
	updates   []System
	setups    []SystemSetup
	renders   []SystemRender
}

func NewSystemStage(label string) *SystemStage {
	return &SystemStage{
		Label:     label,
		Enabled:   true,
		executors: make([]SystemExecutor, 0),
		updates:   make([]System, 0),
		setups:    make([]SystemSetup, 0),
		renders:   make([]SystemRender, 0),
	}
}

func (s *SystemStage) All() []System {
	return s.updates
}

func (s *SystemStage) AddSub(executors ...SystemExecutor) {
	s.executors = append(s.executors, executors...)
}

func (s *SystemStage) Add(systems ...System) {
	for _, system := range systems {
		if setupSystem, ok := system.(SystemSetup); ok {
			s.setups = append(s.setups, setupSystem)
		}
		if renderSystem, ok := system.(SystemRender); ok {
			s.renders = append(s.renders, renderSystem)
		}
	}
	s.updates = append(s.updates, systems...)
}

func (s *SystemStage) Update(frame *SimulationFrame) {
	for _, sub := range s.executors {
		sub.Update(frame)
	}
	for _, system := range s.updates {
		system.Update(frame)
	}
}

func (s *SystemStage) Render(frame *SimulationFrame) {
	for _, sub := range s.executors {
		sub.Render(frame)
	}
	for _, system := range s.renders {
		system.Render(frame)
	}
}

func (s *SystemStage) Setup(sim *Simulation) error {
	for _, sub := range s.executors {
		err := sub.Setup(sim)
		if err != nil {
			return err
		}
	}
	for _, system := range s.setups {
		err := system.Setup(sim)
		if err != nil {
			return err
		}
	}
	return nil
}

type SystemScheduler struct {
	stages []*SystemStage
}

func NewSystemScheduler() *SystemScheduler {
	return &SystemScheduler{
		stages: []*SystemStage{},
	}
}

func (s *SystemScheduler) Stages() []*SystemStage {
	return s.stages
}

func (s *SystemScheduler) ByName(name string) *SystemStage {
	for _, stage := range s.stages {
		if stage.Label == name {
			return stage
		}
	}
	return nil
}

func (s *SystemScheduler) Add(stages ...*SystemStage) {
	s.stages = append(s.stages, stages...)
}

func (s *SystemScheduler) Update(frame *SimulationFrame) {
	for _, stage := range s.stages {
		if !stage.Enabled {
			continue
		}
		stage.Update(frame)
	}
}

func (s *SystemScheduler) Render(frame *SimulationFrame) {
	for _, stage := range s.stages {
		if !stage.Enabled {
			continue
		}
		stage.Render(frame)
	}
}

func (s *SystemScheduler) Setup(sim *Simulation) error {
	for _, stage := range s.stages {
		err := stage.Setup(sim)
		if err != nil {
			return err
		}
	}
	return nil
}
