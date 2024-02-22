package middleware

import (
	"errors"
	"net/http"
)

type handler func(w http.ResponseWriter, r *http.Request) error

func Middleware(h handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		var customError *CustomError
		err := h(w, r)
		if err != nil {
			if errors.As(err, &customError) {
				if errors.Is(err, ErrEntityNotFound) {
					w.WriteHeader(http.StatusNoContent)
					w.Write(ErrEntityNotFound.Marshal())
					return
				}

				if errors.Is(err, ErrUserDuplicate) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(ErrUserDuplicate.Marshal())
					return
				}

				err = err.(*CustomError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(customError.Marshal())
				return
			}

			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err).Marshal())
		}
	}
}
