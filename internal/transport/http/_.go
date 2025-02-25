package http

import (
	"log"
	"net/http"
	"strconv"

	"profilego/internal/domain"
	"profilego/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	//"golang.org/x/tools/go/analysis/passes/printf"
)

type ProfileHandler struct {
	profileService service.ProfileService
}

func NewProfileHandler(profileService service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {
	var profile domain.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile.ProfileID = uuid.New()

	/*if err := h.profileUsecase.CreateProfile(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}*/
	if err := h.profileService.CreateProfile(c.Request.Context(), &profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Obtener el ID desde la URL
	idStr := c.Param("id")

	// Convertir el ID a uuid.UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Llamar al servicio con contexto y UUID
	profile, err := h.profileService.GetProfile(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	// Obtener userId del contexto
	/*userId, exists := c.Get("userId")  // Kike comenta, validacion para un userId que no existe
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo recuperar el userId"})
		return
	}*/

	// Parsear JSON del body
	var updateProfileRequest domain.Profile

	if err := c.ShouldBindJSON(&updateProfileRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	log.Println("📩 Recibido en handler - profileId:", updateProfileRequest.ProfileID)
	// Llamar al servicio para actualizar el perfil
	ctx := c.Request.Context()
	err := h.profileService.UpdateProfile(ctx, &updateProfileRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//c.JSON(http.StatusOK, gin.H{"message": "Perfil actualizado correctamente"})

	c.JSON(http.StatusOK, gin.H{
		"message": "Perfil actualizado correctamente",
	})
}

func (h *ProfileHandler) UpdateFiscalData(c *gin.Context) {
	// Obtener userId del contexto
	/*userId, exists := c.Get("userId")  // Kike comenta, validacion para un userId que no existe
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo recuperar el userId"})
		return
	}*/

	// Parsear JSON del body
	var updateFiscalData domain.Profile

	if err := c.ShouldBindJSON(&updateFiscalData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	log.Println("📩 Recibido en handler - profileId:", updateFiscalData.ProfileID)
	// Llamar al servicio para actualizar el perfil
	ctx := c.Request.Context()
	err := h.profileService.UpdateFiscalData(ctx, &updateFiscalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//c.JSON(http.StatusOK, gin.H{"message": "Perfil actualizado correctamente"})

	c.JSON(http.StatusOK, gin.H{
		"message": "Perfil actualizado correctamente",
	})
}

func (h *ProfileHandler) UpdateProfilePoints(c *gin.Context) {
	// Obtener userId del contexto
	/*userId, exists := c.Get("userId")  // Kike comenta, validacion para un userId que no existe
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo recuperar el userId"})
		return
	}*/

	// Parsear JSON del body
	//var updateProfilePointsRequest domain.Profile

	var req struct {
		ProfileID     string `json:"profileId" binding:"required"`
		ProfilePoints string `json:"profilePoints" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	log.Println("📩 Recibido en handler - profileId:", req.ProfileID)
	// Llamar al servicio para actualizar el perfil

	i, err := strconv.Atoi(req.ProfilePoints)
	if err == nil {
		log.Println("ESTA ENTRANDO ACA", req.ProfileID, req.ProfilePoints, i)
		err = h.profileService.UpdateProfilePoints(c.Request.Context(), req.ProfileID, i)
		if err != nil {
			log.Println("ESTA ENTRANDO ACA?", req.ProfileID, req.ProfilePoints, i)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	//c.JSON(http.StatusOK, gin.H{"message": "Perfil actualizado correctamente"})

	c.JSON(http.StatusOK, gin.H{
		"message": "Perfil actualizado correctamente",
	})
}

func (h *ProfileHandler) UpdateProfileImage(c *gin.Context) {
	//idStr := c.Param("id")
	//userID := c.Param("userId")

	var req struct {
		ProfileID    string `json:"profileId" binding:"required"`
		ProfileImage string `json:"profileImage" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	err := h.profileService.UpdateProfileImage(c.Request.Context(), req.ProfileID, req.ProfileImage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar la imagen"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Imagen actualizada correctamente"})
}

func (h *ProfileHandler) RegisterRoutes(router *gin.RouterGroup) {
	profileGroup := router.Group("/profiles")
	{
		//profileGroup.POST("/:userid/create", h.CreateProfile) Kike Comento momentaneamente, todavia no se recuperar el userid
		profileGroup.POST("/create", h.CreateProfile)
		profileGroup.GET("/:id", h.GetProfile)
		profileGroup.POST("/updateProfile", h.UpdateProfile)
		profileGroup.POST("/updateFiscalData", h.UpdateFiscalData)
		profileGroup.POST("/updateImage", h.UpdateProfileImage)
		profileGroup.POST("/updateProfilePoints", h.UpdateProfilePoints)
	}
}
