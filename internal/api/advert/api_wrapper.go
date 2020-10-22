package advert

import (
	repo_builder "gitlab.com/projectreferral/marketing-api/lib/dynamodb/repo-builder"
	"gitlab.com/projectreferral/marketing-api/lib/rabbitmq"
	"net/http"
)

func TestFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)

	rabbitmq.BroadcastNewAdvert([]byte("Hello from Marketing Service"))
}

//We check for the recaptcha response and proceed
//Covert the response body into appropriate models
//Create a new user using our dynamodb adapter
//A event message it sent to the queues which are consumed by the relevant services
func Create(w http.ResponseWriter, r *http.Request) {

	repo_builder.Advert.CreateAdvert(w,r)
}

func Delete(w http.ResponseWriter, r *http.Request) {

	repo_builder.Advert.DeleteAdvert(w,r)
}

func Get(w http.ResponseWriter, r *http.Request) {

	repo_builder.Advert.GetAdvert(w,r)
}

func Update(w http.ResponseWriter, r *http.Request) {

	repo_builder.Advert.UpdateAdvert(w,r)
}

func GetBatch(w http.ResponseWriter, r *http.Request) {

	repo_builder.Advert.GetBatchAdvert(w,r)
}

func Apply(w http.ResponseWriter, r *http.Request) {

	repo_builder.Advert.Apply(w,r)
}

