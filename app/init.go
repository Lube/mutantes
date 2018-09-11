package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-redis/redis"
	"github.com/lube/mutantes/errors"
)

// Init returns a middleware that prepares the request context and processing environment.
// The middleware will populate RequestContext, handle possible panics and errors from the processing
// handlers, and add an access log entry.
func Init(logger *logrus.Logger, redisClient *redis.Client) routing.Handler {
	return func(rc *routing.Context) error {
		now := time.Now()

		rc.Response = &access.LogResponseWriter{ResponseWriter: rc.Response, Status: http.StatusOK, BytesWritten: 0}

		rs := newRequestScope(now, logger, rc.Request, redisClient)
		rc.Set("Context", rs)

		fault.Recovery(rs.Errorf, convertError)(rc)
		logAccess(rc, rs.Infof, rs.Now())

		return nil
	}
}

// GetRequestScope returns the RequestScope of the current request.
func GetRequestScope(c *routing.Context) RequestScope {
	return c.Get("Context").(RequestScope)
}

func logAccess(c *routing.Context, logFunc access.LogFunc, start time.Time) {
	rw := c.Response.(*access.LogResponseWriter)
	elapsed := float64(time.Now().Sub(start).Nanoseconds()) / 1e6
	requestLine := fmt.Sprintf("%s %s %s", c.Request.Method, c.Request.URL.Path, c.Request.Proto)
	logFunc(`[%.3fms] %s %d %d`, elapsed, requestLine, rw.Status, rw.BytesWritten)
}

func convertError(c *routing.Context, err error) error {
	switch err.(type) {
	case *errors.APIError:
		return err
	case validation.Errors:
		return errors.InvalidData(err.(validation.Errors))
	case routing.HTTPError:
		switch err.(routing.HTTPError).StatusCode() {
		case http.StatusUnauthorized:
			return errors.Unauthorized(err.Error())
		case http.StatusNotFound:
			return errors.NotFound("the requested resource")
		}
	}
	return errors.InternalServerError(err)
}
