package collector

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"reflect"
	"sync"
)

type natCollector struct {
	SessionActiveConnection           *prometheus.Desc
	SessionActiveConnectionWaterLever *prometheus.Desc
	SessionLimitDropConnection        *prometheus.Desc
	SessionNewConnection              *prometheus.Desc
	SessionNewConnectionWaterLever    *prometheus.Desc
	SessionNewLimitDropConnection     *prometheus.Desc
	PPSRateInFromInside               *prometheus.Desc
	PPSRateInFromOutside              *prometheus.Desc
	PPSRateOutToInside                *prometheus.Desc
	PPSRateOutToOutside               *prometheus.Desc
	BWRateInFromInside                *prometheus.Desc
	BWRateInFromOutside               *prometheus.Desc
	BWRateOutToInside                 *prometheus.Desc
	BWRateOutToOutside                *prometheus.Desc
	BytesInFromInside                 *prometheus.Desc
	BytesInFromOutside                *prometheus.Desc
	BytesOutToInside                  *prometheus.Desc
	BytesOutToOutside                 *prometheus.Desc
	sMutex                            sync.Mutex
}

func NewNatCollector() *natCollector {
	return &natCollector{
		SessionActiveConnection: prometheus.NewDesc(
			"aliyun_nat_session_active_connection",
			"SessionActiveConnection，并发连接数，单位 Count",
			[]string{"user_id", "instance_id"},
			nil,
		),
		SessionActiveConnectionWaterLever: prometheus.NewDesc(
			"aliyun_nat_session_active_connection_waterlever",
			"SessionActiveConnectionWaterLever，并发连接水位，单位 %",
			[]string{"user_id", "instance_id"},
			nil,
		),
		SessionLimitDropConnection: prometheus.NewDesc(
			"aliyun_nat_session_limit_drop_connection",
			"SessionLimitDropConnection，并发丢弃连接速率，单位 Count/s",
			[]string{"user_id", "instance_id"},
			nil,
		),
		SessionNewConnection: prometheus.NewDesc(
			"aliyun_nat_session_new_connection",
			"SessionNewConnection，新建连接速率，单位 Count/s",
			[]string{"user_id", "instance_id"},
			nil,
		),
		SessionNewConnectionWaterLever: prometheus.NewDesc(
			"aliyun_nat_session_newconnection_waterlever",
			"SessionNewConnectionWaterLever，新建连接水位，单位 %",
			[]string{"user_id", "instance_id"},
			nil,
		),
		SessionNewLimitDropConnection: prometheus.NewDesc(
			"aliyun_nat_session_newlimit_drop_connection",
			"SessionNewLimitDropConnection，新建丢弃连接速率，单位 Count/s",
			[]string{"user_id", "instance_id"},
			nil,
		),
		PPSRateInFromInside: prometheus.NewDesc(
			"aliyun_nat_ppsrate_in_from_inside",
			"PPSRateInFromInside，从VPC来包速率，单位 Count/s",
			[]string{"user_id", "instance_id"},
			nil,
		),
		PPSRateInFromOutside: prometheus.NewDesc(
			"aliyun_nat_ppsrate_in_from_outside",
			"PPSRateInFromOutside，从公网来包速率，单位 Count/s",
			[]string{"user_id", "instance_id"},
			nil,
		),
		PPSRateOutToInside: prometheus.NewDesc(
			"aliyun_nat_ppsrate_out_to_inside",
			"PPSRateOutToInside，入VPC包速率，单位 Count/s",
			[]string{"user_id", "instance_id"},
			nil,
		),
		PPSRateOutToOutside: prometheus.NewDesc(
			"aliyun_nat_ppsrate_out_to_outside",
			"PPSRateOutToOutside，入公网包速率，单位 Count/s",
			[]string{"user_id", "instance_id"},
			nil,
		),
		BWRateInFromInside: prometheus.NewDesc(
			"aliyun_nat_bwrate_in_from_inside",
			"BWRateInFromInside，从VPC来流量速率，单位 bps",
			[]string{"user_id", "instance_id"},
			nil,
		),
		BWRateInFromOutside: prometheus.NewDesc(
			"aliyun_nat_bwrate_in_from_outside",
			"BWRateInFromOutside，从公网来流量速率，单位 bps",
			[]string{"user_id", "instance_id"},
			nil,
		),
		BWRateOutToInside: prometheus.NewDesc(
			"aliyun_nat_bwrate_out_to_inside",
			"BWRateOutToInside，入VPC流量速率，单位 bps",
			[]string{"user_id", "instance_id"},
			nil,
		),
		BWRateOutToOutside: prometheus.NewDesc(
			"aliyun_nat_bwrate_out_to_outside",
			"BWRateOutToOutside，入公网流量速率，单位 bps",
			[]string{"user_id", "instance_id"},
			nil,
		),
		BytesInFromInside: prometheus.NewDesc(
			"aliyun_nat_bytes_in_from_inside",
			"BytesInFromInside，从VPC来流量，单位 Byte",
			[]string{"user_id", "instance_id"},
			nil,
		),
		BytesInFromOutside: prometheus.NewDesc(
			"aliyun_nat_bytes_in_from_outside",
			"BytesInFromOutside，从公网来流量，单位 Byte",
			[]string{"user_id", "instance_id"},
			nil,
		),
		BytesOutToInside: prometheus.NewDesc(
			"aliyun_nat_bytes_out_to_inside",
			"BytesOutToInside，入VPC流量，单位 Byte",
			[]string{"user_id", "instance_id"},
			nil,
		),
		BytesOutToOutside: prometheus.NewDesc(
			"aliyun_nat_bytes_out_to_outside",
			"BytesOutToOutside，入公网流量，单位 Byte",
			[]string{"user_id", "instance_id"},
			nil,
		),
	}
}

func (n *natCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- n.SessionActiveConnection
	ch <- n.SessionActiveConnectionWaterLever
	ch <- n.SessionLimitDropConnection
	ch <- n.SessionNewConnection
	ch <- n.SessionNewConnectionWaterLever
	ch <- n.SessionNewLimitDropConnection
	ch <- n.PPSRateInFromInside
	ch <- n.PPSRateInFromOutside
	ch <- n.PPSRateOutToInside
	ch <- n.PPSRateOutToOutside
	ch <- n.BWRateInFromInside
	ch <- n.BWRateInFromOutside
	ch <- n.BWRateOutToInside
	ch <- n.BWRateOutToOutside
	ch <- n.BytesInFromInside
	ch <- n.BytesInFromOutside
	ch <- n.BytesOutToInside
	ch <- n.BytesOutToOutside
}

func (n *natCollector) Collect(ch chan<- prometheus.Metric) {
	n.sMutex.Lock()
	defer n.sMutex.Unlock()

	value := reflect.ValueOf(n)
	types := reflect.TypeOf(n)
	for i := 0; i < types.Elem().NumField()-1; i++ {
		metricName := types.Elem().Field(i).Name
		var d interface{}

		response, err := describeMetricLastResponse(metricName, "acs_nat_gateway")

		if err != nil {
			level.Error(logger).Log("msg", err)
			break
		}

		if *response.Body.Code != "200" {
			level.Error(logger).Log("msg", "The result returned by the server is not 200", "code", *response.Body.Code, "metric", metricName)
			fmt.Println(*response)
			break
		}

		err = json.Unmarshal([]byte(*response.Body.Datapoints), &d)
		if err != nil {
			level.Error(logger).Log("msg", err)
		}

		datapoints := d.([]interface{})
		for _, datapoint := range datapoints {
			metricData := datapoint.(map[string]interface{})
			ch <- prometheus.MustNewConstMetric(
				value.Elem().FieldByName(metricName).Interface().(*prometheus.Desc),
				prometheus.GaugeValue,
				metricData["Value"].(float64),
				metricData["userId"].(string),
				metricData["instanceId"].(string),
			)
		}
	}
}
