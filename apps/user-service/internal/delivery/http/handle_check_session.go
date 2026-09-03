package httpHandler

import (
	"log/slog"
	"net/http"
)

func (h *UserHandler) CheckSessionIDHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("X-User-ID")
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userSession, err := h.UC.ValidateSession(r.Context(), cookie.Value)
	if err != nil {
		h.Logger.Error("error checking session id", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !userSession.IsValid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Add("X-User-ID", userSession.UserID.String())
	w.WriteHeader(http.StatusOK)

}
