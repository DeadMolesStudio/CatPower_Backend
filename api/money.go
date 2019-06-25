package api

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"

	db "github.com/go-park-mail-ru/2018_2_DeadMolesStudio/database"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/logger"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/middleware"

	"CatPower/models"
	"CatPower/queries"
)

func (s *Server) GetCategories(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value(middleware.KeyIsAuthenticated).(bool) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	res, err := queries.GetAllCategories(s.dm, r.Context().Value(middleware.KeyUserID).(uint))
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json, err := models.MoneyCategoryList(res).MarshalJSON()
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (s *Server) AddCategory(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value(middleware.KeyIsAuthenticated).(bool) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c := &models.MoneyCategory{}
	if err = c.UnmarshalJSON(body); err != nil || c.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c.User = r.Context().Value(middleware.KeyUserID).(uint)

	if err := queries.AddCategory(s.dm, c); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json, err := c.MarshalJSON()
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (s *Server) AddCategoryPic(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value(middleware.KeyIsAuthenticated).(bool) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rawFor := r.URL.Query().Get("for")
	if rawFor == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	forCat, err := strconv.ParseUint(rawFor, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	path := uploadFile(w, r, "picture", "media/categories/")
	if path == "" {
		return
	}

	if err := queries.AddCategoryPic(s.dm, r.Context().Value(middleware.KeyUserID).(uint), forCat, path); err != nil {
		if err == queries.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	json, err := models.MoneyCategory{ID: uint(forCat), Pic: &path}.MarshalJSON()
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func (s *Server) AddMoney(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value(middleware.KeyIsAuthenticated).(bool) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !r.Context().Value(middleware.KeyIsAuthenticated).(bool) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	op := &models.MoneyOp{}
	if err = op.UnmarshalJSON(body); err != nil || op.Delta == 0 || op.From == op.To {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	op.User = r.Context().Value(middleware.KeyUserID).(uint)
	op.UUID = uuid.NewV4().String()

	for {
		if err := queries.AddMoney(s.dm, op); err != nil {
			switch err {
			case db.ErrUniqueConstraintViolation:
				logger.Info("collision for uuid of money op")
				continue
			case queries.ErrForeignKeyViolation:
				w.WriteHeader(http.StatusNotFound)
				return
			default:
				logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		break
	}
	json, err := op.MarshalJSON()
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (s *Server) AddMoneyPhoto(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value(middleware.KeyIsAuthenticated).(bool) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	forOp := r.URL.Query().Get("for")
	if forOp == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	path := uploadFile(w, r, "photo", "media/ops/")
	if path == "" {
		return
	}

	if err := queries.AddMoneyPic(s.dm, r.Context().Value(middleware.KeyUserID).(uint), forOp, path); err != nil {
		if err == queries.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	json, err := models.MoneyOp{UUID: forOp, Photo: &path}.MarshalJSON()
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(json)
}

func (s *Server) GetHistory(w http.ResponseWriter, r *http.Request) {
	if !r.Context().Value(middleware.KeyIsAuthenticated).(bool) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	since := &time.Time{}
	if raw := r.URL.Query().Get("since"); raw != "" {
		var err error
		if *since, err = time.Parse(time.RFC3339, raw); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	// todo: limit value
	res, err := queries.GetMoneyHistory(s.dm, r.Context().Value(middleware.KeyUserID).(uint), since)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json, err := models.MoneyOpHistory(res).MarshalJSON()
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
