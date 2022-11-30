// Copyright 2022 Authors of spidernet-io
// SPDX-License-Identifier: Apache-2.0

package loadRequest

import (
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"time"
)

type HttpMethod string

const (
	HttpMethodGet     = HttpMethod("GET")
	HttpMethodPost    = HttpMethod("POST")
	HttpMethodPut     = HttpMethod("PUT")
	HttpMethodDelete  = HttpMethod("DELETE")
	HttpMethodConnect = HttpMethod("CONNECT")
	HttpMethodOptions = HttpMethod("OPTIONS")
	HttpMethodPatch   = HttpMethod("PATCH")
	HttpMethodHead    = HttpMethod("HEAD")
)

type HttpRequestData struct {
	Method              HttpMethod
	Url                 string
	Qps                 int
	PerRequestTimeoutMS int
	RequestTimeSecond   int
}

func HttpRequest(req *HttpRequestData) *vegeta.Metrics {
	rate := vegeta.Rate{
		Freq: req.Qps,
		Per:  time.Duration(req.PerRequestTimeoutMS) * time.Millisecond,
	}
	duration := time.Duration(req.RequestTimeSecond) * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: string(req.Method),
		URL:    req.Url,
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		// if Code >= 200 && Code < 400 { m.success++ }
		metrics.Add(res)
	}
	metrics.Close()

	return &metrics
}
