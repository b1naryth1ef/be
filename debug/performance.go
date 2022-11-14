package debug

import (
	"fmt"
	"math"

	"github.com/b1naryth1ef/be/ecs"
	"github.com/inkyblackness/imgui-go/v4"
)

type PerformanceDebugSystem struct {
	fps []float32 `debug:"-"`

	metricNames []string
	metrics     map[string][]float32 `debug:"-"`
}

func NewPerformanceDebugSystem() *PerformanceDebugSystem {
	return &PerformanceDebugSystem{
		fps:     make([]float32, 1),
		metrics: map[string][]float32{},
	}
}

func (d *PerformanceDebugSystem) PushMetric(name string, metric float32) {
	arr, ok := d.metrics[name]
	if !ok {
		d.metricNames = append(d.metricNames, name)
		d.metrics[name] = []float32{metric}
	} else {
		arr = append(arr, metric)
		if len(arr) >= 255 {
			arr = arr[1:]
		}
		d.metrics[name] = arr
	}
}

func (d *PerformanceDebugSystem) PushFPS(fps float32) {
	d.fps = append(d.fps, fps)
	if len(d.fps) > 64 {
		d.fps = d.fps[1:]
	}
}

func (d *PerformanceDebugSystem) Update(frame *ecs.SimulationFrame) {
}

func (d *PerformanceDebugSystem) Render(frame *ecs.SimulationFrame) {
	d.PushMetric("frametime", float32(frame.LastFrameTime))

	open := false
	imgui.SetNextWindowPos(imgui.Vec2{})
	imgui.BeginV("Performance Overlay", &open, imgui.WindowFlagsNoDecoration|imgui.WindowFlagsNoInputs)
	imgui.Text(fmt.Sprintf("%d fps (%.2fms)", int(d.fps[len(d.fps)-1]), float64(frame.LastFrameTime)/1000))
	// imgui.Text("Frame Time: ")
	// imgui.SameLine()
	// imgui.SameLineV(imgui.WindowWidth()-64, -1)
	// imgui.Text(fmt.Sprintf("%.2fms", float64(frame.LastFrameTime)/1000.0))
	// imgui.PlotLinesV("###frametime", d.frameTimes, 0, "", math.MaxFloat32, math.MaxFloat32, imgui.Vec2{
	// 	X: 250,
	// 	Y: 25,
	// })
	// imgui.SameLine()
	// imgui.SameLineV(imgui.WindowWidth()-64, -1)
	// imgui.Text(fmt.Sprintf("%.2fms", float64(frame.LastFrameTime)/1000.0))

	for _, name := range d.metricNames {
		imgui.Text(name + ": ")
		imgui.PlotLinesV("###"+name, d.metrics[name], 0, "", math.MaxFloat32, math.MaxFloat32, imgui.Vec2{
			X: 250,
			Y: 25,
		})
	}

	imgui.End()
}
