package objects

import "github.com/aerostatka/mongodb-example/blog/structs"

type Service interface {
	TestConnection()
	CreatePost(post *structs.Post) (*structs.Post, error)
	ReadPost(id string) (*structs.Post, error)
	DeletePost(id string) (bool, error)
	UpdatePost(post *structs.Post) (*structs.Post, error)
	ListPosts() ([]*structs.Post, error)
}

type StandardService struct {
	rep Repository
}

func NewStandardService(r Repository) *StandardService {
	return &StandardService{
		rep: r,
	}
}

func (s *StandardService) TestConnection() {
	s.rep.TestConnection()
}

func (s *StandardService) CreatePost(post *structs.Post) (*structs.Post, error) {
	return s.rep.CreatePost(post)
}

func (s *StandardService) ReadPost(id string) (*structs.Post, error) {
	return s.rep.ReadPost(id)
}

func (s *StandardService) UpdatePost(post *structs.Post) (*structs.Post, error) {
	return s.rep.UpdatePost(post)
}

func (s *StandardService) DeletePost(id string) (bool, error) {
	return s.rep.DeletePost(id)
}

func (s *StandardService) ListPosts() ([]*structs.Post, error) {
	return s.rep.ListPosts()
}
