package translate

type Locale string

func (l Locale) String() string {
	return string(l)
}

func (l Locale) Is(locales ...Locale) bool {
	for _, locale := range locales {
		if l == locale {
			return true
		}
	}

	return false
}

const (
	LocaleRussian Locale = "ru"
	LocaleEnglish Locale = "en"
	LocaleKazakh  Locale = "kk"
)
