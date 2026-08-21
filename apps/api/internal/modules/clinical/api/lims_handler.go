package api

import (
	"net/http"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type LimsHandler struct {
	server *server.Server
}

func NewLimsHandler(s *server.Server) *LimsHandler {
	return &LimsHandler{server: s}
}

type CreateOrderRequest struct {
	PatientID      string   `json:"patientId"`
	VisitID        string   `json:"visitId,omitempty"`
	ProviderID     string   `json:"providerId,omitempty"`
	TestCatalogIDs []string `json:"testCatalogIds"`
	Priority       string   `json:"priority"`
}

func (h *LimsHandler) CreateOrder(c echo.Context) error {
	var req CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid request body"})
	}

	orderID := uuid.New().String()
	orderNumber := "ORD-" + time.Now().Format("20060102") + "-" + orderID[:6]

	return c.JSON(http.StatusCreated, echo.Map{
		"orderId":     orderID,
		"orderNumber": orderNumber,
		"status":      "pending",
	})
}

type AccessionRequest struct {
	OrderID          string `json:"orderId"`
	SpecimenTypeCode string `json:"specimenTypeCode"`
	Barcode          string `json:"barcode"`
}

func (h *LimsHandler) AccessionSpecimen(c echo.Context) error {
	var req AccessionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid request body"})
	}

	specimenID := uuid.New().String()
	return c.JSON(http.StatusOK, echo.Map{
		"specimenId": specimenID,
		"barcode":    req.Barcode,
		"status":     "accessioned",
	})
}

type ResultItemDTO struct {
	TestCatalogID  string `json:"testCatalogId"`
	ParameterName  string `json:"parameterName"`
	Value          string `json:"value"`
	Unit           string `json:"unit,omitempty"`
	ReferenceRange string `json:"referenceRange,omitempty"`
	Flag           string `json:"flag,omitempty"`
}

type EnterResultsRequest struct {
	OrderID string          `json:"orderId"`
	Results []ResultItemDTO `json:"results"`
}

func (h *LimsHandler) EnterResults(c echo.Context) error {
	var req EnterResultsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid request body"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"orderId":       req.OrderID,
		"recordedCount": len(req.Results),
	})
}

type AuthorizeOrderRequest struct {
	OrderID       string `json:"orderId"`
	SignatureHash string `json:"signatureHash"`
	Notes         string `json:"notes,omitempty"`
}

func (h *LimsHandler) AuthorizeOrder(c echo.Context) error {
	var req AuthorizeOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid request body"})
	}

	authID := uuid.New().String()
	return c.JSON(http.StatusOK, echo.Map{
		"authorizationId": authID,
		"orderId":         req.OrderID,
		"status":          "authorized",
	})
}
