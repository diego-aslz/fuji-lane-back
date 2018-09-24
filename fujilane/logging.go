package fujilane

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func (a *Application) logMiddleware(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path

	c.Next()

	end := time.Now()
	duration := end.Sub(start)
	method := c.Request.Method
	statusCode := c.Writer.Status()

	log.Printf("at=%v status=%d duration=%v ip=%s method=%s path=%s %s\n", end.Format("2006-01-02T15:04:05Z"),
		statusCode, duration, c.ClientIP(), method, path, c.GetString("log-details"))
}

func (a *routeContext) addLog(key, value string) {
	logs := a.context.GetString("log-details")
	if len(logs) > 0 {
		logs += " "
	}
	a.context.Set("log-details", logs+key+"="+value)
}

func (a *routeContext) addLogQuoted(key, value string) {
	a.addLog(key, "\""+value+"\"")
}

func (a *routeContext) addLogError(err error) {
	a.addLogQuoted("error", err.Error())
}

func (a *routeContext) addLogJSON(key string, value interface{}) {
	jsonObj, err := json.Marshal(value)
	if err == nil {
		a.addLog(key, string(jsonObj))
	}
}

type filterableLog interface {
	filterSensitiveInformation() filterableLog
}

func (a *routeContext) addLogFiltered(key string, value filterableLog) {
	a.addLogJSON(key, value.filterSensitiveInformation())
}
