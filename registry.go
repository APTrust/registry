package main

import "github.com/APTrust/registry/app"

// Run the application.
// The app itself moved to app/application.go so we can launch it
// from our test suite as well.
func main() {
	app.Run()
}
