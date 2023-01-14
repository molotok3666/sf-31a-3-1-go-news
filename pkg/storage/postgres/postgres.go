package postgres

import (
	"GoNews/pkg/storage"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DbStorage struct {
	db *pgxpool.Pool
}

func New(connstr string) (*DbStorage, error) {
	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	s := DbStorage{
		db: db,
	}
	return &s, nil
}

func (s *DbStorage) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(),
		`
			SELECT
			p.id,
			p.title,
			p.content,
			p.author_id,
			a.name as author_name,
			p.created_at,
			p.published_at
			FROM posts p
			INNER JOIN authors a on p.author_id = a.id
			ORDER BY p.id
		`,
	)
	if err != nil {
		return nil, err
	}

	var posts []storage.Post

	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.AuthorID,
			&p.AuthorName,
			&p.CreatedAt,
			&p.PublishedAt,
		)

		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func (s *DbStorage) AddPost(p storage.Post) error {
	var id int
	err := s.db.QueryRow(context.Background(),
		`
			INSERT INTO posts (title, content, author_id, created_at, published_at)
			VALUES ($1, $2, $3, $4, $5) RETURNING id;
		`,
		p.Title,
		p.Content,
		p.AuthorID,
		p.CreatedAt,
		p.PublishedAt,
	).Scan(&id)
	return err
}
func (s *DbStorage) UpdatePost(p storage.Post) error {
	var id int
	err := s.db.QueryRow(context.Background(),
		`
			UPDATE posts
			SET title = $1, content = $2, author_id = $3,
			created_at = $4, published_at = $5
			WHERE id = $6
			RETURNING id;
		`,
		p.Title,
		p.Content,
		p.AuthorID,
		p.CreatedAt,
		p.PublishedAt,
		p.ID,
	).Scan(&id)
	return err
}
func (s *DbStorage) DeletePost(p storage.Post) error {
	var id int
	err := s.db.QueryRow(context.Background(),
		`
			DELETE FROM posts WHERE id = $1
			RETURNING id
		`,
		p.ID,
	).Scan(&id)
	return err
}
