package domain

import (
	"time"

	"github.com/google/uuid"
)

// Address representa el modelo de datos para una dirección.
type Address struct {
	AddressID     uuid.UUID `json:"addressId"`
	CP            string    `json:"CP"`
	Street        string    `json:"street"`
	Number        int       `json:"number"`
	Floor         *string   `json:"floor,omitempty"`
	MainAddress   bool      `json:"mainAddress"`
	CreationDate  time.Time `json:"creationDate"`
	UpdatedDate   time.Time `json:"updatedDate"`
	ActiveAddress bool      `json:"activeAddress"`
	IdProfile     uuid.UUID `json:"idProfile"` // Clave foránea a Profile
}

// NewAddress crea una nueva instancia de Address con un ID generado.
func NewAddress(CP, street string, floor *string, mainAddress, activeAddress bool, number int, profileID uuid.UUID) *Address {
	return &Address{
		AddressID:     uuid.New(),
		CP:            CP,
		Street:        street,
		Number:        number,
		Floor:         floor,
		MainAddress:   mainAddress,
		ActiveAddress: activeAddress,
		IdProfile:     profileID,
		CreationDate:  time.Now(),
		UpdatedDate:   time.Now(),
	}
}
