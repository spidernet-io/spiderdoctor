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

func HttpRequest(method HttpMethod, URL string, qps int, PerRequestTimeoutSecond int, RequestTimeSecond int) *vegeta.Metrics {
	rate := vegeta.Rate{
		Freq: qps,
		Per:  time.Duration(PerRequestTimeoutSecond) * time.Second,
	}
	duration := time.Duration(RequestTimeSecond) * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: string(method),
		URL:    URL,
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	return &metrics
}
