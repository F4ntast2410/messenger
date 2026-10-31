package httpHandler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"user-service/internal/delivery/http/dto"
	customErrors "user-service/internal/errors"
)

func (h *UserHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Warn("failed to decode request body", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	loginResult, err := h.UC.LoginUser(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, customErrors.ErrInvalidEmailOrPassword) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Error: "Неверная почта или пароль",
			})
			return
		} else {
			h.Logger.Error("failed to login user", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    loginResult.SessionToken,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusOK)
}
func (h *UserHandler) LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := h.UC.DeleteSession(r.Context(), cookie.Value); err != nil {
		h.Logger.Error("failed to delete session", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusNoContent)
}
