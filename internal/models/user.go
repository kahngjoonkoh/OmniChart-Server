package models

type SignUpRequest struct {
	Email	 string	`json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// This struct contains user information to be returned to frontend
// Currently it only contains username. In later iterations, this
// will be enriched with more information such as avatar
type UserInfo struct {
	Username 		string 	`json:"username"`
	AccessToken 	string 	`json:"accessToken"`
	RefreshToken 	string 	`json:"refreshToken"`
}