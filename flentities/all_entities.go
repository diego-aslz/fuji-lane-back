package flentities

// AllEntities returns a slice with one object per persistent entity
func AllEntities() []interface{} {
	return []interface{}{
		Image{},
		Property{},
		User{},
		Account{},
		City{},
		Country{},
	}
}
