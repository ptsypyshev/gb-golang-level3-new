package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ptsypyshev/gb-golang-level3-new/pkg/pb"
)

const (
	ctxTimeout = 2 * time.Second
)

func newUsersHandler(usersClient usersClient) *usersHandler {
	return &usersHandler{client: usersClient}
}

type usersHandler struct {
	client usersClient
}

func (h *usersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	users, err := h.client.ListUsers(ctx, &pb.Empty{})
	if err != nil {
		slog.Error("cannot get list of Users at GetUsers handler", slog.Any("err", err))
		http.Error(w, "500 - Cannot get Users", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(users)
	if err != nil {
		slog.Error("cannot marshal list of Users to JSON at GetUsers handler", slog.Any("err", err))
		http.Error(w, "500 - Cannot marshal Users", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		slog.Error("cannot write response at GetUsers handler", slog.Any("err", err))
	}
}

func (h *usersHandler) PostUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	var userReq pb.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		slog.Error("cannot decode request body at PostUsers handler", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userReq.Id != "" || userReq.Username == "" || userReq.Password == "" {
		slog.Error("invalid body params at PostUsers handler", slog.Any("err", err))
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	_, err = h.client.CreateUser(ctx, &userReq)
	if err != nil {
		slog.Error("cannot create User at PostUsers handler", slog.Any("err", err))
		http.Error(w, "500 - Cannot create User", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *usersHandler) DeleteUsersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	req := &pb.GetUserRequest{Id: r.PathValue("id")}

	_, err := h.client.GetUser(ctx, req)
	if err != nil {
		slog.Info("cannot get User at DeleteUsersId handler", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("404 - User with ID %s is not found", r.PathValue("id")), http.StatusNotFound)
		return
	}

	delReq := &pb.DeleteUserRequest{Id: r.PathValue("id")}
	_, err = h.client.DeleteUser(ctx, delReq)
	if err != nil {
		slog.Error("cannot delete User at DeleteUsersId handler", slog.Any("err", err))
		http.Error(w, "500 - Cannot create User", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *usersHandler) GetUsersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	req := &pb.GetUserRequest{Id: r.PathValue("id")}

	user, err := h.client.GetUser(ctx, req)
	if err != nil {
		slog.Info("cannot get User at GetUsersId handler", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("404 - User with ID %s is not found", r.PathValue("id")), http.StatusNotFound)
		return
	}

	b, err := json.Marshal(user)
	if err != nil {
		slog.Error("cannot marshal User to JSON at GetUsersId handler", slog.Any("err", err))
		http.Error(w, "500 - Cannot marshal User", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		slog.Error("cannot write response at GetUsersId handler", slog.Any("err", err))
	}
}

func (h *usersHandler) PutUsersId(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	var userReq pb.UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		slog.Error("cannot decode request body at PutUsersId handler", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := &pb.GetUserRequest{Id: r.PathValue("id")}

	user, err := h.client.GetUser(ctx, req)
	if err != nil {
		slog.Error("cannot get User at PutUsersId handler", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("404 - User with ID %s is not found", r.PathValue("id")), http.StatusNotFound)
		return
	}

	updReq := &pb.CreateUserRequest{Id: user.Id}
	if user.Username != userReq.Username {
		updReq.Username = userReq.Username
	}

	if user.Password != userReq.Password {
		updReq.Password = userReq.Password
	}

	_, err = h.client.CreateUser(ctx, updReq)
	if err != nil {
		slog.Error("cannot create User at PutUsersId handler", slog.Any("err", err))
		http.Error(w, "500 - Cannot update User", http.StatusInternalServerError)
		return
	}
}
