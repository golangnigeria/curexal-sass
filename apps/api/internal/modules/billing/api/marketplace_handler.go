package api

import (
	"net/http"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/labstack/echo/v4"
)

type MarketplaceHandler struct {
	server *server.Server
}

func NewMarketplaceHandler(s *server.Server) *MarketplaceHandler {
	return &MarketplaceHandler{server: s}
}

type ModuleDTO struct {
	ID           string  `json:"id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	Description  string  `json:"description,omitempty"`
	PriceMonthly float64 `json:"priceMonthly"`
	IsActive     bool    `json:"isActive"`
}

func (h *MarketplaceHandler) ListModules(c echo.Context) error {
	modules := []ModuleDTO{
		{
			ID:           "mod_clinic_emr",
			Code:         "clinic_emr",
			Name:         "Clinic HMS / EMR",
			Category:     "clinical",
			Description:  "Electronic Medical Records, Encounters, Consultations & Prescriptions",
			PriceMonthly: 25000.00,
			IsActive:     true,
		},
		{
			ID:           "mod_laboratory_lims",
			Code:         "laboratory_lims",
			Name:         "Laboratory LIS / LIMS",
			Category:     "diagnostics",
			Description:  "Phlebotomy Accessioning, Worksheet Entry & Two-Stage Pathologist Verification",
			PriceMonthly: 35000.00,
			IsActive:     true,
		},
		{
			ID:           "mod_referral_marketplace",
			Code:         "referral_marketplace",
			Name:         "B2B Referral Marketplace",
			Category:     "marketplace",
			Description:  "B2B Laboratory Network Referrals & Settlement Ledger",
			PriceMonthly: 15000.00,
			IsActive:     true,
		},
	}
	return c.JSON(http.StatusOK, modules)
}

type SubscribeRequest struct {
	ModuleCode string `json:"moduleCode"`
}

func (h *MarketplaceHandler) Subscribe(c echo.Context) error {
	var req SubscribeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid request body"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"subscriptionId": "sub_01HGB72A9X",
		"licenseKey":     "LIC-CUREXAL-2026-ACTIVE-88F",
		"status":         "active",
	})
}
