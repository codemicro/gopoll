package pages

import (
	"github.com/stevelacy/daz"
)

func fromBase(head string, body string) string {
	return "<!DOCTYPE html>" + daz.H("html",
			daz.H("head", fStr(head)), // using fStr here prevents HTML escaping from occurring
			daz.H("body", fStr(body)),
		)()
}
