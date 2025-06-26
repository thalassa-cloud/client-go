# client-go for the Thalassa Cloud Platform API

> **Note**: This client is currently in alpha state. While we strive to maintain stability, breaking changes may occur as we continue to develop and improve the client.

This is the official Go client library for interacting with the Thalassa Cloud Platform API. It provides a simple and efficient way to integrate Thalassa Cloud services into your Go applications.

## Documentation

For detailed documentation, including API references, guides, and examples, visit [docs.thalassa.cloud](https://docs.thalassa.cloud).

## Installation

To use this client in your Go project, add it to your dependencies:

```bash
go get github.com/thalassa-cloud/client-go
```

## Quick Start

Here's a simple example of how to use the client:

```go
package main

import (
    "context"
    "log"

    "github.com/thalassa-cloud/client-go/pkg/client"
    "github.com/thalassa-cloud/client-go/pkg/thalassa"
)

func main() {
    baseURL := "https://<YOUR BASE URL FOR THE THALASSA API>"

    client, err := thalassa.NewClient(
		client.WithBaseURL(baseURL),
		client.WithOrganisation("organisation-slug-or-identity"),
        client.WithAuthPersonalToken("your-api-key"))
    if err != nil {
        log.Fatal(err)
    }
    ctx := context.Background()
    
    // Use the client to interact with Thalassa Cloud services
    // Example: List your vpcs
    vpcs, err := client.IaaS().ListVpcs(ctx, &iaas.ListVpcsRequest{})
    if err != nil {
        log.Fatal(err)
    }
    // Process vpcs...
}
```

## Examples

### Infrastructure as a Service (IaaS)

```go
// List all VPCs
vpcs, err := client.IaaS().ListVpcs(ctx)
if err != nil {
    log.Fatal(err)
}

for _, vpc := range vpcs {
    fmt.Printf("VPC ID: %s, Name: %s\n", vpc.ID, vpc.Name)
}

// Get details of a specific VPC
vpc, err := client.IaaS().GetVpc(ctx, "vpc-id")
if err != nil {
    log.Fatal(err)
}
```

### Kubernetes Service

```go
// List all Kubernetes clusters
clusters, err := client.Kubernetes().ListKubernetesClusters(ctx)
if err != nil {
    log.Fatal(err)
}

for _, cluster := range clusters {
    fmt.Printf("Cluster Name: %s\n", cluster.Name)
}

// Get details for a specific Kubernetes cluster
cluster, err := client.Kubernetes().GetKubernetesCluster(ctx, "cluster-id")
if err != nil {
    log.Fatal(err)
}
```

### User and Organization Management

```go
// Get information about your organizations
meClient := client.Me()
organizations, err := meClient.ListMyOrganisations(ctx)
if err != nil {
    log.Fatal(err)
}

for _, org := range organizations {
    fmt.Printf("Organization: %s\n", org.Name)
}
```

### Using the Alternative Client Approach

You can also initialize the client components separately:

```go
// Initialize the base client
baseClient, err := client.NewClient(
    client.WithBaseURL(baseURL),
    client.WithOrganisation(organisation),
    client.WithAuthPersonalToken(pat),
)
if err != nil {
    log.Fatal(err)
}

// Create the IaaS service
iaasService, err := iaas.New(baseClient)
if err != nil {
    log.Fatal(err)
}

// Use the IaaS service
vpcs, err := iaasService.ListVpcs(ctx)
if err != nil {
    log.Fatal(err)
}
```

## Contributing

We welcome contributions! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the [MIT License](/LICENSE).
