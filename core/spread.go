package core

import (
	"cykl/midi"
	"fmt"
)

type SpreadEmitter struct{}

func NewSpreadEmitter(midi midi.Midi, direction Direction) *BaseEmitter {
	return &BaseEmitter{
		direction: direction,
		note:      NewNote(midi),
		behavior:  &SpreadEmitter{},
	}
}

func (e *SpreadEmitter) EmitDirections(dir Direction, pulse uint64) Direction {
	return dir
}

func (e *SpreadEmitter) ArmedOnStart() bool {
	return false
}

func (e *SpreadEmitter) Symbol(dir Direction) string {
	return fmt.Sprintf("%s%s", "S", dir.Symbol())
}

func (e *SpreadEmitter) Name() string {
	return "mult"
}

func (e *SpreadEmitter) Color() string {
	return "177"
}
