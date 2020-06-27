package repo_builder

import (
	"gitlab.com/projectreferral/marketing-api/internal/models"
)

func (c *AdvertWrapper) UpdateValue(email string, cr *models.ChangeRequest) error{

	switch cr.Type {
	// string value
	case 1:
		return c.DC.UpdateStringField(cr.Field,email,cr.NewString)
		break
	// map value
	case 2:
		return c.DC.AppendNewMap(cr.Id, email, &cr.NewMap, cr.Field)
		break
		// string value
	case 3:
		return c.DC.UpdateBoolField(cr.Field,email,cr.NewBool)
		break
	}

	return nil
}
