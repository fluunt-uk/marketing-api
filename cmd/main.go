package main

import (
	"gitlab.com/projectreferral/marketing-api/cmd/dep"
	"gitlab.com/projectreferral/marketing-api/configs"
	"gitlab.com/projectreferral/marketing-api/internal/api"
	"gitlab.com/projectreferral/marketing-api/internal/models"
	"gitlab.com/projectreferral/util/util"
	"os"
)

func main() {

	//gets all the necessary configs into our object
	//completes connections
	//assigns connections to repos
	dep.Inject(&util.ServiceConfigs{
		Environment: 	os.Getenv("ENV"),
		Region:       	configs.EU_WEST_2,
		Table:        	configs.TABLE_NAME,
		SearchParam:  	configs.UNIQUE_IDENTIFIER,
		GenericModel: 	models.Advert{},
		Port:		  	configs.PORT,
		BrokerUrl: 		configs.QAPI_URL,
	})

	api.SetupEndpoints()
}
