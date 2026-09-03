package httpHandler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"user-service/internal/delivery/http/dto"
	customErrors "user-service/internal/errors"
)

func (h *UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.UC.RegisterUser(r.Context(), req.Email, req.Username, req.Name, req.Password)
	if err != nil {
		if errors.Is(err, customErrors.ErrUserExists) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Error: "Пользователь уже существует",
			})
			return
		}
		h.Logger.Error("failed to register user", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Error: "Неизвестная ошибка",
		})
		return
	}
	w.WriteHeader(http.StatusCreated)

}
