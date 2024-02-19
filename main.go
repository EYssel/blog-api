// Implemented following this guide
// https://www.jetbrains.com/guide/go/tutorials/rest_api_series/stdlib/

package main

import (
	"blog-app/blogs"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/gosimple/slug"
)

func main() {
	store := blogs.NewMemStore()
	blogsHandler := NewBlogsHandler(store)

	mux := http.NewServeMux()

	mux.Handle("/", &HomeHandler{})
	mux.Handle("/blogs", blogsHandler)
	mux.Handle("/blog", blogsHandler)

	http.ListenAndServe(":8080", mux)
}

type HomeHandler struct{}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World from Go :D"))
}

type BlogsHandler struct {
	store blogStore
}

func (h *BlogsHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var blog blogs.Blog
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	resourceID := slug.Make(blog.Title)

	if err := h.store.Add(resourceID, blog); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BlogsHandler) ListBlogs(w http.ResponseWriter, r *http.Request)  {}
func (h *BlogsHandler) GetBlog(w http.ResponseWriter, r *http.Request)    {}
func (h *BlogsHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {}
func (h *BlogsHandler) DeleteBlog(w http.ResponseWriter, r *http.Request) {}

var (
	BlogRe       = regexp.MustCompile(`^/blogs/*$`)
	BlogReWithID = regexp.MustCompile(`^/blogs/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

func (h *BlogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && BlogRe.MatchString(r.URL.Path):
		h.CreateBlog(w, r)
		return
	case r.Method == http.MethodGet && BlogRe.MatchString(r.URL.Path):
		h.ListBlogs(w, r)
		return
	case r.Method == http.MethodGet && BlogReWithID.MatchString(r.URL.Path):
		h.GetBlog(w, r)
		return
	case r.Method == http.MethodPut && BlogReWithID.MatchString(r.URL.Path):
		h.UpdateBlog(w, r)
		return
	case r.Method == http.MethodDelete && BlogReWithID.MatchString(r.URL.Path):
		h.DeleteBlog(w, r)
		return
	default:
		return
	}
}

func NewBlogsHandler(s blogStore) *BlogsHandler {
	return &BlogsHandler{
		store: s,
	}
}

type blogStore interface {
	Add(name string, blog blogs.Blog) error
	Get(name string) (blogs.Blog, error)
	Update(name string, blog blogs.Blog) error
	List() (map[string]blogs.Blog, error)
	Remove(name string) error
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
