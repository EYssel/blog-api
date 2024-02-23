package blogs

import (
	"errors"

	"github.com/google/uuid"
)

type Blog struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Author   string    `json:"author"`
	Likes    int       `json:"likes"`
	Comments []Comment `json:"comments"`
}

type Comment struct {
	Author      string `json:"author"`
	CommentText string `json:"commentText"`
}

var (
	NotFoundErr = errors.New("not found")
)

type MemStore struct {
	list map[uuid.UUID]Blog
}

func NewMemStore() *MemStore {
	list := make(map[uuid.UUID]Blog)
	return &MemStore{
		list,
	}
}

func (m MemStore) Add(blog Blog) error {
	id := uuid.New()
	blog.ID = id

	m.list[id] = blog
	return nil
}

func (m MemStore) Get(id uuid.UUID) (Blog, error) {

	if val, ok := m.list[id]; ok {
		return val, nil
	}

	return Blog{}, NotFoundErr
}

func (m MemStore) List() (map[uuid.UUID]Blog, error) {

	// var blogs []Blog

	// for _, item := range m.list {
	// 	blogs = append(blogs, item)
	// }

	return m.list, nil
}

func (m MemStore) Update(id uuid.UUID, blog Blog) error {

	if _, ok := m.list[id]; ok {
		m.list[id] = blog
		return nil
	}

	return NotFoundErr
}

func (m MemStore) Remove(id uuid.UUID) error {
	delete(m.list, id)
	return nil
}
