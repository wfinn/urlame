package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

var numberregex = regexp.MustCompile("^\\d+$")
var profilepageregex = regexp.MustCompile("^/(u|user|profile|author|member)s?/[^/]+/?$")
var titleregex = regexp.MustCompile("^[A-Za-z0-9-.]+$")
var langregex = regexp.MustCompile("^(af|af-ZA|ar|ar-AE|ar-BH|ar-DZ|ar-EG|ar-IQ|ar-JO|ar-KW|ar-LB|ar-LY|ar-MA|ar-OM|ar-QA|ar-SA|ar-SY|ar-TN|ar-YE|az|az-AZ|az-AZ|be|be-BY|bg|bg-BG|bs-BA|ca|ca-ES|cs|cs-CZ|cy|cy-GB|da|da-DK|de|de-AT|de-CH|de-DE|de-LI|de-LU|dv|dv-MV|el|el-GR|en|en-AU|en-BZ|en-CA|en-CB|en-GB|en-IE|en-JM|en-NZ|en-PH|en-TT|en-US|en-ZA|en-ZW|eo|es|es-AR|es-BO|es-CL|es-CO|es-CR|es-DO|es-EC|es-ES|es-ES|es-GT|es-HN|es-MX|es-NI|es-PA|es-PE|es-PR|es-PY|es-SV|es-UY|es-VE|et|et-EE|eu|eu-ES|fa|fa-IR|fi|fi-FI|fo|fo-FO|fr|fr-BE|fr-CA|fr-CH|fr-FR|fr-LU|fr-MC|gl|gl-ES|gu|gu-IN|he|he-IL|hi|hi-IN|hr|hr-BA|hr-HR|hu|hu-HU|hy|hy-AM|id|id-ID|is|is-IS|it|it-CH|it-IT|ja|ja-JP|ka|ka-GE|kk|kk-KZ|kn|kn-IN|ko|ko-KR|kok|kok-IN|ky|ky-KG|lt|lt-LT|lv|lv-LV|mi|mi-NZ|mk|mk-MK|mn|mn-MN|mr|mr-IN|ms|ms-BN|ms-MY|mt|mt-MT|nb|nb-NO|nl|nl-BE|nl-NL|nn-NO|ns|ns-ZA|pa|pa-IN|pl|pl-PL|ps|ps-AR|pt|pt-BR|pt-PT|qu|qu-BO|qu-EC|qu-PE|ro|ro-RO|ru|ru-RU|sa|sa-IN|se|se-FI|se-FI|se-FI|se-NO|se-NO|se-NO|se-SE|se-SE|se-SE|sk|sk-SK|sl|sl-SI|sq|sq-AL|sr-BA|sr-BA|sr-SP|sr-SP|sv|sv-FI|sv-SE|sw|sw-KE|syr|syr-SY|ta|ta-IN|te|te-IN|th|th-TH|tl|tl-PH|tn|tn-ZA|tr|tr-TR|tt|tt-RU|ts|uk|uk-UA|ur|ur-PK|uz|uz-UZ|uz-UZ|vi|vi-VN|xh|xh-ZA|zh|zh-CN|zh-HK|zh-MO|zh-SG|zh-TW|zu|zu-ZA)$")
var uuidregex = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")
var hashregex = regexp.MustCompile("^[a-zA-Z0-9]+$")
var hashlens = []int{32, 40, 64, 128}

var exts = []string{".css", ".png", ".jpg", ".jpeg", ".svg", ".gif", ".mp3", ".mp4", ".rss", ".ttf", ".woff", ".woff2", ".eot", ".pdf"}
var paths = []string{"static", "assets", "wp-content", "blog", "blogs", "product", "doc", "docs", "support"}

func main() {
	printNormalized := flag.Bool("print-normalized", false, "print the normalized version of the urls (for debugging)")
	blockPaths := flag.Bool("block-paths", false, "block common paths like /static, /wp-content")
	flag.Usage = func() {
		fmt.Printf("cat urls.txt | %s [OPTIONS]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	seen := map[string]bool{}
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		urlstr := stdin.Text()
		if u, err := url.Parse(urlstr); err == nil {
			if lamefiletype(u) || profilepage(u) || (*blockPaths && lamedir(u)) {
				//skip those that we can be certain are lame
				continue
			}
			normalized := normalizeURL(urlstr)
			if seen[normalized] {
				continue
			} else {
				seen[normalized] = true
			}
			if *printNormalized {
				fmt.Println(normalized)
			} else {
				fmt.Println(urlstr)
			}
		}
	}
}

func lamefiletype(u *url.URL) bool {
	filetype := strings.ToLower(path.Ext(u.Path))
	for _, ext := range exts {
		if filetype == ext {
			return true
		}
	}
	return false
}

func lamedir(u *url.URL) bool {
	for _, part := range strings.Split(u.Path, "/") {
		lower := strings.ToLower(part)
		for _, path := range paths {
			if lower == path {
				return true
			}
		}
	}
	return false
}

func profilepage(u *url.URL) bool {
	if profilepageregex.MatchString(u.Path) {
		return true
	}
	return false
}

func normalizeURL(urlstr string) string {
	//this func needs the original string instead of a url as it may return it unchanged
	if u, err := url.Parse(urlstr); err == nil {
		newvals := url.Values{}
		for key := range u.Query() {
			newvals.Set(key, "!-P-!")
		}
		return newURL(u, normalizePath(u.Path), newvals)
	}
	return urlstr
}

func normalizePath(path string) string {
	normalized := ""
	for _, part := range strings.Split(path, "/") {
		if strings.TrimSpace(part) == "" {
			continue
		}
		normalized += "/" + normalizeItem(part)
	}
	return normalized
}

func normalizeItem(item string) string {
	// it's unlikely that we have urls with !-X-! in them which we would miss here

	if numberregex.MatchString(item) {
		return "!-N-!"
	} else if postitle(item) {
		return "!-T-!"
	} else if hash(item) {
		return "!-H-!"
	} else if langcode(item) {
		return "!-L-!"
	} else if uuid(item) {
		return "!-U-!"
	}
	return item
}

func postitle(str string) bool {
	if !titleregex.MatchString(str) {
		return false
	}
	if len(str) > 10 {
		return true
	}
	return strings.Count(str, "-") > 2
}

func hash(str string) bool {
	if hashregex.MatchString(str) {
		strlen := len(str)
		for _, l := range hashlens {
			if strlen == l {
				return true
			}
		}
	}
	return false
}

func langcode(str string) bool {
	return langregex.MatchString(str)
}

func uuid(str string) bool {
	return uuidregex.MatchString(str)
}

func newURL(old *url.URL, path string, vals url.Values) string {
	return /*ignore scheme*/ cleanHostname(old) + path + "?" + vals.Encode() + "#" + old.Fragment
}

func cleanHostname(u *url.URL) string {
	if u.Port() == "80" || u.Port() == "443" {
		return u.Hostname()
	}
	return u.Host
}
