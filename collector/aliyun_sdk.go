package collector

import (
	"encoding/base64"
	cms20190101 "github.com/alibabacloud-go/cms-20190101/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	vpc20160428 "github.com/alibabacloud-go/vpc-20160428/v2/client"
	"github.com/go-kit/log/level"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	AccessKeyId     = kingpin.Flag("access.keyid", "The aliyun AccessKeyId, use base64 encode").Default("Default").String()
	AccessKeySecret = kingpin.Flag("accesss.key.secret", "The aliyun AccessKeySecret, use base64 encode").Default("Default").String()
	regionId        = kingpin.Flag("region.id", "The aliyun regionid").Default("Default, like cn-zhangjiakou").String()
	endpoint        = kingpin.Flag("endpoint", "The aliyun endpoint").Default("metrics.cn-hangzhou.aliyuncs.com").String()
)

func CreateClient(endpoint *string) (config *openapi.Config) {
	accesskeyid, _ := base64.StdEncoding.DecodeString(*AccessKeyId)
	accesskeysecret, _ := base64.StdEncoding.DecodeString(*AccessKeySecret)
	config = &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: tea.String(string(accesskeyid)),
		// 您的AccessKey Secret
		AccessKeySecret: tea.String(string(accesskeysecret)),
	}
	// 访问的域名
	config.Endpoint = endpoint

	return config
}

func describeMetricLastResponse(metrics string, namespace string) (*cms20190101.DescribeMetricLastResponse, error) {
	config := CreateClient(
		tea.String(*endpoint),
	)
	client := &cms20190101.Client{}
	client, _err := cms20190101.NewClient(config)
	if _err != nil {
		level.Error(logger).Log("msg", _err)
	}

	describeMetricLastRequest := &cms20190101.DescribeMetricLastRequest{
		Namespace:  tea.String(namespace),
		MetricName: tea.String(metrics),
		Period:     tea.String("60"),
	}
	dataResponse, _err := client.DescribeMetricLast(describeMetricLastRequest)

	return dataResponse, _err
}

func describeLoadBalancersResponse() *slb20140515.DescribeLoadBalancersResponse {
	config := CreateClient(
		tea.String("slb."+*regionId+".aliyuncs.com"),
	)
	client := &slb20140515.Client{}
	client, _err := slb20140515.NewClient(config)
	if _err != nil {
		level.Error(logger).Log("msg", _err)
	}

	describeLoadBalancersRequest := &slb20140515.DescribeLoadBalancersRequest{
		RegionId: tea.String(*regionId),
	}
	dataResponse, _err := client.DescribeLoadBalancers(describeLoadBalancersRequest)
	if _err != nil {
		level.Error(logger).Log("msg", _err)
	}

	return dataResponse
}

func describeEipAddressesResponse() *vpc20160428.DescribeEipAddressesResponse {
	config := CreateClient(
		tea.String("vpc."+*regionId+".aliyuncs.com"),
	)
	client := &vpc20160428.Client{}
	client, _err := vpc20160428.NewClient(config)
	if _err != nil {
		level.Error(logger).Log("msg", _err)
	}

	describeEipAddressesRequest := &vpc20160428.DescribeEipAddressesRequest{
		RegionId: tea.String(*regionId),
	}
	dataResponse, _err := client.DescribeEipAddresses(describeEipAddressesRequest)
	if _err != nil {
		level.Error(logger).Log("msg", _err)
	}

	return dataResponse
}
