package fujilane

import (
	"net/http"
)

func (a *Application) routeStatus(c *routeContext) {
	c.success(http.StatusOK, map[string]string{"status": "active"})
}
