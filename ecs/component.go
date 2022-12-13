package ecs

/// NameComponent provides a name for entities
type NameComponent struct {
	Name string
}

func NewNameComponent(name string) *NameComponent {
	return &NameComponent{
		Name: name,
	}
}

/// LabelComponent provides the ability to label entities
type LabelComponent struct {
	Labels map[string]string
}

func NewLabelComponent() *LabelComponent {
	return &LabelComponent{
		Labels: make(map[string]string),
	}
}

type Component = any
