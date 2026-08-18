package middleware

import (
	"errors"
	"net/http"

	"educlypro/shared/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// userCenter mirrors the one column of auth.User this file needs, to avoid
// importing modules/auth — modules/auth (and every other module) already
// imports this middleware package, so the reverse import would cycle. Same
// reasoning as the tokenRecord mirror in auth.go above.
type userCenter struct {
	CenterID *uint `gorm:"column:center_id"`
}

// CenterAcademicSetupComplete reports whether a center's academic structure
// is actually usable: at least one grade, every grade has at least one
// major, and every one of those majors has at least one subject. Queries
// the grades/majors/subjects tables by name rather than importing
// modules/academic, for the same import-cycle reason as userCenter above.
//
// This is the single source of truth for "onboarding done" — both
// RequireAcademicSetup below and auth.Service (for the has_grades field on
// the logged-in user) call into it, so the two can never drift apart.
func CenterAcademicSetupComplete(db *gorm.DB, centerID uint) (bool, error) {
	var gradeCount int64
	if err := db.Table("grades").Where("center_id = ?", centerID).Count(&gradeCount).Error; err != nil {
		return false, err
	}
	if gradeCount == 0 {
		return false, nil
	}

	var gradesMissingMajors int64
	if err := db.Table("grades").
		Where("center_id = ?", centerID).
		Where("id NOT IN (SELECT DISTINCT grade_id FROM majors)").
		Count(&gradesMissingMajors).Error; err != nil {
		return false, err
	}
	if gradesMissingMajors > 0 {
		return false, nil
	}

	var majorsMissingSubjects int64
	err := db.Table("majors").
		Where("grade_id IN (SELECT id FROM grades WHERE center_id = ?)", centerID).
		Where("id NOT IN (SELECT DISTINCT major_id FROM subjects)").
		Count(&majorsMissingSubjects).Error
	if err != nil {
		return false, err
	}

	return majorsMissingSubjects == 0, nil
}

// RequireAcademicSetup blocks a center_owner's request until
// CenterAcademicSetupComplete is true for their center — the server-side
// half of the "can't use the dashboard without at least one grade" gate
// (the other half is the frontend route guard, which only stops
// navigation and can't be trusted on its own). No-op for every other role.
//
// Must run after AuthMiddleware (needs "userID" and "role" in context).
// Must NEVER wrap the /academic/* routes themselves, nor /auth/* — those
// are respectively how a center_owner satisfies this requirement, and how
// the frontend learns whether it's satisfied in the first place.
func RequireAcademicSetup(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, _ := c.Get("role")
		if role, _ := roleVal.(string); role != "center_owner" {
			c.Next()
			return
		}

		userIDVal, exists := c.Get("userID")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		userID, ok := userIDVal.(uint)
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		var user userCenter
		if err := db.Table("users").Select("center_id").Where("id = ?", userID).Take(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
				c.Abort()
				return
			}
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "onboarding_check_failed", "Failed to verify your account. Please try again.")
			c.Abort()
			return
		}
		if user.CenterID == nil {
			utils.ErrorResponseWithCode(c, http.StatusForbidden, "onboarding_required", "Finish setting up your academic structure before continuing.")
			c.Abort()
			return
		}

		complete, err := CenterAcademicSetupComplete(db, *user.CenterID)
		if err != nil {
			utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "onboarding_check_failed", "Failed to verify your account. Please try again.")
			c.Abort()
			return
		}
		if !complete {
			utils.ErrorResponseWithCode(c, http.StatusForbidden, "onboarding_required", "Finish setting up your academic structure before continuing.")
			c.Abort()
			return
		}

		c.Next()
	}
}
