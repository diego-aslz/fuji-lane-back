package flactions

// Action provides an interface for all actions in the system
type Action interface {
	Perform(Context)
}

// ActionCreator creates new instances of actions
type ActionCreator func() Action
