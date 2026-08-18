package academic

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

// GET /academic/grades
func (h *Handler) Tree(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	resp, err := h.service.Tree(c.Request.Context(), ownerUserID)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "academic_tree_failed", "Failed to load your academic setup. Please try again.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "academic tree retrieved", resp)
}

// POST /academic/grades
func (h *Handler) AddGrade(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	var req AddGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.service.AddGrade(c.Request.Context(), ownerUserID, &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrGradeExists):
			utils.ValidationErrorResponse(c, http.StatusConflict, "grade_exists", "This grade already exists.", map[string]string{"name": "validation.taken"})
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "academic_add_grade_failed", "Failed to add grade. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "grade added", resp)
}

// DELETE /academic/grades/:id
func (h *Handler) RemoveGrade(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	gradeID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_grade_id", "Please provide a valid grade id.")
		return
	}

	if err := h.service.RemoveGrade(c.Request.Context(), ownerUserID, uint(gradeID)); err != nil {
		switch {
		case errors.Is(err, ErrGradeNotFound):
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "grade_not_found", "Grade not found.")
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "academic_remove_grade_failed", "Failed to remove grade. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "grade removed", nil)
}

// POST /academic/grades/:id/majors
func (h *Handler) AddMajor(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	gradeID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_grade_id", "Please provide a valid grade id.")
		return
	}

	var req AddMajorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.service.AddMajor(c.Request.Context(), ownerUserID, uint(gradeID), &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrGradeNotFound):
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "grade_not_found", "Grade not found.")
		case errors.Is(err, ErrMajorExists):
			utils.ValidationErrorResponse(c, http.StatusConflict, "major_exists", "This major already exists.", map[string]string{"name": "validation.taken"})
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "academic_add_major_failed", "Failed to add major. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "major added", resp)
}

// DELETE /academic/majors/:id
func (h *Handler) RemoveMajor(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	majorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_major_id", "Please provide a valid major id.")
		return
	}

	if err := h.service.RemoveMajor(c.Request.Context(), ownerUserID, uint(majorID)); err != nil {
		switch {
		case errors.Is(err, ErrMajorNotFound):
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "major_not_found", "Major not found.")
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "academic_remove_major_failed", "Failed to remove major. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "major removed", nil)
}

// POST /academic/majors/:id/subjects
func (h *Handler) AddSubject(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	majorID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_major_id", "Please provide a valid major id.")
		return
	}

	var req AddSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "validation_failed", "Please correct the highlighted fields.", utils.FormatValidationErrors(err))
		return
	}

	resp, err := h.service.AddSubject(c.Request.Context(), ownerUserID, uint(majorID), &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrMajorNotFound):
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "major_not_found", "Major not found.")
		case errors.Is(err, ErrSubjectExists):
			utils.ValidationErrorResponse(c, http.StatusConflict, "subject_exists", "This subject already exists.", map[string]string{"name": "validation.taken"})
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "academic_add_subject_failed", "Failed to add subject. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "subject added", resp)
}

// DELETE /academic/subjects/:id
func (h *Handler) RemoveSubject(c *gin.Context) {
	ownerUserID := c.MustGet("userID").(uint)

	subjectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_subject_id", "Please provide a valid subject id.")
		return
	}

	if err := h.service.RemoveSubject(c.Request.Context(), ownerUserID, uint(subjectID)); err != nil {
		switch {
		case errors.Is(err, ErrSubjectNotFound):
			utils.ErrorResponseWithCode(c, http.StatusNotFound, "subject_not_found", "Subject not found.")
		default:
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "academic_remove_subject_failed", "Failed to remove subject. Please try again.")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "subject removed", nil)
}
