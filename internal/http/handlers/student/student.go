package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/titanSarim/one-rest-api-go/internal/storage"
	"github.com/titanSarim/one-rest-api-go/internal/types"
	"github.com/titanSarim/one-rest-api-go/internal/utils/response"
)

func Create(storage storage.Storage) http.HandlerFunc{
	return func(res http.ResponseWriter, req *http.Request) {

		var student types.Student

		slog.Info("Creating a student...")

		err := json.NewDecoder(req.Body).Decode(&student)

		if errors.Is(err, io.EOF){
			response.WriteJson(res, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return 
		}

		if(err != nil){
			response.WriteJson(res, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(res, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}


		id, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		slog.Info("User created a successfully", slog.String("UserId", fmt.Sprint(id)))

		if(err != nil){
			response.WriteJson(res, http.StatusInternalServerError, response.GeneralError(err))
            return     
		}

		response.WriteJson(res, http.StatusCreated, map[string] int64{"id": id})

	}
}


