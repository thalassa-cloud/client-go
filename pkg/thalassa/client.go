package thalassa

import (
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/pkg/iaas"
	"github.com/thalassa-cloud/client-go/pkg/kubernetesclient"
	"github.com/thalassa-cloud/client-go/pkg/me"
)

type Client interface {
	IaaS() *iaas.Client
	Kubernetes() *kubernetesclient.Client
	Me() *me.Client
}

type thalassaCloudClient struct {
	client client.Client
}

// Option is a function that modifies the Client.
type Option func(*Client) error

// NewClient applies all options, configures authentication, and returns the client.
func NewClient(opts ...client.Option) (Client, error) {
	c, err := client.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &thalassaCloudClient{
		client: c,
	}, nil
}

func (c *thalassaCloudClient) SetOrganisation(organisation string) {
	c.client.SetOrganisation(organisation)
}

func (c *thalassaCloudClient) IaaS() *iaas.Client {
	iaasClient, err := iaas.New(c.client)
	if err != nil {
		panic(err)
	}
	return iaasClient
}

func (c *thalassaCloudClient) Kubernetes() *kubernetesclient.Client {
	kubernetesClient, err := kubernetesclient.New(c.client)
	if err != nil {
		panic(err)
	}
	return kubernetesClient
}

func (c *thalassaCloudClient) Me() *me.Client {
	meClient, err := me.New(c.client)
	if err != nil {
		panic(err)
	}
	return meClient
}
