# About project
This project is created to simplify getting started with web based client server applications

## Installation
```bash
go get -u github.com/kubestaff/web-helper/server
```

## Getting started
```go
package main

import (
"fmt"
"github.com/kubestaff/web-helper/server"
"time"
)

func main() {
	opts := server.Options{
		Port: 34567,
	}
	// we create the simplified web server
	s := server.NewServer(opts)

	// we close the server at the end
	defer s.Stop()

	variables := map[string]string{"%name%": "Max Mustermann"}

	// we output the contents of index.html statically
	s.PrintFile("/", "html/index.html", variables)
	s.PrintFile("/status", "status.html", variables)
	//try URLs like /months?month=1 or just /month
	s.Handle("/months", HandleMonths)

	// we start the webserver don't put any code after it
	s.Start()
}

func HandleMonths(inputs map[string]string) (filename string, placeholders map[string]string) {
	months := map[string]string{
		"1":  "Jan",
		"2":  "Feb",
		"3":  "Mar",
		"4":  "Apr",
		"5":  "Mai",
		"6":  "Jun",
		"7":  "Jul",
		"8":  "Aug",
		"9":  "Sep",
		"10": "Oct",
		"11": "Nov",
		"12": "Dec",
	}

	givenMonthNumber, ok := inputs["month"]
	if !ok {
		givenMonthNumber = time.Now().Format("1")
	}

	output := map[string]string{
		"%month%": fmt.Sprintf("Number %s is month %s", givenMonthNumber, months[givenMonthNumber]),
	}

	return "html/month.html", output
}

```

This will start a webserver which will output this file to your browser. 
When started, have a look at the message `started at http://127.0.0.1:49765`. This will be a URL for your server. 