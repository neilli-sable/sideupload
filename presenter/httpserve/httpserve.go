package httpserve

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// JSON API実行結果を返すヘルパー的関数
func JSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}

// CSVUTF8 CSV出力(UTF-8)
func CSVUTF8(w http.ResponseWriter, fileName string, header []string, data [][]string) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", ContentDisposition(fileName))
	writer := csv.NewWriter(w)
	writer.Write(header)
	writer.WriteAll(data)
}

// CSVShiftJIS CSV出力(Shift-JIS)
func CSVShiftJIS(w http.ResponseWriter, fileName string, header []string, data [][]string) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", ContentDisposition(fileName))

	writer := csv.NewWriter(transform.NewWriter(w, japanese.ShiftJIS.NewEncoder()))
	writer.Write(header)
	writer.WriteAll(data)
}

// Binary PDFファイルなどの出力
func Binary(w http.ResponseWriter, contentType, contentDisposition string, body io.Reader) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", contentDisposition)
	io.Copy(w, body)
}

// Redirect Redirectする(307なのでキャッシュしない)
func Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, 307)
}

// MultipartInput MultipartFormの入力からデータを取り出す
func MultipartInput(r *http.Request, fieldName string) (multipart.File, *multipart.FileHeader, error) {
	file, fileHeader, err := r.FormFile(fieldName)
	if err != nil {
		return nil, nil, err
	}

	return file, fileHeader, nil
}

// GetSchema Schema名の取得
func GetSchema(r *http.Request) string {
	const (
		https = "https"
		http  = "http"
	)

	if r.Header.Get("Cloudfront-Forwarded-Proto") == https {
		// CloudFrontがSSLアクセラレータとして動いているとき、r.TLSはnilだがhttpsでリダイレクト
		return https
	} else if r.TLS == nil {
		// 通常時の非TLS
		return http
	}
	// デフォルトは本番で使うhttpsに
	return https
}

// GetHost Host名の取得
func GetHost(r *http.Request) string {
	if r.URL.IsAbs() {
		return r.URL.Host
	}

	return r.Host
}

// GetOrigin URLのベース部分(プロトコル+Host名)の取得
func GetOrigin(r *http.Request) string {
	host := GetHost(r)
	if r.URL.IsAbs() {
		return r.URL.Scheme + "://" + host
	}

	return GetSchema(r) + "://" + host
}

// ContentDisposition はファイル名をUTF-8でエンコードして、ContentDisposition文字列にセットして返します
func ContentDisposition(fileName string) string {
	encoded := url.QueryEscape(fileName)
	encoded = regexp.MustCompile(`([^%])(\+)`).ReplaceAllString(encoded, "$1%20")
	return `attachment; filename*=UTF-8''` + encoded
}
