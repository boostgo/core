package translate

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"maps"
	"os"
	"path/filepath"

	"github.com/boostgo/core/fsx"
)

const (
	ExtJson = ".json"
	ExtYaml = ".yaml"
)

type Translator struct {
	path      string
	locales   []Locale
	extension string
	texts     map[Locale]map[string]string
}

func NewTranslator(
	path string,
	extension string,
	locales ...Locale,
) *Translator {
	texts := make(map[Locale]map[string]string)
	maps.Copy(texts, defaultTexts)

	return &Translator{
		path: path,
		locales: append([]Locale{
			LocaleEnglish,
			LocaleRussian,
			LocaleKazakh,
		}, locales...),
		extension: extension,
		texts:     texts,
	}
}

func (t *Translator) RegisterLocale(locale Locale) *Translator {
	t.locales = append(t.locales, locale)
	return t
}

func (t *Translator) Read() error {
	for _, locale := range t.locales {
		filePath := filepath.Join(t.path, locale.String()+t.extension)
		if !fsx.FileExist(filePath) {
			continue
		}

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return ErrReadFile.
				SetError(err).
				AddParam("path", filePath)
		}

		translations := make(map[string]string)

		switch t.extension {
		case ExtJson:
			if err = json.Unmarshal(fileContent, &translations); err != nil {
				return err
			}
		case ExtYaml:
			if err = yaml.Unmarshal(fileContent, &translations); err != nil {
				return err
			}
		}

		for key, text := range translations {
			if t.texts[locale] == nil {
				t.texts[locale] = make(map[string]string)
			}

			if _, ok := t.texts[locale][key]; ok {
				continue
			}

			t.texts[locale][key] = text
		}
	}

	return nil
}

func (t *Translator) TextByKey(locale Locale, key string) (string, error) {
	translations, ok := t.texts[locale]
	if !ok {
		return "", ErrLocaleNotFound.AddParam("locale", locale)
	}

	translation, ok := translations[key]
	if !ok {
		return "", ErrKeyNotFound.
			AddParam("key", key).
			AddParam("locale", locale)
	}

	return translation, nil
}
