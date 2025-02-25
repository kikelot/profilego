package http

import (
	"fmt"
	"log"

	//"log"
	"net/http"
	"profilego/internal/domain"

	//"strconv"

	//"profilego/internal/middleware"
	"profilego/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	profileService service.ProfileService
}

func NewProfileHandler(profileService service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	// Extraer userId del middleware (lo extrae del token JWT)
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo recuperar el userId"})
		return
	}

	// Parsear el JSON del body
	var profile domain.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Asignar el userId al perfil antes de crearlo
	profile.UserID = userId.(string)
	profile.ProfileID = uuid.New()

	// Llamar al servicio para crear el perfil
	if err := h.profileService.CreateProfile(c.Request.Context(), userId.(string), &profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Devolver respuesta con perfil creado
	c.JSON(http.StatusCreated, profile)
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userId := c.Param("userId")

	/*if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo recuperar el userId"})
		return
	}*/

	profile, err := h.profileService.GetProfile(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Perfil no encontrado"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userId := c.Param("userId") // 📌 Obtener el userId desde la URL

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId es requerido"})
		return
	}

	var updateProfileRequest domain.Profile

	if err := c.ShouldBindJSON(&updateProfileRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	ctx := c.Request.Context()
	err := h.profileService.UpdateProfile(ctx, userId, &updateProfileRequest) // 🔹 Pasamos userId
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Perfil actualizado correctamente"})
}

func (h *ProfileHandler) UpdateFiscalData(c *gin.Context) {
	// 🔹 Extraer `userId` del token (100% confiable)
	userIdToken, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo recuperar el userId"})
		return
	}

	// 🔹 Extraer `userId` de la URL
	userIdParam := c.Param("userId")
	if userIdParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId es requerido en la URL"})
		return
	}

	// 🔹 Buscar el perfil en la base de datos usando el `userIdParam`
	profile, err := h.profileService.Repo.GetByUserID(c.Request.Context(), userIdParam)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Perfil no encontrado"})
		return
	}

	// 🔹 Validar que el perfil le pertenece al usuario autenticado (comparar con el `userIdToken`)
	if profile.UserID != fmt.Sprintf("%s", userIdToken) {

		c.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos para modificar este perfil"})
		return
	}
	// 🔹 Parsear JSON del body
	var updateFiscalData domain.Profile
	if err := c.ShouldBindJSON(&updateFiscalData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// 🔹 Asegurar que `profileId` y `userId` sean los correctos
	updateFiscalData.ProfileID = profile.ProfileID
	updateFiscalData.UserID = profile.UserID // ⚠ Importante: Usar `profile.UserID`, ya que es el validado.

	// 🔹 Llamar al servicio para actualizar los datos fiscales
	ctx := c.Request.Context()
	err = h.profileService.UpdateFiscalData(ctx, profile.UserID, &updateFiscalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 🔹 Respuesta de éxito
	c.JSON(http.StatusOK, gin.H{"message": "Datos fiscales actualizados correctamente"})
}

func (h *ProfileHandler) UpdateProfileImage(c *gin.Context) {
	userId := c.Param("userId") // Obtener el userId desde la URL

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId es requerido"})
		return
	}

	var req struct {
		ProfileID    string `json:"profileId" binding:"required"`
		ProfileImage string `json:"profileImage" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Llamar al servicio para actualizar la imagen
	err := h.profileService.UpdateProfileImage(c.Request.Context(), userId, req.ProfileID, req.ProfileImage)
	if err != nil {
		if err.Error() == "usuario no encontrado" { // 📌 Ahora sí manejamos el error en el servicio
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar la imagen"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Imagen actualizada correctamente"})
}

func (h *ProfileHandler) UpdateProfilePoints(c *gin.Context) {
	log.Println("ESTA ENTRANDO ACA")
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo recuperar el userId"})
		return
	}
	log.Println("ESTA ENTRANDO ACA")
	var req domain.Profile

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	err := h.profileService.UpdateProfilePoints(c.Request.Context(), userId, &req)
	if err != nil {
		log.Println("ESTA ENTRANDO ACA")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Puntos del perfil actualizados correctamente"})
}

func (h *ProfileHandler) UpdateProfileLevel(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo recuperar el userId"})
		return
	}

	var req domain.Profile

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	err := h.profileService.UpdateProfileLevel(c.Request.Context(), userId, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nivel del perfil actualizado correctamente"})
}

func (h *ProfileHandler) RegisterRoutes(router *gin.RouterGroup) {
	profileGroup := router.Group("/profiles")
	{
		profileGroup.POST("/:userId/create", h.CreateProfile)
		profileGroup.GET("/:userId", h.GetProfile)
		profileGroup.POST("/:userId/updateProfile", h.UpdateProfile)
		profileGroup.POST("/:userId/updateFiscalData", h.UpdateFiscalData)
		profileGroup.POST("/:userId/updateImage", h.UpdateProfileImage)
		profileGroup.POST("/:userId/updateProfilePoints", h.UpdateProfilePoints)
		profileGroup.POST("/:userId/updateProfileLevel", h.UpdateProfileLevel)
	}
}
