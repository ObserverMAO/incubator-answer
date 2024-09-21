package notificationcommonlang

import (
	"github.com/apache/incubator-answer/plugin"
)

type Lang string

func GetLanguage(lang string) Lang {
	switch lang {
	case "zh_CN":
		return LanguageZhCN
	case "en_US":
		return LanguageEnUS
	default:
		return LanguageZhCN
	}
}

type LangSupportService interface {
	TranslateGetCardInfo() string
	TranslateDescription(mixinNotificationMsg *plugin.NotificationMessage) string
	GetLangType() Lang
}

type LangPicker struct {
	pool map[Lang]LangSupportService
}

func NewLangPicker() *LangPicker {
	langPicker := &LangPicker{
		pool: make(map[Lang]LangSupportService),
	}
	langPicker.register(newZhCN())
	langPicker.register(newEnUS())

	return langPicker
}

func (ap *LangPicker) register(a LangSupportService) {
	ap.pool[a.GetLangType()] = a
}

func (ap *LangPicker) Pick(typ Lang) LangSupportService {
	if a, ok := ap.pool[typ]; ok {
		return a
	}
	return ap.pool[LanguageZhCN]
}
