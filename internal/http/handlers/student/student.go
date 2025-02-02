package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

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


func GetById(storage storage.Storage) http.HandlerFunc{
	return func(res http.ResponseWriter, req *http.Request) {

		id := req.PathValue("id")

		slog.Info("gettitng a student", slog.String("id", id))

		parsedId, err := strconv.ParseInt(id, 10, 64)

		if err!= nil || parsedId <= 0 {
            response.WriteJson(res, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id")))
            return
        }

		student, err := storage.GetStudentById(parsedId)

		if err != nil{		
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(res, http.StatusInternalServerError, response.GeneralError(err))
            return
		}

		response.WriteJson(res, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc{
	return func(res http.ResponseWriter, req *http.Request) {

		slog.Info("getting all students")

		students, err := storage.GetStudents()

		if(err != nil){
			slog.Error("error getting users")
            response.WriteJson(res, http.StatusInternalServerError, response.GeneralError(err))
            return
		}

		response.WriteJson(res, http.StatusOK, students)
	}
}