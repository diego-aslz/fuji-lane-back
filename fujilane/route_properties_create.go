package fujilane

import (
	"net/http"
)

func (a *Application) routePropertiesCreate(c *routeContext) {
	property, err := a.propertiesRepository.create(c.currentUser())
	if err != nil {
		c.fatal(err)
		return
	}

	c.respond(http.StatusCreated, property)
}
