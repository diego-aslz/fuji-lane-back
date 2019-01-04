package flactions

const defaultPageSize = 25

// Action provides an interface for all actions in the system
type Action interface {
	Perform()
}
