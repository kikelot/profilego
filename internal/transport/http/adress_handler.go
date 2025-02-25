package http

import (
	"fmt"
	//"log"
	"net/http"
	"profilego/internal/domain"
	"profilego/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdressHandler struct {
	addressService service.AddressService
}

func NewAdressHandler(addressService service.AddressService) *AdressHandler {
	return &AdressHandler{addressService: addressService}
}

// Obtener todas las direcciones de un perfil
func (h *AdressHandler) GetAdressesByProfile(c *gin.Context) {
	fmt.Println("Todos los parámetros recibidos:", c.Params) // 🔍 Debug
	profileID := c.Param("userId")
	fmt.Println("ProfileID recibido:", profileID) // 🔍 Debug

	profileId, err := uuid.Parse(profileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de perfil inválido"})
		return
	}

	adresses, err := h.addressService.GetAddressesByProfile(c.Request.Context(), profileId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, adresses)
}

func (h *AdressHandler) CreateAddress(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId es requerido"})
		return
	}

	var newAddress domain.Address
	if err := c.ShouldBindJSON(&newAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	ctx := c.Request.Context()
	err := h.addressService.CreateAddress(ctx, userId, &newAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dirección creada correctamente"})
}

func (h *AdressHandler) GetAddress(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId es requerido"})
		return
	}

	ctx := c.Request.Context()
	address, err := h.addressService.GetAddress(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if address == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontró una dirección activa"})
		return
	}

	c.JSON(http.StatusOK, address)
}

func (h *AdressHandler) UpdateAddress(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId es requerido"})
		return
	}

	var updateAddress domain.Address
	if err := c.ShouldBindJSON(&updateAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	ctx := c.Request.Context()
	err := h.addressService.UpdateAddress(ctx, userId, &updateAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dirección actualizada correctamente"})
}

func (h *AdressHandler) DeleteAddress(c *gin.Context) {
	userId := c.Param("userId")

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId es requerido"})
		return
	}

	// Leer el body del request
	var requestBody struct {
		AddressID     string `json:"addressId"`
		ActiveAddress bool   `json:"activeAddress"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de request inválido"})
		return
	}

	// Validar que addressId no esté vacío
	if requestBody.AddressID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "addressId es requerido"})
		return
	}

	ctx := c.Request.Context()
	err := h.addressService.DeleteAddress(ctx, requestBody.AddressID, requestBody.ActiveAddress, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dirección eliminada correctamente"})
}

func (h *AdressHandler) RegisterRoutes(router *gin.RouterGroup) {
	adressGroup := router.Group("/address")
	{
		adressGroup.POST("/:userId/createAddress", h.CreateAddress)
		adressGroup.GET("/:userId/getAddress", h.GetAddress)
		adressGroup.POST("/:userId/updateAddress", h.UpdateAddress)
		adressGroup.PUT("/:userId/deleteAddress", h.DeleteAddress)
		adressGroup.GET("/:userId", h.GetAdressesByProfile)
	}
}
