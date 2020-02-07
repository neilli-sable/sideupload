package handler

import (
	"net/http"

	"github.com/neilli-sable/sideupload/infrastructure/adaptor"
	"github.com/neilli-sable/sideupload/infrastructure/setting"
)

// SideUploadHandler ...
type SideUploadHandler struct {
	Handler
	Setting *setting.Setting
}

// Backup ...
func (h *SideUploadHandler) Backup(w http.ResponseWriter, r *http.Request) {
	usecase := adaptor.UsecaseFactory(h.Setting)
	defer usecase.Close()

	targets, err := usecase.GetTargets("")
	if err != nil {
		h.Error(w, err)
		return
	}

	archives, err := usecase.CompressTargets(targets)
	if err != nil {
		h.Error(w, err)
		return
	}

	err = usecase.UploadArchives(archives)
	if err != nil {
		h.Error(w, err)
		return
	}

	h.OK(w, nil)
}

// Clean ...
func (h *SideUploadHandler) Clean(w http.ResponseWriter, r *http.Request) {
	usecase := adaptor.UsecaseFactory(h.Setting)
	defer usecase.Close()

	_, err := usecase.DeleteOldArchives()
	if err != nil {
		h.Error(w, err)
		return
	}

	h.OK(w, nil)
}
