package request

type ReqLogin struct{
	Email 	string	`json:"email" validate: "require" `
	Password string `json:"password" validate: "require" `
}