package ecs

/// NameComponent provides a name for entities
type NameComponent struct {
	Name string
}

/// LabelComponent provides the ability to label entities
type LabelComponent struct {
	Labels map[string]string
}

type Component = any
