package param

import (
	"fmt"
	"strconv"

	"signls/core/common"
	"signls/core/music"
	"signls/ui/util"
)

type Channel struct {
	nodes []common.Node
}

func (c Channel) Name() string {
	return "cha"
}

func (c Channel) Help() string {
	return ""
}

func (c Channel) Display() string {
	if c.nodes[0].(music.Audible).Note().Channel.RandomAmount() != 0 {
		return util.Normalize(
			fmt.Sprintf(
				"%d%+d\u033c",
				c.nodes[0].(music.Audible).Note().Channel.Value()+1,
				c.nodes[0].(music.Audible).Note().Channel.RandomAmount(),
			),
		)
	}
	return fmt.Sprintf("%d", c.nodes[0].(music.Audible).Note().Channel.Value()+1)
}

func (c Channel) Value() int {
	return int(c.nodes[0].(music.Audible).Note().Channel.Value())
}

func (c Channel) AltValue() int {
	return 0
}

func (c Channel) Up() {
	c.Set(c.Value() + 1)
}

func (c Channel) Down() {
	c.Set(c.Value() - 1)
}

func (c Channel) Left() {
	c.SetAlt(c.nodes[0].(music.Audible).Note().Channel.RandomAmount() - 1)
}

func (c Channel) Right() {
	c.SetAlt(c.nodes[0].(music.Audible).Note().Channel.RandomAmount() + 1)
}

func (c Channel) AltUp() {}

func (c Channel) AltDown() {}

func (c Channel) AltLeft() {}

func (c Channel) AltRight() {}

func (c Channel) Set(value int) {
	for _, n := range c.nodes {
		n.(music.Audible).Note().SetChannel(uint8(value))
	}
}

func (c Channel) SetAlt(value int) {
	for _, n := range c.nodes {
		n.(music.Audible).Note().Channel.SetRandomAmount(value)
	}
}

func (c Channel) SetEditValue(input string) {
	value, err := strconv.Atoi(input)
	if err != nil {
		return
	}
	c.Set(value - 1)
}
