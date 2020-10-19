package configs

const (
	LOG_PATH         = "../logs/marketingAPI_log.txt"
	PORT = ":5003"

	/************** DynamoDB configs *************/
	EU_WEST_2        	 = "eu-west-2"
	TABLE_NAME        	= "advert"
	UNIQUE_IDENTIFIER 	= "id"
	/*********************************************/
	/*********** Authentication configs **********/
	AUTH_REGISTER      	= "new_advert"
	AUTH_AUTHENTICATED	= "crud"
	NO_ACCESS         	= "admin_gui"
	AUTHORIZED			= "Authorization"
	/*********************************************/
	/************** RabbitMQ configs *************/
	FANOUT_EXCHANGE 	= "advert.fanout"
	QAPI_URL 			= "http://35.179.11.178:5004"
	/*********************************************/
	/************** AccountAPI configs *************/
	ACCOUNT_API 		= "http://35.179.11.178:5001/account"
	/*********************************************/
)
var (
	/***************************************** CORS configs ****************************************/
	ALLOWED_HEADERS		= []string{"g-recaptcha-response",
		"Authorization", "Content-Type","Origin",
		"Accept", "Accept-Encoding", "Accept-Language",
		"Host", "Connection", "Referer", "Sec-Fetch-Mode",
		"User-Agent", "Access-Control-Request-Headers", "Access-Control-Request-Method: "}

	ALLOWED_METHODS		= []string{"POST", "PUT", "GET", "PATCH"}
	ALLOWED_ORIGINS		= []string{"*"}
	/***********************************************************************************************/

)

