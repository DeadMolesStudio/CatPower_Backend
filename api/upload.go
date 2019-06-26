package api

import (
	"net/http"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/logger"
	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/middleware"

	"CatPower/filesystem"
)

func uploadFile(w http.ResponseWriter, r *http.Request, key, dir string) string {
	if err := r.ParseMultipartForm(5 * (1 << 20)); err != nil { // 5 MB
		if err == http.ErrNotMultipart || err == http.ErrMissingBoundary {
			w.WriteHeader(http.StatusBadRequest)
			return ""
		}
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return ""
	}
	file, fileHeader, err := r.FormFile(key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return ""
	}
	defer file.Close()

	uID := r.Context().Value(middleware.KeyUserID).(uint)
	filename := fileHeader.Filename
	filename, err = filesystem.GetHashedNameForFile(uID, filename)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return ""
	}
	err = filesystem.SaveFile(file, dir, filename)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return ""
	}

	return dir + filename
}
