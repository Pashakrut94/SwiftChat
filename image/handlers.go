package image

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/Pashakrut94/SwiftChat/auth"
	"github.com/Pashakrut94/SwiftChat/handlers"
	"github.com/Pashakrut94/SwiftChat/users"
	"github.com/pkg/errors"
)

func Upload(repo users.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authSession := auth.SessionValue(ctx)
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			handlers.HandleResponseError(w, errors.Wrap(err, "error parsing request body").Error(), http.StatusBadRequest)
			return
		}
		file, header, err := r.FormFile("uploadfile")
		if err != nil {
			handlers.HandleResponseError(w, errors.Wrap(err, "error returning file by provided key").Error(), http.StatusInternalServerError)
			return
		}
		url, err := HandleUpload(repo, file, header, authSession.UserID)
		if err != nil {
			handlers.HandleResponseError(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "File uploaded to %s\n", url)
	}
}

func UploadTemplate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, err := template.ParseFiles("./forms/update.gtpl")
		if err != nil {
			handlers.HandleResponseError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, token)
	}
}
