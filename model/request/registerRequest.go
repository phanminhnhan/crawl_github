package request

type ReqRes struct{
	Fullname string `validate: "require" json:"fullname"`
	Email 	string	`validate: "require" json:"email"`
	Password string `validate: "require" json:"password"`
}