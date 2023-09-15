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
	opts := server.Options{}
	// we create the simplified web server
	s := server.NewServer(opts)

	// we close the server at the end
	defer s.Stop()

	variables := map[string]string{"%name%": "Max Mustermann"}

	// we output the contents of index.html
	s.PrintFile("index.html", variables)

	// we start the webserver don't put any code after it
	s.Start()
}
```

This will start a webserver which will output this file to your browser. 
When started, have a look at the message `started at http://127.0.0.1:49765`. This will be a URL for your server. 