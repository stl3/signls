// Package filesystem provides interfaces and serializable structures that
// allows saving/loading grid state to/from json files.
package filesystem

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"signls/core/common"
	"signls/core/music"
	"signls/core/music/meta"
	"signls/core/theory"
)

const (
	defaultTempo                = 120.
	defaultRootKey theory.Key   = 60 // Middle C
	defaultScale   theory.Scale = theory.CHROMATIC
	defaultSize                 = 20
	maxGrids                    = 32
)

// Bank holds a slice of grids in memory
type Bank struct {
	mu sync.Mutex

	Grids    []Grid `json:"grids"`
	Active   int    `json:"active"`
	filename string
}

// Grid holds a grid in memory
type Grid struct {
	Nodes []Node  `json:"nodes"`
	Tempo float64 `json:"tempo"`

	Height int `json:"height"`
	Width  int `json:"width"`

	Device string `json:"device"`

	Key   uint8  `json:"key"`
	Scale uint16 `json:"scale"`

	SendClock     bool `json:"send_clock"`
	SendTransport bool `json:"send_transport"`
}

// NewGrid creates a new grid with default values.
func NewGrid() Grid {
	return Grid{
		Nodes:  []Node{},
		Height: defaultSize,
		Width:  defaultSize,
		Tempo:  defaultTempo,
		Key:    uint8(defaultRootKey),
		Scale:  uint16(defaultScale),
	}
}

// IsEmpty returns true if the grid is empty (no nodes).
func (g Grid) IsEmpty() bool {
	return len(g.Nodes) == 0
}

// Node represents a grid node that is json serializable.
type Node struct {
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Device    string `json:"device"`
	Note      Note   `json:"note"`
	Type      string `json:"type"`
	Direction int    `json:"direction"`
	Muted     bool   `json:"muted"`

	Params map[string]Param `json:"params"`
}

type Note struct {
	Key          Key                    `json:"key"`
	Channel      Param                  `json:"channel"`
	Velocity     Param                  `json:"velocity"`
	Length       Param                  `json:"length"`
	Probability  int                    `json:"probability"`
	Controls     []CC                   `json:"controls"`
	MetaCommands map[string]MetaCommand `json:"meta_commands"`
}

func NewNote(n music.Note) Note {
	controls := make([]CC, len(n.Controls))
	for i, c := range n.Controls {
		controls[i] = NewCC(*c)
	}
	metaCmds := make(map[string]MetaCommand, len(n.MetaCommands))
	for _, c := range n.MetaCommands {
		metaCmds[c.Name()] = NewMetaCommand(c)
	}
	return Note{
		Key:          NewKey(*n.Key),
		Channel:      NewParam(*n.Channel),
		Velocity:     NewParam(*n.Velocity),
		Length:       NewParam(*n.Length),
		Probability:  int(n.Probability),
		Controls:     controls,
		MetaCommands: metaCmds,
	}
}

type Key struct {
	Key    int
	Amount int
	Silent bool
}

func NewKey(key music.KeyValue) Key {
	return Key{
		Key:    int(key.BaseValue()),
		Amount: key.RandomAmount(),
		Silent: key.IsSilent(),
	}
}

type CC struct {
	Type       int   `json:"type"`
	Controller int   `json:"controller"`
	Value      Param `json:"value"`
}

func NewCC(cc music.CC) CC {
	return CC{
		Type:       int(cc.Type),
		Controller: int(cc.Controller),
		Value:      NewParam(*cc.Value),
	}
}

type MetaCommand struct {
	Active bool  `json:"active"`
	Value  Param `json:"value"`
}

func NewMetaCommand(cmd meta.Command) MetaCommand {
	return MetaCommand{
		Active: cmd.Active(),
		Value:  NewParam(*cmd.Value()),
	}
}

type Param struct {
	Value  int
	Amount int
}

func NewParam[T uint8 | int](p common.ControlValue[T]) Param {
	return Param{
		Value:  int(p.Value()),
		Amount: p.RandomAmount(),
	}
}

// New creates and loads a new bank from a given file.
func New(filename string) *Bank {
	grids := make([]Grid, maxGrids)
	for k := range grids {
		grids[k] = NewGrid()
	}
	bank := &Bank{
		filename: filename,
		Grids:    grids,
	}
	bank.Read(filename)
	return bank
}

// ActiveGrid returns the active grid from the bank.
func (b *Bank) ActiveGrid() Grid {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Grids[b.Active]
}

// ClearGrid clears a given grid.
func (b *Bank) ClearGrid(nb int) {
	b.Grids[nb] = NewGrid()
}

// Filename returns the bank filename.
func (b *Bank) Filename() string {
	return strings.TrimSuffix(b.filename, filepath.Ext(b.filename))
}

// Save saves a grid to the active slot and writes.
func (b *Bank) Save(grid Grid) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Grids[b.Active] = grid
	b.Write()
}

// Write serializes the Bank and writes it to a file.
func (b *Bank) Write() {
	content, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(b.filename, content, 0o644)
	if err != nil {
		log.Fatal(err)
	}
}

// Read reads a json and unmarshal its content to the Bank..
func (b *Bank) Read(filename string) {
	f, err := os.Open(filename)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return
	} else if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	content, _ := io.ReadAll(f)
	err = json.Unmarshal(content, b)
	if err != nil {
		log.Fatal(err)
	}
	b.Grids = b.Grids[:cap(b.Grids)]
}
