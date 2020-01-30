package handler

import (
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/neilli-sable/sideupload/presenter/httpserve"
)

// Handler 共通Handler
type Handler struct{}

// OK 正常に結果を返す
func (*Handler) OK(w http.ResponseWriter, body interface{}) {
	type Response struct {
		Body   interface{}           `json:"body"`
		Errors []string `json:"errors"`
	}
	res := &Response{
		Body:   body,
		Errors: nil,
	}
	httpserve.JSON(w, res, http.StatusOK)
}

// Error エラーとして結果を返す
func (*Handler) Error(w http.ResponseWriter, err error) {
	httpserve.JSON(w, nil, 500)
}

// Binary バイナリを返す
func (*Handler) Binary(w http.ResponseWriter, contentType, fileName string, body io.Reader) {
	httpserve.Binary(w, contentType, contentDisposition(fileName), body)
}

// ContentDisposition はファイル名をUTF-8でエンコードして、ContentDisposition文字列にセットして返します
func contentDisposition(fileName string) string {
	encoded := url.QueryEscape(fileName)
	encoded = regexp.MustCompile(`([^%])(\+)`).ReplaceAllString(encoded, "$1%20")
	return `attachment; filename*=UTF-8''` + encoded
}
