package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type Metrics struct {
	metrics map[string]* prometheus.Desc
	mutex sync.Mutex
}

func newGlobalMetric(namespace string, metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, labels, nil)
}


/**
 * 工厂方法：NewMetrics
 * 功能：初始化指标信息，即Metrics结构体
 */
func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		metrics: map[string]*prometheus.Desc{
			"root_disk_used_persent": newGlobalMetric(namespace, "root_disk_used_persent", "The description of root_disk_used_persent", []string{"host"}),
			"my_gauge_metric": newGlobalMetric(namespace, "my_gauge_metric","The description of my_gauge_metric", []string{"host"}),
		},
	}
}

/**
 * 接口：Describe
 * 功能：传递结构体中的指标描述符到channel
 */
func (c *Metrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

/**
 * 接口：Collect
 * 功能：抓取最新的数据，传递给channel
 */
func (c *Metrics) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()  // 加锁
	defer c.mutex.Unlock()

	mockCounterMetricData, mockGaugeMetricData := c.GenerateMockData()
	for host, currentValue := range mockCounterMetricData {
		ch <-prometheus.MustNewConstMetric(c.metrics["root_disk_used_persent"], prometheus.CounterValue, float64(currentValue), host)
	}
	for host, currentValue := range mockGaugeMetricData {
		ch <-prometheus.MustNewConstMetric(c.metrics["my_gauge_metric"], prometheus.GaugeValue, float64(currentValue), host)
	}
}


/**
 * 函数：GenerateMockData
 * 功能：生成模拟数据
 */
func (c *Metrics) GenerateMockData() (mockCounterMetricData map[string]int, mockGaugeMetricData map[string]int) {
	res := cmdres("df -Th | awk 'NR==2{print $6}' | awk -F \"%\" '{print $1}'")
	b,err := strconv.Atoi(res)
	if err != nil {
		fmt.Printf(err.Error())
	}
	mockCounterMetricData = map[string]int{
		cmdres("hostname"): b,
	}
	mockGaugeMetricData = map[string]int{
		"test": int(rand.Int31n(10)),
	}
	return
}

func cmdres(command string) string {
	out, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		log.Println("err: " + err.Error())
		return ""
	}
	return string(out)
}

