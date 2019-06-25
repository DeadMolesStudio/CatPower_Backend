package api

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/session"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/logger"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/middleware"

	"CatPower/models"
	"CatPower/queries"
)

func comparePasswords(hashed string, clean string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(clean))
	switch err {
	case nil:
		return true, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return false, nil
	default:
		return false, err
	}
}

func loginUser(s *session.SessionManager, w http.ResponseWriter, userID uint) error {
	sessionID, err := s.Create(userID)
	if err != nil {
		logger.Error(err)
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(30 * 24 * time.Hour),
		// Secure:   true,
		// HttpOnly: true,
	})

	return nil
}

func (s *Server) logIn(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(middleware.KeyIsAuthenticated).(bool) {
		// user has already logged in
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	u := &models.UserProfile{}
	if err = u.UnmarshalJSON(body); err != nil || !govalidator.IsEmail(u.Email) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbResponse, err := queries.GetUserPassword(s.dm, u.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	passwordsMatch, err := comparePasswords(dbResponse.Password, u.Password)
	if err != nil {
		logger.Errorf("compare passwords error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if passwordsMatch {
		if err = loginUser(s.sm, w, dbResponse.UserID); err != nil {
			logger.Errorf("cannot login user %d: %s", dbResponse.UserID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.Info("user %s with id %d email %s logged in", dbResponse.Username, dbResponse.UserID, dbResponse.Email)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
}

func (s *Server) logOut(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value(middleware.KeyIsAuthenticated).(bool) {
		// user has already logged out
		return
	}
	if err := s.sm.Delete(r.Context().Value(middleware.KeySessionID).(string)); err != nil { // but we continue
		logger.Error(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		// Secure:   true,
		// HttpOnly: true,
	})
}
