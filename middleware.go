package main

import (
	"fmt"
	"net/http"

	"github.com/rubacodes/boot-dev-rss/internal/auth"
	"github.com/rubacodes/boot-dev-rss/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// pre execution, nice
		apikey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error:%s", err))
			return
		}
		user, err := apiCfg.DB.GetuserByApiKey(r.Context(), apikey)
		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("User not found: %v", err))
			return
		}
		handler(w, r, user)
	}

}
