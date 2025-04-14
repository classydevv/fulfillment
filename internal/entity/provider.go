package entity

type Provider struct {
	Id ProviderId `json:"id" example:"kuper"`
	Name string `json:"name" example:"Купер"`
}

type ProviderId string