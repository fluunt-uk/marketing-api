package repo_builder

import (
	"gitlab.com/projectreferral/marketing-api/internal/models"
)

func (a *AdvertWrapper) UpdateValue(email string, cr *models.ChangeRequest) error{

	switch cr.Type {
	// string value
	case 1:
		return a.DC.UpdateStringField(cr.Field,email,cr.NewString)
		break
	// map value
	case 2:
		return a.DC.AppendNewMap(cr.Id, email, &cr.NewMap, cr.Field)
		break
		// string value
	case 3:
		return a.DC.UpdateBoolField(cr.Field,email,cr.NewBool)
		break
	}

	return nil
}
