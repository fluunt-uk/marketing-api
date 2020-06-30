package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gitlab.com/projectreferral/marketing-api/configs"
	"gitlab.com/projectreferral/util/pkg/security"
	"io/ioutil"
	"log"
	"net/http"
)

func SetupEndpoints() {
	_router := mux.NewRouter()

	_router.HandleFunc("/test", TestFunc)

	_router.HandleFunc("/apply", security.WrapHandlerWithSpecialAuth(Apply, configs.AUTH_AUTHENTICATED)).Methods("POST")
	_router.HandleFunc("/adverts", security.WrapHandlerWithSpecialAuth(CreateAdvert, configs.AUTH_AUTHENTICATED)).Methods("PUT")
	_router.HandleFunc("/adverts", security.WrapHandlerWithSpecialAuth(DeleteAdvert, configs.AUTH_AUTHENTICATED)).Methods("DELETE")
	_router.HandleFunc("/adverts", security.WrapHandlerWithSpecialAuth(UpdateAdvert, configs.AUTH_AUTHENTICATED)).Methods("PATCH")
	_router.HandleFunc("/adverts", security.WrapHandlerWithSpecialAuth(GetAdvert, configs.AUTH_AUTHENTICATED)).Methods("GET")
	_router.HandleFunc("/adverts/query", security.WrapHandlerWithSpecialAuth(GetBatchAdverts, "")).Methods("GET")
	_router.HandleFunc("/adverts/apply", security.WrapHandlerWithSpecialAuth(Apply, "")).Methods("POST")
	_router.HandleFunc("/log", displayLog).Methods("GET")

	c := cors.New(cors.Options{
		AllowedMethods: []string{"POST"},
		AllowedOrigins: []string{"*"},
		AllowCredentials: true,
		AllowedHeaders: []string{"g-recaptcha-response", "Authorization", "Content-Type","Origin","Accept", "Accept-Encoding", "Accept-Language", "Host", "Connection", "Referer", "Sec-Fetch-Mode", "User-Agent", "Access-Control-Request-Headers", "Access-Control-Request-Method: "},
		OptionsPassthrough: true,
	})

	handler := c.Handler(_router)
	log.Fatal(http.ListenAndServe(configs.PORT,handler))
}


func displayLog(w http.ResponseWriter, r *http.Request){
	b, err := ioutil.ReadFile(configs.LOG_PATH)

	if err != nil {
		fmt.Println(err.Error()) //output to main
		w.WriteHeader(http.StatusInternalServerError)
	}else{
		w.Write(b)
	}
}