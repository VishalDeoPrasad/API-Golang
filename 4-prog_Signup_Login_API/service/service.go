package service

import (
	"errors"
	"golang/models"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Conn struct {
	db *gorm.DB
}

// NewService is the constructor for the Conn struct.
func NewConn(db *gorm.DB) (*Conn, error) {
	if db == nil {
		return nil, errors.New("please provide a valid connection")
	}

	s := &Conn{db: db}
	return s, nil
}

func (c *Conn) AutoMigrate() error {
	err := c.db.AutoMigrate(&models.User{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform database migration")
		return err
	}
	return nil
}

func (c *Conn) CreateUser(nu models.NewUserReq) (models.User, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("hash password is not generating.")
	}

	u1 := models.User{
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: string(hashPass),
	}

	err = c.db.Create(&u1).Error
	if err != nil {
		log.Error().Err(err).Str("user_name", u1.Name).Msg("Failed to create user")
		return u1, err
	}
	return u1, nil
}

func (c *Conn) Authentication(login models.LoginReq) error {
	email := login.Email
	password := login.Password

	var user models.User
	tx := c.db.Where("email = ?", email).First(&user)
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("Email Not Found:")
		return tx.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Error().Err(err).Msg("password didn't match:")
		return err
	}

	return nil
}
