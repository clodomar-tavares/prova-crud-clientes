package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cliente struct {
	Id             primitive.ObjectID `json:"id,omitempty"`
	Nome           string             `json:"nome,omitempty" validate:"required"`
	CPF            string             `json:"cpf,omitempty" validate:"required"`
	DataNascimento string             `json:"dataNascimento,omitempty" validate:"required"`
	Endereco       string             `json:"endereco,omitempty" validate:"required"`
	Bairro         string             `json:"bairro,omitempty" validate:"required"`
	Municipio      string             `json:"municipio,omitempty" validate:"required"`
	Estado         string             `json:"estado,omitempty" validate:"required"`
}
