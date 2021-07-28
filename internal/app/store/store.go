package store

// Storer defines interface for app's stores
type Storer interface {
	Posts() IPostRepository
	Categories() ICategoryRepository
	Materials() IMaterialRepository
	MatCategories() IMatCategoryRepository
	Users() IUserRepository
	Services() IServiceRepository
	Pages() IPageRepository
}
