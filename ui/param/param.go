package param

import (
	"cykl/core/common"
	"cykl/core/field"
	"cykl/core/music"
)

type Values map[int]string

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
	}
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
