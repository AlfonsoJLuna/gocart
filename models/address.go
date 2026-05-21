package models

type Address struct {
	FullName    string	`json:"full_name"`
	CompanyName	string 	`json:"company_name"`
	Address1    string 	`json:"address_1"`
	Address2    string 	`json:"address_2"`
	PostalCode  string 	`json:"postal_code"`
	City        string 	`json:"city"`
}
