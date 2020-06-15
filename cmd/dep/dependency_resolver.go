package dep

import (
	"gitlab.com/projectreferral/marketing-api/lib/dynamodb/repo-builder"
	"gitlab.com/projectreferral/marketing-api/lib/rabbitmq"
	"gitlab.com/projectreferral/util/client"
	"gitlab.com/projectreferral/util/pkg/dynamodb"
	"log"
)

//methods that are implemented on util
//and will be used
type ConfigBuilder interface{
	LoadEnvConfigs()
	LoadDynamoDBConfigs() *dynamodb.Wrapper
	LoadRabbitMQConfigs() *client.DefaultQueueClient
}

//internal specific configs are loaded at runtime
//takes in a object(implemented interface) of type ServiceConfigs
func Inject(builder ConfigBuilder) {

	//load the env into the object
	builder.LoadEnvConfigs()

	//setup dynamo library
	//TODO:shall the dynamo configs injected here? or in the main?
	dynamoClient := builder.LoadDynamoDBConfigs()
	//connect to the instance
	log.Println("Connecting to dynamo client")
	dynamoClient.DefaultConnect()

	//dependency injection to our resource
	//we inject the dynamo client
	//shared client, therefore shared in between all the repos
	LoadAdvertRepo(&repo_builder.AdvertWrapper{
		DC: dynamoClient,
	})

	rabbitMQClient := builder.LoadRabbitMQConfigs()

	LoadRabbitMQClient(rabbitMQClient)
}

//variable injected with the interface methods
func LoadAdvertRepo (r repo_builder.AdvertBuilder){
	log.Println("Injecting Advert Repo")
	repo_builder.Advert = r
}

func LoadRabbitMQClient(c client.QueueClient){
	log.Println("Injecting RabbitMQ Client")
	rabbitmq.Client = c
}

