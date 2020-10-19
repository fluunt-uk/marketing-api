package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gitlab.com/projectreferral/marketing-api/configs"
	"gitlab.com/projectreferral/marketing-api/internal/api/advert"
	"gitlab.com/projectreferral/util/pkg/security"
	"io/ioutil"
	"log"
	"net/http"
)

func SetupEndpoints() {
	_router := mux.NewRouter()

	_router.HandleFunc("/test", advert.TestFunc)

	_router.HandleFunc("/apply", security.WrapHandlerWithSpecialAuth(advert.Apply, configs.AUTH_AUTHENTICATED)).Methods("POST")
	_router.HandleFunc("/advert", security.WrapHandlerWithSpecialAuth(advert.Create, configs.AUTH_AUTHENTICATED)).Methods("PUT")
	_router.HandleFunc("/advert", security.WrapHandlerWithSpecialAuth(advert.Delete, configs.AUTH_AUTHENTICATED)).Methods("DELETE")
	_router.HandleFunc("/advert", security.WrapHandlerWithSpecialAuth(advert.Update, configs.AUTH_AUTHENTICATED)).Methods("PATCH")
	_router.HandleFunc("/advert", security.WrapHandlerWithSpecialAuth(advert.Get, configs.AUTH_AUTHENTICATED)).Methods("GET")
	_router.HandleFunc("/advert/query", security.WrapHandlerWithSpecialAuth(advert.GetBatch, "")).Methods("GET")
	_router.HandleFunc("/advert/apply", security.WrapHandlerWithSpecialAuth(advert.Apply, "")).Methods("POST")
	_router.HandleFunc("/log", displayLog).Methods("GET")

	handler := withCORS().Handler(_router)
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

func withCORS() *cors.Cors{
	c := cors.New(cors.Options{
		AllowedMethods: configs.ALLOWED_METHODS,
		AllowedOrigins: configs.ALLOWED_ORIGINS,
		AllowCredentials: true,
		AllowedHeaders: configs.ALLOWED_HEADERS,
		OptionsPassthrough: true,
	})

	return c
}