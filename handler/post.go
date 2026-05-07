package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test-post-article/model"
	"test-post-article/router"
	"test-post-article/services"
)

type GetAllResponse struct {
	Data   []model.Post `json:"data"`
	Total  int          `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

type PostHandler struct {
	Service *services.PostService
}

func NewPostHandler(svc *services.PostService) *PostHandler {
	return &PostHandler{Service: svc}
}

func getIdFromParam(w http.ResponseWriter, r *http.Request) *int {
	idStr, ok := r.Context().Value(router.ContextKey("id")).(string)
	if !ok {
		http.Error(w, "missing id", http.StatusBadRequest)
		return nil
	}

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return nil
	}

	return &id
}

func getPaginationParams(w http.ResponseWriter, r *http.Request) (limit int, offset int, ok bool) {
	limitStr, hasLimit := r.Context().Value(router.ContextKey("limit")).(string)
	offsetStr, hasOffset := r.Context().Value(router.ContextKey("offset")).(string)

	if !hasLimit || !hasOffset {
		http.Error(w, "missing pagination params", http.StatusBadRequest)
		return 0, 0, false
	}

	limit, err := strconv.Atoi(limitStr)

	if err != nil {
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return 0, 0, false
	}

	offset, err = strconv.Atoi(offsetStr)

	if err != nil {
		http.Error(w, "invalid offset", http.StatusBadRequest)
		return 0, 0, false
	}

	return limit, offset, true
}

func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := h.Service.GetAll()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) GetAllPaginate(w http.ResponseWriter, r *http.Request) {
	limit, offset, ok := getPaginationParams(w, r)

	if !ok {
		return
	}

	posts, total, err := h.Service.GetAllPaginate(limit, offset)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := GetAllResponse{
		Data:   posts,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	json.NewDecoder(r.Body).Decode(&post)

	err := h.Service.Create(&post)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id := getIdFromParam(w, r)

	if id == nil {
		return
	}

	post, err := h.Service.GetById(*id)

	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := getIdFromParam(w, r)

	if id == nil {
		return
	}

	var post model.Post
	json.NewDecoder(r.Body).Decode(&post)
	post.ID = *id

	err := h.Service.Update(*id, &post)

	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := getIdFromParam(w, r)

	if id == nil {
		return
	}

	err := h.Service.Delete(*id)

	if err != nil {
		if err.Error() == "not found" {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
