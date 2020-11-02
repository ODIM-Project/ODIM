package redfish

import (
	"strings"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
)

type Translator struct {
	Dictionaries *config.URLTranslation
}

func (u *Translator) ODIMToRedfish(data string) string {
	translated := data
	for k, v := range u.Dictionaries.SouthBoundURL {
		translated = strings.Replace(translated, k, v, -1)
	}
	return translated
}

func (u *Translator) RedfishToODIM(data string) string {
	translated := data
	for k, v := range u.Dictionaries.NorthBoundURL {
		translated = strings.Replace(translated, k, v, -1)
	}
	return translated
}
