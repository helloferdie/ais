package libvalidator

import (
	"ais/lib/libresponse"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	// Register validator to detect struct with tag "json"
	validate.RegisterTagNameFunc(func(f reflect.StructField) string {
		json := f.Tag.Get("json")
		loc := f.Tag.Get("loc")
		if strings.HasSuffix(loc, ".") {
			loc += json
		}
		return json + "|" + loc
	})
}

// Validate - Validate obj request and return response with error
func Validate(obj interface{}) (*libresponse.Response, error) {
	resp := libresponse.GetDefault()
	err := validate.Struct(obj)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return resp.ErrorValidation(), err
		}

		// Loop error list
		tmp := map[string]libresponse.Error{}
		for _, err := range err.(validator.ValidationErrors) {
			tag := strings.Split(err.Field(), "|")
			json := tag[0]
			loc := tag[1]

			e := libresponse.Error{}
			e.Error = "validation.error." + err.Tag()
			tmp[json] = e

			if resp.Error == "" {
				resp.Error = e.Error + "_var"
				resp.ErrorVar = map[string]interface{}{
					"var": "." + loc,
				}
			}
		}

		resp.Code = 422
		resp.Message = "validation.error.input"
		if len(tmp) > 0 {
			resp.Data = tmp
		}
	}
	return resp, err
}
