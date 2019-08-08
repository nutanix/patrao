package upgradeagent

type mockUpstreamClient struct {
}

// NewMockUpstreamClient returns a new instance of mock upstream service client
func NewMockUpstreamClient() UpstreamClient {
	return mockUpstreamClient{}
}

// RequestUpgrade implements the transport layer between Patrao and mocked upstream service
func (mock mockUpstreamClient) RequestUpgrade(solutionInfo LocalSolutionInfo) (*UpstreamResponseUpgradeInfo, bool) {
	return nil, false
}
