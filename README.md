# aptTransport

Go API for creating transports for Debian's `apt` package manager(and it's relatives, etc)

## Usage:

By default, the transport will use Go's default HTTP client, which will result in a transport which works for HTTP and HTTPS clearnet clients.
This client is as simple as:

```Go
package main

import "github.com/eyedeekay/apttransport"

func main() {
	transport := &apttransport.AptMethod{}
	transport.Main = transport.DefaultMain
	transport.Main()
}
```

You can also hook into the process by creating a custom client.
A client must implement the `AptClient` interface, which must have a `Get` function which returns an `*http.Response`.
Of course, an HTTP Client works, here is an example which uses a SOCKS5 proxy:

```Go
package main

import "github.com/eyedeekay/apttransport"

func main() {
	transport := &apttransport.AptMethod{}
	transport.Main = transport.DefaultMain
    proxy, err := SOCKS5("tcp", "127.0.0.1:9050", nil, nil)
    if err != nil {
        panic(err)
    }
    transport.Client := &http.Client{
		Timeout: time.Duration(6) * time.Minute,
		Transport: &http.Transport{
			Dial:                  proxy.Dial,
		},
		CheckRedirect: nil,
	}
    transport.Main()
}
```

That is simply the type used for the response.
The content of the `Get` function must compose an `*http.Response` to return.
TODO: example which constructs an HTTP response from a bittorrent download.
