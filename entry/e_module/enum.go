package e_module

type Module int

const (
	None Module = iota
	NodeObject
	Alarm
)

func (m Module) String() string {
	return [...]string{"", "node_object", "alarm"}[m]
}
