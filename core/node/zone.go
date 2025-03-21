package node

import (
	"signls/core/common"
	"signls/core/music"
	"signls/midi"
)

type ZoneEmitter struct{}

func NewZoneEmitter(midi midi.Midi, device *midi.Device, direction common.Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi, device),
		behavior:  &ZoneEmitter{},
	}
}

func (e *ZoneEmitter) EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction {
	return dir
}

func (e *ZoneEmitter) ShouldPropagate() bool {
	return true
}

func (e *ZoneEmitter) ArmedOnStart() bool {
	return false
}

func (e *ZoneEmitter) Copy() common.EmitterBehavior {
	return &ZoneEmitter{}
}

func (e *ZoneEmitter) Symbol() string {
	return "Z"
}

func (e *ZoneEmitter) Name() string {
	return "zone"
}

func (e *ZoneEmitter) Color() string {
	return "197"
}

func (e *ZoneEmitter) Reset() {}
