package httpHandler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"user-service/internal/delivery/http/dto"

	"github.com/google/uuid"
)

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := h.UC.GetUserByID(r.Context(), userID)
	if err != nil {
		h.Logger.Error("failed to get user by id", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	meResponse := dto.ToDTOMeResponse(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(meResponse)
}

func (h *UserHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Logger.Warn("failed to parse user_id from query to uuid", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.UC.GetUserByID(r.Context(), userID)
	if err != nil {
		h.Logger.Error("failed to get user by id",
			slog.String("error", err.Error()),
		)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	dtoUser := dto.ToDTOSearchResult(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtoUser)

}

func (h *UserHandler) GetUsersByIDsHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.GetUserByIDsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.Logger.Warn("failed to decode user_ids", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	users, err := h.UC.GetUsersByIDs(r.Context(), req.UserIDs)
	if err != nil {
		h.Logger.Error("failed to get users by ids",
			slog.String("error", err.Error()),
		)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	var dtoUsers []dto.SearchUsersResponse
	for _, u := range users {
		dtoUsers = append(dtoUsers, *dto.ToDTOSearchResult(u))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtoUsers)

}
