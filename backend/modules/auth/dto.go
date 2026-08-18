package auth

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=150"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

type UserResponse struct {
	ID       uint            `json:"id"`
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Role     string          `json:"role"`
	Center   *CenterResponse `json:"center,omitempty"`
	// HasGrades is true once the user's center has a usable academic
	// structure: at least one grade, every grade has at least one major,
	// and every major has at least one subject. Only meaningful for
	// center-scoped users; always false otherwise. The frontend uses this
	// to gate center_owner access to the dashboard until onboarding is
	// complete.
	HasGrades bool `json:"has_grades"`
}

type CenterResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type UpdateProfileRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email,max=150"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6,max=72"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	Token        string       `json:"token,omitempty"`
	RefreshToken string       `json:"refreshToken,omitempty"`
}
