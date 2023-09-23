package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	config := newAppConfig()

	keycloakUrl, err := url.Parse(config.Keycloak.Host)
	if err != nil {
		fmt.Sprintln("Error parsing keycloak url", config.Keycloak.Host, "cannot be parsed")
		os.Exit(1)
	}

	server := http.NewServeMux()
	rp := httputil.NewSingleHostReverseProxy(keycloakUrl)

	server.HandleFunc("/oauth/authorize", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Header().Set("Location", fmt.Sprintf("%s?%s", config.Keycloak.AuthUrl, r.URL.RawQuery))
		return
	})

	server.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		r.Host = keycloakUrl.Host
		r.URL.Host = keycloakUrl.Host
		r.URL.Scheme = keycloakUrl.Scheme
		r.URL.Path = config.Keycloak.TokenUrl

		rp.ServeHTTP(w, r)
	})

	server.HandleFunc("/api/v4/user", func(w http.ResponseWriter, r *http.Request) {
		r.Host = keycloakUrl.Host
		r.URL.Host = keycloakUrl.Host
		r.URL.Scheme = keycloakUrl.Scheme
		r.URL.Path = config.Keycloak.UserInfoUrl

		rp.ServeHTTP(w, r)
	})

	err = http.ListenAndServe(":80", server)
	if err != nil {
		panic(err)
	}
}
