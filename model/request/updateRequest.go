package request

type ReqUpdate struct{
	Email 	string	`json:"email" validate: "require" `
	Fullname string `json:"fullname" validate: "require" `
}
