package acg

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	tmpl *template.Template
	err  error
)

const postPerPage = 15

func init() {
	tmpl = template.Must(template.ParseGlob("internal/*/views/*.gohtml"))
}

func (s *Server) handleHomePage() http.HandlerFunc {
	type homepage struct {
		Page     *models.Page
		Services []*models.Service
		Posts    []*models.Post
	}

	return func(w http.ResponseWriter, r *http.Request) {
		page, err := s.store.Pages().FindByURL(r.URL.Path)
		if err != nil {
			s.logger.Logf("[DEBUG] page: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		services, err := s.store.Services().FindAll(bson.M{"deleted": false})
		if err != nil {
			s.logger.Logf("[DEBUG] services: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		posts, err := s.store.Posts().Aggregate(mongo.Pipeline{
			bson.D{{"$lookup", bson.D{{"from", "categories"}, {"localField", "category_id"}, {"foreignField", "_id"}, {"as", "category_slug"}}}},
			bson.D{{"$match", bson.D{{"deleted", false}}}},
			bson.D{{"$project", bson.D{{"title", 1}, {"snippet", 1}, {"postimg", 1}, {"time", 1}, {"slug", 1}, {"category_slug", "$category_slug.slug"}}}},
			bson.D{{"$sort", bson.D{{"time", -1}}}},
			bson.D{{"$limit", 3}},
			bson.D{{"$unwind", "$category_slug"}}})
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		err = tmpl.ExecuteTemplate(w, "index.gohtml", &homepage{
			Page:     page,
			Services: services,
			Posts:    posts,
		})
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}
}

func (s *Server) handleAboutPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aboutpage, err := s.store.Pages().FindByURL(r.URL.Path)
		if err != nil {
			s.logger.Logf("[DEBUG] page: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		tmpl.ExecuteTemplate(w, "singlepage.gohtml", aboutpage)
	}
}

func (s *Server) handlePostsPage() http.HandlerFunc {
	type postsPage struct {
		Page          *models.Page
		Posts         []*models.Post
		Categories    []*models.Category
		Pagination    []helpers.PaginationLink
		NumberOfPages int
		CurrentPage   string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var pageNumber uint64

		page, err := s.store.Pages().FindByURL(r.URL.Path)
		if err != nil {
			s.logger.Logf("[DEBUG] page: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
		}

		pNum := r.URL.Query().Get("page")
		if pNum != "" {
			pageNumber, err = strconv.ParseUint(pNum, 10, 64)
			if err != nil {
				s.logger.Logf("[DEBUG] %v\n", err)
				http.Redirect(w, r, "/posts", http.StatusSeeOther)
				return
			}
		} else {
			pageNumber = 1
		}

		numOfPosts, err := s.store.Posts().Count(bson.D{{"deleted", false}})
		if err != nil {
			s.logger.Logf("[DEBUG] page: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
		}

		// Calculate maximum number of pages
		maxPageNumber := numOfPosts / postPerPage

		// Fix number maximum number of pages for odd value
		if numOfPosts%postPerPage != 0 {
			maxPageNumber++
		}

		// If pageNumber out of maximum
		if pageNumber > uint64(maxPageNumber) {
			pageNumber = uint64(maxPageNumber)
		}

		// Number of posts to skip
		numOfSkip := (pageNumber - 1) * postPerPage

		// Find posts with joining information from categories colleciton
		posts, err := s.store.Posts().Aggregate(mongo.Pipeline{
			bson.D{{"$lookup", bson.D{{"from", "categories"}, {"localField", "category_id"}, {"foreignField", "_id"}, {"as", "category_slug"}}}},
			bson.D{{"$match", bson.D{{"deleted", false}}}},
			bson.D{{"$project", bson.D{{"title", 1}, {"snippet", 1}, {"postimg", 1}, {"time", 1}, {"slug", 1}, {"category_slug", "$category_slug.slug"}}}},
			bson.D{{"$sort", bson.D{{"time", -1}}}},
			bson.D{{"$skip", numOfSkip}},
			bson.D{{"$limit", postPerPage}},
			bson.D{{"$unwind", "$category_slug"}}})
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		// Generate pagination slice
		pagination := helpers.GeneratePagination(uint(pageNumber), uint(maxPageNumber))

		categories, err := s.store.Categories().FindAll(bson.M{"deleted": false})
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		err = tmpl.ExecuteTemplate(w, "posts.gohtml", &postsPage{
			Page:          page,
			Posts:         posts,
			Categories:    categories,
			Pagination:    pagination,
			NumberOfPages: int(maxPageNumber),
			CurrentPage:   strconv.Itoa(int(pageNumber)),
		})
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			// http.Redirect(w, r, "/404", http.StatusSeeOther)
		}
	}
}

func (s *Server) handleSinglePostPage() http.HandlerFunc {
	type singlePost struct {
		*models.Post
		CategoryName string
		CategoryURL  string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		category, err := s.store.Categories().FindBySlug(chi.URLParam(r, "categorySlug"))
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		post, err := s.store.Posts().FindBySlug(chi.URLParam(r, "postSlug"))
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		buf := &bytes.Buffer{}

		err = tmpl.ExecuteTemplate(buf, "singlepost.gohtml", &singlePost{
			Post:         post,
			CategoryName: category.Title,
			CategoryURL:  category.URL(),
		})
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		io.Copy(w, buf)
	}
}

func (s *Server) handleSingleCategoryPage() http.HandlerFunc {
	type categoryPage struct {
		Page            *models.Page
		Posts           []*models.Post
		CurrentCategory primitive.ObjectID
		Categories      []*models.Category
		Pagination      []helpers.PaginationLink
		NumberOfPages   int
		CurrentPage     string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var pageNumber uint64

		pNum := r.URL.Query().Get("page")
		if pNum != "" {
			pageNumber, err = strconv.ParseUint(pNum, 10, 64)
			if err != nil {
				s.logger.Logf("[DEBUG] %v\n", err)
				http.Redirect(w, r, "/posts", http.StatusSeeOther)
				return
			}
		} else {
			pageNumber = 1
		}

		category, err := s.store.Categories().FindBySlug(chi.URLParam(r, "categorySlug"))
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		numOfPosts, err := s.store.Posts().Count(bson.D{{"deleted", false}, {"category_id", category.ID}})
		if err != nil {
			s.logger.Logf("[DEBUG] page: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
		}

		if pageNumber == 0 {
			s.logger.Logf("[DEBUG] %v\n", pageNumber)
		}

		// Calculate maximum number of pages
		maxPageNumber := numOfPosts / postPerPage

		// Fix number maximum number of pages for odd value
		if numOfPosts%postPerPage != 0 {
			maxPageNumber++
		}

		// If pageNumber out of maximum
		if pageNumber > uint64(maxPageNumber) {
			pageNumber = uint64(maxPageNumber)
		}

		// Number of posts to skip
		numOfSkip := (pageNumber - 1) * postPerPage

		// Find posts with joining information from categories colleciton
		posts, err := s.store.Posts().Aggregate(mongo.Pipeline{
			bson.D{{"$lookup", bson.D{{"from", "categories"}, {"localField", "category_id"}, {"foreignField", "_id"}, {"as", "category_slug"}}}},
			bson.D{{"$match", bson.D{{"deleted", false}, {"category_id", category.ID}}}},
			bson.D{{"$project", bson.D{{"title", 1}, {"snippet", 1}, {"postimg", 1}, {"time", 1}, {"slug", 1}, {"category_slug", "$category_slug.slug"}}}},
			bson.D{{"$sort", bson.D{{"time", -1}}}},
			bson.D{{"$skip", numOfSkip}},
			bson.D{{"$limit", postPerPage}},
			bson.D{{"$unwind", "$category_slug"}}})
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		// Generate pagination slice
		pagination := helpers.GeneratePagination(uint(pageNumber), uint(maxPageNumber))

		categories, err := s.store.Categories().FindAll(bson.M{"deleted": false})
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		buf := &bytes.Buffer{}

		err = tmpl.ExecuteTemplate(buf, "category.gohtml", &categoryPage{
			Page: &models.Page{
				Title:    category.Title,
				Subtitle: category.Subtitle,
				MetaDesc: category.MetaDesc,
				URL:      category.URL(),
			},
			Posts:           posts,
			CurrentCategory: category.ID,
			Categories:      categories,
			Pagination:      pagination,
			NumberOfPages:   int(maxPageNumber),
			CurrentPage:     strconv.Itoa(int(pageNumber)),
		})
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
		}

		io.Copy(w, buf)
	}
}

func (s *Server) handleMaterialsPage() http.HandlerFunc {
	type matcat struct {
		models.MatCategory
		Materials []models.Material
	}

	type materialsPage struct {
		Page    *models.Page
		MatCats []*models.MaterialShow
	}

	return func(w http.ResponseWriter, r *http.Request) {
		page, err := s.store.Pages().FindByURL(r.URL.Path)
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		mats, err := s.store.MatCategories().Aggregate(mongo.Pipeline{
			{
				{
					Key: "$lookup", Value: bson.D{
						{Key: "from", Value: "materials"},
						{Key: "let", Value: bson.D{
							{Key: "matcat_id", Value: "$_id"},
						}},
						{Key: "pipeline", Value: bson.A{
							bson.D{
								{Key: "$match", Value: bson.D{
									{Key: "deleted", Value: false},
								}},
							},
							bson.D{
								{Key: "$sort", Value: bson.D{
									{Key: "time", Value: -1},
								}},
							},
							bson.D{
								{Key: "$match", Value: bson.D{
									{Key: "$expr", Value: bson.D{
										{Key: "$eq", Value: bson.A{"$matcategory_id", "$$matcat_id"}},
									}},
								}},
							}, bson.D{
								{Key: "$limit", Value: 3},
							},
						}},
						{Key: "as", Value: "materials"},
					},
				},
			},
			{
				{
					Key: "$match", Value: bson.D{
						{Key: "deleted", Value: false},
					},
				},
			},
			{
				{
					Key: "$project", Value: bson.D{
						{Key: "_id", Value: 1},
						{Key: "title", Value: 1},
						{Key: "slug", Value: 1},
						{Key: "desc", Value: 1},
						{Key: "materials", Value: 1},
					},
				},
			},
		})
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			// s.error(w, r, http.StatusInternalServerError, err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
			return
		}

		buf := &bytes.Buffer{}

		err = tmpl.ExecuteTemplate(buf, "materials.gohtml", materialsPage{
			Page:    page,
			MatCats: mats,
		})
		if err != nil {
			s.logger.Logf("[DEBUG] %v\n", err)
			http.Redirect(w, r, "/404", http.StatusSeeOther)
		}

		io.Copy(w, buf)
	}
}

func (s *Server) handleServicesPage() http.HandlerFunc {
	type servicepage struct {
		Page     *models.Page
		Services []*models.Service
	}

	return func(w http.ResponseWriter, r *http.Request) {
		page, err := s.store.Pages().FindByURL(r.URL.Path)
		if err != nil {
			s.logger.Logf("[DEBUG] page: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		services, err := s.store.Services().FindAll(bson.M{"deleted": false})
		if err != nil {
			s.logger.Logf("[DEBUG] services: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		tmpl.ExecuteTemplate(w, "services.gohtml", &servicepage{
			Page:     page,
			Services: services,
		})
	}
}

func (s *Server) handleContactsPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contactspage, err := s.store.Pages().FindByURL(r.URL.Path)
		if err != nil {
			s.logger.Logf("[DEBUG] page: %v\n", err)
			http.Redirect(w, r, "/404", http.StatusNotFound)
		}

		tmpl.ExecuteTemplate(w, "singlepage.gohtml", contactspage)
	}
}
