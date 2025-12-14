package main

import (
	"github.com/MhmdEagel/lms-usti-be/env"
	"github.com/MhmdEagel/lms-usti-be/model"
	"github.com/MhmdEagel/lms-usti-be/router"
)

func main() {
	model.ConnectDatabase()
	r := router.InitRouter()
	r.Run(env.DEFAULT_PORT)

}
