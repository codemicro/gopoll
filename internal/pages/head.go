package pages

type headBuilder []func() string

func newHeadBuilder() *headBuilder {
	return new(headBuilder)
}

func (h *headBuilder) Add(x ...func() string) *headBuilder {
	*h = append(*h, x...)
	return h
}

func (h *headBuilder) Render() (o string) {
	for _, y := range *h {
		o += y()
	}
	return
}
