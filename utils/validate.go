package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
)

func ValidateRequest[T interface{}](r *http.Request, w http.ResponseWriter, isPatch bool) (req T, ok bool) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorUnableToReadRequestBody())
		ok = false
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorUnableToParseRequestBody())
		ok = false
		return
	}

	rt := reflect.TypeOf(req)

	errs := make([]interface{}, 0)
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		reqSlice := reflect.ValueOf(req)
		for i := 0; i < reqSlice.Len(); i++ {
			if validateErrs, ok := ValidateStruct(reqSlice.Index(i).Interface(), isPatch); !ok {
				errs = append(errs, map[string]interface{}{
					"index":  i,
					"errors": validateErrs,
				})
			}
		}
	default:
		if validateErrs, ok := ValidateStruct(req, isPatch); !ok {
			errs = append(errs, map[string]interface{}{
				"errors": validateErrs,
			})
		}
	}

	if len(errs) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewErrorInvalidRequestBody(errs))
		ok = false
		return
	}

	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		reqSlice := reflect.ValueOf(req)
		if reqSlice.Len() == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(NewErrorEmptyRequestBody())
			ok = false
			return
		}
	default:
		if reflect.DeepEqual(req, reflect.Zero(rt).Interface()) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(NewErrorEmptyRequestBody())
			ok = false
			return
		}
	}

	ok = true
	return
}
