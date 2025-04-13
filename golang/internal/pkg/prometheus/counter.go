package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	// 누적 요청 수.
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "flow",
			Subsystem: "myapp",
			Name:      "request",
			Help:      "accumulate count",
		},
		[]string{},
	)
	// 누적 성공 건수.
	SuccessCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "flow",
			Subsystem: "call",
			Name:      "success",
			Help:      "accumulate count",
		},
		[]string{},
	)
	// 누적 실패 건수.
	FailCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "flow",
			Subsystem: "call",
			Name:      "fail",
			Help:      "accumulate count",
		},
		[]string{},
	)
)

func init() {
	prometheus.Register(RequestCount)
	prometheus.Register(SuccessCount)
	prometheus.Register(FailCount)
}
