package pages

import (
	"io/ioutil"
	"testing"
)

func TestSimplePage(t *testing.T) {
	x, _ := SimplePage("Abi")
	_ = ioutil.WriteFile("out.html", []byte(x), 0644)
}
