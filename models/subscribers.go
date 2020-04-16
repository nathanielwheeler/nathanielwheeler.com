package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	// Not directly used, but needed to help gorm communicate with postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound : Indicates that a resource does not exist within postgres
	ErrNotFound = errors.New("models: resource not found")
	// ErrInvalidID : Returned when an invalid ID is provided to a method like Delete.
	ErrInvalidID = errors.New("models: ID provided was invalid")
)

// Subscriber : Model for people that want updates from my website.
type Subscriber struct {
	gorm.Model
	Email string `gorm:"not null;type:varchar(100);unique_index"`
	MonthlyUpdate bool `gorm:"default:true"`
	EveryUpdate bool `gorm:"default:false"`
}



// SubsService : Processes the logic for subscribers
type SubsService struct {
	db *gorm.DB
}

// NewSubsService : constructor for SubsService
func NewSubsService(connectionStr string) (*SubsService, error) {
	db, err := gorm.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &SubsService{
		db: db,
	}, nil
}

// #region SERVICE METHODS

// Create : Creates the provided subscriber and fills provided data fields
func (ss *SubsService) Create(sub *Subscriber) error {
	return ss.db.Create(sub).Error
}

// Update : Changes subscriber preferences
func (ss *SubsService) Update(sub *Subscriber) error {
	return ss.db.Save(sub).Error
}

// Delete : Removes the subscriber identified by the id
func (ss *SubsService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	sub := Subscriber{Model: gorm.Model{ID: id}}
	return ss.db.Delete(&sub).Error
}

// Close : Shuts down the connection to database
func (ss *SubsService) Close() error {
	return ss.db.Close()
}

// #endregion