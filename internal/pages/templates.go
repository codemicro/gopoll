package pages

//go:generate qtc -dir=./../../web/templates -skipLineComments

func fStr(str string) func() string {
	return func() string {
		return str
	}
}
