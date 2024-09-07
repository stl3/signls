package param

import (
	"cykl/core/common"
	"cykl/core/field"
	"cykl/core/music"
	"cykl/core/node"
)

type Param interface {
	Name() string
	Value() int
	Display() string
	Set(value int)
	Increment()
	Decrement()
	Left()
	Right()
}

func NewParamsForNodes(grid *field.Grid, nodes []common.Node) []Param {
	if len(nodes) == 0 {
		return []Param{}
	} else if _, ok := nodes[0].(*node.TeleportEmitter); ok && len(nodes) == 1 {
		return []Param{
			Destination{
				nodes:  nodes,
				width:  grid.Width,
				height: grid.Height,
			},
		}
	} else if nodes[0].Name() == "quota" && len(nodes) == 1 {
		return append([]Param{
			Threshold{nodes: nodes},
		}, DefaultEmitterParams(grid, nodes)...)
	}

	emitters := filterNodes[*node.Emitter](nodes)
	return DefaultEmitterParams(grid, emitters)
}

func DefaultEmitterParams(grid *field.Grid, nodes []common.Node) []Param {
	return []Param{
		Key{
			nodes: nodes,
			keys:  music.AllKeysInScale(grid.Key, grid.Scale),
			root:  grid.Key,
			mode:  KeyMode{nodes: nodes, modes: music.AllNoteBehaviors()},
		},
		Velocity{nodes: nodes},
		Length{nodes: nodes},
		Channel{nodes: nodes},
	}
}

func NewParamsForGrid(grid *field.Grid) []Param {
	return []Param{
		Root{grid: grid},
		Scale{grid: grid, scales: music.AllScales()},
	}
}

func Get(name string, params []Param) Param {
	for _, p := range params {
		if p.Name() == name {
			return p
		}
	}
	return params[0]
}

func filterNodes[T any](nodes []common.Node) []common.Node {
	filteredNodes := []common.Node{}
	for _, n := range nodes {
		if _, ok := n.(T); !ok {
			continue
		}
		filteredNodes = append(filteredNodes, n)
	}
	return filteredNodes
}
