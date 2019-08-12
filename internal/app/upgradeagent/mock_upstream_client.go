package upgradeagent

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type mockUpstreamClient struct {
	context *cli.Context
}

// NewMockUpstreamClient returns a new instance of mock upstream service client
func NewMockUpstreamClient(ctx *cli.Context) UpstreamClient {
	return &mockUpstreamClient{context: ctx}
}

// RequestUpgrade implements the transport layer between Patrao and mocked upstream service
func (mock *mockUpstreamClient) RequestUpgrade(solutionInfo LocalSolutionInfo) (*UpstreamResponseUpgradeInfo, error) {
	var currentPath string

	mockURL := mock.context.GlobalString(UpstreamName) + UpstreamGetUpgrade
	upgradeInfo := NewUpstreamResponseUpgradeInfo()

	if mockURL[len(mockURL)-1:] != "/" {
		currentPath = mockURL + "/" + solutionInfo.GetName()
	} else {
		currentPath = mockURL + solutionInfo.GetName()
	}

	resp, err := http.Get(currentPath)
	if nil != err {
		log.Error(err)
		return NewUpstreamResponseUpgradeInfo(), err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		log.Error(err)
		return NewUpstreamResponseUpgradeInfo(), err
	}

	err = json.Unmarshal([]byte(body), &upgradeInfo)

	if nil != err {
		log.Error(err)
		return NewUpstreamResponseUpgradeInfo(), err
	}

	return upgradeInfo, nil
}
