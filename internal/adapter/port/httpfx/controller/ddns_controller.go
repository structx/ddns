package controller

import (
	"go.uber.org/zap"

	"githhub.com/structx/ddns/internal/core/domain"
)

// DDNS
type DDNS struct {
	log     *zap.SugaredLogger
	service domain.DDNS
}

// NewDDNS
func NewDDNS(logger *zap.Logger, ddns domain.DDNS) *DDNS {
	return &DDNS{
		log:     logger.Sugar().Named("DdnsController"),
		service: ddns,
	}
}
