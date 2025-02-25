package domain

import (
	"time"

	"github.com/google/uuid"
)

// Profile representa el modelo de datos para un perfil de usuario.
type Profile struct {
	ProfileID       uuid.UUID `json:"profileId"`
	UserID          string    `json:"userId"`
	ProfileImage    *string   `json:"profileImage,omitempty"`
	ProfileName     string    `json:"profileName"`
	ProfileLevel    int       `json:"profileLevel"`
	ProfilePoints   int       `json:"profilePoints"`
	ProfileMail     string    `json:"profileMail"`
	Phone           string    `json:"phone"`
	CUIL            string    `json:"CUIL"`
	FiscalAdress    string    `json:"fiscalAdress"`
	FiscalCondition string    `json:"fiscalCondition"`
	IIBB            string    `json:"IIBB"`
	CreationDate    time.Time `json:"creationDate"`
	UpdatedDate     time.Time `json:"updatedDate"`
}

// NewProfile crea una nueva instancia de Profile con un ID generado.
func NewProfile(userID, profileName, profileMail, phone, CUIL, fiscalAdress, fiscalCondition, IIBB string, profileLevel, profilePoints int, profileImage *string) *Profile {
	return &Profile{
		ProfileID:       uuid.New(),
		UserID:          userID,
		ProfileName:     profileName,
		ProfileMail:     profileMail,
		Phone:           phone,
		CUIL:            CUIL,
		FiscalAdress:    fiscalAdress,
		FiscalCondition: fiscalCondition,
		IIBB:            IIBB,
		ProfileLevel:    profileLevel,
		ProfilePoints:   profilePoints,
		ProfileImage:    profileImage,
		CreationDate:    time.Now(),
		UpdatedDate:     time.Now(),
	}
}
