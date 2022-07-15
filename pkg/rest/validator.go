package rest

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslate "github.com/go-playground/validator/v10/translations/en"
)

// StructValidator is a validator with automatic errors translations.
type StructValidator struct {
	uni   *ut.UniversalTranslator
	valid *validator.Validate
}

// NewStructValidator returns new validator with default english locale.
func NewStructValidator() StructValidator {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	valid := validator.New()
	valid.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	_ = entranslate.RegisterDefaultTranslations(valid, trans)
	_ = valid.RegisterTranslation(
		"fqdn",
		trans,
		registerTranslation("fqdn", "{0} must be valid FQDN", false),
		translateFunc,
	)
	_ = valid.RegisterTranslation(
		"uuid_rfc4122",
		trans,
		registerTranslation("uuid_rfc4122", "{0} must be valid UUID", false),
		translateFunc,
	)
	return StructValidator{
		uni:   uni,
		valid: valid,
	}
}

func registerTranslation(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, override); err != nil {
			return
		}
		return
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		return fe.(error).Error()
	}
	return t
}

// Validate validates structure according attached validation tags.
func (s StructValidator) Validate(ctx context.Context, value any) error {
	if err := s.valid.StructCtx(ctx, value); err != nil {
		if fieldsErr, ok := err.(validator.ValidationErrors); ok {
			trans, _ := s.uni.GetTranslator("en")
			// TODO(eartemov): Report all available errors not only
			// the first one. To implement this I need to change the
			// error format first.
			fErr := fieldsErr[0]
			return BadRequestErrorf("%s: %s", fErr.Namespace(), fErr.Translate(trans))
		}
		return err
	}
	return nil
}
