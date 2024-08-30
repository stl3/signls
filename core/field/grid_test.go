package field

import (
	"fmt"
	"testing"

	"cykl/core/common"
	"cykl/core/node"
	"cykl/midi"
)

var benchmarks = []struct {
	size int
}{
	{size: 20},
	{size: 100},
	{size: 200},
	{size: 300},
	{size: 400},
	{size: 500},
}

func BenchmarkGrid(b *testing.B) {
	midi := &midi.Mock{}
	for _, v := range benchmarks {
		b.Run(fmt.Sprintf("grid_size_%dx%d", v.size, v.size), func(b *testing.B) {
			grid := NewGrid(v.size, v.size, midi)
			grid.AddNode(node.NewBangEmitter(midi, common.DOWN|common.RIGHT, true), 7, 7)
			grid.AddNode(node.NewSpreadEmitter(midi, common.DOWN), 11, 7)
			grid.AddNode(node.NewSpreadEmitter(midi, common.LEFT), 11, 11)
			grid.AddNode(node.NewSpreadEmitter(midi, common.UP), 7, 11)
			grid.AddNode(node.NewBangEmitter(midi, common.RIGHT, true), 7, 2)
			grid.AddNode(node.NewSpreadEmitter(midi, common.LEFT), 12, 2)
			grid.AddNode(node.NewBangEmitter(midi, common.RIGHT, true), 7, 3)
			grid.AddNode(node.NewSpreadEmitter(midi, common.LEFT), 9, 3)
			for i := 0; i < b.N; i++ {
				grid.Update()
			}
		})
	}
}
