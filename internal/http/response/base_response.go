package response

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/ffajarpratama/pos-wash-api/pkg/constant"
	"github.com/ffajarpratama/pos-wash-api/pkg/custom_error"
	"github.com/ffajarpratama/pos-wash-api/pkg/custom_validator"
)

const (
	CONTENT_TYPE_HEADER  = "Content-Type"
	CONTENT_DESC_HEADER  = "Content-Description"
	CONTENT_DISPO_HEADER = "Content-Disposition"

	CONTENT_TYPE_JSON         = "application/json"
	CONTENT_TYPE_PDF          = "application/pdf"
	CONTENT_TYPE_OCTET_STREAM = "aplication/octet-stream"

	CONTENT_DESC_FILE_TRANSFER = "File Transfer"
)

type JsonResponse struct {
	Success bool           `json:"success"`
	Paging  *PagingJSON    `json:"paging"`
	Data    interface{}    `json:"data"`
	Error   *ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Code    int      `json:"code"`
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

type PagingJSON struct {
	Page      int   `json:"page"`
	PerPage   int   `json:"per_page"`
	Count     int64 `json:"count"`
	PageCount int   `json:"page_count"`
	Next      bool  `json:"next"`
	Prev      bool  `json:"prev"`
}

func OK(w http.ResponseWriter, data interface{}) {
	w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JsonResponse{
		Data:    data,
		Success: true,
	})
}

func Paging(w http.ResponseWriter, list interface{}, page, perPage int, cnt int64) {
	var paging *PagingJSON

	total := calculateTotalPage(cnt, perPage)
	if page > 0 {
		paging = &PagingJSON{
			Page:      page,
			PerPage:   perPage,
			Count:     cnt,
			PageCount: total,
			Next:      hasNext(page, total),
			Prev:      hasPrev(page),
		}
	}

	w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JsonResponse{
		Success: true,
		Paging:  paging,
		Data:    list,
		Error:   nil,
	})
}

func Error(w http.ResponseWriter, err error) {
	v, isValidationErr := err.(custom_validator.ValidatorError)
	if isValidationErr {
		w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_JSON)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonResponse{
			Error: &ErrorResponse{
				Code:    v.Code,
				Status:  v.Status,
				Message: v.Message,
				Details: v.Details,
			},
		})

		return
	}

	e, isCustomErr := err.(*custom_error.CustomError)
	if !isCustomErr {
		if err != nil && !errors.Is(err, context.Canceled) {
			fmt.Println(err.Error(), "[unhandled-error]")
		}

		w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_JSON)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JsonResponse{
			Error: &ErrorResponse{
				Code:    constant.DefaultUnhandledError,
				Status:  http.StatusInternalServerError,
				Message: constant.HTTPStatusText(http.StatusInternalServerError),
			},
		})

		return
	}

	httpCode := http.StatusInternalServerError
	internalCode := constant.DefaultUnhandledError
	msg := constant.HTTPStatusText(http.StatusInternalServerError)

	if e.ErrorContext != nil && e.ErrorContext.HTTPCode > 0 {
		httpCode = e.ErrorContext.HTTPCode
		internalCode = constant.InteralResponseCodeMap[httpCode]
		msg = constant.HTTPStatusText(httpCode)

		if e.ErrorContext.Message != "" {
			msg = e.ErrorContext.Message
		}
	}

	w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_JSON)
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(JsonResponse{
		Error: &ErrorResponse{
			Code:    internalCode,
			Status:  httpCode,
			Message: msg,
		},
	})
}

func UnauthorizedError(w http.ResponseWriter) {
	w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_JSON)
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(JsonResponse{
		Data:    nil,
		Success: false,
		Error: &ErrorResponse{
			Code:    constant.DefaultUnauthorizedError,
			Status:  http.StatusUnauthorized,
			Message: constant.HTTPStatusText(http.StatusUnauthorized),
		},
	})
}

func BinaryExcel(w http.ResponseWriter, filename string, b *bytes.Buffer) {
	filename = fmt.Sprintf("%s.xlsx", filename)
	w.Header().Set(CONTENT_DESC_HEADER, CONTENT_DESC_FILE_TRANSFER)
	w.Header().Set(CONTENT_DISPO_HEADER, "attachment; filename="+filename)
	w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_OCTET_STREAM)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(b.Bytes())
	if err != nil {
		panic(err)
	}
}

func BinaryPdf(w http.ResponseWriter, filename string, b *bytes.Buffer) {
	filename = fmt.Sprintf("%s.pdf", filename)
	w.Header().Set(CONTENT_DESC_HEADER, CONTENT_DESC_FILE_TRANSFER)
	w.Header().Set(CONTENT_DISPO_HEADER, "attachment; filename="+filename)
	w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_PDF)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(b.Bytes())
	if err != nil {
		panic(err)
	}
}

func BinaryCsv(w http.ResponseWriter, filename string, b *bytes.Buffer) {
	filename = fmt.Sprintf("%s.csv", filename)
	w.Header().Set(CONTENT_DESC_HEADER, CONTENT_DESC_FILE_TRANSFER)
	w.Header().Set(CONTENT_DISPO_HEADER, "attachment; filename="+filename)
	w.Header().Set(CONTENT_TYPE_HEADER, CONTENT_TYPE_OCTET_STREAM)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(b.Bytes())
	if err != nil {
		panic(err)
	}
}

func hasNext(currentPage, totalPages int) bool {
	return currentPage < totalPages
}

func hasPrev(currentPage int) bool {
	return currentPage > 1
}

func calculateTotalPage(cnt int64, limit int) (total int) {
	return int(math.Ceil(float64(cnt) / float64(limit)))
}
