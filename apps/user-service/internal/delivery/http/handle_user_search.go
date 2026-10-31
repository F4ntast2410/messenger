package httpHandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"user-service/internal/delivery/http/dto"

	"github.com/google/uuid"
)

func (h *UserHandler) SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	queryParams := r.URL.Query()
	query := queryParams.Get("query")
	limitStr := queryParams.Get("limit")
	cursorStr := queryParams.Get("cursor")

	limit := 10
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	var cursor *string
	if cursorStr != "" {
		cursor = &cursorStr
	}
	users, err := h.UC.SearchUsers(r.Context(), query, cursor, limit, userID)
	if err != nil {
		h.Logger.Error("failed to search users", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	searchResult := make([]dto.SearchUsersResponse, 0, limit)
	for _, u := range users {
		searchResult = append(searchResult, *dto.ToDTOSearchResult(&u))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(searchResult)
}
