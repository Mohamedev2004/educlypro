package teachers

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

// GET /teachers
func (h *Handler) List(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_teacher_page", "Please provide a valid page number.", map[string]string{"page": "validation.invalid"})
		return
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_teacher_per_page", "Please provide a valid page size.", map[string]string{"per_page": "validation.invalid"})
		return
	}
	if _, ok := AllowedPerPage[perPage]; !ok {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_teacher_per_page", "Please provide a valid page size.", map[string]string{"per_page": "validation.invalid"})
		return
	}

	sort := c.DefaultQuery("sort", "created_at")
	if _, ok := AllowedSort[sort]; !ok {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_teacher_sort", "Please provide a valid sort field.", map[string]string{"sort": "validation.invalid"})
		return
	}

	direction := c.DefaultQuery("direction", "desc")
	if direction != "asc" && direction != "desc" {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_teacher_direction", "Please provide a valid sort direction.", map[string]string{"direction": "validation.invalid"})
		return
	}

	results, err := h.service.List(c.Request.Context(), ownerUserID, ListParams{
		Page:      page,
		PerPage:   perPage,
		Search:    c.Query("q"),
		Sort:      sort,
		Direction: direction,
	})
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "teacher_list_failed", "Failed to load teachers. Please try again.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "teachers retrieved", results)
}

// GET /teachers/:slug
func (h *Handler) GetBySlug(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)
	slug := c.Param("slug")

	resp, err := h.service.GetBySlug(c.Request.Context(), ownerUserID, slug)
	if err != nil {
		if errors.Is(err, ErrTeacherNotFound) {
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "teacher_not_found", "Teacher not found.")
			return
		}
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "teacher_get_failed", "Failed to load teacher. Please try again.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "teacher retrieved", resp)
}

// POST /teachers
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
		case errors.Is(err, ErrEmailTaken):
			utils.ValidationErrorResponse(c, http.StatusConflict, "email_taken", "Email is already in use.", map[string]string{"email": "validation.taken"})
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "teacher_create_failed", "Failed to create teacher. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "teacher created", resp)
}

// PUT /teachers/:id
func (h *Handler) Update(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	teacherID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_teacher_id", "Please provide a valid teacher id.")
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.service.Update(c.Request.Context(), ownerUserID, uint(teacherID), &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrTeacherNotFound):
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "teacher_not_found", "Teacher not found.")
		case errors.Is(err, ErrEmailTaken):
			utils.ValidationErrorResponse(c, http.StatusConflict, "email_taken", "Email is already in use.", map[string]string{"email": "validation.taken"})
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "teacher_update_failed", "Failed to update teacher. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "teacher updated", resp)
}

// PUT /teachers/:id/subjects
func (h *Handler) UpdateSubjects(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	teacherID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_teacher_id", "Please provide a valid teacher id.")
		return
	}

	var req UpdateSubjectsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.service.UpdateSubjects(c.Request.Context(), ownerUserID, uint(teacherID), &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrTeacherNotFound):
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "teacher_not_found", "Teacher not found.")
		case errors.Is(err, ErrInvalidSubjectSelection):
			utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_subject_selection", "One or more selected subjects are invalid.", map[string]string{"subject_ids": "validation.invalid"})
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "teacher_update_failed", "Failed to update teacher. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "teacher subjects updated", resp)
}

// PUT /teachers/:id/classes
func (h *Handler) UpdateClasses(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	teacherID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_teacher_id", "Please provide a valid teacher id.")
		return
	}

	var req UpdateClassesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.service.UpdateClasses(c.Request.Context(), ownerUserID, uint(teacherID), &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrTeacherNotFound):
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "teacher_not_found", "Teacher not found.")
		case errors.Is(err, ErrInvalidClassSelection):
			utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_class_selection", "One or more selected classes are invalid.", map[string]string{"class_ids": "validation.invalid"})
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "teacher_update_failed", "Failed to update teacher. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "teacher classes updated", resp)
}

// DELETE /teachers/:id
func (h *Handler) Delete(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	teacherID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_teacher_id", "Please provide a valid teacher id.")
		return
	}

	if err := h.service.Delete(c.Request.Context(), ownerUserID, uint(teacherID)); err != nil {
		switch {
		case errors.Is(err, ErrTeacherNotFound):
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "teacher_not_found", "Teacher not found.")
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "teacher_delete_failed", "Failed to delete teacher. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "teacher deleted", nil)
}
