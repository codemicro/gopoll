package pages

import (
	"github.com/codemicro/gopoll/internal/pages/internal/templates"
	"github.com/codemicro/gopoll/internal/pages/internal/templates/parts"
)

func Homepage() string {
	head := newHeadBuilder()
	head.Add(fStr(parts.StandardHead("GoPoll")))
	body := templates.Homepage()
	return fromBase(head.Render(), body)
}

