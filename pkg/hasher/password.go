package hasher

type PasswordHasher interface {
	Hash(password string) string
	CheckPassword(hash, password string) bool
}

type BCRYTHasher struct{}

func NewBCRYTHasher() *BCRYTHasher {
	return &BCRYTHasher{}
}

func (b *BCRYTHasher) Hash(password string) string {

}

func (b *BCRYTHasher) CheckPassword(hash, password string) bool {

}
