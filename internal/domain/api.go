// Copyright (c) 2025 EterLine (Andrew)
// This file is part of fstmon.
// Licensed under the MIT License. See the LICENSE file for details.

package domain

import (
	"context"
	"net"
	"net/http"
	"net/netip"
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
	SourceIP string `json:"source_ip"`
	ClientIP string `json:"client_ip"`
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

func InitRequestInfo(r *http.Request, ext IpExtractor) RequestInfo {
	ip, err := ext.ExtractIP(r)
	if err != nil {
		panic(err)
	}

	return RequestInfo{
		ClientIP: ip.String(),
		SourceIP: r.RemoteAddr,
	}
}
