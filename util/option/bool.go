package option

var (
	True  = newBool(true)
	False = newBool(false)
)

func newBool(value bool) *bool {
	return &value
}
