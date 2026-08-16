package staff

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

// GET /staff
func (h *Handler) List(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_staff_page", "Please provide a valid page number.", map[string]string{"page": "validation.invalid"})
		return
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_staff_per_page", "Please provide a valid page size.", map[string]string{"per_page": "validation.invalid"})
		return
	}
	if _, ok := AllowedPerPage[perPage]; !ok {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_staff_per_page", "Please provide a valid page size.", map[string]string{"per_page": "validation.invalid"})
		return
	}

	role := c.Query("role")
	if role != "" && role != RoleScanner && role != RoleReceptionist {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_staff_role", "Please provide a valid role filter.", map[string]string{"role": "validation.invalid"})
		return
	}

	sort := c.DefaultQuery("sort", "created_at")
	if _, ok := AllowedSort[sort]; !ok {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_staff_sort", "Please provide a valid sort field.", map[string]string{"sort": "validation.invalid"})
		return
	}

	direction := c.DefaultQuery("direction", "desc")
	if direction != "asc" && direction != "desc" {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_staff_direction", "Please provide a valid sort direction.", map[string]string{"direction": "validation.invalid"})
		return
	}

	results, err := h.service.List(c.Request.Context(), ownerUserID, ListParams{
		Page:      page,
		PerPage:   perPage,
		Role:      role,
		Search:    c.Query("q"),
		Sort:      sort,
		Direction: direction,
	})
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "staff_list_failed", "Failed to load staff. Please try again.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "staff retrieved", results)
}

// POST /staff
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
		case errors.Is(err, ErrCenterMismatch):
			utils.ErrorResponseWithCode(c, http.StatusForbidden, "center_mismatch", "You can only add staff to your own center.")
		case errors.Is(err, ErrEmailTaken):
			utils.ValidationErrorResponse(c, http.StatusConflict, "email_taken", "Email is already in use.", map[string]string{"email": "validation.taken"})
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "staff_create_failed", "Failed to create staff member. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "staff created", resp)
}

// PUT /staff/:id
func (h *Handler) Update(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	staffID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_staff_id", "Please provide a valid staff id.")
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.service.Update(c.Request.Context(), ownerUserID, uint(staffID), &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrCenterMismatch):
			utils.ErrorResponseWithCode(c, http.StatusForbidden, "center_mismatch", "You can only manage staff in your own center.")
		case errors.Is(err, ErrStaffNotFound):
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "staff_not_found", "Staff member not found.")
		case errors.Is(err, ErrEmailTaken):
			utils.ValidationErrorResponse(c, http.StatusConflict, "email_taken", "Email is already in use.", map[string]string{"email": "validation.taken"})
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "staff_update_failed", "Failed to update staff member. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "staff updated", resp)
}

// DELETE /staff/:id
func (h *Handler) Delete(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	staffID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_staff_id", "Please provide a valid staff id.")
		return
	}

	if err := h.service.Delete(c.Request.Context(), ownerUserID, uint(staffID)); err != nil {
		switch {
		case errors.Is(err, ErrStaffNotFound):
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "staff_not_found", "Staff member not found.")
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "staff_delete_failed", "Failed to delete staff member. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "staff deleted", nil)
}
