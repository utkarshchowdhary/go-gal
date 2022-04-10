package main

import (
	"database/sql"
	"os"
)

const pageSize = 25

type ImageStore interface {
	FindAll(offset int) ([]Image, error)
	Find(id string) (*Image, error)
	FindAllByUser(user *User, offset int) ([]Image, error)
	Save(image *Image) error
}

var globalImageStore ImageStore

type DbImageStore struct {
	db *sql.DB
}

func NewDbImageStore(filepath string) (ImageStore, error) {
	err := os.MkdirAll(filepath, 0660)
	if err != nil {
		return nil, err
	}

	return &DbImageStore{
		globalPostgresDb,
	}, nil
}

func (store *DbImageStore) FindAll(offset int) ([]Image, error) {
	var images []Image

	rows, err := store.db.Query(
		`
		SELECT id, user_id, description, location, size, created_at
		FROM images
		ORDER BY created_at DESC
		LIMIT $1
		OFFSET $2
		`,
		pageSize,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var image Image
		err := rows.Scan(
			&image.Id,
			&image.UserId,
			&image.Description,
			&image.Location,
			&image.Size,
			&image.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	return images, nil
}

func (store *DbImageStore) Find(id string) (*Image, error) {
	row := store.db.QueryRow(
		`
		SELECT id, user_id, description, location, size, created_at
		FROM images
		WHERE id = $1
		`,
		id,
	)

	var image Image
	err := row.Scan(
		&image.Id,
		&image.UserId,
		&image.Description,
		&image.Location,
		&image.Size,
		&image.CreatedAt,
	)
	return &image, err
}

func (store *DbImageStore) FindAllByUser(user *User, offset int) ([]Image, error) {
	var images []Image

	rows, err := store.db.Query(
		`
		SELECT id, user_id, description, location, size, created_at
		FROM images
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
		OFFSET $3
		`,
		user.Id,
		pageSize,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var image Image
		err := rows.Scan(
			&image.Id,
			&image.UserId,
			&image.Description,
			&image.Location,
			&image.Size,
			&image.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	return images, nil
}

func (store *DbImageStore) Save(image *Image) error {
	_, err := store.db.Exec(
		`
		INSERT INTO images (id, user_id, description, location, size, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE 
		SET description = EXCLUDED.description
		`,
		image.Id,
		image.UserId,
		image.Description,
		image.Location,
		image.Size,
		image.CreatedAt,
	)
	return err
}
