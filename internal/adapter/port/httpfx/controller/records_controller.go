package controller

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"

	"githhub.com/structx/ddns/internal/core/domain"
)

// Records
type Records struct {
	log     *zap.SugaredLogger
	service domain.DDNS
}

// NewRecords
func NewRecords(logger *zap.Logger, service domain.DDNS) *Records {
	return &Records{
		log:     logger.Sugar().Named("RecordsController"),
		service: service,
	}
}

// RegisterRoutesV1
func (rc *Records) RegisterRoutesV1(r chi.Router) {

	rr := chi.NewRouter()

	rr.Put("/", rc.Upsert)

	r.Mount("/records", rr)
}

// UpsertRecordPayload
type UpsertRecordPayload struct {
	RecordType string `json:"record_type"`
	Root       string `json:"root"`
	Content    string `json:"content"`
	TTL        int64  `json:"ttl"`
}

// UpsertRecordParams
type UpsertRecordParams struct {
	Payload *UpsertRecordPayload `json:"payload"`
}

// Render
func (up *UpsertRecordParams) Bind(r *http.Request) error {

	if up.Payload == nil {
		return errors.New("missing request payload")
	}

	return nil
}

// RecordPayload
type RecordPayload struct{}

// UpsertRecordResponse
type UpsertRecordResponse struct{}

// Upsert
func (rc *Records) Upsert(w http.ResponseWriter, r *http.Request) {

	var params UpsertRecordParams
	err := render.Bind(r, &params)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	var record domain.Record

	switch params.Payload.RecordType {
	case "A":
		record = &domain.A{}
	case "CNAME":
		record = &domain.CName{}
	default:
		render.Render(w, r, ErrInvalidRequest(errors.New("invalid record type provided")))
	}

	err = rc.service.AddOrUpdateRecord(r.Context(), record)
	if err != nil {
		rc.log.Errorf("unable to add or update record %v", err)
		render.Render(w, r, ErrInternalServerError)
		return
	}
}
