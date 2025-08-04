package dto

type UpdateUserRequest struct {
	ID    string `json:"id" binding:"required"`
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

type DeleteUserRequest struct {
	ID string `json:"id" binding:"required"`
}
