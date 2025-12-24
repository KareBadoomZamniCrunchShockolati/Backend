package middleware

import (
	domainLoc "challenge-app/internal/domain/localization"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const TranslatorContextKey = "translator"

type LocalizationMiddleware struct {
	translator domainLoc.Translator
}

func NewLocalizationMiddleware(translator domainLoc.Translator) *LocalizationMiddleware {
	return &LocalizationMiddleware{translator: translator}
}

func (lm *LocalizationMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := getLocale(c.Request)
		tr := lm.translator.GetTranslator(locale)
		c.Set(TranslatorContextKey, tr)
		c.Next()
	}
}

func getLocale(r *http.Request) string {
	lang := r.Header.Get("Accept-Language")
	lang = strings.TrimSpace(lang)
	if lang == "" {
		return "en_US"
	}

	lang = strings.ReplaceAll(lang, "-", "_")
	lang = strings.Split(lang, ",")[0]
	lang = strings.Split(lang, ";")[0]
	lang = strings.TrimSpace(lang)

	switch {
	case strings.HasPrefix(strings.ToLower(lang), "fa"):
		return "fa_IR"
	case strings.HasPrefix(strings.ToLower(lang), "en"):
		return "en_US"
	default:
		return "en_US"
	}
}