package fujilane

import (
	"net/http"
)

func (a *Application) routePropertiesCreate(c *routeContext) {
	v, _ := c.context.Get("current-user")
	user := v.(*User)

	property, err := a.propertiesRepository.create(user)
	if err != nil {
		c.fail(http.StatusInternalServerError, err)
		return
	}

	c.respond(http.StatusCreated, property)
}
