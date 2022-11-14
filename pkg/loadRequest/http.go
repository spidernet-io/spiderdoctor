package loadRequest

import (
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"time"
)

func HttpRequest(URL string, qps int, PerRequestTimeoutSecond int, RequestTimeSecond int) *vegeta.Metrics {
	rate := vegeta.Rate{
		Freq: qps,
		Per:  time.Duration(PerRequestTimeoutSecond) * time.Second,
	}
	duration := time.Duration(RequestTimeSecond) * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
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
