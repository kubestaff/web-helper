package main

import "github.com/kubestaff/web-helper/server"

func main() {
	opts := server.Options{
		Port: 34567,
	}
	// we create the simplified web server
	s := server.NewServer(opts)

	// we close the server at the end
	defer s.Stop()

	variables := map[string]string{"%name%": "Max Mustermann"}

	// we output the contents of index.html
	s.PrintFile("/", "index.html", variables)
	s.PrintFile("/status", "status.html", variables)

	// we start the webserver don't put any code after it
	s.Start()
}
