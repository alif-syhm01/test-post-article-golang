package repositories

import (
	"database/sql"
	"errors"
	"test-post-article/model"
	"time"
)

type PostRepository struct {
	DB *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) Create(post *model.Post) error {
	now := time.Now()
	post.CreatedDate = now
	post.UpdatedDate = now

	query := "INSERT INTO posts (title, content, category, created_date, updated_date, status) VALUES (?, ?, ?, ?, ?, ?)"

	result, err := r.DB.Exec(query, post.Title, post.Content, post.Category, post.CreatedDate, post.UpdatedDate, post.Status)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	post.ID = int(id)
	return nil
}

func (r *PostRepository) GetAll() ([]model.Post, error) {
	query := "SELECT id, title, content, category, created_date, updated_date, status FROM posts"

	rows, err := r.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var p model.Post

		rows.Scan(&p.ID, &p.Title, &p.Content, &p.Category, &p.CreatedDate, &p.UpdatedDate, &p.Status)

		posts = append(posts, p)
	}

	return posts, nil
}

func (r *PostRepository) GetAllPaginate(limit, offset int) ([]model.Post, int, error) {
	var total int

	err := r.DB.QueryRow("SELECT COUNT(*) FROM posts").Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	query := "SELECT id, title, content, category, created_date, updated_date, status FROM posts LIMIT ? OFFSET ?"

	rows, err := r.DB.Query(query, limit, offset)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	var posts []model.Post

	for rows.Next() {
		var p model.Post

		rows.Scan(&p.ID, &p.Title, &p.Content, &p.Category, &p.CreatedDate, &p.UpdatedDate, &p.Status)

		posts = append(posts, p)
	}

	return posts, total, nil
}

func (r *PostRepository) GetById(id int) (*model.Post, error) {
	query := "SELECT id, title, content, category, created_date, updated_date, status FROM posts WHERE id = ?"

	row := r.DB.QueryRow(query, id)
	var p model.Post

	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Category, &p.CreatedDate, &p.UpdatedDate, &p.Status)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PostRepository) Update(id int, post *model.Post) error {
	now := time.Now()
	post.UpdatedDate = now

	query := "UPDATE posts SET title=?, content=?, category=?, updated_date=?, status=? WHERE id=?"

	result, err := r.DB.Exec(query, post.Title, post.Content, post.Category, post.UpdatedDate, post.Status, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("not found")
	}

	return nil
}

func (r *PostRepository) Delete(id int) error {
	query := "DELETE FROM posts where id = ?"

	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("not found")
	}

	return nil
}
