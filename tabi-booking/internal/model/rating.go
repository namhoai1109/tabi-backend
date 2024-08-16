package model

type Rating struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	UserID   int    `json:"user_id"`
	BranchID int    `json:"branch_id"`
	Rating   int    `json:"rating"`
	Comment  string `json:"comment"`
	Base

	Branch *Branch `gorm:"foreignKey:BranchID"`
	User   *User   `gorm:"foreignKey:UserID"`
}

type RatingResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}

func (s *Rating) ToRatingResponse() *RatingResponse {
	username := ""
	if s.User != nil && s.User.Account != nil {
		username = s.User.Account.Username
	}

	return &RatingResponse{
		ID:        s.ID,
		Username:  username,
		Rating:    s.Rating,
		Comment:   s.Comment,
		CreatedAt: s.CreatedAt.Format("2006-01-02"),
	}
}
