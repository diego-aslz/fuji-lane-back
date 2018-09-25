package flentities

// AllEntities returns a slice with one object per persistent entity
func AllEntities() []interface{} {
	return []interface{}{
		Property{},
		User{},
		Account{},
		Country{},
	}
}
