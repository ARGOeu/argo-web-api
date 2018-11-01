package hbase

import (
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/tsuna/gohbase"
)

// CreateClient to hbase
func CreateClient(hbaseCfg config.HbaseConfig) gohbase.Client {
	client := gohbase.NewClient(hbaseCfg.ZkQuorum)
	return client
}

// CloseClient
func CloseClient(client gohbase.Client) {
	client.Close()
}
