package pages

import (
	"github.com/codemicro/gopoll/internal/pages/internal/templates"
	"github.com/codemicro/gopoll/internal/pages/internal/templates/parts"
)

func NewPoll() string {
	head := newHeadBuilder()
	head.Add(fStr(parts.StandardHead("Hello!")))
	body := templates.NewPoll()
	return fromBase(head.Render(), body)
}
