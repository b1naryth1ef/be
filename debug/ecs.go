package debug

import (
	"fmt"
	"reflect"

	"github.com/b1naryth1ef/be/ecs"
	"github.com/inkyblackness/imgui-go/v4"
	"github.com/nwillc/genfuncs"
	"github.com/nwillc/genfuncs/container"
)

type Debuggable interface {
	Debug()
}

type entityWithName struct {
	Id   ecs.EntityId
	Name *ecs.NameComponent `ecs:"optional"`
}

var entityWithNameQuery = ecs.NewQuery[entityWithName]()

type ECSDebugSystem struct {
	openEntityWindows map[ecs.EntityId]struct{}

	selectedSystem ecs.System
}

func NewECSDebugSystem() *ECSDebugSystem {
	return &ECSDebugSystem{
		openEntityWindows: map[ecs.EntityId]struct{}{},
	}
}

func (d *ECSDebugSystem) Update(frame *ecs.SimulationFrame) {
	d.renderDebugWindow(frame.Sim)
}

func (d *ECSDebugSystem) Render(frame *ecs.SimulationFrame) {}

func (d *ECSDebugSystem) renderDebugWindow(sim *ecs.Simulation) {
	open := true
	imgui.BeginV("ECS Debugger", &open, imgui.WindowFlagsAlwaysAutoResize)
	imgui.BeginTabBar("ecs-debugger-pane")

	if imgui.BeginTabItem("Entities") {
		d.renderEntities(sim)
		imgui.EndTabItem()
	}

	if imgui.BeginTabItem("Systems") {
		d.renderSystems(sim)
		imgui.EndTabItem()
	}

	imgui.EndTabBar()
	imgui.End()

	d.renderOpenEntities(sim)
}

func (d *ECSDebugSystem) renderSystems(sim *ecs.Simulation) {
	selectedSystemName := ""
	if d.selectedSystem != nil {
		selectedSystemName = reflect.TypeOf(d.selectedSystem).Elem().Name()
	}
	if imgui.BeginCombo("Selected System", selectedSystemName) {
		systems := sim.Executor.All()
		for _, system := range systems {
			if imgui.SelectableV(
				reflect.TypeOf(system).Elem().Name(),
				d.selectedSystem != system,
				0,
				imgui.Vec2{X: 0, Y: 0},
			) {
				d.selectedSystem = system
			}
		}

		imgui.EndCombo()
	}

	if d.selectedSystem != nil {
		RenderStruct(d.selectedSystem)
		if dbg, ok := d.selectedSystem.(Debuggable); ok {
			dbg.Debug()
		}
	}
}

func (d *ECSDebugSystem) renderOpenEntities(sim *ecs.Simulation) {
	for entityId := range d.openEntityWindows {
		open := true
		// imgui.SetNextWindowSize(imgui.Vec2{})
		imgui.BeginV(fmt.Sprintf("Entity %v Details", entityId), &open, 0)
		components := sim.Storage.Get(entityId)

		if !open {
			delete(d.openEntityWindows, entityId)
		}

		container.GSlice[interface{}](components).SortBy(func(a interface{}, b interface{}) bool {
			return genfuncs.OrderedLessThan(reflect.TypeOf(a).Elem().Name())(reflect.TypeOf(b).Elem().Name())
		})

		for _, componentData := range components {
			componentType := reflect.TypeOf(componentData).Elem()
			if imgui.CollapsingHeaderV(componentType.Name(), imgui.TreeNodeFlagsDefaultOpen) {
				RenderStruct(componentData)
				if dbg, ok := componentData.(Debuggable); ok {
					dbg.Debug()
				}
			}
		}

		imgui.End()
	}
}

func (d *ECSDebugSystem) renderEntities(sim *ecs.Simulation) {
	if imgui.BeginTableV("entities", 3, imgui.TableFlagsBorders, imgui.Vec2{}, 0.0) {
		imgui.TableSetupColumn("ID")
		imgui.TableSetupColumn("Name")
		imgui.TableSetupColumn("Options")
		imgui.TableHeadersRow()

		iter := entityWithNameQuery.Execute(sim)
		iter.Sort()

		for iter.Next() {
			imgui.TableNextRow()
			imgui.TableNextColumn()
			imgui.Text(fmt.Sprintf("%v", iter.Item.Id))

			imgui.TableNextColumn()
			if iter.Item.Name != nil {
				imgui.Text(iter.Item.Name.Name)
			} else {
				imgui.Text("")
			}

			imgui.TableNextColumn()
			if imgui.Button(fmt.Sprintf("View###%v", iter.Item.Id)) {
				d.openEntityWindows[iter.Item.Id] = struct{}{}
			}
		}
		imgui.EndTable()
	}
}
