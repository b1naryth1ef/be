package ecs

import (
	"log"
	"reflect"
	"strings"
	"unsafe"

	"github.com/nwillc/genfuncs"
	"github.com/nwillc/genfuncs/container"
	"github.com/viant/xunsafe"
)

// Query which returns all entities in the simulation
type AllEntities struct {
	Id EntityId
}

type QueryComponent struct {
	Index uint16
	Type  reflect.Type
}

// Query abstracts away fetching entities based on their archetype.
type Query[T any] struct {
	queryComponents  []reflect.Type
	absentComponents []reflect.Type
	components       []reflect.Type
	fields           []*xunsafe.Field
	entityId         *xunsafe.Field
}

// Read a single entity from the given storage into a pointer towards the inner query type. This is useful for reading entities into archetypes.
func (q *Query[T]) Read(storage EntityStorage, id EntityId, target unsafe.Pointer) {
	if q.entityId != nil {
		q.entityId.SetUint32(target, id)
	}

	for index, componentType := range q.components {
		field := q.fields[index]
		value := storage.GetComponent(id, componentType)

		field.SetValue(target, value)
	}
}

func (q *Query[T]) Get(sim *Simulation, id EntityId) *T {
	var result T
	ptr := &result
	q.Read(sim.Storage, id, unsafe.Pointer(ptr))
	return ptr
}

func (q *Query[T]) Execute(sim *Simulation) *QueryResultIterator[T] {
	return q.ExecuteStorage(sim.Storage)
}

func (q *Query[T]) ExecuteStorage(storage EntityStorage) *QueryResultIterator[T] {
	res := &QueryResultIterator[T]{
		ids:     storage.FindAll(q.queryComponents, q.absentComponents),
		index:   0,
		storage: storage,
		query:   q,
	}
	res.ptr = unsafe.Pointer(&res.Item)
	return res
}

func NewQuery[T any]() *Query[T] {
	var query T
	queryType := reflect.TypeOf(query)
	if queryType.Kind() != reflect.Struct {
		log.Panicf("invalid query type: %v", queryType)
		return nil
	}

	result := &Query[T]{
		components:       []reflect.Type{},
		queryComponents:  []reflect.Type{},
		absentComponents: []reflect.Type{},
		fields:           []*xunsafe.Field{},
		entityId:         nil,
	}

	for fieldIdx := 0; fieldIdx < queryType.NumField(); fieldIdx++ {
		field := queryType.Field(fieldIdx)
		if !field.IsExported() {
			continue
		}

		optional := false
		absent := false
		skipped := false

		tags := strings.Split(field.Tag.Get("ecs"), ",")
		for _, tag := range tags {
			if tag == "-" {
				skipped = true
				break
			} else if tag == "optional" {
				optional = true
			} else if tag == "absent" {
				absent = true
			}
		}

		if skipped {
			continue
		}

		if field.Type == entityIdType {
			if result.entityId != nil {
				log.Panicf("multiple entity id fields in query %v", query)
				return nil
			}

			result.entityId = xunsafe.FieldByIndex(queryType, fieldIdx)
			continue
		}

		if absent {
			result.absentComponents = append(result.absentComponents, field.Type)
			continue
		}

		if !optional {
			result.queryComponents = append(result.queryComponents, field.Type)
		}

		result.components = append(result.components, field.Type)
		result.fields = append(result.fields, xunsafe.FieldByIndex(queryType, fieldIdx))
	}

	return result
}

type QueryResultIterator[T any] struct {
	Item T

	storage EntityStorage
	query   *Query[T]
	ptr     unsafe.Pointer
	ids     []EntityId
	index   uint32
}

// Sorts the underlying entity index for this query, ensuring entities are iterated in ascending order by id
func (q *QueryResultIterator[T]) Sort() {
	container.GSlice[uint32](q.ids).SortBy(func(a EntityId, b EntityId) bool {
		return genfuncs.OrderedLessThan(a)(b)
	})
}

func (q *QueryResultIterator[T]) Get() *T {
	if q.index >= uint32(len(q.ids)) {
		return nil
	}

	var result T
	ptr := &result
	q.query.Read(q.storage, q.ids[q.index], unsafe.Pointer(ptr))
	return ptr
}

func (q *QueryResultIterator[T]) Next() bool {
	if q.index >= uint32(len(q.ids)) {
		return false
	}

	id := q.ids[q.index]
	q.query.Read(q.storage, id, q.ptr)
	q.index += 1
	return true
}

func (q *QueryResultIterator[T]) ToList() []T {
	result := make([]T, len(q.ids))
	for idx := range result {
		q.query.Read(q.storage, q.ids[idx], unsafe.Pointer(&result[idx]))
	}

	return result
}

func (q *QueryResultIterator[T]) First() (T, bool) {
	var result T
	if len(q.ids) == 0 {
		return result, false
	}
	id := q.ids[0]
	q.query.Read(q.storage, id, unsafe.Pointer(&result))
	return result, true
}
