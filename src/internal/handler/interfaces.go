package handler

import "net/http"

type CountryHandler interface {
	ConvertCountry(w http.ResponseWriter, r *http.Request)
	Health(w http.ResponseWriter, r *http.Request)
}
