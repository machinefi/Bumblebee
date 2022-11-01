package routes

import (
	"github.com/machinefi/Bumblebee/kit/httptransport"
	"github.com/machinefi/Bumblebee/kit/httptransport/swagger"
	"github.com/machinefi/Bumblebee/kit/kit"
)

var RootRouter = kit.NewRouter(httptransport.BasePath("/demo"))

func init() {
	RootRouter.Register(swagger.Router)
}
