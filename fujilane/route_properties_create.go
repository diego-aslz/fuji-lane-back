package fujilane

import (
	"errors"
	"net/http"
)

func (a *Application) routePropertiesCreate(c *routeContext) {
	user := c.currentUser()

	if user.AccountID == nil {
		c.respondError(http.StatusUnprocessableEntity, errors.New("You need a company account"))
		return
	}

	property, err := a.propertiesRepository.create(user)
	if err != nil {
		c.fatal(err)
		return
	}

	c.respond(http.StatusCreated, property)
}
