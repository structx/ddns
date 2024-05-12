package controller

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/structx/ddns/internal/core/domain"
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

// RecordPayload
type RecordPayload struct {
	RecordType string `json:"record_type"`
	Root       string `json:"root"`
	Content    string `json:"content"`
	TTL        int64  `json:"ttl"`
}

// UpsertRecordParams
type UpsertRecordParams struct {
	Payload *RecordPayload `json:"payload"`
}

// Render
func (up *UpsertRecordParams) Bind(r *http.Request) error {

	if up.Payload == nil {
		return errors.New("missing request payload")
	}

	return nil
}

// UpsertRecordResponse
type UpsertRecordResponse struct {
	Payload *RecordPayload `json:"payload"`
}

// Render
func (upr *UpsertRecordResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusCreated)
	return nil
}

// NewUpsertRecordResponse
func NewUpsertRecordResponse(payload *RecordPayload) *UpsertRecordResponse {
	return &UpsertRecordResponse{Payload: payload}
}

// Upsert
func (rc *Records) Upsert(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var params UpsertRecordParams
	err := render.Bind(r, &params)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	var record domain.Record

	switch params.Payload.RecordType {
	case "A":
		record = domain.NewARecord(params.Payload.Root, params.Payload.Content, params.Payload.TTL)
	case "CNAME":
		record = domain.NewCNameRecord(params.Payload.Root, params.Payload.Content, params.Payload.TTL)
	default:
		_ = render.Render(w, r, ErrInvalidRequest(errors.New("invalid record type provided")))
		return
	}

	err = rc.service.AddOrUpdateRecord(ctx, record)
	if err != nil {
		rc.log.Errorf("unable to add or update record %v", err)
		_ = render.Render(w, r, ErrInternalServerError)
		return
	}

	_ = render.Render(w, r, NewUpsertRecordResponse(params.Payload))
}
