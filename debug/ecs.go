package debug

import (
	"fmt"
	"reflect"

	"github.com/b1naryth1ef/be/ecs"
	"github.com/inkyblackness/imgui-go/v4"
	"github.com/nwillc/genfuncs"
	"github.com/nwillc/genfuncs/container"
)

type Debugable interface {
	Debug()
}

type entityWithName struct {
	Id   ecs.EntityId
	Name *ecs.NameComponent `ecs:"optional"`
}

var entityWithNameQuery = ecs.NewQuery[entityWithName]()

type ECSDebugWindow struct {
	Open bool

	selectedStage  *ecs.SystemStage
	selectedSystem ecs.System
}

func NewECSDebugWindow() *ECSDebugWindow {
	return &ECSDebugWindow{
		Open: true,
	}
}

func (e *ECSDebugWindow) Render(sim *ecs.Simulation) {
	if imgui.BeginV("ECS Debugger", &e.Open, imgui.WindowFlagsAlwaysAutoResize) {
		imgui.BeginTabBar("ecs-debugger-pane")

		if imgui.BeginTabItem("Entities") {
			e.renderEntities(sim)
			imgui.EndTabItem()
		}

		if imgui.BeginTabItem("Systems") {
			e.renderSystems(sim)
			imgui.EndTabItem()
		}

		imgui.EndTabBar()
	}
	imgui.End()
}

func (e *ECSDebugWindow) renderSystemStage(sim *ecs.Simulation, stage *ecs.SystemStage) {
	for _, subStage := range stage.SubStages {
		if imgui.TreeNodeV(subStage.Label, imgui.TreeNodeFlagsFramed) {
			e.renderSystemStage(sim, subStage)
			imgui.TreePop()
		}
	}

	for _, system := range stage.All() {
		systemName := reflect.TypeOf(system).Elem().Name()
		if imgui.TreeNodeV(systemName, imgui.TreeNodeFlagsFramed) {
			RenderStruct(system)

			if d, ok := system.(Debugable); ok {
				d.Debug()
			}

			imgui.TreePop()
		}
	}
}

func (e *ECSDebugWindow) renderSystems(sim *ecs.Simulation) {
	scheduler, ok := sim.Executor.(*ecs.SystemScheduler)
	if !ok {
		return
	}

	stages := scheduler.Stages()
	for _, stage := range stages {
		if imgui.TreeNodeV(stage.Label, imgui.TreeNodeFlagsFramed) {
			e.renderSystemStage(sim, stage)
			imgui.TreePop()
		}
	}
}

func (e *ECSDebugWindow) renderEntities(sim *ecs.Simulation) {
	if imgui.BeginTableV("entities", 3, imgui.TableFlagsBorders|imgui.TableFlagsResizable, imgui.Vec2{}, 0.0) {
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
				sim.AddEntity(
					NewECSDebugEntityWindow(iter.Item.Id),
					&ecs.NameComponent{
						Name: fmt.Sprintf("ECS Debug Entity %v", iter.Item.Id),
					})
			}
		}
		imgui.EndTable()
	}
}

type ECSDebugEntityWindow struct {
	Id   ecs.EntityId
	Open bool
}

func NewECSDebugEntityWindow(id ecs.EntityId) *ECSDebugEntityWindow {
	return &ECSDebugEntityWindow{
		Id:   id,
		Open: true,
	}
}

func (e *ECSDebugEntityWindow) Render(sim *ecs.Simulation) {
	imgui.BeginV(fmt.Sprintf("Entity %v Details", e.Id), &e.Open, 0)
	components := sim.Storage.Get(e.Id)

	if len(components) == 0 && !sim.Storage.Has(e.Id) {
		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{255, 0, 0, 255})
		imgui.Text("Entity Deleted")
		imgui.PopStyleColor()
	}

	container.GSlice[interface{}](components).SortBy(func(a interface{}, b interface{}) bool {
		return genfuncs.OrderedLessThan(reflect.TypeOf(a).Elem().Name())(reflect.TypeOf(b).Elem().Name())
	})

	for _, componentData := range components {
		componentType := reflect.TypeOf(componentData).Elem()
		if imgui.CollapsingHeaderV(componentType.Name(), imgui.TreeNodeFlagsDefaultOpen) {
			RenderStruct(componentData)
			if dbg, ok := componentData.(Debugable); ok {
				dbg.Debug()
			}
		}
	}

	imgui.End()
}

type ECSDebugSystem struct {
	openEntityWindows map[ecs.EntityId]struct{}

	selectedStage  *ecs.SystemStage
	selectedSystem ecs.System
}

func NewECSDebugSystem() *ECSDebugSystem {
	return &ECSDebugSystem{
		openEntityWindows: map[ecs.EntityId]struct{}{},
	}
}

func (d *ECSDebugSystem) Update(frame *ecs.SimulationFrame) {
}

func (d *ECSDebugSystem) Render(frame *ecs.SimulationFrame) {
	debugWindows := ecs.NewQuery[struct {
		ecs.EntityId
		*ECSDebugWindow
	}]().Execute(frame.Sim)
	for debugWindows.Next() {
		if debugWindows.Item.Open {
			debugWindows.Item.Render(frame.Sim)
		} else {
			frame.Sim.DeleteEntity(debugWindows.Item.EntityId)
		}
	}

	debugEntities := ecs.NewQuery[struct {
		ecs.EntityId
		*ECSDebugEntityWindow
	}]().Execute(frame.Sim)
	for debugEntities.Next() {
		if debugEntities.Item.Open {
			debugEntities.Item.Render(frame.Sim)
		} else {
			frame.Sim.DeleteEntity(debugEntities.Item.EntityId)
		}
	}
}
