package service

import (
	"context"
	"errors"
	"time"

	"profilego/internal/domain"
	"profilego/internal/repository"

	"github.com/google/uuid"
)

// AddressService maneja la lógica de negocio para direcciones.
type AddressService struct {
	Repo        repository.AddressRepository
	ProfileRepo repository.ProfileRepository
}

// NewAddressService crea una nueva instancia de AddressService.
func NewAddressService(repo repository.AddressRepository, profileRepo repository.ProfileRepository) *AddressService {
	return &AddressService{
		Repo:        repo,
		ProfileRepo: profileRepo,
	}
}

// CreateAddress crea una nueva dirección.
/*func (s *AddressService) CreateAddress(ctx context.Context, address *domain.Address) (*domain.Address, error){
	// Validación básica
	if address.Street == "" || address.Number == 0 {
		return errors.New("la calle y el número son obligatorios")
	}

	address.AddressID = uuid.New()
	address.CreationDate = time.Now()
	address.UpdatedDate = time.Now()

	return s.Repo.CreateAddress(ctx, address)
}*/

func (s *AddressService) CreateAddress(ctx context.Context, userId string, address *domain.Address) error {
	profile, err := s.ProfileRepo.GetByUserID(ctx, userId)
	if err != nil || profile == nil {
		return errors.New("perfil no encontrado")
	}

	address.AddressID = uuid.New()
	address.IdProfile = profile.ProfileID
	address.ActiveAddress = true
	address.CreationDate = time.Now()
	address.UpdatedDate = time.Now()

	return s.Repo.CreateAddress(ctx, address)
}

func (s *AddressService) GetAddress(ctx context.Context, userId string) (*domain.Address, error) {
	profile, err := s.ProfileRepo.GetByUserID(ctx, userId)
	if err != nil || profile == nil {
		return nil, errors.New("perfil no encontrado")
	}

	return s.Repo.GetAddress(ctx, profile.ProfileID)
}

func (s *AddressService) UpdateAddress(ctx context.Context, userId string, address *domain.Address) error {
	profile, err := s.ProfileRepo.GetByUserID(ctx, userId)
	if err != nil || profile == nil {
		return errors.New("perfil no encontrado")
	}

	address.IdProfile = profile.ProfileID
	address.UpdatedDate = time.Now()

	return s.Repo.UpdateAddress(ctx, address)
}

func (s *AddressService) DeleteAddress(ctx context.Context, addressId string, activeAddress bool, userId string) error {
	profile, err := s.ProfileRepo.GetByUserID(ctx, userId)
	if err != nil || profile == nil {
		return errors.New("perfil no encontrado")
	}

	return s.Repo.DeleteAddress(ctx, addressId, activeAddress, profile.ProfileID)
}

// GetAddressesByProfile obtiene todas las direcciones de un perfil.
func (s *AddressService) GetAddressesByProfile(ctx context.Context, profileID uuid.UUID) ([]domain.Address, error) {
	return s.Repo.GetAddressesByProfile(ctx, profileID)
}
