package fox

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Database string `json:"database"`
		User     string `json:"user"`
		Password string `json:"password"`
	}
	Addr string `json:"addr"`

}

type ConfigService interface {
	GetConfig() Config
}

type User struct {
	ID         int    `json:"id"`
	FirstName  string `json:"firstName" db:"firstName"`
	LastName   string `json:"lastName" db:"lastName"`
	Email      string `json:"email"`
	IsVerified bool   `json:"isVerified" db:"isVerified"`
}

type UserService interface {
	Create(u *User) error
	Update(u *User) error
	List() ([]User, error)
	ListOne(id int) (User, error)
	Verify(id int) error
}