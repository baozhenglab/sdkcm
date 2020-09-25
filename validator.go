package sdkcm

import (
	"log"
	"reflect"
	"regexp"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/iancoleman/strcase"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

var trans ut.Translator
var validate *validator.Validate

// Execute validate form
func ExecuteValidator(form interface{}) error {
	return validate.Struct(form)
}

//Load Register validator for variable Validator and load custom error
func LoadValidator() {
	translator := en.New()
	validate = validator.New()
	uni := ut.New(translator, translator)

	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	var found bool
	trans, found = uni.GetTranslator("en")

	if !found {
		log.Fatal("translator not found")
	}
	registerItemValidator(validate)
	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Fatal(err)
	}
	registerTranslationValidator(validate)

}

//GetErrors used to Format message error into map key string
func GetErrors(err validator.ValidationErrors) (errors map[string][]string) {

	errors = map[string][]string{}
	for _, errItem := range err {
		if _, ok := errors[errItem.Field()]; ok {
			errors[strcase.ToSnake(errItem.Field())] = append(errors[errItem.Field()], errItem.Translate(trans))
		} else {
			errors[strcase.ToSnake(errItem.Field())] = []string{errItem.Translate(trans)}
		}

	}
	return errors
}

func registerItemValidator(Validator *validator.Validate) {
	// Validator.RegisterValidation("umail", func(fl validator.FieldLevel) bool {
	// 	email := fl.Field().String()
	// 	if email == "" {
	// 		return true
	// 	}
	// 	user := models.User{}
	// 	err := orm.FindOneByQuery(&user, map[string]interface{}{"email": email})
	// 	if err != nil {
	// 		return gorm.IsRecordNotFoundError(err)
	// 	}
	// 	return false

	// })

	Validator.RegisterValidation("isbase64", func(fl validator.FieldLevel) bool {
		base64 := Base64(fl.Field().String())
		if base64 == "" {
			return true
		}
		return base64.IsBase64(OptionsIsBase64{MimeRequired: true})
	})

	Validator.RegisterValidation("isphone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		r := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
		return r.MatchString(phone)
	})

	Validator.RegisterValidation("iscode", func(fl validator.FieldLevel) bool {
		code := fl.Field().String()
		r := regexp.MustCompile(`[A-Z\d]+$`)
		return r.MatchString(code)
	})

	Validator.RegisterValidation("mineimage", func(fl validator.FieldLevel) bool {
		base64 := fl.Field().String()
		if base64 == "" {
			return true
		}
		r, _ := regexp.Compile(`^data:image\/([\s\S]+);base64,([\s\S]+)`)
		matches := r.FindAllStringSubmatch(base64, -1)
		if len(matches) != 1 {
			return false
		}
		allowType := []string{"png", "jpeg", "jpg", "svg"}
		if len(matches[0]) == 3 {
			for _, v := range allowType {
				if v == matches[0][1] {
					return true
				}
			}
		}
		return false
	})

	Validator.RegisterValidation("integer", func(fl validator.FieldLevel) bool {
		if fl.Field().Kind() == reflect.Int || fl.Field().Kind() == reflect.Uint {
			return true
		}
		value := fl.Field().String()
		if value == "" {
			return true
		}
		r, _ := regexp.Compile(`^[0-9]*$`)
		return r.MatchString(value)
	})

	Validator.RegisterValidation("uinteger", func(fl validator.FieldLevel) bool {
		if fl.Field().Kind() == reflect.Uint {
			return true
		}
		value := fl.Field().String()
		if value == "" {
			return true
		}
		r, _ := regexp.Compile(`^[0-9]*$`)
		return r.MatchString(value)
	})
	Validator.RegisterValidation("integer32", func(fl validator.FieldLevel) bool {
		if fl.Field().Kind() == reflect.Int32 {
			return true
		}
		value := fl.Field().String()
		if value == "" {
			return true
		}
		r, _ := regexp.Compile(`^[0-9]*$`)
		return r.MatchString(value)
	})

	Validator.RegisterValidation("date", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}
		r, _ := regexp.Compile(`^\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d(?:\.\d+)?(?:Z|\+[0-2]\d(?:\:[0-5]\d)?)?$`)
		return r.MatchString(value)
	})
	Validator.RegisterValidation("mine", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		mine := fl.Param()
		if value == "" {
			return true
		}
		r := regexp.MustCompile(`^data:([\s\S]+);base64,([\s\S]+)`)
		matches := r.FindAllStringSubmatch(value, -1)
		if len(matches) != 1 {
			return false
		} else if len(matches[0]) != 3 {
			return false
		} else if matches[0][1] != mine {
			return false
		}
		return true

	})
	Validator.RegisterValidation("passw", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}
		r, _ := regexp.Compile(`^(.*[A-Z].*)(.*[a-z].*)(.*\d.*)`)
		return r.MatchString(value)
	})
}

func registerTranslationValidator(Validator *validator.Validate) {
	Validator.RegisterTranslation("umail", trans, func(ut ut.Translator) error {
		return ut.Add("umail", "{0} is already exists", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("umail", fe.Field())
		return t
	})
	Validator.RegisterTranslation("iscode", trans, func(ut ut.Translator) error {
		return ut.Add("iscode", "{0} is incorrect", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("iscode", fe.Field())
		return t
	})
	Validator.RegisterTranslation("isphone", trans, func(ut ut.Translator) error {
		return ut.Add("isphone", "{0} is incorrect", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("isphone", fe.Field())
		return t
	})
	Validator.RegisterTranslation("isbase64", trans, func(ut ut.Translator) error {
		return ut.Add("isbase64", "{0} is not base64", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("isbase64", fe.Field())
		return t
	})
	Validator.RegisterTranslation("mineimage", trans, func(ut ut.Translator) error {
		return ut.Add("mineimage", "{0} is not mine type base64", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("mineimage", fe.Field())
		return t
	})
	Validator.RegisterTranslation("date", trans, func(ut ut.Translator) error {
		return ut.Add("date", "{0} is not type date string", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("date", fe.Field())
		return t
	})
	Validator.RegisterTranslation("integer", trans, func(ut ut.Translator) error {
		return ut.Add("integer", "{0} is not type integer", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("integer", fe.Field())
		return t
	})
	Validator.RegisterTranslation("uinteger", trans, func(ut ut.Translator) error {
		return ut.Add("uinteger", "{0} is not type uint", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("uinteger", fe.Field())
		return t
	})
	Validator.RegisterTranslation("integer32", trans, func(ut ut.Translator) error {
		return ut.Add("integer32", "{0} is not type integer32", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("integer32", fe.Field())
		return t
	})

	Validator.RegisterTranslation("passw", trans, func(ut ut.Translator) error {
		return ut.Add("passw", "Password not correct", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("passw", fe.Field())
		return t
	})
	Validator.RegisterTranslation("mine", trans, func(ut ut.Translator) error {
		return ut.Add("mine", "{0} required mine type {1}", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("mine", fe.Field(), fe.Param())
		return t
	})
}
