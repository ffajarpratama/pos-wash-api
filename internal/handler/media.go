package handler

import (
	"errors"
	"net/http"

	"github.com/ffajarpratama/pos-wash-api/internal/http/request"
	"github.com/ffajarpratama/pos-wash-api/internal/http/response"
	"github.com/ffajarpratama/pos-wash-api/pkg/constant"
	"github.com/ffajarpratama/pos-wash-api/pkg/custom_error"
)

func (h *handler) CreateMedia(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, constant.FILE_UPLOAD_MAX_SIZE)
	if err := r.ParseMultipartForm(constant.FILE_UPLOAD_MAX_SIZE); err != nil {
		err = custom_error.SetCustomError(&custom_error.ErrorContext{
			HTTPCode: http.StatusBadRequest,
			Message:  "file size exceeds 2MB",
		})

		response.Error(w, err)
		return
	}

	purpose := r.FormValue("purpose")
	file, header, err := r.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			err = custom_error.SetCustomError(&custom_error.ErrorContext{
				HTTPCode: http.StatusBadRequest,
				Message:  "file cannot be empty",
			})

			response.Error(w, err)
			return
		}

		response.Error(w, err)
		return
	}

	mimetype := header.Header.Get("Content-Type")

	if !constant.MimetypeWhitelist[mimetype] {
		err = custom_error.SetCustomError(&custom_error.ErrorContext{
			HTTPCode: http.StatusBadRequest,
			Message:  "file extension not allowed",
		})

		response.Error(w, err)
		return
	}

	req := &request.CreateMedia{
		File:      file,
		Header:    header,
		Filename:  header.Filename,
		Mimetype:  mimetype,
		Purpose:   constant.UploadPurpose(purpose),
		AssetType: constant.GetAssetType(mimetype),
	}

	res, err := h.uc.CreateMedia(r.Context(), req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.OK(w, res)
}
