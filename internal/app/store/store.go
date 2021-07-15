package store

// Storer defines interface for app's stores
type Storer interface {
	Posts() IPostRepository
	Categories() ICategoryRepository
	// Pages()
	// Users()
}
