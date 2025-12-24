package localization

import (
	"fmt"
	"strings"
	"sync"

	domainLoc "challenge-app/internal/domain/localization"
	"github.com/go-playground/locales/en_US"
	"github.com/go-playground/locales/fa_IR"
	ut "github.com/go-playground/universal-translator"
)

var translationMap = make(map[string]map[string]string)

type translationService struct {
	universalTranslator *ut.UniversalTranslator
	mu                  sync.RWMutex
}

func NewTranslationService() domainLoc.Translator {
	service := &translationService{
		universalTranslator: createUniversalTranslator(),
	}
	service.loadAndAddTranslations()
	return service
}

func (t *translationService) GetTranslator(locale string) domainLoc.TranslatorInstance {
	t.mu.RLock()
	defer t.mu.RUnlock()
	normalized := normalizeLocale(locale)
	translator, found := t.universalTranslator.GetTranslator(normalized)
	if !found {
		translator, _ = t.universalTranslator.GetTranslator("fa_IR") 
	}
	return &translatorInstanceImpl{translator: translator}
}

type translatorInstanceImpl struct {
	translator ut.Translator
}

func (t *translatorInstanceImpl) Translate(key string, params ...string) (string, error) {
	return t.translator.T(key, params...)
}

func (t *translatorInstanceImpl) Locale() string {
	return t.translator.Locale()
}

func createUniversalTranslator() *ut.UniversalTranslator {
	en := en_US.New()
	fa := fa_IR.New()
	return ut.New(en, en, fa)
}

func (t *translationService) loadAndAddTranslations() {
	addTranslations("fa_IR", Persian, t.universalTranslator)
	addTranslations("en_US", English, t.universalTranslator)
}

func addTranslations(locale string, translations map[string]interface{}, universalTranslator *ut.UniversalTranslator) {
	translator, found := universalTranslator.GetTranslator(locale)
	if !found {
		panic(fmt.Errorf("translator for locale %s not found", locale))
	}

	flattenedTranslations := loadTranslations(locale, translations)
	for key, translation := range flattenedTranslations {
		_ = translator.Add(key, translation, true)
	}
}

func loadTranslations(locale string, translations map[string]interface{}) map[string]string {
	if translations, ok := translationMap[locale]; ok {
		return translations
	}
	flattenedTranslations := make(map[string]string)
	flattenMap("", translations, flattenedTranslations)
	translationMap[locale] = flattenedTranslations
	return flattenedTranslations
}

func flattenMap(prefix string, input map[string]interface{}, output map[string]string) {
	for k, v := range input {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}
		switch value := v.(type) {
		case map[string]interface{}:
			flattenMap(fullKey, value, output)
		case string:
			output[fullKey] = value
		default:
		}
	}
}

func normalizeLocale(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "fa_IR"
	}

	first := raw
	if idx := strings.Index(first, ","); idx >= 0 {
		first = first[:idx]
	}
	if idx := strings.Index(first, ";"); idx >= 0 {
		first = first[:idx]
	}
	first = strings.TrimSpace(first)
	first = strings.ReplaceAll(first, "-", "_")

	low := strings.ToLower(first)
	if low == "fa" || strings.HasPrefix(low, "fa_") {
		return "fa_IR"
	}
	if low == "en" || strings.HasPrefix(low, "en_") {
		return "en_US"
	}
	return "fa_IR"
}
