package hbase

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gcfg.v1"
)

// This is a utility suite struct used in tests (see pkg "testify")
type hbaseTestSuite struct {
	suite.Suite
	cfg config.Config
}

// Setup the Test Environment
// This function runs before any test and setups the environment
func (suite *hbaseTestSuite) SetupTest() {

	log.SetOutput(ioutil.Discard)

	const testConfig = `
    [server]
    bindip = ""
    port = 8080
    maxprocs = 4
    cache = false
    lrucache = 700000000
    gzip = true
    [mongodb]
    host = "127.0.0.1"
    port = 27017
    db = "argo_core_test_mongo"
    [hbase]
    zkquorum="localhost"
    `

	_ = gcfg.ReadStringInto(&suite.cfg, testConfig)

}

func (suite *hbaseTestSuite) TestFindAndProject() {
	suite.Equal("localhost", suite.cfg.Hbase.ZkQuorum)

}

func (suite *hbaseTestSuite) TearDownSuite() {

}

func TestHbaseTestSuite(t *testing.T) {
	suite.Run(t, new(hbaseTestSuite))
}
