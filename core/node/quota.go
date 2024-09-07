package node

import (
	"fmt"

	"cykl/core/common"
	"cykl/core/music"
	"cykl/midi"
)

type QuotaEmitter struct {
	threshold int
	count     int
}

func NewQuotaEmitter(midi midi.Midi, direction common.Direction) *Emitter {
	return &Emitter{
		direction: direction,
		note:      music.NewNote(midi),
		behavior:  &QuotaEmitter{},
	}
}

func (e *QuotaEmitter) EmitDirections(dir common.Direction, inDir common.Direction, pulse uint64) common.Direction {
	e.count++
	if e.count < e.threshold {
		return common.NONE
	}
	e.count = 0
	return dir
}

func (e *QuotaEmitter) ArmedOnStart() bool {
	return false
}

func (e *QuotaEmitter) Symbol(dir common.Direction) string {
	return fmt.Sprintf("%s%s", "Q", dir.Symbol())
}

func (e *QuotaEmitter) Name() string {
	return "quota"
}

func (e *QuotaEmitter) Color() string {
	return "197"
}

func (e *QuotaEmitter) Threshold() int {
	return e.threshold
}

func (e *QuotaEmitter) SetThreshold(threshold int) {
	e.threshold = threshold
}
