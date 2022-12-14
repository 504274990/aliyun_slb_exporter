package collector

import (
	"encoding/json"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/promlog"
	"reflect"
	"sync"
)

var (
	logger = promlog.New(&promlog.Config{})
)

type slbCollector struct {
	ActiveConnection                 *prometheus.Desc
	MaxConnection                    *prometheus.Desc
	NewConnection                    *prometheus.Desc
	PacketRX                         *prometheus.Desc
	PacketTX                         *prometheus.Desc
	TrafficRXNew                     *prometheus.Desc
	TrafficTXNew                     *prometheus.Desc
	InactiveConnection               *prometheus.Desc
	HeathyServerCount                *prometheus.Desc
	UnhealthyServerCount             *prometheus.Desc
	DropConnection                   *prometheus.Desc
	DropPacketRX                     *prometheus.Desc
	DropPacketTX                     *prometheus.Desc
	DropTrafficRX                    *prometheus.Desc
	DropTrafficTX                    *prometheus.Desc
	InstanceDropConnection           *prometheus.Desc
	InstanceDropPacketRX             *prometheus.Desc
	InstanceDropPacketTX             *prometheus.Desc
	InstanceDropTrafficRX            *prometheus.Desc
	InstanceDropTrafficTX            *prometheus.Desc
	InstanceActiveConnection         *prometheus.Desc
	InstanceInactiveConnection       *prometheus.Desc
	InstanceMaxConnection            *prometheus.Desc
	InstanceMaxConnectionUtilization *prometheus.Desc
	InstanceNewConnection            *prometheus.Desc
	InstanceNewConnectionUtilization *prometheus.Desc
	InstancePacketRX                 *prometheus.Desc
	InstancePacketTX                 *prometheus.Desc
	InstanceTrafficRX                *prometheus.Desc
	InstanceTrafficTX                *prometheus.Desc
	InstanceTrafficTXUtilization     *prometheus.Desc
	sMutex                           sync.Mutex
}

func NewSlbCollector() *slbCollector {
	return &slbCollector{
		ActiveConnection: prometheus.NewDesc(
			"aliyun_slb_active_connection",
			"ActiveConnection???TCP???????????????????????? Count",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		MaxConnection: prometheus.NewDesc(
			"aliyun_slb_max_connection",
			"MaxConnection????????????????????????????????? Count/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		NewConnection: prometheus.NewDesc(
			"aliyun_slb_new_connection",
			"NewConnection???TCP???????????????????????? Count",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		PacketRX: prometheus.NewDesc(
			"aliyun_slb_packet_RX",
			"PacketRX???????????????????????????????????? Count/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		PacketTX: prometheus.NewDesc(
			"aliyun_slb_packet_TX",
			"PacketTX???????????????????????????????????? Count/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		TrafficRXNew: prometheus.NewDesc(
			"aliyun_slb_traffic_rxnew",
			"TrafficRXNew???????????????????????? bit/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		TrafficTXNew: prometheus.NewDesc(
			"aliyun_slb_traffic_txnew",
			"TrafficTXNew???????????????????????? bit/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InactiveConnection: prometheus.NewDesc(
			"aliyun_slb_inactive_connection",
			"InactiveConnection???????????????????????????????????? Count",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		HeathyServerCount: prometheus.NewDesc(
			"aliyun_slb_heathy_servercount",
			"HeathyServerCount???????????????ECS????????????????????? Count",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		UnhealthyServerCount: prometheus.NewDesc(
			"aliyun_slb_unhealthy_servercount",
			"UnhealthyServerCount???????????????ECS????????????????????? Count",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		DropConnection: prometheus.NewDesc(
			"aliyun_slb_drop_connection",
			"DropConnection??????????????????????????????????????? Count/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		DropPacketRX: prometheus.NewDesc(
			"aliyun_slb_drop_packet_RX",
			"DropPacketRX??????????????????????????????????????? Count/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		DropPacketTX: prometheus.NewDesc(
			"aliyun_slb_drop_packet_TX",
			"DropPacketTX??????????????????????????????????????? Count/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		DropTrafficRX: prometheus.NewDesc(
			"aliyun_slb_drop_traffic_RX",
			"DropTrafficRX????????????????????????bit???????????? bit/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		DropTrafficTX: prometheus.NewDesc(
			"aliyun_slb_drop_traffic_TX",
			"DropTrafficTX????????????????????????bit???????????? bit/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstanceDropConnection: prometheus.NewDesc(
			"aliyun_slb_instance_drop_connection",
			"InstanceDropConnection??????????????????????????????????????? Count/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstanceDropPacketRX: prometheus.NewDesc(
			"aliyun_slb_instance_drop_packet_RX",
			"InstanceDropPacketRX??????????????????????????????????????? Count/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstanceDropPacketTX: prometheus.NewDesc(
			"aliyun_slb_instance_drop_packet_TX",
			"InstanceDropPacketTX??????????????????????????????????????? Count/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstanceDropTrafficRX: prometheus.NewDesc(
			"aliyun_slb_instance_drop_traffic_RX",
			"InstanceDropTrafficRX????????????????????????bit???????????? bit/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstanceDropTrafficTX: prometheus.NewDesc(
			"aliyun_slb_instance_drop_traffic_TX",
			"InstanceDropTrafficTX????????????????????????bit???????????? bit/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstanceActiveConnection: prometheus.NewDesc(
			"aliyun_slb_instance_active_connection",
			"InstanceActiveConnection???????????????????????????Count/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstanceInactiveConnection: prometheus.NewDesc(
			"aliyun_slb_instance_inactive_connection",
			"InstanceInactiveConnection????????????????????????????????????Count/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstanceMaxConnection: prometheus.NewDesc(
			"aliyun_slb_instance_maxconnection",
			"InstanceMaxConnection????????????????????????????????????????????? Count/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstanceMaxConnectionUtilization: prometheus.NewDesc(
			"aliyun_slb_instance_maxconnection_utilization",
			"InstanceMaxConnectionUtilization???????????????????????????????????? %",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstanceNewConnection: prometheus.NewDesc(
			"aliyun_slb_instance_new_connection",
			"InstanceNewConnection??????????????????????????????????????? Count/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstanceNewConnectionUtilization: prometheus.NewDesc(
			"aliyun_slb_instance_newconnection_utilization",
			"InstanceNewConnectionUtilization???????????????????????????????????? %",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstancePacketRX: prometheus.NewDesc(
			"aliyun_slb_instance_packet_RX",
			"InstancePacketRX????????????????????????????????? Count/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstancePacketTX: prometheus.NewDesc(
			"aliyun_slb_instance_packet_TX",
			"InstancePacketTX????????????????????????????????? Count/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstanceTrafficRX: prometheus.NewDesc(
			"aliyun_slb_instance_traffic_RX",
			"InstanceTrafficRX??????????????????bit???????????? bit/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstanceTrafficTX: prometheus.NewDesc(
			"aliyun_slb_instance_traffic_TX",
			"InstanceTrafficTX??????????????????bit???????????? bit/s",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
		InstanceTrafficTXUtilization: prometheus.NewDesc(
			"aliyun_slb_instance_traffic_TX_utilization",
			"InstanceTrafficTXUtilization??????????????????????????????????????? %",
			[]string{"user_id", "instance_id", "port", "vip", "instance_name"},
			nil,
		),
	}
}

func (s *slbCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.ActiveConnection
	ch <- s.ActiveConnection
	ch <- s.MaxConnection
	ch <- s.NewConnection
	ch <- s.PacketRX
	ch <- s.PacketTX
	ch <- s.TrafficRXNew
	ch <- s.TrafficTXNew
	ch <- s.InactiveConnection
	ch <- s.HeathyServerCount
	ch <- s.UnhealthyServerCount
	ch <- s.DropConnection
	ch <- s.DropPacketRX
	ch <- s.DropPacketTX
	ch <- s.DropTrafficRX
	ch <- s.DropTrafficTX
	ch <- s.InstanceDropConnection
	ch <- s.InstanceDropPacketRX
	ch <- s.InstanceDropPacketTX
	ch <- s.InstanceDropTrafficRX
	ch <- s.InstanceDropTrafficTX
	ch <- s.InstanceActiveConnection
	ch <- s.InstanceInactiveConnection
	ch <- s.InstanceMaxConnection
	ch <- s.InstanceMaxConnectionUtilization
	ch <- s.InstanceNewConnection
	ch <- s.InstanceNewConnectionUtilization
	ch <- s.InstancePacketRX
	ch <- s.InstancePacketTX
	ch <- s.InstanceTrafficRX
	ch <- s.InstanceTrafficTX
}

func (s *slbCollector) Collect(ch chan<- prometheus.Metric) {
	s.sMutex.Lock()
	defer s.sMutex.Unlock()

	slbInstanceMap := make(map[string]string)
	slbs := describeLoadBalancersResponse().Body.LoadBalancers.LoadBalancer
	for _, v := range slbs {
		slbInstanceMap[*v.LoadBalancerId] = *v.LoadBalancerName
	}

	value := reflect.ValueOf(s)
	types := reflect.TypeOf(s)
	for i := 0; i < types.Elem().NumField()-1; i++ {
		metricName := types.Elem().Field(i).Name
		var d interface{}

		response, err := describeMetricLastResponse(metricName, "acs_slb_dashboard")

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
			_, okPort := metricData["port"]
			_, okVip := metricData["vip"]

			if okPort && okVip {
				ch <- prometheus.MustNewConstMetric(
					value.Elem().FieldByName(metricName).Interface().(*prometheus.Desc),
					prometheus.GaugeValue,
					metricData["Average"].(float64),
					metricData["userId"].(string),
					metricData["instanceId"].(string),
					metricData["port"].(string),
					metricData["vip"].(string),
					slbInstanceMap[metricData["instanceId"].(string)],
				)
			} else if okPort {
				ch <- prometheus.MustNewConstMetric(
					value.Elem().FieldByName(metricName).Interface().(*prometheus.Desc),
					prometheus.GaugeValue,
					metricData["Average"].(float64),
					metricData["userId"].(string),
					metricData["instanceId"].(string),
					metricData["port"].(string),
					"",
					slbInstanceMap[metricData["instanceId"].(string)],
				)
			} else if okVip {
				ch <- prometheus.MustNewConstMetric(
					value.Elem().FieldByName(metricName).Interface().(*prometheus.Desc),
					prometheus.GaugeValue,
					metricData["Average"].(float64),
					metricData["userId"].(string),
					metricData["instanceId"].(string),
					"",
					metricData["vip"].(string),
					slbInstanceMap[metricData["instanceId"].(string)],
				)
			} else {
				ch <- prometheus.MustNewConstMetric(
					value.Elem().FieldByName(metricName).Interface().(*prometheus.Desc),
					prometheus.GaugeValue,
					metricData["Average"].(float64),
					metricData["userId"].(string),
					metricData["instanceId"].(string),
					"",
					"",
					slbInstanceMap[metricData["instanceId"].(string)],
				)
			}
		}

	}
}
