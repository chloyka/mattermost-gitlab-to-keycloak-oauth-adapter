package main

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	config := newAppConfig()
	log, err := NewDefaultAppLogger(config.Logger.Level)
	if err != nil {
		panic(err)
	}
	enableDebug := log.Level == LoggerLevelDebug

	keycloakUrl, err := url.Parse(config.Keycloak.Host)
	if err != nil {
		log.Error("error parsing keycloak url", zap.String("url", config.Keycloak.Host), zap.Error(err))
		os.Exit(1)
	}

	server := http.NewServeMux()
	rp := httputil.NewSingleHostReverseProxy(keycloakUrl)

	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w = newResponseWriterWithStatusCode(w, enableDebug)
		defer log.Info("request", zap.String("method", r.Method), zap.String("url", r.URL.String()), zap.Int("status", w.(*responseWriterWithStatusCode).statusCode))

		if r.URL.Path == "/oauth/authorize" {
			w.Header().Set("Location", fmt.Sprintf("%s%s?%s", config.Keycloak.Host, config.Keycloak.AuthUrl, r.URL.RawQuery))
			w.WriteHeader(http.StatusTemporaryRedirect)
			log.Debug("requesting authorize", zap.String("url", r.URL.String()), zap.String("redirectUrl", w.Header().Get("Location")))
			return
		}

		r.Host = keycloakUrl.Host
		r.URL.Host = keycloakUrl.Host
		r.URL.Scheme = keycloakUrl.Scheme

		if r.URL.Path == "/oauth/token" {
			r.URL.Path = config.Keycloak.TokenUrl

			log.Debug("requesting token", zap.String("url", r.URL.String()))
			rp.ServeHTTP(w, r)
			log.Debug("token response", zap.String("url", r.URL.String()), zap.Int("status", w.(*responseWriterWithStatusCode).statusCode), zap.String("body", string(w.(*responseWriterWithStatusCode).body)))

			w.WriteHeader(http.StatusOK)
			return
		}

		if r.URL.Path == "/api/v4/user" {
			r.URL.Path = config.Keycloak.UserInfoUrl

			log.Debug("requesting user info", zap.String("url", r.URL.String()))
			rp.ServeHTTP(w, r)
			log.Debug("user info response", zap.String("url", r.URL.String()), zap.Int("status", w.(*responseWriterWithStatusCode).statusCode), zap.String("body", string(w.(*responseWriterWithStatusCode).body)))

			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		return
	})

	err = http.ListenAndServe(":80", server)
	if err != nil {
		panic(err)
	}
}
