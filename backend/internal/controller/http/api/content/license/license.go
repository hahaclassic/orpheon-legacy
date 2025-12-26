package license_ctrl

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ctxclaims "github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/claims"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/license"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
)

type LicenseController struct {
	licenseService license.LicenseService
	authMiddleware gin.HandlerFunc
}

func NewLicenseController(licenseService license.LicenseService, authMiddleware gin.HandlerFunc) *LicenseController {
	return &LicenseController{
		licenseService: licenseService,
		authMiddleware: authMiddleware,
	}
}

func (c *LicenseController) RegisterRoutes(router *gin.RouterGroup) {
	licenses := router.Group("/licenses")
	{
		licenses.GET("/:id", c.GetLicense)
		licenses.GET("", c.GetAllLicenses)

		protected := licenses.Group("/")
		protected.Use(c.authMiddleware)
		{
			protected.POST("", c.CreateLicense)
			protected.PUT("/:id", c.UpdateLicense)
			protected.DELETE("/:id", c.DeleteLicense)
		}
	}
}

// GetLicense godoc
// @Summary Get license information
// @Description Get license information by ID
// @Tags licenses
// @Produce json
// @Param id path string true "License ID"
// @Success 200 {object} entity.License
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/licenses/{id} [get]
func (c *LicenseController) GetLicense(ctx *gin.Context) {
	licenseID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid license ID"})
		return
	}

	license, err := c.licenseService.GetLicenseByID(ctx.Request.Context(), licenseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get license"})
		return
	}

	if license == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "License not found"})
		return
	}

	ctx.JSON(http.StatusOK, license)
}

// GetAllLicenses godoc
// @Summary Get all licenses
// @Description Get all licenses
// @Tags licenses
// @Produce json
// @Success 200 {array} entity.License
// @Failure 500 {object} gin.H
// @Router /api/v1/licenses [get]
func (c *LicenseController) GetAllLicenses(ctx *gin.Context) {
	licenses, err := c.licenseService.GetAllLicenses(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all licenses"})
		return
	}

	ctx.JSON(http.StatusOK, licenses)
}

// CreateLicense godoc
// @Summary Create a new license
// @Description Create a new license with the provided information
// @Tags licenses
// @Accept json
// @Produce json
// @Param license body entity.License true "License information"
// @Security BearerAuth
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/licenses [post]
func (c *LicenseController) CreateLicense(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var license entity.License
	if err := ctx.ShouldBindJSON(&license); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := c.licenseService.CreateLicense(ctx.Request.Context(), claims, &license); err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to create license"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create license"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "License created successfully"})
}

// UpdateLicense godoc
// @Summary Update license information
// @Description Update license information by ID
// @Tags licenses
// @Accept json
// @Produce json
// @Param id path string true "License ID"
// @Param license body entity.License true "License information"
// @Security BearerAuth
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/licenses/{id} [put]
func (c *LicenseController) UpdateLicense(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	licenseID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid license ID"})
		return
	}

	var license entity.License
	if err := ctx.ShouldBindJSON(&license); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	license.ID = licenseID

	if err := c.licenseService.UpdateLicense(ctx.Request.Context(), claims, &license); err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to update license"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update license"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "License updated successfully"})
}

// DeleteLicense godoc
// @Summary Delete license
// @Description Delete license by ID
// @Tags licenses
// @Produce json
// @Param id path string true "License ID"
// @Security BearerAuth
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/licenses/{id} [delete]
func (c *LicenseController) DeleteLicense(ctx *gin.Context) {
	claims := ctxclaims.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	licenseID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid license ID"})
		return
	}

	if err := c.licenseService.DeleteLicense(ctx.Request.Context(), claims, licenseID); err != nil {
		if errors.Is(err, commonerr.ErrForbidden) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to delete license"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete license"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "License deleted successfully"})
}
