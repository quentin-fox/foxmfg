package fox

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Database string `json:"database"`
		User     string `json:"user"`
		Password string `json:"password"`
	}
	Addr       string `json:"addr"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
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
	Hash       string `json:"hash,omitempty" db:"hash"`
}

type UserService interface {
	Create(u *User) error
	Update(u *User) error
	List() ([]User, error)
	ListOne(id int) (User, error)
	ListWithHash(id int) (User, error)
	Verify(id int) error
}

type AuthService interface {
	GenerateHash(password string) (hash string, err error)
	ValidatePassword(hash string, password string) (bool, error)
	IssueJWT(u User) (string, error)
	VerifyJWT(tokenStr string) (bool, error)
}
