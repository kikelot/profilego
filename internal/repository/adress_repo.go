package repository

import (
	"context"
	"database/sql"

	//"errors"

	"profilego/internal/domain"

	"github.com/google/uuid"
)

// AddressRepository define los métodos para interactuar con la base de datos.
type AddressRepository struct {
	DB *sql.DB
}

// NewAddressRepository crea una nueva instancia del repositorio.
func NewAddressRepository(db *sql.DB) *AddressRepository {
	return &AddressRepository{DB: db}
}

func (r *AddressRepository) CreateAddress(ctx context.Context, address *domain.Address) error {
	query := `INSERT INTO address (addressId, CP, street, number, floor, mainAddress, creationDate, updatedDate, activeAddress, idProfile)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.DB.ExecContext(ctx, query,
		address.AddressID,
		address.CP,
		address.Street,
		address.Number,
		address.Floor,
		address.MainAddress,
		address.CreationDate,
		address.UpdatedDate,
		address.ActiveAddress,
		address.IdProfile,
	)
	return err
}

func (r *AddressRepository) GetAddress(ctx context.Context, idProfile uuid.UUID) (*domain.Address, error) {
	var address domain.Address
	query := `SELECT addressId, CP, street, number, floor, mainAddress, creationDate, updatedDate, activeAddress, idProfile 
	          FROM address WHERE idProfile = $1 AND activeAddress = TRUE`

	err := r.DB.QueryRowContext(ctx, query, idProfile).Scan(
		&address.AddressID,
		&address.CP,
		&address.Street,
		&address.Number,
		&address.Floor,
		&address.MainAddress,
		&address.CreationDate,
		&address.UpdatedDate,
		&address.ActiveAddress,
		&address.IdProfile,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &address, err
}

func (r *AddressRepository) UpdateAddress(ctx context.Context, address *domain.Address) error {
	query := `UPDATE address SET CP = $1, street = $2, number = $3, floor = $4, mainAddress = $5, updatedDate = $6
			WHERE idProfile = $7 AND activeAddress = TRUE`

	_, err := r.DB.ExecContext(ctx, query,
		address.CP,
		address.Street,
		address.Number,
		address.Floor,
		address.MainAddress,
		address.UpdatedDate,
		address.IdProfile,
	)
	return err
}

func (r *AddressRepository) DeleteAddress(ctx context.Context, addressId string, activeAddress bool, idprofile uuid.UUID) error {
	query := `UPDATE address SET activeaddress = $2
			WHERE addressid = $1 AND idprofile =$3`
	_, err := r.DB.ExecContext(ctx, query, addressId, activeAddress, idprofile)
	return err
}

// GetAddressesByProfile obtiene todas las direcciones de un perfil.
func (r *AddressRepository) GetAddressesByProfile(ctx context.Context, IdProfile uuid.UUID) ([]domain.Address, error) {
	query := `SELECT addressId, CP, street, number, floor, mainAddress,
			creationDate, updatedDate, activeAddress, idprofile
			FROM address WHERE idProfile= $1`

	rows, err := r.DB.QueryContext(ctx, query, IdProfile)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []domain.Address
	for rows.Next() {
		var address domain.Address
		err := rows.Scan(
			&address.AddressID, &address.CP, &address.Street, &address.Number, &address.Floor,
			&address.MainAddress, &address.CreationDate, &address.UpdatedDate,
			&address.ActiveAddress, &address.IdProfile,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}
