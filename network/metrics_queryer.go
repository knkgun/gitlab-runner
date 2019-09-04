package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	prometheusV1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
	"gitlab.com/gitlab-org/gitlab-runner/common"
)

var metricsArtifactOptions = common.ArtifactsOptions{
	BaseName: "monitor.log",
	Format:   "raw",
	Type:     "monitor",
	ExpireIn: "10000000",
}

type MetricsQueryer struct {
	metricQueries []string
	queryInterval time.Duration
	network       common.Network
	labelName     string
	log           func() *logrus.Entry
}

func (mq *MetricsQueryer) Query(
	ctx context.Context,
	prometheusAddress string,
	labelValue string,
	startTime time.Time,
	endTime time.Time,
) (map[string][]model.SamplePair, error) {
	// create prometheus client from server address in config
	clientConfig := api.Config{Address: prometheusAddress}
	prometheusClient, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	// create a prometheus api from the client config
	prometheusAPI := prometheusV1.NewAPI(prometheusClient)
	// specify the range used for the PromQL query
	queryRange := prometheusV1.Range{
		Start: startTime,
		End:   endTime,
		Step:  mq.queryInterval,
	}

	metrics := make(map[string][]model.SamplePair)
	// use config file to pull metrics from prometheus range queries
	for _, metricQuery := range mq.metricQueries {
		selector := fmt.Sprintf("%s=\"%s\"", mq.labelName, labelValue)
		query := strings.ReplaceAll(metricQuery, "{selector}", selector)
		result, err := prometheusAPI.QueryRange(ctx, query, queryRange)
		if err != nil {
			return nil, err
		}

		// check for a result and pull first
		if result == nil || result.(model.Matrix).Len() == 0 {
			continue
		}

		// save first result set values at metric
		metrics[query] = (result.(model.Matrix)[0]).Values
	}

	return metrics, nil
}

func (mq *MetricsQueryer) Upload(
	metrics map[string][]model.SamplePair,
	jobCredentials *common.JobCredentials,
) error {
	// convert metrics sample pairs to JSON
	output, err := json.Marshal(metrics)
	if err != nil {
		fmt.Errorf("Failed to marshall metrics into json for artifact upload")
		return err
	}

	// upload JSON to GitLab as monitor.log artifact
	reader := bytes.NewReader(output)
	mq.network.UploadRawArtifacts(*jobCredentials, reader, metricsArtifactOptions)
	return nil
}

func NewMetricsQueryer(
	queryMetrics common.QueryMetricsConfig,
	labelName string,
	network common.Network,
) (*MetricsQueryer, error) {
	queryIntervalDuration, err := time.ParseDuration(queryMetrics.QueryInterval)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse query interval from config")
	}

	return &MetricsQueryer{
		metricQueries: queryMetrics.MetricQueries,
		queryInterval: queryIntervalDuration,
		labelName:     labelName,
		network:       network,
	}, nil
}
