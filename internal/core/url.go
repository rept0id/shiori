package core

import (
	"fmt"
	nurl "net/url"
	"sort"
	"strings"
)

// queryEncodeWithoutEmptyValues is a copy of `values.Encode` but checking if the queryparam
// value is empty to prevent sending the = symbol empty which breaks in some servers.
func queryEncodeWithoutEmptyValues(v nurl.Values) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		keyEscaped := nurl.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			if v != "" {
				buf.WriteByte('=')
				buf.WriteString(nurl.QueryEscape(v))
			}
		}
	}
	return buf.String()
}

// Parse parses a URL. If it fails, it tries again with "https://" prefix.
func Parse(url string) (*nurl.URL, error) {
	urlParsed, err := nurl.Parse(url)
	if err != nil || urlParsed.Scheme == "" || urlParsed.Hostname() == "" {
		if strings.HasPrefix(url, "https://") {
			return nil, fmt.Errorf("URL is not valid")
		}

		return Parse("https://" + url)
	}

	return urlParsed, err
}

// RemoveUTMParams removes the UTM parameters from URL.
func RemoveUTMParams(url string) (string, error) {
	// Parse string URL
	tmp, err := Parse(url)
	if err != nil {
		return url, err
	}

	// Remove UTM queries
	queries := tmp.Query()
	for key := range queries {
		if strings.HasPrefix(key, "utm_") {
			queries.Del(key)
		}
	}

	tmp.RawQuery = queryEncodeWithoutEmptyValues(queries)
	return tmp.String(), nil
}
