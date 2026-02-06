package main

import (
	"encoding/json"
	"net/http"

	"github.com/WadeGulbrandsen/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid ApiKey", err)
		return
	}

	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	switch params.Event {
	case "user.upgraded":
		user_id, err := uuid.Parse(params.Data.UserID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}
		_, err = cfg.dbQueries.UpgradeToRed(r.Context(), user_id)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		w.WriteHeader(http.StatusNoContent)
		return
	}
}
