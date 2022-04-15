package main

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		longPath, prs := pathsToUrls[request.RequestURI]
		if prs {
			http.Redirect(writer, request, longPath, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(writer, request)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type PathAndURL struct {
	path string `yaml:"path"`
	url  string `yaml:"url"`
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var t []PathAndURL
	err := yaml.Unmarshal(yml, &t)
	if err != nil {
		return nil, err
	}
	fmt.Println("array is ", t)
	urlMap := make(map[string]string)
	for i := 0; i < len(t); i++ {
		urlMap[t[i].path] = t[i].url
	}
	fmt.Println("map", urlMap)
	return MapHandler(urlMap, fallback), nil
}
