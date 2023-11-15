package models

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Conn is our main struct, including the database instance for working with data.
type Conn struct {
	// db is an instance of the SQLite database.
	db *gorm.DB
}

// NewService is the constructor for the Conn struct.
func NewConn(db *gorm.DB) (*Conn, error) {
	// We check if the database instance is nil, which would indicate an issue.
	if db == nil {
		return nil, errors.New("please provide a valid connection")
	}
	// We initialize our service with the passed database instance.
	s := &Conn{db: db}
	return s, nil
}

// Define the function CreatInventory, which belongs to the struct 'Conn'.
// This function takes in 3 parameters: a context `ctx` of type `Context`, `ni` of type `NewInventory`, and `userId` of type `uint`.
// This function will return an `Inventory` and an `error`.

func (s *Conn) CreatInventory(ctx context.Context, ni NewInventory, userId uint) (Inventory, error) {
	// Create a new 'Inventory' struct named 'inv'.
	// Initialize it with parameters from the 'NewInventory' struct and the `userId` passed to the function.
	inv := Inventory{
		ItemName:    ni.ItemName,
		Quantity:    ni.Quantity,
		Category:    ni.Category,
		UserId:      userId,
		CostPerItem: ni.CostPerItem,
	}

	// Create a new database transaction using `ctx` as the context.
	// Within this transaction, create a new row in the database for the 'inv' struct.
	tx := s.db.WithContext(ctx).Create(&inv)

	// If there's an error with the database transaction.
	if tx.Error != nil {
		// Return an empty 'Inventory' struct and the error.
		return Inventory{}, tx.Error
	}

	// If there was no error with the database transaction, return 'inv' and nil as the error.
	return inv, nil
}

func (s *Conn) ViewInventory(ctx context.Context, userId string) ([]Inventory, float64, error) {
	var inv = make([]Inventory, 0, 10)
	tx := s.db.WithContext(ctx).Where("user_id = ?", userId)
	err := tx.Find(&inv).Error
	if err != nil {
		return nil, 0, err
	}

	totalCost, err := CalculateTotalCost(inv, "shirts")
	if err != nil {
		return nil, 0, err
	}
	return inv, totalCost, nil

}

func CalculateTotalCost(inventories []Inventory, category string) (float64, error) {
	if category == "" {
		return 0, errors.New("category doesn't exist")
	}
	if inventories == nil {
		return 0, errors.New("inventory not found")
	}
	// Compute the total cost
	var totalCost float64
	for _, inventory := range inventories {
		totalCost += inventory.CostPerItem * float64(inventory.Quantity)
	}

	return totalCost, nil
}

// CreateUser is a method that creates a new user record in the database.
func (s *Conn) CreateUser(ctx context.Context, nu NewUser) (User, error) {

	// We hash the user's password for storage in the database.
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generating password hash: %w", err)
	}

	// We prepare the User record.
	u := User{
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: string(hashedPass),
	}

	// We attempt to create the new User record in the database.
	err = s.db.Create(&u).Error
	if err != nil {
		return User{}, err
	}

	// Successfully created the record, return the user.
	return u, nil
}

// Authenticate is a method that checks a user's provided email and password against the database.
func (s *Conn) Authenticate(ctx context.Context, email, password string) (jwt.RegisteredClaims,
	error) {

	// We attempt to find the User record where the email
	// matches the provided email.
	var u User
	tx := s.db.Where("email = ?", email).First(&u)
	if tx.Error != nil {
		return jwt.RegisteredClaims{}, tx.Error
	}

	// We check if the provided password matches the hashed password in the database.
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}

	// Successful authentication! Generate JWT claims.
	c := jwt.RegisteredClaims{
		Issuer:    "service project",
		Subject:   strconv.FormatUint(uint64(u.ID), 10),
		Audience:  jwt.ClaimStrings{"students"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	// And return those claims.
	return c, nil
}

func (s *Conn) AutoMigrate() error {
	//if s.db.Migrator().HasTable(&User{}) {
	//	return nil
	//}
	err := s.db.Migrator().DropTable(&User{}, &Inventory{})
	if err != nil {
		return err
	}

	// AutoMigrate function will ONLY create tables, missing columns and missing indexes, and WON'T change existing column's type or delete unused columns
	err = s.db.Migrator().AutoMigrate(&User{}, &Inventory{})
	if err != nil {
		// If there is an error while migrating, log the error message and stop the program
		return err
	}
	return nil
}
