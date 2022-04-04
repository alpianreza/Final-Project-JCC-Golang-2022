package models

import (
	"html"
	"strings"
	"time"

	"finalproject/utils/token"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type (
	// User
	Users struct {
		ID        uint      `gorm:"primaryKey,column:id" json:"id"`
		FullName  string    `gorm:"not null;" json:"full_name"`
		Username  string    `gorm:"not null;unique" json:"username"`
		Email     string    `gorm:"not null;unique" json:"email"`
		Password  string    `gorm:"not null;" json:"password"`
		Role      string    `gorm:"not null;" json:"role"`
		CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
		Posts     *[]Posts  `json:"posts" gorm:"foreignKey:user_id"`
	}

	Guest struct {
		ID        uint      `gorm:"primaryKey,column:id" json:"id"`
		FullName  string    `gorm:"not null;" json:"full_name"`
		Username  string    `gorm:"not null;unique" json:"username"`
		Role      string    `gorm:"not null;" json:"role"`
		CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
		Posts     *[]Posts  `json:"posts" gorm:"foreignKey:user_id"`
	}
)

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string, db *gorm.DB) (string, error) {

	var err error

	u := Users{}

	err = db.Model(Users{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(u.ID, u.Role)

	if err != nil {
		return "", err
	}

	return token, nil

}

func (u *Users) SaveUser(db *gorm.DB) (*Users, error) {
	//turn password into hash
	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if errPassword != nil {
		return &Users{}, errPassword
	}
	u.Password = string(hashedPassword)
	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	var err error = db.Create(&u).Error
	if err != nil {
		return &Users{}, err
	}
	return u, nil
}
