package user

import "time"

type userProfileResponse struct {
	ShortID string `json:"shortId"`
	// Social network links
	Linkedin    *string `json:"linkedin"`
	Email       *string `json:"email"`
	Whatsapp    *string `json:"whatsapp"`
	Medium      *string `json:"medium"`
	TwitterX    *string `json:"twitterX"`
	Website     *string `json:"website"`
	Description *string `json:"description"`
	Nickname    *string `json:"nickname"`

	// Non available fort the moment
	Image *string `json:"image"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
