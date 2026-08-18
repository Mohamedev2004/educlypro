package centers

import (
	"errors"
	"net/http"
	"strconv"

	"educlypro/modules/staff"
	"educlypro/shared/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GET /centers
func (h *Handler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_centers_page", "Please provide a valid page number.", map[string]string{"page": "validation.invalid"})
		return
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_centers_per_page", "Please provide a valid page size.", map[string]string{"per_page": "validation.invalid"})
		return
	}
	if _, ok := AllowedPerPage[perPage]; !ok {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_centers_per_page", "Please provide a valid page size.", map[string]string{"per_page": "validation.invalid"})
		return
	}

	sort := c.DefaultQuery("sort", "created_at")
	if _, ok := AllowedSort[sort]; !ok {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_centers_sort", "Please provide a valid sort field.", map[string]string{"sort": "validation.invalid"})
		return
	}

	direction := c.DefaultQuery("direction", "desc")
	if direction != "asc" && direction != "desc" {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_centers_direction", "Please provide a valid sort direction.", map[string]string{"direction": "validation.invalid"})
		return
	}

	results, err := h.service.List(c.Request.Context(), ListParams{
		Page:      page,
		PerPage:   perPage,
		Search:    c.Query("q"),
		Sort:      sort,
		Direction: direction,
	})
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "centers_list_failed", "Failed to load centers. Please try again.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "centers retrieved", results)
}

// POST /centers
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "centers_create_failed", "Failed to create center. Please try again.")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "center created", resp)
}

// GET /centers/:slug
func (h *Handler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	resp, err := h.service.GetDetail(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, ErrCenterNotFound) {
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "center_not_found", "Center not found.")
			return
		}
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "centers_get_failed", "Failed to load center. Please try again.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "center retrieved", resp)
}

// POST /centers/:slug/owner
func (h *Handler) AddOwner(c *gin.Context) {
	slug := c.Param("slug")

	var req CreateOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.service.AddOwner(c.Request.Context(), slug, &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrCenterNotFound):
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "center_not_found", "Center not found.")
		case errors.Is(err, ErrOwnerAlreadyAssigned):
			utils.ErrorResponseWithCode(c, http.StatusConflict, "owner_already_assigned", "This center already has an owner.")
		case errors.Is(err, ErrEmailTaken):
			utils.ValidationErrorResponse(c, http.StatusConflict, "email_taken", "Email is already in use.", map[string]string{"email": "validation.taken"})
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "centers_add_owner_failed", "Failed to add owner. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "owner added", resp)
}

// POST /centers/:slug/staff
func (h *Handler) AddStaff(c *gin.Context) {
	slug := c.Param("slug")

	var req AddStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	actorUserID := c.MustGet("userID").(uint)

	resp, err := h.service.AddStaff(c.Request.Context(), slug, actorUserID, &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrCenterNotFound):
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "center_not_found", "Center not found.")
		case errors.Is(err, staff.ErrEmailTaken):
			utils.ValidationErrorResponse(c, http.StatusConflict, "email_taken", "Email is already in use.", map[string]string{"email": "validation.taken"})
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "centers_add_staff_failed", "Failed to add staff member. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "staff added", resp)
}

// GET /centers/:slug/subcenters
func (h *Handler) ListSubCenters(c *gin.Context) {
	slug := c.Param("slug")

	resp, err := h.service.ListSubCenters(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, ErrCenterNotFound) {
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "center_not_found", "Center not found.")
			return
		}
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "centers_list_subcenters_failed", "Failed to load sub-centers. Please try again.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "sub-centers retrieved", resp)
}
