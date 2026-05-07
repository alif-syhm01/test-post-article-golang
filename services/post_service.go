package services

import (
	"errors"
	"test-post-article/model"
	"test-post-article/repositories"
)

type PostService struct {
	Repo *repositories.PostRepository
}

func NewPostService(repo *repositories.PostRepository) *PostService {
	return &PostService{Repo: repo}
}

func postValidation(post *model.Post) error {
	if post.Title == "" {
		return errors.New("title is required.")
	}

	if post.Category == "" {
		return errors.New("category is required.")
	}

	if post.Content == "" {
		return errors.New("content is required.")
	}

	if post.Status == "" {
		return errors.New("status is required.")
	}

	if len(post.Title) < 20 {
		return errors.New("title must be at least 20 characters.")
	}

	if len(post.Content) < 200 {
		return errors.New("content must be at least 200 characters.")
	}

	if len(post.Category) < 3 {
		return errors.New("category must be at least 3 characters.")
	}

	if post.Status != "publish" && post.Status != "draft" && post.Status != "trash" {
		return errors.New("status must be publish, draft, or trash.")
	}

	return nil
}

func (s *PostService) Create(post *model.Post) error {
	if err := postValidation(post); err != nil {
		return err
	}

	return s.Repo.Create(post)
}

func (s *PostService) GetAll() ([]model.Post, error) {
	return s.Repo.GetAll()
}

func (s *PostService) GetAllPaginate(limit, offset int) ([]model.Post, int, error) {
	return s.Repo.GetAllPaginate(limit, offset)
}

func (s *PostService) GetById(id int) (*model.Post, error) {
	return s.Repo.GetById(id)
}

func (s *PostService) Update(id int, post *model.Post) error {
	if err := postValidation(post); err != nil {
		return err
	}

	return s.Repo.Update(id, post)
}

func (s *PostService) Delete(id int) error {
	return s.Repo.Delete(id)
}
