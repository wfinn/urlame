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

var numberregex = regexp.MustCompile("\\d+(\\.\\d+)?")
var profilepageregex = regexp.MustCompile("(?i)/(u|user|profile|author|member|referral)s?/[^/]+/?")
var titleregex = regexp.MustCompile("[A-Za-z0-9-.]+")
var langregex = buildlangregex()
var uuidregex = regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}")
var hashregex = regexp.MustCompile("[a-zA-Z0-9]{32,40,64,128}")
var equivalenceregexes = buildquivalences()

var exts = []string{".css", ".png", ".jpg", ".jpeg", ".svg", ".gif", ".mp3", ".mp4", ".rss", ".ttf", ".woff", ".woff2", ".eot", ".pdf", ".m4v", ".ogv", ".webm"}
var paths = []string{"wp-content", "blog", "blogs", "product", "doc", "docs", "support"}

var langcodes = []string{"af", "af-ZA", "ar", "ar-AE", "ar-BH", "ar-DZ", "ar-EG", "ar-IQ", "ar-JO", "ar-KW", "ar-LB", "ar-LY", "ar-MA", "ar-OM", "ar-QA", "ar-SA", "ar-SY", "ar-TN", "ar-YE", "az", "az-AZ", "az-AZ", "be", "be-BY", "bg", "bg-BG", "bs-BA", "ca", "ca-ES", "cs", "cs-CZ", "cy", "cy-GB", "da", "da-DK", "de", "de-AT", "de-CH", "de-DE", "de-LI", "de-LU", "dv", "dv-MV", "el", "el-GR", "en", "en-AU", "en-BZ", "en-CA", "en-CB", "en-GB", "en-IE", "en-JM", "en-NZ", "en-PH", "en-TT", "en-US", "en-ZA", "en-ZW", "eo", "es", "es-AR", "es-BO", "es-CL", "es-CO", "es-CR", "es-DO", "es-EC", "es-ES", "es-ES", "es-GT", "es-HN", "es-MX", "es-NI", "es-PA", "es-PE", "es-PR", "es-PY", "es-SV", "es-UY", "es-VE", "et", "et-EE", "eu", "eu-ES", "fa", "fa-IR", "fi", "fi-FI", "fo", "fo-FO", "fr", "fr-BE", "fr-CA", "fr-CH", "fr-FR", "fr-LU", "fr-MC", "gl", "gl-ES", "gu", "gu-IN", "he", "he-IL", "hi", "hi-IN", "hr", "hr-BA", "hr-HR", "hu", "hu-HU", "hy", "hy-AM", "id", "id-ID", "is", "is-IS", "it", "it-CH", "it-IT", "ja", "ja-JP", "ka", "ka-GE", "kk", "kk-KZ", "kn", "kn-IN", "ko", "ko-KR", "kok", "kok-IN", "ky", "ky-KG", "lt", "lt-LT", "lv", "lv-LV", "mi", "mi-NZ", "mk", "mk-MK", "mn", "mn-MN", "mr", "mr-IN", "ms", "ms-BN", "ms-MY", "mt", "mt-MT", "nb", "nb-NO", "nl", "nl-BE", "nl-NL", "nn-NO", "ns", "ns-ZA", "pa", "pa-IN", "pl", "pl-PL", "ps", "ps-AR", "pt", "pt-BR", "pt-PT", "qu", "qu-BO", "qu-EC", "qu-PE", "ro", "ro-RO", "ru", "ru-RU", "sa", "sa-IN", "se", "se-FI", "se-FI", "se-FI", "se-NO", "se-NO", "se-NO", "se-SE", "se-SE", "se-SE", "sk", "sk-SK", "sl", "sl-SI", "sq", "sq-AL", "sr-BA", "sr-BA", "sr-SP", "sr-SP", "sv", "sv-FI", "sv-SE", "sw", "sw-KE", "syr", "syr-SY", "ta", "ta-IN", "te", "te-IN", "th", "th-TH", "tl", "tl-PH", "tn", "tn-ZA", "tr", "tr-TR", "tt", "tt-RU", "ts", "uk", "uk-UA", "ur", "ur-PK", "uz", "uz-UZ", "uz-UZ", "vi", "vi-VN", "xh", "xh-ZA", "zh", "zh-CN", "zh-HK", "zh-MO", "zh-SG", "zh-TW", "zu", "zu-zA"}

// Equivalences are explained in README.md
// Modify to include target specific words, which can be normalized to reduce results
// left side must contain a unique string of the pattern !-FOO-!, this is used internally as replacement
var equivalences = map[string][]string{
	//"!-TESLA-!" : {"model-3","model-y", "..."},
	// langcodes are too small and need to be treated seperately
}

var dotstar = regexp.MustCompile(".*")
var paramregexes = map[string]*regexp.Regexp{
	"utm_source":   dotstar,
	"utm_medium":   dotstar,
	"utm_campaign": dotstar,
	"utm_content":  dotstar,
	"utm_term":     dotstar,
	"redirect":     regexp.MustCompile("no"),
	//TODO version, v, cb, cache ...
}

func buildlangregex() *regexp.Regexp {
	// langcodes are currently only completely matched, might remove ^ and $ for the longer ones?
	reg := "(?i)^("
	for i, lang := range langcodes {
		reg += strings.Replace(lang, "-", "[-_]", 1) // match en-US & en_US
		if i < len(langcodes) {
			reg += "|"
		}
	}
	reg += ")$"
	return regexp.MustCompile(reg)
}

func buildquivalences() map[string]*regexp.Regexp {
	regexes := map[string]*regexp.Regexp{}
	for replacement, eqwords := range equivalences {
		regexes[replacement] = buildeqregex(eqwords)
	}
	return regexes
}

func buildeqregex(equivalentwords []string) *regexp.Regexp {
	reg := "("
	for i, lang := range langcodes {
		reg += strings.Replace(lang, "-", "[-_]", 1) // match en-US & en_US
		if i < len(langcodes) {
			reg += "|"
		}
	}
	reg += ")"
	return regexp.MustCompile(reg)
}

func main() {
	printNormalized := flag.Bool("print-normalized", false, "print the normalized version of the urls (for debugging)")
	flag.Usage = func() {
		fmt.Printf("%s [OPTIONS] < urls.txt > less_urls.txt\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	seen := map[string]bool{}
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		urlstr := stdin.Text()
		if u, err := url.Parse(urlstr); err == nil {
			if lamefiletype(u) || profilepage(u) || lamedir(u) {
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
		for i, lamepath := range paths {
			if i > 2 {
				//this is so we match /en-US/blog but not /api/v1/edit/blog
				return false
			}
			if lower == lamepath {
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
			if !lameparam(key, u.Query().Get(key)) {
				// ignoring lame params, if we see /foo and /foo?utm_source=bar we only list the first
				newvals.Set(normalizeItem(key), "!-P-!")
			}
		}
		return newURL(u, normalizePath(u.Path), newvals)
	}
	return urlstr
}

func lameparam(key, val string) bool {
	if paramregexes[key] != nil {
		return paramregexes[key].MatchString(val)
	}
	return false
}

func normalizePath(path string) string {
	normalized := ""
	split := strings.Split(path, "/")
	fileName := split[len(split)-1]
	split = split[:len(split)-1]
	for _, part := range split {
		if strings.TrimSpace(part) == "" {
			continue
		}
		normalized += "/" + normalizeItem(part)
	}
	if fileName != "" { // this makes /foo/bar?x=y and /foo/bar/?x=y equivalent
		// e.g. 123.json // FIXME: after the changes in normalizeItem are finished, this shouldn't needed anymore
		lastInd := strings.LastIndex(fileName, ".")
		if lastInd == -1 {
			normalized += "/" + normalizeItem(fileName)
		} else {
			normalized += "/" + normalizeItem(fileName[:lastInd]) + "." + fileName[lastInd+1:]
		}
	}
	return normalized
}

func normalizeItem(item string) string {
	if len(item) > 10 && titleregex.MatchString(item) {
		return "!-T-!"
	}
	orig := item
	item = hashregex.ReplaceAllString(item, "!-H-!")
	item = uuidregex.ReplaceAllString(item, "!-U-!")
	item = langregex.ReplaceAllString(item, "!-L-!")
	if orig == item {
		// only apply `numberregex` if hash / UUID wasn't found, might be too generic otherwise
		item = numberregex.ReplaceAllString(item, "!-N-!")
	}
	return item
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
