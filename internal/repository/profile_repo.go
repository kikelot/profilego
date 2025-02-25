package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	//"time"

	"profilego/internal/domain"

	"github.com/google/uuid"
)

// ProfileRepository define los métodos para interactuar con la base de datos.
type ProfileRepository struct {
	DB *sql.DB
}

/*=============================RABBIT=====================================
==========================================================================*/
// UpdateProfilePoints actualiza los puntos en la base de datos
func (r *ProfileRepository) UpdateProfilePoints(ctx context.Context, userId string, profile *domain.Profile) error {

	query := `UPDATE profile SET profilePoints = profilePoints + $1 WHERE userId = $2`
	result, err := r.DB.ExecContext(ctx, query, profile.ProfilePoints, userId)
	if err != nil {
		//log.Println("ERROR al ejecutar consulta:", err)
		return err
	}

	// Verificar si se actualizó alguna fila
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		//log.Println("⚠ No se encontró un perfil con userId:", userId)
		return errors.New("usuario no encontrado")
	}

	//log.Println("✅ Puntos actualizados correctamente para userId:", userId)
	return nil
}

// UpdateProfileLevel actualiza el nivel en la base de datos
func (r *ProfileRepository) UpdateProfileLevel(ctx context.Context, userId string, profile *domain.Profile) error {

	if profile.ProfilePoints > 1000 { // si vino a actualizar el nivel, y los puntos son mayores a 1000, reseteo los puntos
		profile.ProfilePoints = 0
	}

	query := `UPDATE profile SET profileLevel = $1, profilepoints =$3  WHERE userId = $2`
	result, err := r.DB.ExecContext(ctx, query,
		profile.ProfileLevel,
		userId,
		profile.ProfilePoints,
	)
	if err != nil {
		//log.Println("ERROR al ejecutar consulta:", err)
		return err
	}

	// Verificar si se actualizó alguna fila
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		//log.Println("⚠ No se encontró un perfil con userId:", userId)
		return errors.New("usuario no encontrado")
	}

	log.Println("✅ Nivel actualizado correctamente para userId:", userId)
	return nil
}

/*=============================ENDRABBIT=====================================
==========================================================================*/

// NewProfileRepository crea una nueva instancia del repositorio.
func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{DB: db}
}

// CreateProfile inserta un nuevo perfil en la base de datos. // Kike comento este es el Create original que insertaba todos los datos
/*func (r *ProfileRepository) CreateProfile(ctx context.Context, profile *domain.Profile) error {
	query := `
		INSERT INTO profile (
			profileId, userId, profileImage, profileName, profileLevel, profilePoints,
			profileMail, phone, CUIL, fiscalAdress, fiscalCondition, IIBB,
			creationDate, updatedDate
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err := r.DB.ExecContext(ctx, query,
		profile.ProfileID, profile.UserID, profile.ProfileImage, profile.ProfileName, profile.ProfileLevel,
		profile.ProfilePoints, profile.ProfileMail, profile.Phone, profile.CUIL, profile.FiscalAdress,
		profile.FiscalCondition, profile.IIBB, profile.CreationDate, profile.UpdatedDate,
	)
	return err
}*/

func (r *ProfileRepository) CreateProfile(ctx context.Context, profile *domain.Profile) error {
	query := `
		INSERT INTO profile (
			profileId, userId, profileName,  
			profileMail, phone,
			creationDate, updatedDate
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.DB.ExecContext(ctx, query,
		profile.ProfileID, profile.UserID, profile.ProfileName, profile.ProfileMail, profile.Phone, profile.CreationDate, profile.UpdatedDate,
	)
	return err
}

// GetProfile obtiene un perfil por su ID.
func (r *ProfileRepository) GetProfile(ctx context.Context, userId string) (*domain.Profile, error) {
	query := `SELECT profileId, userId, profileImage, profileName, profileLevel, profilePoints,
			profileMail, phone, CUIL, fiscalAdress, fiscalCondition, IIBB, 
			creationDate, updatedDate FROM profile WHERE userId = $1`

	row := r.DB.QueryRowContext(ctx, query, userId)

	var profile domain.Profile
	var (
		profileImage    sql.NullString
		profileMail     sql.NullString
		phone           sql.NullString
		cuil            sql.NullString
		fiscalAdress    sql.NullString
		fiscalCondition sql.NullString
		iibb            sql.NullString
	)

	err := row.Scan(
		&profile.ProfileID, &profile.UserID, &profileImage, &profile.ProfileName, &profile.ProfileLevel,
		&profile.ProfilePoints, &profileMail, &phone, &cuil, &fiscalAdress,
		&fiscalCondition, &iibb, &profile.CreationDate, &profile.UpdatedDate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No se encontró el perfil
		}
		//log.Println("❌ Error al recuperar el perfil:", err)
		return nil, err
	}

	// Convertir `sql.NullString` a `string` evitando `NULL`

	profile.ProfileMail = safeString(profileMail)
	profile.Phone = safeString(phone)
	profile.CUIL = safeString(cuil)
	profile.FiscalAdress = safeString(fiscalAdress)
	profile.FiscalCondition = safeString(fiscalCondition)
	profile.IIBB = safeString(iibb)

	return &profile, err
}

// Función auxiliar para convertir `sql.NullString` a `string`
func safeString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// UpdateProfile actualiza un perfil existente.
func (r *ProfileRepository) UpdateProfile(ctx context.Context, profileId string, profile *domain.Profile) error {
	query := `
		UPDATE profile SET 
			profileName = $1,profileMail = $2, phone = $3, updatedDate =$4
		WHERE profileid = $5`

	//log.Println("🔍 Ejecutando SQL Update con profileId:", profile.ProfileID)
	//log.Println("🔍 Ejecutando SQL Update con profileId:", profile.UpdatedDate)
	result, err := r.DB.ExecContext(ctx, query,
		profile.ProfileName,
		profile.ProfileMail,
		profile.Phone,
		profile.UpdatedDate,
		profile.ProfileID,
	)
	if err != nil {
		//log.Println("Error en la ejecución de la consulta SQL:", err) // Agrega este log
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		//log.Println("⚠️ No se actualizó ninguna fila. ¿El profileId existe?")
		return fmt.Errorf("no se encontró el perfil con profileId: %s", profile.ProfileID)
	}

	//log.Println("✅ Perfil actualizado correctamente:", profile.ProfileID)
	return nil
}

func (r *ProfileRepository) UpdateFiscalData(ctx context.Context, userid string, profile *domain.Profile) error { // este va a ser el updateFisicalData
	query := `
		UPDATE profile SET
			CUIL = $1, fiscalAdress = $2,fiscalCondition = $3, IIBB = $4, updatedDate = $5
		WHERE profileid = $6`

	_, err := r.DB.ExecContext(ctx, query,
		profile.CUIL, profile.FiscalAdress, profile.FiscalCondition, profile.IIBB, profile.UpdatedDate, profile.ProfileID,
	)
	return err
}

func (r *ProfileRepository) UpdateProfileImage(ctx context.Context, userId, profileID, profileImage string) error {
	query := `UPDATE profile SET profileImage = $1 WHERE userId = $2`
	_, err := r.DB.ExecContext(ctx, query, profileImage, userId)
	if err != nil {
		//log.Println("❌ Error actualizando la imagen de perfil:", err)
		return err
	}
	return nil
}

// DeleteProfile elimina un perfil por su ID.
func (r *ProfileRepository) DeleteProfile(ctx context.Context, profileID uuid.UUID) error {
	query := `DELETE FROM profile WHERE profileId = $1`
	_, err := r.DB.ExecContext(ctx, query, profileID)
	return err
}

func (r *ProfileRepository) GetByUserID(ctx context.Context, userId string) (*domain.Profile, error) {

	var profile domain.Profile
	query := "SELECT profileid, userid, profilename, profilemail, phone, ProfilePoints, ProfileLevel FROM profile WHERE userId = $1"

	err := r.DB.QueryRowContext(ctx, query, userId).Scan(
		&profile.ProfileID,
		&profile.UserID,
		&profile.ProfileName,
		&profile.ProfileMail,
		&profile.Phone,
		&profile.ProfilePoints, //kike agrego estos dos para hacer la sincronizacion de puntos
		&profile.ProfileLevel,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Si no se encuentra, devolvemos nil
		}
		return nil, err
	}

	return &profile, nil
}
