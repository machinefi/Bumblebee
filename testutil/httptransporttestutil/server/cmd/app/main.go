package main

import (
	"github.com/machinefi/Bumblebee/kit/httptransport"
	"github.com/machinefi/Bumblebee/kit/kit"
	"github.com/machinefi/Bumblebee/testutil/httptransporttestutil/server/cmd/app/routes"
)

func main() {
	ht := &httptransport.HttpTransport{
		Port: 8080,
	}
	ht.SetDefault()

	kit.Run(routes.RootRouter, ht)
}
