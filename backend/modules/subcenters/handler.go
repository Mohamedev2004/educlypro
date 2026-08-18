package subcenters

import (
	"educlypro/shared/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GET /subcenters
func (h *Handler) List(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	resp, err := h.service.List(c.Request.Context(), ownerUserID)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "subcenters_list_failed", "Failed to load sub-centers. Please try again.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "sub-centers retrieved", resp)
}

// POST /subcenters
func (h *Handler) Create(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.service.Create(c.Request.Context(), ownerUserID, &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrSubCenterExists):
			utils.ValidationErrorResponse(c, http.StatusConflict, "subcenter_exists", "A sub-center with this name already exists.", map[string]string{"name": "validation.taken"})
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "subcenters_create_failed", "Failed to create sub-center. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "sub-center created", resp)
}

// PUT /subcenters/:id
func (h *Handler) Update(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	subCenterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_subcenter_id", "Please provide a valid sub-center id.")
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.service.Update(c.Request.Context(), ownerUserID, uint(subCenterID), &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrSubCenterNotFound):
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "subcenter_not_found", "Sub-center not found.")
		case errors.Is(err, ErrSubCenterExists):
			utils.ValidationErrorResponse(c, http.StatusConflict, "subcenter_exists", "A sub-center with this name already exists.", map[string]string{"name": "validation.taken"})
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "subcenters_update_failed", "Failed to update sub-center. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "sub-center updated", resp)
}

// DELETE /subcenters/:id
func (h *Handler) Delete(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	subCenterID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_subcenter_id", "Please provide a valid sub-center id.")
		return
	}

	if err := h.service.Delete(c.Request.Context(), ownerUserID, uint(subCenterID)); err != nil {
		switch {
		case errors.Is(err, ErrSubCenterNotFound):
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "subcenter_not_found", "Sub-center not found.")
		case errors.Is(err, ErrSubCenterHasStaff):
			utils.ErrorResponseWithCode(c, http.StatusConflict, "subcenter_has_staff", "Reassign or remove its staff before deleting this sub-center.")
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "subcenters_delete_failed", "Failed to delete sub-center. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "sub-center deleted", nil)
}
