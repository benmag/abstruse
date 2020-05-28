package repository

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/jkuri/abstruse/internal/pkg/auth"
	"github.com/jkuri/abstruse/internal/server/db/model"
	"go.uber.org/zap"
)

// UserRepository iface.
type UserRepository interface {
	Find(ID uint64) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Login(email, password string) (string, error)
}

// DBUserRepository struct.
type DBUserRepository struct {
	logger *zap.Logger
	db     *gorm.DB
}

// NewDBUserRepository func.
func NewDBUserRepository(logger *zap.Logger, db *gorm.DB) UserRepository {
	return &DBUserRepository{
		logger: logger.With(zap.String("type", "UserRepository")),
		db:     db,
	}
}

// Find returns user by id or error if not exists.
func (r *DBUserRepository) Find(ID uint64) (*model.User, error) {
	user := &model.User{}
	if err := r.db.Model(user).Where("id = ?", ID).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindByEmail returns user by email or error if not exists.
func (r *DBUserRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := r.db.Model(user).Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Login func
func (r *DBUserRepository) Login(email, password string) (string, error) {
	user, err := r.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if checkPassword(user, password) {
		jsonwebtoken, err := generateJWT(user)
		if err != nil {
			return "", err
		}
		return jsonwebtoken, nil
	}
	return "", fmt.Errorf("invalid password")
}

func checkPassword(u *model.User, password string) bool {
	return auth.CheckPasswordHash(password, u.Password)
}

func generateJWT(u *model.User) (string, error) {
	user := auth.UserJWT{
		ID:       strconv.Itoa(int(u.ID)),
		Email:    u.Email,
		Fullname: u.Fullname,
		Avatar:   u.Avatar,
	}
	return auth.GenerateJWT(user)
}
