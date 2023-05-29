package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/clodomar/prova/configs"
	"github.com/clodomar/prova/models"
	"github.com/clodomar/prova/responses"
	"github.com/clodomar/prova/services"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var clienteCollection *mongo.Collection = configs.GetCollection(configs.DB, "clientes")
var validate = validator.New()

func CreateCliente() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var cliente models.Cliente
		defer cancel()

		//validate the request body
		if err := json.NewDecoder(r.Body).Decode(&cliente); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.ClienteResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&cliente); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.ClienteResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		NewClient := models.Cliente{
			Id:             primitive.NewObjectID(),
			Nome:           cliente.Nome,
			CPF:            cliente.CPF,
			DataNascimento: cliente.DataNascimento,
			Endereco:       cliente.Endereco,
			Bairro:         cliente.Bairro,
			Municipio:      cliente.Municipio,
			Estado:         cliente.Estado,
		}
		//result, err := clienteCollection.InsertOne(ctx, newUser)
		// if err != nil {
		// 	rw.WriteHeader(http.StatusInternalServerError)
		// 	response := responses.ClienteResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
		// 	json.NewEncoder(rw).Encode(response)
		// 	return
		// }

		message, err := json.Marshal(NewClient)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.ClienteResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return

		}
		mq := services.NewRabbitMQ()
		if err := mq.SendMessage(message); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.ClienteResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return

		}

		rw.WriteHeader(http.StatusCreated)
		response := responses.ClienteResponse{Status: http.StatusCreated, Message: "success"}
		json.NewEncoder(rw).Encode(response)
	}
}

func GetCliente() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["clienteId"]
		var user models.Cliente
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		err := clienteCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.ClienteResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.ClienteResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}}
		json.NewEncoder(rw).Encode(response)
	}
}

func EditACliente() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["userId"]
		var cliente models.Cliente
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		//validate the request body
		if err := json.NewDecoder(r.Body).Decode(&cliente); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.ClienteResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&cliente); validationErr != nil {
			rw.WriteHeader(http.StatusBadRequest)
			response := responses.ClienteResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		update := bson.M{"nome": cliente.Nome, "cpf": cliente.CPF, "datanascimento": cliente.DataNascimento, "endereco": cliente.Endereco, "bairro": cliente.Bairro, "municipio": cliente.Municipio, "estado": cliente.Estado}

		result, err := clienteCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.ClienteResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//get updated user details
		var updatedUser models.Cliente
		if result.MatchedCount == 1 {
			err := clienteCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.ClienteResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
				return
			}
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.ClienteResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}}
		json.NewEncoder(rw).Encode(response)
	}
}

func DeleteAUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["clienteId"]
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		result, err := clienteCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.ClienteResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		if result.DeletedCount < 1 {
			rw.WriteHeader(http.StatusNotFound)
			response := responses.ClienteResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		rw.WriteHeader(http.StatusOK)
		response := responses.ClienteResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}}
		json.NewEncoder(rw).Encode(response)
	}
}

func GetAllCliente() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users []models.Cliente
		defer cancel()

		results, err := clienteCollection.Find(ctx, bson.M{})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.ClienteResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(rw).Encode(response)
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.Cliente
			if err = results.Decode(&singleUser); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.ClienteResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(rw).Encode(response)
			}

			users = append(users, singleUser)
		}
	}
}
