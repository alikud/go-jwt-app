package handler

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func (h *Handler) RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.NewUUID()
		var requestID = id
		//https://medium.com/golangspec/globally-unique-key-for-context-value-in-golang-62026854b48f
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestID, id)))
	})
}
func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		username, password, ok := request.BasicAuth()
		a := "admin"
		b := "123"
		if ok {
			if username == a && password == b {
				next.ServeHTTP(writer, request)
			} else {
				writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			}
		}
	})
}

func (h *Handler) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		tokenString := request.Header["Authorization"]
		if len(tokenString) == 0 {
			NewErrorResponse(writer, request, errors.New("Empty token"), 404)
			//TODO BUG! IF TOKEN EMPTY SEE ONLY 404 status code, but not error response
		} else {
			token := strings.Split(request.Header["Authorization"][0], " ")
			//fmt.Println(token[1])
			id, _ := h.Service.ParseToken(token[1])
			//fmt.Println(id)
			//fmt.Println(err)

			var UserID = "UserID"
			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), UserID, id)))

			//TODO if user exist put them to request context
		}
	})
}
