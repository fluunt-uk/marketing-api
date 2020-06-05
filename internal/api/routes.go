package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.com/projectreferral/marketing-api/configs"
	"gitlab.com/projectreferral/util/pkg/security"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	log.Fatal(http.ListenAndServe(configs.PORT, _router))
}


func displayLog(w http.ResponseWriter, r *http.Request){

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	b, _ := ioutil.ReadFile(path + "/logs/marketingAPI_log.txt")

	w.Write(b)
}
