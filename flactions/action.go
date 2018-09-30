package flactions

// Action provides an interface for all actions in the system
type Action interface {
	Perform(Context)
}
