package repo_builder

import (
	"encoding/json"
	"fmt"
	"gitlab.com/projectreferral/marketing-api/configs"
	"gitlab.com/projectreferral/marketing-api/internal/models"
	"gitlab.com/projectreferral/marketing-api/lib/rabbitmq"
	"gitlab.com/projectreferral/util/pkg/dynamodb"
	"gitlab.com/projectreferral/util/pkg/http_lib"
	_ "gitlab.com/projectreferral/util/pkg/http_lib"
	"gitlab.com/projectreferral/util/pkg/security"
	_ "gitlab.com/projectreferral/util/pkg/security"
	"io/ioutil"
	"net/http"
	"net/url"
	_ "os"
	"strconv"
)

type AdvertWrapper struct {
	//dynamo client
	DC *dynamodb.Wrapper
}

//implement only the necessary methods for each repository
//available to be consumed by the API
type AdvertBuilder interface {
	GetAdvert(http.ResponseWriter, *http.Request)
	GetBatchAdvert(http.ResponseWriter, *http.Request)
	UpdateAdvert(http.ResponseWriter, *http.Request)
	CreateAdvert(http.ResponseWriter, *http.Request)
	DeleteAdvert(http.ResponseWriter, *http.Request)
	Apply(http.ResponseWriter, *http.Request)
}

//interface with the implemented methods will be injected in this variable
var Advert AdvertBuilder

// get advertid, check if the ad exists, update user with adverts applied, update user count for advert
//TODO:
// Check if advert is premium or not, if it is then the user (if not premium) should not be allowed to apply to it,
// the option should be greyed out.
func (a *AdvertWrapper) Apply(w http.ResponseWriter, r *http.Request) {

	var ad models.Advert // id from body
	var ap models.Application

	// Get id from body
	errDecode := dynamodb.DecodeToMap(r.Body, &ad)

	if !HandleError(errDecode, w, false) {

		// Check if ad exists
		i, err := a.DC.GetItem(ad.Uuid)

		_ = dynamodb.Unmarshal(i, &ad)
		_ = dynamodb.Unmarshal(i, &ap)

		if !HandleError(err, w, true) {

			appli_err := addApplicant(ad, a)

			if appli_err != nil {
				fmt.Println(appli_err)
				return
			}

			b, _ := json.Marshal(models.ChangeRequest{
				Field:  "applications",
				Id:     ap.Uuid,
				NewMap: ap,
				Type:   2,
			})

			// Make a patch request to accounts to update applications
			res, _ := http_lib.Patch(configs.ACCOUNT_API, b, map[string]string{configs.AUTHORIZED: r.Header.Get(configs.AUTHORIZED)})

			if res.StatusCode == 200 {

				userCount := convertToInt(w, &ad.UserCount)

				if userCount < convertToInt(w, &ad.MaxUsers) {
					newUserCount := userCount + 1

					cr := &models.ChangeRequest{
						Field:     "users_applied",
						NewString: strconv.Itoa(newUserCount),
						Type:      1,
					}

					// Update adverts to increase count by 1
					err := a.UpdateValue(ad.Uuid, cr)

					if !HandleError(err, w, false) {

						b, je := json.Marshal(&cr)
						if je != nil {
							http.Error(w, "Error parsing body", http.StatusBadRequest)
							return
						}

						w.Write([]byte("Success with payload :"))
						w.Write(b)
						w.WriteHeader(http.StatusOK)
						return
					}
				}
				//reached max users for advert
				w.WriteHeader(http.StatusPreconditionFailed)
				w.Write([]byte("Validation failed: User Count exceeds Max Users"))
				return
			}

			errorBody, errParse := ioutil.ReadAll(res.Body)

			if errParse != nil {
				http.Error(w, "Error parsing body", http.StatusBadRequest)
				return
			}

			http.Error(w, string(errorBody), res.StatusCode)
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	return
}

func addApplicant(ad models.Advert, a *AdvertWrapper) error {

	// Get user that is applying to ad
	var applyingUser models.User
	email := security.GetClaimsOfJWT().Subject
	t, _ := a.DC.GetItem(email)
	_ = dynamodb.Unmarshal(t, applyingUser)

	err := a.DC.AppendNewMap(applyingUser.Uuid, ad.Uuid, applyingUser, "applicants")

	if err != nil {
		fmt.Println("Failed to add applicant")
		return err
	}
	fmt.Println("Success: applicant added to ad.")
	return nil
}

//We check for the recaptcha response and proceed
//Covert the response body into appropriate models
//Create a new user using our dynamodb adapter
//A event message it sent to the queues which are consumed by the relevant services
func (a *AdvertWrapper) CreateAdvert(w http.ResponseWriter, r *http.Request) {
	var ad models.Advert

	dynamoAttr, errDecode := dynamodb.DecodeToDynamoAttribute(r.Body, &ad)

	if !HandleError(errDecode, w, false) {

		err := a.DC.CreateItem(dynamoAttr)

		if !HandleError(err, w, false) {
			w.WriteHeader(http.StatusOK)

			b, _ := json.Marshal(&ad)
			go rabbitmq.BroadcastNewAdvert(b)
		}
	}
}

func (a *AdvertWrapper) DeleteAdvert(w http.ResponseWriter, r *http.Request) {
	var am models.Advert

	id := GetQueryString(r.URL.Query(), "id", w)

	if id != "" {
		errDelete := a.DC.DeleteItem(id)

		if !HandleError(errDelete, w, false) {

			//Check item still exists
			result, err := a.DC.GetItem(am.Uuid)

			//error thrown, record not found
			if !HandleError(err, w, true) {
				http.Error(w, result.GoString(), 302)
			}
		}
	}
}

func (a *AdvertWrapper) GetAdvert(w http.ResponseWriter, r *http.Request) {
	var am models.Advert

	dynamodb.DecodeToMap(r.Body, &am)

	id := GetQueryString(r.URL.Query(), "id", w)

	if id != "" {
		//TODO:perhaps better to get from query string
		result, err := a.DC.GetItem(id)

		if !HandleError(err, w, true) {
			dynamodb.Unmarshal(result, &am)

			b, err := json.Marshal(&am)

			if !HandleError(err, w, false) {

				w.Write(b)
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}

//Creating a new user with same ID replaces the record
//Temporary solution
func (a *AdvertWrapper) UpdateAdvert(w http.ResponseWriter, r *http.Request) {
	var cr models.ChangeRequest

	advertID := GetQueryString(r.URL.Query(), "id", w)

	if advertID != "" {

		err := a.UpdateValue(advertID, &cr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetQueryString(m url.Values, q string, w http.ResponseWriter) string {

	idKeys, ok := m[q]

	if !ok {
		w.Write([]byte("Url Param are missing"))
		w.WriteHeader(http.StatusBadRequest)
		return ""
	}

	advertID := idKeys[0]
	if len(advertID) < 1 {
		w.Write([]byte("Url Param are missing"))
		w.WriteHeader(http.StatusBadRequest)
		return ""
	}

	return advertID
}

func (a *AdvertWrapper) GetBatchAdvert(w http.ResponseWriter, r *http.Request) {
	var am models.Advert

	dynamodb.DecodeToMap(r.Body, &am)

	by := GetQueryString(r.URL.Query(), "by", w)

	if by != "" {
		result, err := a.DC.GetAll(by)

		if !HandleError(err, w, true) {

			b, err := json.Marshal(&result)

			if !HandleError(err, w, false) {

				w.Write(b)
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}

func convertToInt(w http.ResponseWriter, value *string) int {

	//TODO: Change data type to INT in DynamoDB
	i, err := strconv.Atoi(*value)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	return i
}
