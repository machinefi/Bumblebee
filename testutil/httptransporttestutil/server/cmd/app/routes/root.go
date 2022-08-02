package routes

import (
	"github.com/iotexproject/Bumblebee/kit/httptransport"
	"github.com/iotexproject/Bumblebee/kit/httptransport/swagger"
	"github.com/iotexproject/Bumblebee/kit/kit"
)

var RootRouter = kit.NewRouter(httptransport.BasePath("/demo"))

func init() {
	RootRouter.Register(swagger.Router)
}
