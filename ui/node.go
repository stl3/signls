package ui

import (
	"log"

	"signls/core/common"
	"signls/core/music"
	"signls/core/node"
	"signls/ui/param"
	"signls/ui/util"

	"github.com/charmbracelet/lipgloss"
)

var (
	gridStyle = lipgloss.NewStyle().
			Background(lipgloss.AdaptiveColor{Light: "254", Dark: "234"})
	cursorStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("190")).
			Foreground(lipgloss.Color("0"))
	teleportDestinationStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("160")).
					Foreground(lipgloss.Color("15"))
	selectionStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("238")).
			Foreground(lipgloss.Color("244"))
	emitterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15"))
	mutedEmitterStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("247")).
				Foreground(lipgloss.Color("236"))
	activeEmitterStyle = lipgloss.NewStyle().
				Background(lipgloss.AdaptiveColor{Light: "0", Dark: "15"}).
				Foreground(lipgloss.AdaptiveColor{Light: "15", Dark: "0"})
)

func (m mainModel) inSelectionRange(x, y int) bool {
	return x >= m.cursorX &&
		x <= m.selectionX &&
		y >= m.cursorY &&
		y <= m.selectionY
}

func (m mainModel) renderNode(n common.Node, x, y int) string {
	// render cursor
	isCursor := false
	if x == m.cursorX && y == m.cursorY && m.mode != BANK {
		isCursor = true
	}

	isTeleportDestination := false
	if m.mode == EDIT && len(m.params) > 0 {
		p, ok := m.activeParam().(param.Destination)
		if ok {
			destinationX, destinationY := p.Position()
			isTeleportDestination = (destinationX == x && destinationY == y)
		}
	}

	// render grid
	teleportDestinationSymbol := node.HoleDestinationSymbol
	if n == nil && isCursor {
		return cursorStyle.Render("  ")
	} else if n == nil && isTeleportDestination && !m.blink && m.mode != BANK {
		return cursorStyle.Render(teleportDestinationSymbol)
	} else if n == nil && isTeleportDestination && (m.blink || m.mode == BANK) {
		return teleportDestinationStyle.Render(teleportDestinationSymbol)
	} else if n == nil && m.inSelectionRange(x, y) && m.mode != BANK {
		return selectionStyle.Render("..")
	} else if n == nil {
		if (x+y)%2 == 0 {
			return "  "
		}
		return gridStyle.Render("  ")
	}

	// render node
	switch t := n.(type) {
	case common.Movable:
		if isCursor {
			return cursorStyle.Render("  ")
		}
		return activeEmitterStyle.Render("  ")
	case music.Audible:
		symbol := util.Normalize(n.Symbol())

		if isCursor && m.mode != EDIT {
			return cursorStyle.Render(symbol)
		} else if isTeleportDestination && m.mode == EDIT && m.blink {
			return teleportDestinationStyle.Render(teleportDestinationSymbol)
		} else if isCursor && m.mode == EDIT && m.blink {
			return cursorStyle.Render(symbol)
		} else if n.Activated() && t.Muted() {
			return activeEmitterStyle.Render(symbol)
		} else if t.Muted() {
			return mutedEmitterStyle.Render(symbol)
		} else if n.Activated() {
			return activeEmitterStyle.
				Foreground(lipgloss.Color(n.Color())).
				Render(symbol)
		} else {
			return emitterStyle.
				Background(lipgloss.Color(n.Color())).
				Render(symbol)
		}
	case *node.HoleEmitter:
		symbol := n.Symbol()

		if isCursor && m.mode != EDIT {
			return cursorStyle.Render(symbol)
		} else if isCursor && m.mode == EDIT && m.blink {
			return cursorStyle.Render(symbol)
		} else if n.Activated() {
			return activeEmitterStyle.
				Foreground(lipgloss.Color(n.Color())).
				Render(symbol)
		} else {
			return emitterStyle.
				Background(lipgloss.Color(n.Color())).
				Render(symbol)
		}
	default:
		log.Fatalf("cannot render node: %+v", t)
		return ""
	}
}
