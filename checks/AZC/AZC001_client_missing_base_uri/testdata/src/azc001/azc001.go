package azc001

type Options struct {
	SubscriptionId          string
	ResourceManagerEndpoint string
}

type FoosClient struct{}

func NewFoosClient(subscriptionId string) FoosClient { return FoosClient{} }

func NewFoosClientWithBaseURI(endpoint string, subscriptionId string) FoosClient {
	return FoosClient{}
}

// Should be flagged: client created without an explicit base URI
func badClient(o Options) FoosClient {
	return NewFoosClient(o.SubscriptionId) // want `Azure SDK clients should be created with NewFoosClientWithBaseURI and the resource manager endpoint explicitly specified`
}

// Should NOT be flagged: client created with the resource manager endpoint
func goodClient(o Options) FoosClient {
	return NewFoosClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
}

// Should NOT be flagged: not an options struct subscription id
func goodOtherCall(subscriptionId string) FoosClient {
	return NewFoosClient(subscriptionId)
}
