# Robtex API

A simple [Robtex API](https://www.robtex.com/api/) client.

## Examples

### Free JSON API

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/electrologue/robtex"
)

func main() {
	client := robtex.New("")

	ctx := context.Background()

	ipQueryResponse, err := client.IPQuery(ctx, "199.19.54.1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", ipQueryResponse)

	asQueryResponse, err := client.ASQuery(ctx, "1234")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", asQueryResponse)

	pDNSForward, err := client.PassiveDNSForward(ctx, "a.iana-servers.net")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", pDNSForward)

	pDNSReverse, err := client.PassiveDNSReverse(ctx, "199.43.132.53")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", pDNSReverse)
}
```

### Pro JSON API

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/electrologue/robtex"
)

func main() {
	client := robtex.New("my_api_key")

	ctx := context.Background()

	ipQueryResponse, err := client.IPQuery(ctx, "199.19.54.1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", ipQueryResponse)

	asQueryResponse, err := client.ASQuery(ctx, "1234")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", asQueryResponse)

	pDNSForward, err := client.PassiveDNSForward(ctx, "a.iana-servers.net")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", pDNSForward)

	pDNSReverse, err := client.PassiveDNSReverse(ctx, "199.43.132.53")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", pDNSReverse)
}
```
