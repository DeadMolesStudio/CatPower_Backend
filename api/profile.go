package api

import (
	"io/ioutil"
	"net/http"

	"github.com/asaskevich/govalidator"
	db "github.com/go-park-mail-ru/2018_2_DeadMolesStudio/database"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/logger"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/middleware"
	"golang.org/x/crypto/bcrypt"

	"CatPower/models"
	"CatPower/queries"
)

func validateUsername(dm *db.DatabaseManager, s string) ([]models.ProfileError, error) {
	var errors []models.ProfileError

	if !govalidator.StringLength(s, "4", "20") {
		errors = append(errors, models.ProfileError{
			Field: "username",
			Text:  "Username must be at least 4 characters and no more than 20 characters.",
		})
		return errors, nil
	}

	exists, err := queries.CheckExistenceOfUsername(dm, s)
	if err != nil {
		logger.Error(err)
		return errors, err
	}
	if exists {
		errors = append(errors, models.ProfileError{
			Field: "username",
			Text:  "This username is already taken.",
		})
	}

	return errors, nil
}

func validateEmail(dm *db.DatabaseManager, s string) ([]models.ProfileError, error) {
	var errors []models.ProfileError

	if !govalidator.IsEmail(s) {
		errors = append(errors, models.ProfileError{
			Field: "email",
			Text:  "Invalid email.",
		})
		return errors, nil
	}

	exists, err := queries.CheckExistenceOfEmail(dm, s)
	if err != nil {
		logger.Error(err)
		return errors, err
	}
	if exists {
		errors = append(errors, models.ProfileError{
			Field: "email",
			Text:  "This email is already taken.",
		})
	}

	return errors, nil
}

func validatePassword(s string) []models.ProfileError {
	var errors []models.ProfileError

	if !govalidator.StringLength(s, "4", "32") {
		errors = append(errors, models.ProfileError{
			Field: "password",
			Text:  "Password must be at least 4 characters and no more than 32 characters.",
		})
	}

	return errors
}

func validateFields(dm *db.DatabaseManager, u *models.UserProfile) ([]models.ProfileError, error) {
	var errors []models.ProfileError

	valErrors, dbErr := validateUsername(dm, u.Username)
	if dbErr != nil {
		return []models.ProfileError{}, dbErr
	}
	errors = append(errors, valErrors...)

	valErrors, dbErr = validateEmail(dm, u.Email)
	if dbErr != nil {
		return []models.ProfileError{}, dbErr
	}
	errors = append(errors, valErrors...)
	errors = append(errors, validatePassword(u.Password)...)

	return errors, nil
}

func hashAndSalt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (s *Server) getProfile(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value(middleware.KeyIsAuthenticated).(bool) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	profile, err := queries.GetUserProfileByID(s.dm, r.Context().Value(middleware.KeyUserID).(uint))
	if err != nil {
		switch err {
		case queries.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	json, err := profile.MarshalJSON()
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (s *Server) signUp(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	u := &models.UserProfile{}
	if err = u.UnmarshalJSON(body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if u.Username == "" || u.Email == "" || u.Password == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	fieldErrors, err := validateFields(s.dm, u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(fieldErrors) != 0 {
		sendList := models.ProfileErrorList{Errors: fieldErrors}
		json, err := sendList.MarshalJSON()
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write(json)
	} else {
		u.Password, err = hashAndSalt(u.Password)
		if err != nil {
			logger.Errorf("hash and salt password error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err := queries.CreateNewUser(s.dm, u)
		u.Password = ""
		if err != nil {
			if err == db.ErrUniqueConstraintViolation {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = loginUser(s.sm, w, u.UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Infof("New user %s with id %v, email %v logged in", u.Username, u.UserID, u.Email)

		json, err := u.MarshalJSON()
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}
