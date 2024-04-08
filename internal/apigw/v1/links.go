package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ptsypyshev/gb-golang-level3-new/pkg/api/apiv1"
	"github.com/ptsypyshev/gb-golang-level3-new/pkg/pb"
)

func newLinksHandler(linksClient linksClient) *linksHandler {
	return &linksHandler{client: linksClient}
}

type linksHandler struct {
	client linksClient
}

func (h *linksHandler) GetLinks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	links, err := h.client.ListLinks(ctx, &pb.Empty{})
	if err != nil {
		slog.Error("cannot get Links at GetLinks handler", slog.Any("err", err))
		http.Error(w, "500 - Cannot get Links", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(links)
	if err != nil {
		slog.Error("cannot marshal list of Links to JSON at GetLinks handler", slog.Any("err", err))
		http.Error(w, "500 - Cannot marshal Links", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		slog.Error("cannot write response at GetLinks handler", slog.Any("err", err))
	}
}

func (h *linksHandler) PostLinks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	var linkReq apiv1.LinkCreate
	err := json.NewDecoder(r.Body).Decode(&linkReq)
	if err != nil {
		slog.Error("cannot decode request body at PostLinks handler", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if linkReq.Id != "" || linkReq.Url == "" {
		slog.Error("invalid body params at PostLinks handler", slog.Any("err", err))
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	req := &pb.CreateLinkRequest{
		Id:     linkReq.Id,
		Images: linkReq.Images,
		Tags:   linkReq.Tags,
		Title:  linkReq.Title,
		UserId: linkReq.UserId,
		Url:    linkReq.Url,
	}

	_, err = h.client.CreateLink(ctx, req)
	if err != nil {
		slog.Error("cannot create Link at PostLinks handler", slog.Any("err", err))
		http.Error(w, "500 - Cannot create Link", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *linksHandler) DeleteLinksId(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	req := &pb.GetLinkRequest{Id: r.PathValue("id")}

	_, err := h.client.GetLink(ctx, req)
	if err != nil {
		slog.Info("cannot get Link at DeleteLinksId handler", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("404 - Link with ID %s is not found", r.PathValue("id")), http.StatusNotFound)
		return
	}

	delReq := &pb.DeleteLinkRequest{Id: r.PathValue("id")}
	_, err = h.client.DeleteLink(ctx, delReq)
	if err != nil {
		slog.Error("cannot delete Link at DeleteLinksId handler", slog.Any("err", err))
		http.Error(w, "500 - Cannot delete Link", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *linksHandler) GetLinksId(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	req := &pb.GetLinkRequest{Id: r.PathValue("id")}

	link, err := h.client.GetLink(ctx, req)
	if err != nil {
		slog.Info("cannot get Link at GetLinksId handler", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("404 - Link with ID %s is not found", r.PathValue("id")), http.StatusNotFound)
		return
	}

	b, err := json.Marshal(link)
	if err != nil {
		slog.Error("cannot marshal Link to JSON at GetLinksId handler", slog.Any("err", err))
		http.Error(w, "500 - Cannot marshal Link", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		slog.Error("cannot write response at GetLinksId handler", slog.Any("err", err))
	}
}

func (h *linksHandler) PutLinksId(w http.ResponseWriter, r *http.Request, id string) {
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	var linkReq apiv1.LinkCreate
	err := json.NewDecoder(r.Body).Decode(&linkReq)
	if err != nil {
		slog.Error("cannot decode request body at PutLinksId handler", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updReq := &pb.UpdateLinkRequest{
		Id:     linkReq.Id,
		Title:  linkReq.Title,
		Url:    linkReq.Url,
		Images: linkReq.Images,
		Tags:   linkReq.Tags,
		UserId: linkReq.UserId,
	}

	_, err = h.client.UpdateLink(ctx, updReq)
	if err != nil {
		slog.Error("cannot updaet Link at PutLinksId handler", slog.Any("err", err))
		http.Error(w, "500 - Cannot update Link", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *linksHandler) GetLinksUserUserID(w http.ResponseWriter, r *http.Request, userID string) {
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	req := &pb.GetLinksByUserId{UserId: r.PathValue("userID")}

	links, err := h.client.GetLinkByUserID(ctx, req)
	if err != nil {
		slog.Info("cannot get Links by UserID at GetLinksUserUserID handler", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if len(links.Links) == 0 {
		slog.Info("no Links found by UserID at GetLinksUserUserID handler", slog.Any("err", err))
		http.Error(w, fmt.Sprintf("404 - Links for user with ID %s are not found", r.PathValue("id")), http.StatusNotFound)
		return
	}

	b, err := json.Marshal(links)
	if err != nil {
		slog.Error("cannot marshal Links to JSON at GetLinksUserUserID handler", slog.Any("err", err))
		http.Error(w, "500 - Cannot marshal Links", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		slog.Error("cannot write response at GetLinksUserUserID handler", slog.Any("err", err))
	}
}
