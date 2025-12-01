// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.
package middleware

import (
	"net/http"
	"net/netip"
)

type IpExtractor interface {
	ExtractIP(*http.Request) (client netip.Addr, remote netip.AddrPort, err error)
}

type BearerTester interface {
	TestBearer(string) bool
}

type SubnetAllower interface {
	InAllowedSubnets(netip.Addr) bool
}

// type CtxKeyApi int

// const (
// 	RequestInfoKey CtxKeyApi = iota
// )

// type ResponseAPI[T any] struct {
// 	Code    int    `json:"code,omitempty"`
// 	Data    T      `json:"data,omitempty"`
// 	Message string `json:"message,omitempty"`
// }

// type RequestInfo struct {
// 	Source    netip.AddrPort
// 	Client    netip.Addr
// 	startedAt time.Time
// }

// func (i RequestInfo) RequestCreated() time.Time {
// 	return i.startedAt
// }

// func (i RequestInfo) RequestDuration() time.Duration {
// 	return time.Since(i.startedAt)
// }

// func (i RequestInfo) ToContext(ctx context.Context) context.Context {
// 	return context.WithValue(ctx, RequestInfoKey, i)
// }

// func RequestInfoFromContext(ctx context.Context) (RequestInfo, bool) {
// 	i, ok := ctx.Value(RequestInfoKey).(RequestInfo)
// 	return i, ok
// }

// func InitRequestInfo(r *http.Request, ext IpExtractor) RequestInfo {
// 	client, src, err := ext.ExtractIP(r)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return RequestInfo{
// 		Client:    client,
// 		Source:    src,
// 		startedAt: time.Now(),
// 	}
// }
