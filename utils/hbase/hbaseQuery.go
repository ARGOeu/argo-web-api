package hbase

import (
	"context"
	"fmt"

	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/tsuna/gohbase/filter"
	"github.com/tsuna/gohbase/hrpc"
)

const metricTblName = "metrics"
const statusTblName = "status"
const dataColumn = "data"

// Convert bytearray values to strings
func toString(b []byte) string {
	return string(b[:len(b)])
}

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
	scanRequest, _ := hrpc.NewScanStr(context.Background(), fullTblName, hrpc.Filters(pFilter))
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

// queryStatusMetrics queries all status metrics from hbase
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

	scanRequest, err := hrpc.NewScanStr(context.Background(), table, hrpc.Filters(pFilter), hrpc.NumberOfRows(100000))
	if err != nil {
		return nil, err
	}
	scanRsp, err := client.Scan(scanRequest)

	return scanRsp, err
}

// queryStatusEndpoints queries all status endpoints from hbase
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
	scanRequest, err := hrpc.NewScanStr(context.Background(), table, hrpc.Filters(pFilter), hrpc.NumberOfRows(100000))
	if err != nil {
		return nil, err
	}
	scanRsp, err := client.Scan(scanRequest)
	return scanRsp, err
}

// queryStatusServices queries all status services from hbase
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
	scanRequest, err := hrpc.NewScanStr(context.Background(), table, hrpc.Filters(pFilter), hrpc.NumberOfRows(100000))
	if err != nil {
		return nil, err
	}
	scanRsp, err := client.Scan(scanRequest)
	return scanRsp, err
}

// queryStatusGroups queries all status endpoint groups from hbase
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
	scanRequest, err := hrpc.NewScanStr(context.Background(), table, hrpc.Filters(pFilter), hrpc.NumberOfRows(100000))
	if err != nil {
		return nil, err
	}
	fmt.Println(filterStr)
	scanRsp, err := client.Scan(scanRequest)
	return scanRsp, err
}
