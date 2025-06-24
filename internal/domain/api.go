package domain

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CtxKeyApi int

const (
	RequestInfoKey CtxKeyApi = iota
)

type ResponseAPI[T any] struct {
	Code    int    `json:"code,omitempty"`
	Data    T      `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type RequestInfo struct {
	RequestID   uuid.UUID `json:"request_id"`
	RequestTime time.Time `json:"request_time"`
	SourceIP    string    `json:"source_ip"`
}

func (i RequestInfo) IP() net.IP {
	return net.ParseIP(i.SourceIP)
}

func (i RequestInfo) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, RequestInfoKey, i)
}

func RequestInfoFromContext(ctx context.Context) (RequestInfo, bool) {
	i, ok := ctx.Value(RequestInfoKey).(RequestInfo)
	return i, ok
}

func InitRequestInfo(r *http.Request) RequestInfo {

	addr := r.Header.Get("X-Real-IP")
	if addr == "" {
		addr = r.RemoteAddr
	}

	return RequestInfo{
		RequestID:   uuid.New(),
		RequestTime: time.Now(),
		SourceIP:    addr,
	}
}
