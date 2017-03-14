package main

//go:generate sqlboiler mysql

import (
	"github.com/skizhak/echodemo/resources"
)

func main() {
	resources.RunServer(resources.ControllerService())
}
