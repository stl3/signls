package music

import (
	"fmt"
	"math/rand"
	"time"

	"signls/core/common"
	"signls/core/music/meta"
	"signls/core/theory"
	"signls/midi"
)

// Constants defining default values for note properties and their limits.
const (
	defaultKey      theory.Key = 60 // Middle C
	defaultChannel  uint8      = 0
	defaultVelocity uint8      = 100
	defaultLength   uint8      = uint8(common.PulsesPerStep)

	defaultCCNumbers int = 8

	maxVelocity    uint8 = 127
	maxLength      uint8 = 127 // 127 is treated as infinity
	minLength      uint8 = 1
	maxChannel     uint8 = 15
	maxProbability uint8 = 100
)

var lastUsedChannel uint8 = defaultChannel

// Note represents a midi note.
type Note struct {
	midi midi.Midi

	rand *rand.Rand

	Device      *DeviceValue
	Key         *KeyValue
	Channel     *common.ControlValue[uint8]
	Velocity    *common.ControlValue[uint8]
	Length      *common.ControlValue[uint8]
	Probability uint8

	Controls     []*CC
	MetaCommands []meta.Command

	pulse     uint64 // Internal pulse counter to manage note length.
	triggered bool
}

// NewNote initializes a new Note with default settings and the provided MIDI interface.
func NewNote(midi midi.Midi, device *midi.Device) *Note {
	source := rand.NewSource(time.Now().UnixNano())
	ccs := make([]*CC, defaultCCNumbers)
	for i := range ccs {
		ccs[i] = NewCC(midi, SilentControlType)
	}
	cmds := []meta.Command{
		meta.NewTempoCommand(),
		meta.NewBankCommand(),
		meta.NewRootCommand(),
		meta.NewScaleCommand(),
	}
	deviceValue := DeviceValue{
		GridDevice: device,
	}
	return &Note{
		midi:         midi,
		rand:         rand.New(source),
		Device:       &deviceValue,
		Key:          NewKeyValue(defaultKey),
		Channel:      common.NewControlValue[uint8](lastUsedChannel, 0, maxChannel),
		Velocity:     common.NewControlValue[uint8](defaultVelocity, 0, maxVelocity),
		Length:       common.NewControlValue[uint8](defaultLength, minLength, maxLength),
		Probability:  maxProbability,
		Controls:     ccs,
		MetaCommands: cmds,
	}
}

// Copy creates a copy of the note.
func (n Note) Copy() *Note {
	newDevice := *n.Device
	newKey := *n.Key
	newChannel := *n.Channel
	newVelocity := *n.Velocity
	newLength := *n.Length
	source := rand.NewSource(time.Now().UnixNano())
	newControls := make([]*CC, defaultCCNumbers)
	for i, c := range n.Controls {
		newControls[i] = c.Copy()
	}
	newCmds := make([]meta.Command, len(n.MetaCommands))
	for i, c := range n.MetaCommands {
		newCmds[i] = c.Copy()
	}
	return &Note{
		midi:         n.midi,
		rand:         rand.New(source),
		Device:       &newDevice,
		Key:          &newKey,
		Channel:      &newChannel,
		Velocity:     &newVelocity,
		Length:       &newLength,
		Probability:  n.Probability,
		Controls:     newControls,
		MetaCommands: newCmds,
	}
}

// Tick advances the internal pulse counter and stops the note if it exceeds its length.
func (n *Note) Tick() {
	if !n.triggered {
		return
	}
	n.pulse++

	// Stop the note if its duration is complete.
	if n.Length.Last() < maxLength && n.pulse >= uint64(n.Length.Last()) {
		n.Stop()
	}
}

// TransposeAndPlay triggers the note with a specific root and scale, resetting internal state.
func (n *Note) TransposeAndPlay(root theory.Key, scale theory.Scale) {
	if n.Key.IsSilent() {
		return
	}

	if n.Probability < maxProbability &&
		uint8(rand.Int31n((100))) >= n.Probability {
		return
	}

	n.Transpose(root, scale)
	n.Stop()
	n.midi.NoteOn(
		n.Device.Get(),
		n.Channel.Computed(),
		uint8(n.Key.Computed(root, scale)),
		n.Velocity.Computed(),
	)
	n.Length.Computed() // Just trigger length computation

	for _, control := range n.Controls {
		control.Send(n.Device.Get(), n.Channel.Last())
	}

	for _, cmd := range n.MetaCommands {
		cmd.Execute()
	}

	n.triggered = true
	n.pulse = 0
}

// Play just triggers the note. Used for note preview.
func (n *Note) Play() {
	if n.Key.IsSilent() {
		return
	}

	n.Stop()
	n.midi.NoteOn(
		n.Device.Get(),
		n.Channel.Value(),
		uint8(n.Key.Value()),
		n.Velocity.Value(),
	)

	n.triggered = true
	n.pulse = 0
}

// Silence silences the note channel
func (n *Note) Silence() {
	n.midi.Silence(n.Device.Get(), n.Channel.Value())
	n.triggered = false
	n.pulse = 0
}

// Stop sends a MIDI Note Off message and resets the triggered state.
func (n *Note) Stop() {
	n.midi.NoteOff(n.Device.Get(), n.Channel.Last(), uint8(n.Key.Last()))
	n.triggered = false
	n.pulse = 0
}

// Transpose transposes current key for a given root and scale.
func (n *Note) Transpose(root theory.Key, scale theory.Scale) {
	n.Key.SetNext(n.Key.key.Transpose(root, scale, n.Key.interval), root)
}

// SetKey sets the next key to play.
func (n *Note) SetKey(key theory.Key, root theory.Key) {
	n.Key.SetNext(key, root)
	if !n.triggered {
		n.Key.Set(n.Key.Value())
	}
}

// SetVelocity updates the velocity of the note.
func (n *Note) SetVelocity(velocity uint8) {
	n.Velocity.Set(velocity)
}

// SetLength updates the length of the note.
func (n *Note) SetLength(length uint8) {
	n.Length.Set(length)
}

// SetChannel updates the MIDI channel of the note.
func (n *Note) SetChannel(channel uint8) {
	n.Channel.Set(channel)
	lastUsedChannel = channel
}

// ClockDivision returns the pulses per step and steps per quarter note,
// which might be used for timing or synchronization purposes.
func (n *Note) ClockDivision() (int, int) {
	return common.PulsesPerStep, common.StepsPerQuarterNote
}

// Symbol returns the string symbol associated with the current note.
func (n *Note) Symbol() string {
	var keySymbol, paramSymbol string
	if n.Key.silent {
		keySymbol = "\u0353"
	} else if n.Key.amount != 0 {
		keySymbol = "\u033c"
	}
	controlsTypeNb := 0
	for _, c := range n.Controls {
		if c.Type != SilentControlType {
			controlsTypeNb += 1
			break
		}
	}
	for _, c := range n.MetaCommands {
		if c.Active() {
			controlsTypeNb += 1
			break
		}
	}
	switch controlsTypeNb {
	case 1:
		paramSymbol = "\u0307"
	case 2:
		paramSymbol = "\u0308"
	}
	return fmt.Sprintf("%s%s", keySymbol, paramSymbol)
}

// Midi returns the Midi interface.
func (n *Note) Midi() midi.Midi {
	return n.midi
}
