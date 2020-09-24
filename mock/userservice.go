package mock

import "github.com/quentin-fox/fox"

type UserService struct {
	CreateFn func(u *fox.User) error
	CreateInvoked bool

	UpdateFn func(u *fox.User) error
	UpdateInvoked bool

	ListFn func() ([]fox.User, error)
	ListInvoked bool

	ListOneFn func(id int) (fox.User, error)
	ListOneInvoked bool

	VerifyFn func(id int) error
	VerifyInvoked bool

	ListWithHashFn func(id int) (fox.User, error)
	ListWithHashInvoked bool
}

func (s *UserService) Create(u *fox.User) error {
	s.CreateInvoked = true
	return s.CreateFn(u)
}

func (s *UserService) Update(u *fox.User) error {
	s.UpdateInvoked = true
	return s.UpdateFn(u)
}

func (s *UserService) List() ([]fox.User, error) {
	s.ListInvoked = true
	return s.ListFn()
}

func (s *UserService) ListOne(id int) (fox.User, error) {
	s.ListOneInvoked = true
	return s.ListOneFn(id)
}

func (s *UserService) Verify(id int) error {
	s.VerifyInvoked = true
	return s.VerifyFn(id)
}

func (s *UserService) ListWithHash(id int) (fox.User ,error) {
	s.ListWithHashInvoked = true
	return s.ListWithHashFn(id)
}
