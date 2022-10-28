package collector

import (
	"encoding/json"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"reflect"
	"strings"
	"sync"
)

type eipCollector struct {
	NetRxRate             *prometheus.Desc
	NetRxPkgsRate         *prometheus.Desc
	NetTxRate             *prometheus.Desc
	NetTxPkgsRate         *prometheus.Desc
	OutRatelimitDropSpeed *prometheus.Desc
	NetInRatePercentage   *prometheus.Desc
	NetOutRatePercentage  *prometheus.Desc
	sMutex                sync.Mutex
}

func NewEipCollector() *eipCollector {
	return &eipCollector{
		NetRxRate: prometheus.NewDesc(
			"aliyun_eip_net_rx_rate",
			"net_rx.rate,流入带宽，单位 bit/s",
			[]string{"user_id", "instance_id", "ip"},
			nil,
		),
		NetRxPkgsRate: prometheus.NewDesc(
			"aliyun_eip_net_rx_pkgs_rate",
			"net_rxPkgs.rate,流入包速率，单位 Packets/s",
			[]string{"user_id", "instance_id", "ip"},
			nil,
		),
		NetTxRate: prometheus.NewDesc(
			"aliyun_eip_net_tx_rate",
			"net_tx.rate,流出带宽，单位 bit/s",
			[]string{"user_id", "instance_id", "ip"},
			nil,
		),
		NetTxPkgsRate: prometheus.NewDesc(
			"aliyun_eip_net_tx_pkgs_rate",
			"net_txPkgs.rate,流出包速率，单位 Packets/s",
			[]string{"user_id", "instance_id", "ip"},
			nil,
		),
		OutRatelimitDropSpeed: prometheus.NewDesc(
			"aliyun_eip_out_rate_limit_drop_speed",
			"out_ratelimit_drop_speed,限速丢包速率，单位 Packets/s",
			[]string{"user_id", "instance_id", "ip"},
			nil,
		),
		NetInRatePercentage: prometheus.NewDesc(
			"aliyun_eip_net_in_rate_percentage",
			"net_in.rate_percentage,网络流入带宽利用率，单位 %",
			[]string{"user_id", "instance_id", "ip"},
			nil,
		),
		NetOutRatePercentage: prometheus.NewDesc(
			"aliyun_eip_net_out_rate_percentage",
			"net_out.rate_percentage,网络流出带宽利用率，单位 %",
			[]string{"user_id", "instance_id", "ip"},
			nil,
		),
	}
}

func (e *eipCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.NetRxRate
	ch <- e.NetRxPkgsRate
	ch <- e.NetTxRate
	ch <- e.NetTxPkgsRate
	ch <- e.OutRatelimitDropSpeed
	ch <- e.NetInRatePercentage
	ch <- e.NetOutRatePercentage
}

func (e *eipCollector) Collect(ch chan<- prometheus.Metric) {
	e.sMutex.Lock()
	defer e.sMutex.Unlock()

	eipInstanceMap := make(map[string]string)
	eips := describeEipAddressesResponse().Body.EipAddresses.EipAddress
	for _,v := range eips {
		eipInstanceMap[*v.AllocationId] = *v.IpAddress
	}

	value := reflect.ValueOf(e)
	types := reflect.TypeOf(e)
	for i := 0; i < types.Elem().NumField()-1; i++ {
		metricName := strings.Split(value.Elem().FieldByIndex([]int{i, 1}).String(), ",")[0]
		eName := types.Elem().Field(i).Name

		var d interface{}

		response, err := describeMetricLastResponse(metricName, "acs_vpc_eip")

		if err != nil {
			level.Error(logger).Log("msg", err)
			break
		}

		if *response.Body.Code != "200" {
			level.Error(logger).Log("msg", "The result returned by the server is not 200", "code", *response.Body.Code, "metric", metricName)
			break
		}

		err = json.Unmarshal([]byte(*response.Body.Datapoints), &d)
		if err != nil {
			level.Error(logger).Log("msg", err)
		}

		datapoints := d.([]interface{})
		for _, datapoint := range datapoints {
			metricData := datapoint.(map[string]interface{})
			_, ok := metricData["Average"]
			if ok {
				ch <- prometheus.MustNewConstMetric(
					value.Elem().FieldByName(eName).Interface().(*prometheus.Desc),
					prometheus.GaugeValue,
					metricData["Average"].(float64),
					metricData["userId"].(string),
					metricData["instanceId"].(string),
					eipInstanceMap[metricData["instanceId"].(string)],
				)
			} else {
				ch <- prometheus.MustNewConstMetric(
					value.Elem().FieldByName(eName).Interface().(*prometheus.Desc),
					prometheus.GaugeValue,
					metricData["Value"].(float64),
					metricData["userId"].(string),
					metricData["instanceId"].(string),
					eipInstanceMap[metricData["instanceId"].(string)],
				)
			}
		}
	}
}
