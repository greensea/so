package common

import (
	"os"
	"strings"

	"github.com/Xuanwo/go-locale"
)

// Current language, only 'zh' or 'en' now
var lang string

// Get current language. Only return 'zh' or 'en' now
func Lang() string {
	if lang == "" {

		langTag, err := locale.Detect()
		if err != nil {
			return "en"
		}

		langBase, _, _ := langTag.Raw()
		if langBase.String() == "zh" {
			lang = "zh"
		} else {
			lang = "en"
		}
	}

	return lang

}

func SoEndpoint(path string) string {
	endpoint := os.Getenv("SO_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://api-so.pingflash.com"
	}
	endpoint = strings.TrimSuffix(endpoint, "/")

	return endpoint + path
}
