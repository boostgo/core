package echo

import (
	"errors"

	"github.com/boostgo/core/echox"
	"github.com/boostgo/core/errorx"
	"github.com/boostgo/core/translate"

	"github.com/labstack/echo/v4"
)

func FailureMiddleware(
	translator *translate.Translator,
	localeHeaderName string,
) echox.FailureMiddleware {
	return func(ctx echo.Context, statusCode int, err error) {
		var converted *errorx.Error
		if !errors.As(err, &converted) {
			return
		}

		locale := translate.Locale(echox.Header(ctx, localeHeaderName).String(translate.LocaleRussian.String()))
		text, err := translator.TextByKey(locale, converted.Message())
		if err != nil {
			return
		}

		_ = converted.NoCopy()
		_ = converted.SetLocaleMessage(text)
	}
}
