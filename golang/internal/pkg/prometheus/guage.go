package prometheus

import "github.com/prometheus/client_golang/prometheus"

// 현재 점유 트랜잭션 수
var TotalTxUse = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "flow",
		Subsystem: "myapp",
		Name:      "total_tx_use",
		Help:      "current use",
	},
	[]string{},
)

func init() {
	prometheus.Register(TotalTxUse)
}
