package pages

import (
	"github.com/codemicro/gopoll/internal/pages/internal/templates"
	"github.com/codemicro/gopoll/internal/pages/internal/templates/parts"
)

func SimplePage(name string) string {

	head := newHeadBuilder()
	head.Add(fStr(parts.StandardHead("Hello!")))
	body := templates.Simple(name)

	return fromBase(head.Render(), body)

}
