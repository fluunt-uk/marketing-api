package repo_builder

import (
	"gitlab.com/projectreferral/util/pkg/dynamodb"
	"net/http"
)

//Custom made error
func HandleError(err error, w http.ResponseWriter, isCustom bool) bool {

	if err != nil {
		//we can return the error with specific response code and reason
		if isCustom {
			e := err.(*dynamodb.ErrorString)
			http.Error(w, e.Reason, e.Code)
			return true
		}

		//default error
		http.Error(w, err.Error(), 400)
		return true
	}
	return false
}
