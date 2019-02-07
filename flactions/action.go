package flactions

const defaultPageSize = 25

// Action provides an interface for all actions in the system
type Action interface {
	Perform()
}

// Provider is a function that when called with a context returns a new action instance
type Provider func(Context) Action
