package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	//"strings"
	"time"

	"profilego/internal/domain"
	"profilego/internal/repository"

	//"profilego/internal/transport/mq"

	"github.com/google/uuid"
)

type Publisher interface { //Interfaz para comunicarme con publisher y romper la dependencia directa (error de ciclo infinito de importaciones)
	PublishProfilePoints(profileid string, profilepoints int) error
	PublishMessage(queueName string, message []byte) error
}

type ProfileService struct {
	Repo      repository.ProfileRepository
	Publisher Publisher // Usa la interfaz en lugar de `mq`
}

// NewProfileService crea una nueva instancia de ProfileService.
func NewProfileService(repo repository.ProfileRepository) *ProfileService {
	return &ProfileService{Repo: repo}
}

func NewProfileRabbitService(repo repository.ProfileRepository, publisher Publisher) *ProfileService {
	return &ProfileService{
		Repo:      repo,
		Publisher: publisher,
	}
}

// CreateProfile crea un nuevo perfil.
func (s *ProfileService) CreateProfile(ctx context.Context, userId string, profile *domain.Profile) error {
	// PRIMERO VALIDO QUE EXISTA EL USUARIO
	existingProfile, err := s.Repo.GetByUserID(ctx, userId)
	if err != nil {
		//log.Println(" Error al buscar userId:", err)
		return errors.New("error al buscar el perfil del usuario")
	}

	// 🔹 Si ya existe un perfil para este `userId`, rechazar la creación
	if existingProfile != nil {
		//log.Println("El usuario ya tiene un perfil - userId:", userId)
		return errors.New("el usuario ya tiene un perfil creado")
	}

	//DESPUES LAS VALIDACIONES BASICAS DE LOS DATOS

	if profile.ProfileName == "" || profile.ProfileMail == "" {
		return errors.New("el nombre y el email del perfil son obligatorios")
	}
	if profile.Phone == "" || len(profile.Phone) < 10 {
		return errors.New("el número de telefono es obligatorio y debe tener al menos 10 digitos")
	}

	if _, err := strconv.Atoi(profile.Phone); err != nil {
		return errors.New("el numero de teléfono debe ser numérico")
	}

	profile.ProfileID = uuid.New()
	profile.UserID = userId // Kike por ahora le pongo uno aleatorio pq todavia no se como recuperar el correcto
	profile.CreationDate = time.Now()
	profile.UpdatedDate = time.Now()

	return s.Repo.CreateProfile(ctx, profile)
}

// GetProfile obtiene un perfil por ID.
func (s *ProfileService) GetProfile(ctx context.Context, userId string) (*domain.Profile, error) {

	existingProfile, err := s.Repo.GetByUserID(ctx, userId)
	if err != nil {
		//log.Println(" Error al buscar userId:", err)
		return nil, errors.New("error al buscar el perfil del usuario")
	}

	if existingProfile == nil {
		//log.Println("⚠ No se encontró un perfil con userId:", userId)
		return nil, errors.New("usuario no encontrado")
	}

	return s.Repo.GetProfile(ctx, userId)
}

// UpdateProfile actualiza los datos de un perfil.
func (s *ProfileService) UpdateProfile(ctx context.Context, userId string, profile *domain.Profile) error {

	/*profile.UpdatedDate = time.Now()
	return s.Repo.UpdateProfile(ctx, profile)*/
	// Validación básica

	//VALIDO QUE EXISTA EL USERID ANTES DE ACTUALIZAR EL REGISTRO
	existingProfile, err := s.Repo.GetByUserID(ctx, userId)
	if err != nil {
		//log.Println("❌ Error al buscar userId:", err)
		return errors.New("error al buscar el perfil del usuario")
	}

	if existingProfile == nil {
		//log.Println("⚠ No se encontró un perfil con userId:", userId)
		return errors.New("usuario no encontrado")
	}
	// FIN VALIDA USER

	if profile.ProfileName == "" || profile.ProfileMail == "" {
		return errors.New("el nombre y el email del perfil son obligatorios")
	}
	if profile.Phone == "" || len(profile.Phone) < 10 {
		return errors.New("el número de telefono es obligatorio y debe tener al menos 10 digitos")
	}

	if _, err := strconv.Atoi(profile.Phone); err != nil {
		return errors.New("el numero de teléfono debe ser numérico")
	}
	//log.Println("⏳ Iniciando actualización del perfil:", profile.ProfileID)

	profile.UpdatedDate = time.Now()

	err = s.Repo.UpdateProfile(ctx, userId, profile)
	if err != nil {
		log.Println("❌ Error al actualizar el perfil en la base de datos:", err)
	} else {
		log.Println("✅ Perfil actualizado correctamente:", profile.ProfileID)
	}

	return err

}

// UpdateProfile actualiza los datos fiscales de un perfil.
func (s *ProfileService) UpdateFiscalData(ctx context.Context, userId string, profile *domain.Profile) error {

	/*profile.UpdatedDate = time.Now()
	return s.Repo.UpdateProfile(ctx, profile)*/
	//VALIDO QUE EXISTA EL USERID ANTES DE ACTUALIZAR EL REGISTRO
	existingProfile, err := s.Repo.GetByUserID(ctx, userId)
	if err != nil {
		//log.Println("❌ Error al buscar userId:", err)
		return errors.New("error al buscar el perfil del usuario")
	}

	if existingProfile == nil {
		//log.Println("⚠ No se encontró un perfil con userId:", userId)
		return errors.New("usuario no encontrado")
	}
	// Validación básica
	if profile.CUIL == "" || profile.FiscalAdress == "" || profile.FiscalCondition == "" || profile.IIBB == "" {
		return errors.New("los datos fiscales son obligatorios")
	}
	if len(profile.CUIL) < 11 {
		return errors.New("el CUIL debe tener al menos 11 caracteres")
	}

	if _, err := strconv.Atoi(profile.CUIL); err != nil {
		return errors.New("el CUIL debe ser un valor numérico")
	}
	if _, err := strconv.Atoi(profile.IIBB); err != nil {
		return errors.New("IIBB debe ser un valor numérico")
	}
	//log.Println("⏳ Iniciando actualización del perfil:", profile.ProfileID)

	profile.UpdatedDate = time.Now()
	err = s.Repo.UpdateFiscalData(ctx, userId, profile)
	if err != nil {
		log.Println("❌ Error al actualizar el perfil en la base de datos:", err)
	} else {
		log.Println("✅ Perfil actualizado correctamente:", profile.ProfileID)
	}

	return err

}

func (s *ProfileService) UpdateProfileImage(ctx context.Context, userId string, profileID string, profileImage string) error {

	//VALIDO QUE EXISTA EL USERID ANTES DE ACTUALIZAR EL REGISTRO
	existingProfile, err := s.Repo.GetByUserID(ctx, userId)
	if err != nil {
		//log.Println("❌ Error al buscar userId:", err)
		return errors.New("error al buscar el perfil del usuario")
	}

	if existingProfile == nil {
		//log.Println("⚠ No se encontró un perfil con userId:", userId)
		return errors.New("usuario no encontrado")
	}

	//Validacion básica
	if profileImage == "" {
		return errors.New("debe subir un archivo")
	}
	return s.Repo.UpdateProfileImage(ctx, userId, profileID, profileImage)
}

// DeleteProfile elimina un perfil por su ID.
func (s *ProfileService) DeleteProfile(ctx context.Context, profileID uuid.UUID) error {
	return s.Repo.DeleteProfile(ctx, profileID)
}

// RABBIT
func (s *ProfileService) UpdateProfilePoints(ctx context.Context, userId string, profile *domain.Profile) error {
	//VALIDO QUE EXISTA EL USERID ANTES DE ACTUALIZAR EL REGISTRO
	existingProfile, err := s.Repo.GetByUserID(ctx, userId)
	if err != nil {
		log.Println("❌ Error al buscar userId:", err)
		return errors.New("error al buscar el perfil del usuario")
	}

	if existingProfile == nil {
		//log.Println("⚠ No se encontró un perfil con userId:", userId)
		return errors.New("usuario no encontrado")
	}

	// Validar que el profileId pertenece al userId
	if existingProfile.UserID != userId {
		//log.Println("❌ El profileId no pertenece al userId:", userId)
		return errors.New("perfil no pertenece al usuario")
	}

	// Actualizar puntos en la base de datos...
	err = s.Repo.UpdateProfilePoints(ctx, userId, profile)
	if err != nil {
		return err
	}

	// 📢 Serializar el mensaje a JSON
	message, err := json.Marshal(map[string]interface{}{
		"userId":        userId,
		"profileId":     profile.ProfileID,
		"profilePoints": profile.ProfilePoints,
	})
	if err != nil {
		return err
	}

	if s.Publisher == nil {
		//log.Println("❌ ERROR: s.Publisher es nil, la conexión con RabbitMQ no está inicializada")
		return fmt.Errorf("error: publisher no inicializado")
	}
	// 📢 Publicar evento en RabbitMQ
	errRabbit := s.Publisher.PublishMessage("direct_profile", message)

	log.Println("❌ ESTA LLEGANDO ACA2?")

	if errRabbit != nil {
		//log.Println("❌ Error: s.Publisher es nil, la conexión no está inicializada")
		//log.Println(message)
		return fmt.Errorf("error: publisher no inicializado")
	}

	return err
}

// UpdateProfileLevel actualiza el nivel de un perfil
func (s *ProfileService) UpdateProfileLevel(ctx context.Context, userId string, profile *domain.Profile) error {
	// Verificar que el usuario existe en la base de datos
	existingProfile, err := s.Repo.GetByUserID(ctx, userId)
	if err != nil {
		log.Println("❌ Error al buscar userId:", err)
		return errors.New("error al buscar el perfil del usuario")
	}

	if existingProfile == nil {
		log.Println("⚠ No se encontró un perfil con userId:", userId)
		return errors.New("usuario no encontrado")
	}

	// Validar que el profileId pertenece al userId
	if existingProfile.UserID != userId {
		log.Println("❌ El profileId no pertenece al userId:", userId)
		return errors.New("perfil no pertenece al usuario")
	}

	// Actualizar nivel en la base de datos...
	log.Println("📩 Actualizando nivel - userId:", userId, "profileLevel:", profile.ProfileLevel)
	err = s.Repo.UpdateProfileLevel(ctx, userId, profile)
	if err != nil {
		return err
	}

	// 📢 Serializar el mensaje a JSON Kike comento esta parte para no publicar otro mensaje cuando subio el level y entra en bucle
	/*message, err := json.Marshal(map[string]interface{}{
		"userId":       userId,
		"profileId":    profile.ProfileID,
		"profileLevel": profile.ProfileLevel,
	})
	if err != nil {
		return err
	}*/

	if s.Publisher == nil {
		log.Println("❌ ERROR: RabbitMQ no está inicializado")
		return fmt.Errorf("error: publisher no inicializado")
	}

	// 📢 Publicar evento en RabbitMQ   Kike comento esta parte para no publicar otro mensaje cuando subio el level y entra en bucle
	//errRabbit := s.Publisher.PublishMessage("direct_profile", message)

	/*if errRabbit != nil {
		//log.Println("❌ Error publicando en RabbitMQ")
		return fmt.Errorf("error al publicar mensaje en RabbitMQ")
	}*/

	return nil
}
