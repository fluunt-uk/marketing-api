package api

import (
	"github.com/gorilla/mux"
	"gitlab.com/projectreferral/marketing-api/configs"
	"gitlab.com/projectreferral/util/pkg/security"
	"log"
	"net/http"
)

func SetupEndpoints() {
	_router := mux.NewRouter()

	_router.HandleFunc("/test", TestFunc)

	_router.HandleFunc("/adverts", security.WrapHandlerWithSpecialAuth(CreateAdvert, configs.AUTH_AUTHENTICATED)).Methods("PUT")
	_router.HandleFunc("/adverts", security.WrapHandlerWithSpecialAuth(DeleteAdvert, configs.AUTH_AUTHENTICATED)).Methods("DELETE")
	_router.HandleFunc("/adverts", security.WrapHandlerWithSpecialAuth(UpdateAdvert, configs.AUTH_AUTHENTICATED)).Methods("PATCH")
	_router.HandleFunc("/adverts", security.WrapHandlerWithSpecialAuth(GetAdvert, configs.AUTH_AUTHENTICATED)).Methods("GET")
	_router.HandleFunc("/adverts/query", security.WrapHandlerWithSpecialAuth(GetBatchAdverts, "")).Methods("GET")

	log.Fatal(http.ListenAndServe(configs.PORT, _router))
}
