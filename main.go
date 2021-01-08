package main

import (
	"sandcard/application"
	"sandcard/controller"
)

func main() {
	var (
		a = application.New()
		c = controller.New(a)
	)
	c.Serve(57428)
}
