package mongostore

import "github.com/the-NZA/acg-nikolaev/internal/app/models"

// PostRepository implements IPostRepository
type PostRepository struct {
	store *Store
}

// Create save new post
func (p PostRepository) Create(*models.Post) error {
	return nil
}

// Find lookup post by it slug
func (p PostRepository) Find(slug string) (*models.Post, error) {
	return nil, nil
}

// Delete marks post as deleted
func (p PostRepository) Delete(*models.Post) error {
	return nil
}
