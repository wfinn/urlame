// urlame -
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

// maybe numbers should be surrounded by special chars, or be at least a certain amount of digits?
var numberregex = regexp.MustCompile("\\d+(\\.\\d+)?")
var profilepageregex = regexp.MustCompile("(?i)/(u|user|profile|author|member|referral)s?/[^/]+/?")
var titleregex = regexp.MustCompile("^(/[^/]+)?/[A-Za-z0-9.]+-[A-Za-z0-9.]+-[A-Za-z0-9.\\-]+$")
var langregex = buildlangregex()
var uuidregex = regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}")
var hashregex = regexp.MustCompile("[a-zA-Z0-9]{32,40,64,128}")
var equivalenceregexes = buildquivalences()

// This is one of the areas where urlame is very opinionated
// This is what *I* consider lame
var exts = []string{".js", ".css", ".png", ".jpg", ".jpeg", ".svg", ".gif", ".ico", ".bmp", ".rss", ".mp3", ".mp4", ".ttf", ".woff", ".woff2", ".eot", ".pdf", ".m4v", ".ogv", ".webm"}
var paths = []string{"static", "wp-content", "blog", "blogs", "product", "doc", "docs", "support", "news", "article", "fonts"}

// .html here means it matches with many static file exts (staticexts)
var files = []string{"index.html", "robots.txt", "contact.html", "home.html", "impressum.html"}
var staticexts = regexp.MustCompile("^\\.(htm|html|php|cgi)$")

var langcodes = []string{"af", "af-ZA", "ar", "ar-AE", "ar-BH", "ar-DZ", "ar-EG", "ar-IQ", "ar-JO", "ar-KW", "ar-LB", "ar-LY", "ar-MA", "ar-OM", "ar-QA", "ar-SA", "ar-SY", "ar-TN", "ar-YE", "az", "az-AZ", "az-AZ", "be", "be-BY", "bg", "bg-BG", "bs-BA", "ca", "ca-ES", "cs", "cs-CZ", "cy", "cy-GB", "da", "da-DK", "de", "de-AT", "de-CH", "de-DE", "de-LI", "de-LU", "dv", "dv-MV", "el", "el-GR", "en", "en-AU", "en-BZ", "en-CA", "en-CB", "en-GB", "en-IE", "en-JM", "en-NZ", "en-PH", "en-TT", "en-US", "en-ZA", "en-ZW", "eo", "es", "es-AR", "es-BO", "es-CL", "es-CO", "es-CR", "es-DO", "es-EC", "es-ES", "es-ES", "es-GT", "es-HN", "es-MX", "es-NI", "es-PA", "es-PE", "es-PR", "es-PY", "es-SV", "es-UY", "es-VE", "et", "et-EE", "eu", "eu-ES", "fa", "fa-IR", "fi", "fi-FI", "fo", "fo-FO", "fr", "fr-BE", "fr-CA", "fr-CH", "fr-FR", "fr-LU", "fr-MC", "gl", "gl-ES", "gu", "gu-IN", "he", "he-IL", "hi", "hi-IN", "hr", "hr-BA", "hr-HR", "hu", "hu-HU", "hy", "hy-AM", "id", "id-ID", "is", "is-IS", "it", "it-CH", "it-IT", "ja", "ja-JP", "ka", "ka-GE", "kk", "kk-KZ", "kn", "kn-IN", "ko", "ko-KR", "kok", "kok-IN", "ky", "ky-KG", "lt", "lt-LT", "lv", "lv-LV", "mi", "mi-NZ", "mk", "mk-MK", "mn", "mn-MN", "mr", "mr-IN", "ms", "ms-BN", "ms-MY", "mt", "mt-MT", "nb", "nb-NO", "nl", "nl-BE", "nl-NL", "nn-NO", "ns", "ns-ZA", "pa", "pa-IN", "pl", "pl-PL", "ps", "ps-AR", "pt", "pt-BR", "pt-PT", "qu", "qu-BO", "qu-EC", "qu-PE", "ro", "ro-RO", "ru", "ru-RU", "sa", "sa-IN", "se", "se-FI", "se-FI", "se-FI", "se-NO", "se-NO", "se-NO", "se-SE", "se-SE", "se-SE", "sk", "sk-SK", "sl", "sl-SI", "sq", "sq-AL", "sr-BA", "sr-BA", "sr-SP", "sr-SP", "sv", "sv-FI", "sv-SE", "sw", "sw-KE", "syr", "syr-SY", "ta", "ta-IN", "te", "te-IN", "th", "th-TH", "tl", "tl-PH", "tn", "tn-ZA", "tr", "tr-TR", "tt", "tt-RU", "ts", "uk", "uk-UA", "ur", "ur-PK", "uz", "uz-UZ", "uz-UZ", "vi", "vi-VN", "xh", "xh-ZA", "zh", "zh-CN", "zh-HK", "zh-MO", "zh-SG", "zh-TW", "zu", "zu-zA"}

// TODO instead of ignoring these params when comparing, how about just removing them? shouldn't affect anything and be much better
var dotstar = regexp.MustCompile(".*")
var paramregexes = map[string]*regexp.Regexp{
	"utm_source":   dotstar,
	"utm_medium":   dotstar,
	"utm_campaign": dotstar,
	"utm_content":  dotstar,
	"utm_term":     dotstar,
	"redirect":     regexp.MustCompile("no"),
	// TODO version, v, cb, cache, ref=[usernameregex], gclid, fbclid, aid, referrer=[usernameregex,urlregex], affiliate
	// Tracking: _hsenc, _hsmi, __hssc, __hstc, hsCtaTracking, msclkid, mkt_tok, yclid, yadclid  ...
}

func buildlangregex() *regexp.Regexp {
	// langcodes are currently only completely matched, might remove ^ and $ for the longer ones?
	reg := "(?i)^("
	for i, lang := range langcodes {
		reg += strings.Replace(lang, "-", "[-_]", 1) // match en-US & en_US
		if i < len(langcodes)-1 {
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
	for i, word := range equivalentwords {
		reg += strings.Replace(word, "-", "[-_]", 1)
		if i < len(equivalentwords)-1 {
			reg += "|"
		}
	}
	reg += ")"
	return regexp.MustCompile(reg)
}

// --- setup of vars ends ---

func main() {
	printNormalized := flag.Bool("print-normalized", false, "print the normalized version of the urls (for debugging)")
	flag.Usage = func() {
		fmt.Printf("%s [OPTIONS] < urls.txt > less_urls.txt\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	runurlame(os.Stdin, os.Stdout, *printNormalized)
}

func runurlame(input io.Reader, output io.Writer, printNormalized bool) error {
	seen := map[string]bool{}
	stdin := bufio.NewScanner(input)
	for stdin.Scan() {
		urlstr := stdin.Text()
		if u, err := url.Parse(urlstr); err == nil && len(urlstr) > 1 {
			if lamefiletype(u) || profilepage(u) || lamedir(u) {
				//skip those that we can be certain are lame
				continue
			}
			normalized := normalizeURL(urlstr)
			if seen[normalized] {
				continue
			} else {
				seen[normalized] = true
				seen[urldecode(normalized)] = true
			}
			if printNormalized {
				fmt.Fprintf(output, "%s\n", normalized)
			} else {
				fmt.Fprintf(output, "%s\n", urlstr)
			}
		}
	}
	return nil
}

// detects lame file extensions
func lamefiletype(u *url.URL) bool {
	filetype := strings.ToLower(path.Ext(u.Path))
	for _, ext := range exts {
		if filetype == ext {
			return true
		}
	}
	return false
}

// detects if one of the first path segments is lame
func lamedir(u *url.URL) bool {
	for i, part := range strings.Split(u.Path, "/") {
		lower := strings.ToLower(part)
		for _, lamepath := range paths {
			if i > 2 {
				//this is so we match /en-US/blog but not /api/v1/edit/blog
				return false
			}
			if lower == lamepath {
				return true
			}
		}
	}
	return titleregex.MatchString(u.Path)
}

// sees if u.Path matches profilepageregex
func profilepage(u *url.URL) bool {
	if profilepageregex.MatchString(u.Path) {
		return true
	}
	return false
}

// ...
func normalizeURL(urlstr string) string {
	//this func needs the original string instead of a url as it may return it unchanged
	if u, err := url.Parse(urlstr); err == nil {
		newvals := url.Values{}
		for key := range u.Query() {
			if !lameparam(key, u.Query().Get(key)) {
				// ignoring lame params, if we see /foo and /foo?utm_source=bar we only list the first
				newvals.Set(normalizeItem(key), "!-P-!")
				//TODO, replace !-P-! with a pattern of inputs
			}
		}
		return newURL(u, normalizePath(u.Path), newvals)
	}
	return urlstr
}

// detects lame parameter names, like utm_source
func lameparam(key, val string) bool {
	if paramregexes[key] != nil {
		return paramregexes[key].MatchString(val)
	}
	return false
}

// splits path into segments and normalizes them, the last segment (filename) is a special case
func normalizePath(path string) string {
	normalized := ""
	split := strings.Split(strings.TrimRight(path, "/"), "/")
	file := split[len(split)-1] // TODO filepath.Split exists
	segments := split[:len(split)-1]
	for _, part := range segments {
		if strings.TrimSpace(part) == "" {
			continue
		}
		normalized += "/" + normalizeItem(part)
	}
	// normalFilename is feature to remove common filenames, e.g. contact.html robots.txt etc
	// By doing it here and not just blocking them we don't miss any directory or host, but we can ignore lame files
	// (this avoids  filtering in cases where https://neverseenbefore.host/robots.txt is the only URL of neverseenbefore.host for example)
	return normalized + "/" + normalFilename(file)
}

// removes lame filenames, so we can ignore them safely
func normalFilename(filename string) string {
	for _, lamefile := range files {
		if lamefile == filename || staticexts.ReplaceAllString(filename, ".html") == lamefile {
			return ""
		}
	}
	return normalizeItem(filename)
}

// normalizes a single item, e.g. path segment
func normalizeItem(item string) string {
	orig := item
	item = applyequivalences(item)
	item = hashregex.ReplaceAllString(item, "!-H-!")
	item = uuidregex.ReplaceAllString(item, "!-U-!")
	item = langregex.ReplaceAllString(item, "!-L-!")
	if orig == item {
		// only apply `numberregex` if hash / UUID wasn't found, might be too generic otherwise
		item = numberregex.ReplaceAllString(item, "!-N-!")
	}
	return item
}

// experimental feature to define target specific "equivalent" words, see equivalences.go
func applyequivalences(item string) string {
	for replacement, regex := range equivalenceregexes {
		item = regex.ReplaceAllString(item, "!-"+replacement+"-!")
	}
	return item
}

// this builds the URL "pattern" (/ normalized URL) we use for comparison internally
func newURL(old *url.URL, path string, vals url.Values) string {
	return /*ignore scheme*/ cleanHostname(old) + path + "?" + vals.Encode() + "#" + old.Fragment
}

// this function removes default port info, so http://foo & http://foo:80 are equal, but http://foo:123 is different
func cleanHostname(u *url.URL) string {
	if u.Port() == "80" || u.Port() == "443" {
		return u.Hostname()
	}
	return u.Host
}

func urldecode(str string) string {
	if decoded, err := url.QueryUnescape(str); err == nil {
		return decoded
	}
	return str
}
