package node

import (
	"cykl/core/common"
)

type TeleportEmitter struct {
	direction common.Direction
	pulse     uint64
}

func NewTeleportEmitter(direction common.Direction, pulse uint64) *TeleportEmitter {
	return &TeleportEmitter{
		direction: direction,
		pulse:     pulse,
	}
}

func (s *TeleportEmitter) Direction() common.Direction {
	return s.direction
}

func (s *TeleportEmitter) SetDirection(dir common.Direction) {
	s.direction = dir
}

func (e *TeleportEmitter) ArmedOnStart() bool {
	return false
}

func (s *TeleportEmitter) Symbol() string {
	return "  "
}

func (s *TeleportEmitter) Name() string {
	return "telep"
}

func (s *TeleportEmitter) Color() string {
	return "15"
}

func (s *TeleportEmitter) updated(pulse uint64) bool {
	return s.pulse == pulse
}
