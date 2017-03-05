package main

import (
	"github.com/skizhak/echodemo/resources"
)

func main() {
	resources.RunServer(&resources.Controller{})
}
