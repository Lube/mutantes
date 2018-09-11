package app

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

// RequestScope contains the application-specific information that are carried around in a request.
type RequestScope interface {
	Logger
	RequestID() string
	DB() *redis.Client
	Now() time.Time
}

type requestScope struct {
	Logger
	client    *redis.Client
	now       time.Time
	requestID string
}

func (rs *requestScope) DB() *redis.Client {
	return rs.client
}

func (rs *requestScope) RequestID() string {
	return rs.requestID
}

func (rs *requestScope) Now() time.Time {
	return rs.now
}

func newRequestScope(now time.Time, logger *logrus.Logger, request *http.Request, redisClient *redis.Client) RequestScope {
	l := NewLogger(logger, logrus.Fields{})
	requestID := request.Header.Get("X-Request-Id")
	if requestID != "" {
		l.SetField("RequestID", requestID)
	}
	return &requestScope{
		client:    redisClient,
		Logger:    l,
		now:       now,
		requestID: requestID,
	}
}
