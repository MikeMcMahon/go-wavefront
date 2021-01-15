# Golang Wavefront Client

Golang SDK for interacting with the Wavefront v2 API, and sending metrics through a Wavefront proxy.

## Usage

### API Client

Presently support for:
 * Querying
 * Searching
 * Dashboard Management
 * Alert (and Alert Target) Management
 * Events Management

Please see the [examples](examples) directory for an example on how to use each, or check out the [documentation](https://godoc.org/github.com/WavefrontHQ/go-wavefront-management-api).

```Go
package main

import (
	"fmt"

	wavefront "github.com/WavefrontHQ/go-wavefront-management-api"
)

func main() {
	client, err := wavefront.NewClient(
		&wavefront.Config{
			Address: "test.wavefront.com",
			Token:   "xxxx-xxxx-xxxx-xxxx-xxxx",
		},
	)
	if err != nil {
		panic(err)
	}

	query := client.NewQuery(
		wavefront.NewQueryParams(`ts("cpu.load.1m.avg", dc=dc1)`),
	)

	result, err := query.Execute()
	if err != nil {
		panic(err)
	}

	fmt.Println(result.TimeSeries[0].Label)
	fmt.Println(result.TimeSeries[0].DataPoints[0])
}
```

### Writer

Writer has full support for metric tagging etc.

Again, see [examples](examples) for a more detailed explanation.

```Go
package main

import (
    "log"
    "os"

    wavefront "github.com/WavefrontHQ/go-wavefront-management-api/writer"
)

func main() {
    source, _ := os.Hostname()

    wf, err := wavefront.NewWriter("wavefront-proxy.example.com", 2878, source, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer wf.Close()

    wf.Write(wavefront.NewMetric("something.very.good.count", 33))
}
```

## Contributing

Pull requests are welcomed.

If you'd like to contribute to this project, please raise an issue and indicate that you'd like to take on the work prior to submitting a pull request.
