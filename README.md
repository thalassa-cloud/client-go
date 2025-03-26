# client-go for the Thalassa Cloud Platform API

> **Note**: This client is currently in alpha state. While we strive to maintain stability, breaking changes may occur as we continue to develop and improve the client.

This is the official Go client library for interacting with the Thalassa Cloud Platform API. It provides a simple and efficient way to integrate Thalassa Cloud services into your Go applications.

## Documentation

For detailed documentation, including API references, guides, and examples, visit [docs.thalassa.cloud](https://docs.thalassa.cloud).

## Installation

To use this client in your Go project, add it to your dependencies:

```bash
go get github.com/thalassa-cloud/go-client
```

## Quick Start

Here's a simple example of how to use the client:

```go
package main

import (
    "context"
    "log"
    "github.com/thalassa-cloud/go-client"
)

func main() {
    client, err := thalassa.NewClient(thalassa.WithAuthPersonalToken("your-api-key"))
    if err != nil {
        log.Fatal(err)
    }
    ctx := context.Background()
    
    // Use the client to interact with Thalassa Cloud services
    // Example: List your vpcs
    vpcs, err := client.ListVpcs(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    // Process vpcs...
}
```

## Contributing

We welcome contributions! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
