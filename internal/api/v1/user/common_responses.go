package user

import "time"

type userProfileResponse struct {
	Status   string `json:"status"`
	ShortID  string `json:"shortId,omitemtpy"`
	Owner    bool   `json:"owner"`
	CanClaim bool   `json:"canClaim"`

	// Social network links
	Linkedin    *string `json:"linkedin,omitemtpy"`
	Email       *string `json:"email,omitemtpy"`
	Whatsapp    *string `json:"whatsapp,omitemtpy"`
	Medium      *string `json:"medium,omitemtpy"`
	TwitterX    *string `json:"twitterX,omitemtpy"`
	Website     *string `json:"website,omitemtpy"`
	Description *string `json:"description,omitempty"`
	Nickname    *string `json:"nickname,omitemtpy"`

	// Non available fort the moment
	//Image string `json:"image"`

	CreatedAt time.Time  `json:"createdAt,omitemtpy"`
	UpdatedAt *time.Time `json:"updatedAt,omitemtpy"`
}
