package notifications

import (
	"educlypro/shared/utils"
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

// GET /notifications
func (h *Handler) List(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_notification_page", "Please provide a valid page number.", map[string]string{"page": "validation.invalid"})
		return
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_notification_per_page", "Please provide a valid page size.", map[string]string{"per_page": "validation.invalid"})
		return
	}

	if _, ok := AllowedPerPage[perPage]; !ok {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_notification_per_page", "Please provide a valid page size.", map[string]string{"per_page": "validation.invalid"})
		return
	}

	filter := c.DefaultQuery("filter", FilterUnread)
	if filter != FilterRead && filter != FilterUnread {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_notification_filter", "Please provide a valid notification filter.", map[string]string{"filter": "validation.invalid"})
		return
	}

	results, err := h.service.List(c.Request.Context(), userID, ListParams{
		Page:    page,
		PerPage: perPage,
		Filter:  filter,
	})
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "notifications_list_failed", "Failed to load notifications. Please try again.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "notifications retrieved", results)
}

// GET /notifications/unread-count
func (h *Handler) UnreadCount(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	count, err := h.service.UnreadCount(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "notifications_unread_count_failed", "Failed to load unread notifications count.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "unread count retrieved", gin.H{"count": count})
}

// PATCH /notifications/:id/read
func (h *Handler) MarkRead(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_notification_id", "Please provide a valid notification id.")
		return
	}

	if err := h.service.MarkRead(c.Request.Context(), userID, uint(id)); err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "notification_mark_read_failed", "Failed to mark the notification as read.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "notification marked as read", nil)
}

// PATCH /notifications/read-all
func (h *Handler) MarkAllRead(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	if err := h.service.MarkAllRead(c.Request.Context(), userID); err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "notifications_mark_all_read_failed", "Failed to mark all notifications as read.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "all notifications marked as read", nil)
}

// DELETE /notifications/:id
func (h *Handler) Delete(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusBadRequest, "invalid_notification_id", "Please provide a valid notification id.")
		return
	}

	if err := h.service.Delete(c.Request.Context(), userID, uint(id)); err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "notification_delete_failed", "Failed to delete the notification.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "notification deleted", nil)
}

// DELETE /notifications/read
func (h *Handler) DeleteAllRead(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	if err := h.service.DeleteAllRead(c.Request.Context(), userID); err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "notifications_delete_all_read_failed", "Failed to delete read notifications.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "all read notifications deleted", nil)
}
