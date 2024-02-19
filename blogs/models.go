package blogs

import (
	"errors"
)

type Blog struct {
	ID       string    `json:"id"`
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
	list map[string]Blog
}

func NewMemStore() *MemStore {
	list := make(map[string]Blog)
	return &MemStore{
		list,
	}
}

func (m MemStore) Add(title string, blog Blog) error {
	m.list[title] = blog
	return nil
}

func (m MemStore) Get(title string) (Blog, error) {

	if val, ok := m.list[title]; ok {
		return val, nil
	}

	return Blog{}, NotFoundErr
}

func (m MemStore) List() (map[string]Blog, error) {
	return m.list, nil
}

func (m MemStore) Update(title string, blog Blog) error {

	if _, ok := m.list[title]; ok {
		m.list[title] = blog
		return nil
	}

	return NotFoundErr
}

func (m MemStore) Remove(title string) error {
	delete(m.list, title)
	return nil
}
