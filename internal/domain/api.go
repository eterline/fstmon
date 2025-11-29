// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package domain

import (
	"context"
	"net/http"
	"net/netip"
	"time"
)

type CtxKeyApi int

const (
	RequestInfoKey CtxKeyApi = iota
)

type IpExtractor interface {
	ExtractIP(r *http.Request) (netip.Addr, error)
}

type ResponseAPI[T any] struct {
	Code    int    `json:"code,omitempty"`
	Data    T      `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type RequestInfo struct {
	Source    netip.AddrPort
	Client    netip.Addr
	startedAt time.Time
}

func (i RequestInfo) RequestDuration() time.Duration {
	return time.Since(i.startedAt)
}

func (i RequestInfo) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, RequestInfoKey, i)
}

func RequestInfoFromContext(ctx context.Context) (RequestInfo, bool) {
	i, ok := ctx.Value(RequestInfoKey).(RequestInfo)
	return i, ok
}

func InitRequestInfo(r *http.Request, ext IpExtractor) RequestInfo {
	ip, err := ext.ExtractIP(r)
	if err != nil {
		panic(err)
	}

	ap, _ := netip.ParseAddrPort(r.RemoteAddr)

	return RequestInfo{
		Client:    ip,
		Source:    ap,
		startedAt: time.Now(),
	}
}
