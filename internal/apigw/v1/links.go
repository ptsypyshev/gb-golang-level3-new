package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/ptsypyshev/gb-golang-level3-new/pkg/pb"
)

func newLinksHandler(linksClient linksClient) *linksHandler {
	return &linksHandler{client: linksClient}
}

type linksHandler struct {
	client linksClient
}

func (h *linksHandler) GetLinks(w http.ResponseWriter, r *http.Request) {
	// TODO implement me - implemented
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	links, err := h.client.ListLinks(ctx, &pb.Empty{})
	if err != nil {
		http.Error(w, "500 - Cannot get Links", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(links)
	if err != nil {
		http.Error(w, "500 - Cannot marshal Links", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (h *linksHandler) PostLinks(w http.ResponseWriter, r *http.Request) {
	// TODO implement me - implemented
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	var linkReq pb.CreateLinkRequest
	err := json.NewDecoder(r.Body).Decode(&linkReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if linkReq.Id != "" || linkReq.Url == "" {
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	_, err = h.client.CreateLink(ctx, &linkReq)
	if err != nil {
		http.Error(w, "500 - Cannot create Link", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *linksHandler) DeleteLinksId(w http.ResponseWriter, r *http.Request, id string) {
	// TODO implement me - implemented
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	req := &pb.GetLinkRequest{Id: r.PathValue("id")}

	_, err := h.client.GetLink(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("404 - Link with ID %s is not found", r.PathValue("id")), http.StatusNotFound)
		return
	}

	delReq := &pb.DeleteLinkRequest{Id: r.PathValue("id")}
	_, err = h.client.DeleteLink(ctx, delReq)
	if err != nil {
		http.Error(w, "500 - Cannot create Link", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *linksHandler) GetLinksId(w http.ResponseWriter, r *http.Request, id string) {
	// TODO implement me - implemented
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	req := &pb.GetLinkRequest{Id: r.PathValue("id")}

	link, err := h.client.GetLink(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("404 - Link with ID %s is not found", r.PathValue("id")), http.StatusNotFound)
		return
	}

	b, err := json.Marshal(link)
	if err != nil {
		http.Error(w, "500 - Cannot marshal Link", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (h *linksHandler) PutLinksId(w http.ResponseWriter, r *http.Request, id string) {
	// TODO implement me - implemented
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	var linkReq pb.UpdateLinkRequest
	err := json.NewDecoder(r.Body).Decode(&linkReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := &pb.GetLinkRequest{Id: r.PathValue("id")}

	link, err := h.client.GetLink(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("404 - Link with ID %s is not found", r.PathValue("id")), http.StatusNotFound)
		return
	}

	updReq := &pb.UpdateLinkRequest{Id: link.Id}
	if link.Title != linkReq.Title {
		updReq.Title = linkReq.Title
	}

	if link.Url != linkReq.Url {
		updReq.Url = linkReq.Url
	}

	if slices.Compare(link.Images, linkReq.Images) == 0 {
		updReq.Images = linkReq.Images
	}

	if slices.Compare(link.Tags, linkReq.Tags) == 0 {
		updReq.Tags = linkReq.Tags
	}

	if link.UserId != linkReq.UserId {
		updReq.UserId = linkReq.UserId
	}

	_, err = h.client.UpdateLink(ctx, updReq)
	if err != nil {
		http.Error(w, "500 - Cannot update Link", http.StatusInternalServerError)
		return
	}
}

func (h *linksHandler) GetLinksUserUserID(w http.ResponseWriter, r *http.Request, userID string) {
	// TODO implement me - implemented
	ctx, cancel := context.WithTimeout(r.Context(), ctxTimeout)
	defer cancel()

	req := &pb.GetLinksByUserId{UserId: r.PathValue("userID")}

	links, err := h.client.GetLinkByUserID(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if len(links.Links) == 0 {
		http.Error(w, fmt.Sprintf("404 - Links for user with ID %s are not found", r.PathValue("id")), http.StatusNotFound)
		return
	}

	b, err := json.Marshal(links)
	if err != nil {
		http.Error(w, "500 - Cannot marshal Links", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
