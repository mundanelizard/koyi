package models

type PasswordHistory struct {
	ID        *string `json:"id"`
	UserId    *string `json:"userId"`
	Password  *string `json:"password"`
	CreatedAt *string `json:"createdAt"`
}

type EmailHistory struct {
	ID        *string `json:"id"`
	UserId    *string `json:"userId"`
	Email     *string `json:"password"`
	CreatedAt *string `json:"createdAt"`
}

type PhoneNumberHistory struct {
	ID          *string `json:"id"`
	UserId      *string `json:"userId"`
	PhoneNumber *string `json:"password"`
	CreatedAt   *string `json:"createdAt"`
}
