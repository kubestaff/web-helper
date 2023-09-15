package main

import "github.com/kubestaff/web-helper/server"

func main() {
	// we create the simplified web server
	s := server.NewServer()

	// we close the server at the end
	defer s.Stop()

	variables := map[string]string{"%name%": "Max Mustermann"}

	// we output the contents of index.html
	s.PrintFile("index.html", variables)

	// we start the webserver don't put any code after it
	s.Start()
}
