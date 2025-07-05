package internal

import (
	"encoding/json"
	"net/http"
)

// Handler struct wires everything
type Handler struct {
	User *UserService
	Chat *ChatService
	mux  *http.ServeMux
}

func NewHandler(user *UserService, chat *ChatService) *Handler {
	h := &Handler{
		User: user,
		Chat: chat,
		mux:  http.NewServeMux(),
	}
	h.routes()
	return h
}

func (h *Handler) routes() {
	h.mux.HandleFunc("/signup", h.signup)
	h.mux.HandleFunc("/login", h.login)
	h.mux.HandleFunc("/dm/send", AuthMiddleware(h.User, h.dmSend))
	h.mux.HandleFunc("/dm/history", AuthMiddleware(h.User, h.dmHistory))
	h.mux.HandleFunc("/group/create", AuthMiddleware(h.User, h.groupCreate))
	h.mux.HandleFunc("/group/send", AuthMiddleware(h.User, h.groupSend))
	h.mux.HandleFunc("/group/history", AuthMiddleware(h.User, h.groupHistory))
	h.mux.HandleFunc("/broadcast/send", AuthMiddleware(h.User, h.broadcastSend))
	h.mux.HandleFunc("/broadcast/history", AuthMiddleware(h.User, h.broadcastHistory))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) { h.mux.ServeHTTP(w, r) }
func (h *Handler) respond(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

// --- USER ---

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	var req struct{ Username, Password string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}
	if err := h.User.Register(r.Context(), req.Username, req.Password); err != nil {
		h.respond(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusCreated, map[string]string{"message": "user registered"})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req struct{ Username, Password string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}
	token, err := h.User.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		h.respond(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusOK, map[string]string{"token": token})
}

// --- DM ---

func (h *Handler) dmSend(w http.ResponseWriter, r *http.Request) {
	username := UsernameFromCtx(r.Context())
	var req struct{ To, Content string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}
	if err := h.Chat.SendDirectMessage(r.Context(), username, req.To, req.Content); err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusOK, map[string]string{"message": "sent"})
}

func (h *Handler) dmHistory(w http.ResponseWriter, r *http.Request) {
	username := UsernameFromCtx(r.Context())
	to := r.URL.Query().Get("user")
	if to == "" {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "missing user"})
		return
	}
	msgs, err := h.Chat.GetDirectMessages(r.Context(), username, to)
	if err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusOK, msgs)
}

// --- GROUP ---

func (h *Handler) groupCreate(w http.ResponseWriter, r *http.Request) {
	var req struct{ Group string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}
	if req.Group == "" {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "missing group name"})
		return
	}
	if err := h.Chat.CreateGroup(r.Context(), req.Group); err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusCreated, map[string]string{"message": "group created"})
}

func (h *Handler) groupSend(w http.ResponseWriter, r *http.Request) {
	username := UsernameFromCtx(r.Context())
	var req struct{ Group, Content string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}
	if err := h.Chat.SendGroupMessage(r.Context(), username, req.Group, req.Content); err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusOK, map[string]string{"message": "sent"})
}

func (h *Handler) groupHistory(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	if group == "" {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "missing group"})
		return
	}
	msgs, err := h.Chat.GetGroupMessages(r.Context(), group)
	if err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusOK, msgs)
}

// --- BROADCAST ---

func (h *Handler) broadcastSend(w http.ResponseWriter, r *http.Request) {
	username := UsernameFromCtx(r.Context())
	var req struct{ Content string }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}
	if err := h.Chat.Broadcast(r.Context(), username, req.Content); err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusOK, map[string]string{"message": "broadcasted"})
}

func (h *Handler) broadcastHistory(w http.ResponseWriter, r *http.Request) {
	msgs, err := h.Chat.GetBroadcasts(r.Context())
	if err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusOK, msgs)
}
