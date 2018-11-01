package hbase

import (
	"context"
	"errors"
	"time"

	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/tsuna/gohbase/filter"
	"github.com/tsuna/gohbase/hrpc"
)

const metricTblName = "metrics"
const statusTblName = "status"
const dataColumn = "data"
const hbaseTimeout = 3

// Convert bytearray values to strings
func toString(b []byte) string {
	return string(b[:len(b)])
}

//CellsToMap converts hbase data to map representation
func CellsToMap(cells []*hrpc.Cell) map[string]string {
	result := make(map[string]string)
	for _, cell := range cells {
		key := toString(cell.Qualifier)
		value := toString(cell.Value)
		result[key] = value
	}

	return result
}

// queryRawMetrics querys a raw metric from hbase
func queryRawMetrics(hcfg config.HbaseConfig, tenant string, hostname string, service string, metric string, ts string) (map[string]string, error) {
	client := CreateClient(hcfg)
	defer CloseClient(client)
	result := make(map[string]string)

	fullTblName := tenant + ":" + metricTblName
	rowPrefix := hostname + "|" + service + "|" + metric + "|" + ts
	pFilter := filter.NewPrefixFilter([]byte(rowPrefix))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(hbaseTimeout)*time.Second)
	defer cancel()
	scanRequest, _ := hrpc.NewScanStr(ctx, fullTblName, hrpc.Filters(pFilter))
	scanRsp, err := client.Scan(scanRequest)

	if err != nil {
		return result, err
	}

	if len(scanRsp) > 0 {
		r1 := scanRsp[0]
		result = CellsToMap(r1.Cells)
	}

	return result, nil
}

// QueryStatusMetrics queries all status metrics from hbase
func QueryStatusMetrics(hcfg config.HbaseConfig, tenant string, report string, date string, group string, service string, endpoint string, metric string) ([]*hrpc.Result, error) {

	client := CreateClient(hcfg)
	defer CloseClient(client)
	tp := "metric"
	table := tenant + ":" + "status"
	filterStr := report + "|" + tp + "|" + date + "|" + group + "|" + service + "|" + endpoint
	if metric != "" {
		filterStr = filterStr + "|" + metric
	}
	pFilter := filter.NewPrefixFilter([]byte(filterStr))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(hbaseTimeout)*time.Second)
	defer cancel()
	scanRequest, err := hrpc.NewScanStr(ctx, table, hrpc.Filters(pFilter), hrpc.NumberOfRows(100000))
	if err != nil {
		return nil, err
	}
	err = errors.New("hbase timeout")
	scanRsp, err := client.Scan(scanRequest)

	return scanRsp, err
}

// QueryStatusEndpoints queries all status endpoints from hbase
func QueryStatusEndpoints(hcfg config.HbaseConfig, tenant string, report string, date string, group string, service string, endpoint string) ([]*hrpc.Result, error) {
	client := CreateClient(hcfg)
	defer CloseClient(client)
	tp := "endpoint"
	table := tenant + ":" + "status"
	filterStr := report + "|" + tp + "|" + date + "|" + group + "|" + service
	if endpoint != "" {
		filterStr = filterStr + "|" + endpoint
	}

	pFilter := filter.NewPrefixFilter([]byte(filterStr))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(hbaseTimeout)*time.Second)
	defer cancel()
	scanRequest, err := hrpc.NewScanStr(ctx, table, hrpc.Filters(pFilter), hrpc.NumberOfRows(100000))
	if err != nil {
		return nil, err
	}
	err = errors.New("hbase timeout")
	scanRsp, err := client.Scan(scanRequest)
	return scanRsp, err
}

// QueryStatusServices queries all status services from hbase
func QueryStatusServices(hcfg config.HbaseConfig, tenant string, report string, date string, group string, service string) ([]*hrpc.Result, error) {
	client := CreateClient(hcfg)
	defer CloseClient(client)
	tp := "service"
	table := tenant + ":" + "status"
	filterStr := report + "|" + tp + "|" + date + "|" + group
	if service != "" {
		filterStr = filterStr + "|" + service
	}
	pFilter := filter.NewPrefixFilter([]byte(filterStr))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(hbaseTimeout)*time.Second)
	defer cancel()
	scanRequest, err := hrpc.NewScanStr(ctx, table, hrpc.Filters(pFilter), hrpc.NumberOfRows(100000))
	if err != nil {
		return nil, err
	}
	err = errors.New("hbase timeout")
	scanRsp, err := client.Scan(scanRequest)
	return scanRsp, err
}

// QueryStatusGroups queries all status endpoint groups from hbase
func QueryStatusGroups(hcfg config.HbaseConfig, tenant string, report string, date string, group string) ([]*hrpc.Result, error) {
	client := CreateClient(hcfg)
	defer CloseClient(client)
	tp := "endpoint_group"
	table := tenant + ":" + "status"
	filterStr := report + "|" + tp + "|" + date
	if group != "" {
		filterStr = filterStr + "|" + group
	}

	pFilter := filter.NewPrefixFilter([]byte(filterStr))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(hbaseTimeout)*time.Second)
	defer cancel()
	scanRequest, err := hrpc.NewScanStr(ctx, table, hrpc.Filters(pFilter), hrpc.NumberOfRows(100000))
	if err != nil {
		return nil, err
	}

	err = errors.New("hbase timeout")
	scanRsp, err := client.Scan(scanRequest)

	return scanRsp, err
}
