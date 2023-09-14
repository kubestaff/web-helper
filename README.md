# About project
This project is created to simplify getting started with web based client server applications

## Installation
```bash
go get -u github.com/kubestaff/web-helper/server
```

## Getting started
```go
package main

import "github.com/kubestaff/web-helper/server"

func main() {
	s := server.NewServer()

	defer s.Stop()

	s.PrintFile("index.html")

	s.Start()
}
```

This will start a webserver which will output this file to your browser. 
When started, have a look at the message `started at http://127.0.0.1:49765`. This will be a URL for your server. 