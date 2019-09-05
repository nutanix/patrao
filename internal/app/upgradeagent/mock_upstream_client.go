package upgradeagent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
		return NewUpstreamResponseUpgradeInfo(), fmt.Errorf("MockUpstreamClient::RequestUpgrade() [%v]", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return NewUpstreamResponseUpgradeInfo(), fmt.Errorf("MockUpstreamClient::RequestUpgrade() [%v]", err)
	}

	err = json.Unmarshal([]byte(body), &upgradeInfo)

	if nil != err {
		return NewUpstreamResponseUpgradeInfo(), fmt.Errorf("MockUpstreamClient::RequestUpgrade() [%v]", err)
	}

	return upgradeInfo, nil
}
