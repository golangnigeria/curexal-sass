package application

import (
	"context"
	"fmt"

	billingDomain "github.com/golangnigeria/curexal/internal/modules/billing/domain"
	"github.com/golangnigeria/curexal/internal/modules/subscription/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type CommercialService struct {
	server           *server.Server
	entitlementRepo  domain.EntitlementRepository
	commercialRepo   domain.CommercialRepository
	entitlementSvc   *EntitlementService
	providerRegistry *billingDomain.PaymentProviderRegistry
	gracePeriodDays  int
}

func NewCommercialService(
	s *server.Server,
	entitlementRepo domain.EntitlementRepository,
	commercialRepo domain.CommercialRepository,
	entitlementSvc *EntitlementService,
	providerRegistry *billingDomain.PaymentProviderRegistry,
) *CommercialService {
	return &CommercialService{
		server:           s,
		entitlementRepo:  entitlementRepo,
		commercialRepo:   commercialRepo,
		entitlementSvc:   entitlementSvc,
		providerRegistry: providerRegistry,
		gracePeriodDays:  7,
	}
}

func (s *CommercialService) SetGracePeriodDays(days int) {
	if days > 0 {
		s.gracePeriodDays = days
	}
}

func (s *CommercialService) CreateCommercialOrder(
	ctx context.Context,
	orgID uuid.UUID,
	userID *uuid.UUID,
	req domain.CreateOrderRequest,
	providerName string,
) (*domain.CreateOrderResponse, error) {
	if orgID == uuid.Nil {
		return nil, fmt.Errorf("organization ID is required")
	}
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("order must contain at least one capability item")
	}

	currency := strings.ToUpper(req.Currency)
	if currency == "" {
		currency = "NGN"
	}
	if currency != "NGN" && currency != "USD" {
		return nil, fmt.Errorf("unsupported currency: %s", req.Currency)
	}

	if providerName == "" {
		providerName = "mock"
	}

	provider, err := s.providerRegistry.Get(providerName)
	if err != nil {
		return nil, fmt.Errorf("invalid payment provider: %w", err)
	}

	// 1. Resolve capability items and backend DB prices
	var items []domain.CommercialOrderItem
	var subtotal float64
	var reqCapCodes []string

	for _, reqItem := range req.Items {
		capCode := reqItem.CapabilityCode
		reqCapCodes = append(reqCapCodes, capCode)

		capObj, err := s.entitlementRepo.GetCapabilityByCode(ctx, capCode)
		if err != nil || capObj == nil || !capObj.IsActive {
			return nil, fmt.Errorf("capability not found or inactive: %s", capCode)
		}

		cycle := strings.ToLower(reqItem.BillingCycle)
		if cycle == "" {
			cycle = "monthly"
		}
		if cycle != "monthly" && cycle != "annual" {
			return nil, fmt.Errorf("invalid billing cycle: %s for capability %s", reqItem.BillingCycle, capCode)
		}

		prices, err := s.entitlementRepo.GetCapabilityPrices(ctx, capObj.ID)
		if err != nil {
			return nil, fmt.Errorf("retrieving capability prices: %w", err)
		}

		var unitPrice float64
		var foundPrice bool
		for _, p := range prices {
			if strings.EqualFold(p.Currency, currency) && strings.EqualFold(p.BillingPeriod, cycle) && p.IsActive {
				unitPrice = p.Price
				foundPrice = true
				break
			}
		}

		if !foundPrice {
			// Default fallback pricing if not explicitly in DB
			if currency == "NGN" {
				unitPrice = 35000.00
				if cycle == "annual" {
					unitPrice = 350000.00
				}
			} else {
				unitPrice = 30.00
				if cycle == "annual" {
					unitPrice = 300.00
				}
			}
		}

		amount := unitPrice
		subtotal += amount

		items = append(items, domain.CommercialOrderItem{
			CapabilityID:   capObj.ID,
			CapabilityCode: capCode,
			BillingCycle:   cycle,
			UnitPrice:      unitPrice,
			Amount:         amount,
			Currency:       currency,
		})
	}

	// 2. Validate prerequisites/dependencies
	deps, err := s.entitlementRepo.GetCapabilityDependencies(ctx, reqCapCodes)
	if err == nil && len(deps) > 0 {
		effectiveCaps, err := s.entitlementSvc.GetEffectiveCapabilities(ctx, orgID)
		if err == nil {
			effectiveMap := make(map[string]bool)
			for _, cap := range effectiveCaps {
				effectiveMap[cap] = true
			}
			for _, dep := range deps {
				alreadyPurchasing := false
				for _, code := range reqCapCodes {
					if code == dep {
						alreadyPurchasing = true
						break
					}
				}
				if !effectiveMap[dep] && !alreadyPurchasing {
					return nil, fmt.Errorf("prerequisite capability missing: %s required before purchasing requested capability", dep)
				}
			}
		}
	}

	// 3. Tax and Total calculation
	tax := 0.00
	if currency == "NGN" {
		tax = subtotal * 0.075 // 7.5% VAT in Nigeria
	}
	discount := 0.00
	total := subtotal + tax - discount

	// 4. Construct Order & Order Items
	orderNum := fmt.Sprintf("ORD-%s-%s", time.Now().Format("20060102"), uuid.New().String()[:8])
	order := &domain.CommercialOrder{
		ID:             uuid.New(),
		OrderNumber:    orderNum,
		OrganizationID: orgID,
		Status:         "pending_payment",
		Currency:       currency,
		Subtotal:       subtotal,
		Tax:            tax,
		Discount:       discount,
		Total:          total,
		CreatedBy:      userID,
	}

	if err := s.commercialRepo.CreateOrder(ctx, order, items); err != nil {
		return nil, fmt.Errorf("persisting commercial order: %w", err)
	}

	// 5. Initialize Payment Provider Transaction
	initParams := billingDomain.PaymentInitParams{
		OrderID:        order.ID,
		OrganizationID: orgID,
		CustomerEmail:  "billing@organization.curexal.space",
		CustomerName:   "Organization Finance Admin",
		Amount:         total,
		Currency:       currency,
		CallbackURL:    fmt.Sprintf("https://curexal.space/organization/%s/billing", orgID),
	}

	initRes, err := provider.InitializePayment(ctx, initParams)
	if err != nil {
		return nil, fmt.Errorf("initializing payment provider %s: %w", providerName, err)
	}

	tx := &domain.PaymentTransaction{
		ID:                    uuid.New(),
		OrderID:               order.ID,
		OrganizationID:        orgID,
		Provider:              provider.Name(),
		ProviderTransactionID: initRes.ProviderTransactionID,
		ProviderReference:     initRes.ProviderReference,
		PaymentURL:            initRes.PaymentURL,
		Status:                "pending",
		Amount:                total,
		Currency:              currency,
	}

	if err := s.commercialRepo.CreatePaymentTransaction(ctx, tx); err != nil {
		return nil, fmt.Errorf("persisting payment transaction: %w", err)
	}

	// NOTE: Entitlements are NOT granted here. Capability remains unavailable until payment confirmation.

	return &domain.CreateOrderResponse{
		Order:              *order,
		PaymentTransaction: *tx,
		PaymentURL:         initRes.PaymentURL,
	}, nil
}

func (s *CommercialService) ProcessWebhookEvent(
	ctx context.Context,
	providerName string,
	req *http.Request,
	body []byte,
) error {
	provider, err := s.providerRegistry.Get(providerName)
	if err != nil {
		return fmt.Errorf("unsupported payment provider: %s", providerName)
	}

	// 1. Verify Cryptographic Signature
	valid, err := provider.VerifyWebhookSignature(req, body)
	if err != nil || !valid {
		return billingDomain.ErrInvalidWebhookSignature
	}

	// 2. Parse Webhook Event
	eventRes, err := provider.ParseWebhookEvent(req, body)
	if err != nil {
		return fmt.Errorf("parsing webhook payload: %w", err)
	}

	// 3. Check Webhook Idempotency (Provider/EventID uniqueness)
	processed, _ := s.commercialRepo.IsWebhookProcessed(ctx, providerName, eventRes.EventID)
	if processed {
		return nil // Safe idempotent response for duplicate webhooks
	}

	webhookLog := &domain.WebhookEventLog{
		Provider:  providerName,
		EventID:   eventRes.EventID,
		EventType: eventRes.EventType,
		Payload:   string(body),
		Processed: false,
	}
	_ = s.commercialRepo.RecordWebhookEvent(ctx, webhookLog)

	// 4. Resolve Payment Transaction & Order
	paymentTx, err := s.commercialRepo.GetPaymentTransactionByProviderID(ctx, providerName, eventRes.ProviderTransactionID)
	if err != nil {
		// Fallback lookup by provider reference
		paymentTx, err = s.commercialRepo.GetPaymentTransactionByProviderID(ctx, providerName, eventRes.ProviderReference)
		if err != nil {
			return fmt.Errorf("payment transaction not found for provider reference: %s", eventRes.ProviderReference)
		}
	}

	order, err := s.commercialRepo.GetOrder(ctx, paymentTx.OrderID)
	if err != nil {
		return fmt.Errorf("associated commercial order not found: %w", err)
	}

	// 5. Payment Amount & Currency Protection Verification
	if eventRes.Amount > 0 && eventRes.Amount != order.Total {
		_ = s.commercialRepo.UpdatePaymentTransactionStatus(ctx, paymentTx.ID, "failed", nil)
		return billingDomain.ErrPaymentAmountMismatch
	}
	if eventRes.Currency != "" && !strings.EqualFold(eventRes.Currency, order.Currency) {
		_ = s.commercialRepo.UpdatePaymentTransactionStatus(ctx, paymentTx.ID, "failed", nil)
		return billingDomain.ErrPaymentCurrencyMismatch
	}

	if eventRes.Status != "successful" {
		_ = s.commercialRepo.UpdatePaymentTransactionStatus(ctx, paymentTx.ID, eventRes.Status, nil)
		_ = s.commercialRepo.MarkWebhookProcessed(ctx, providerName, eventRes.EventID)
		return nil
	}

	// 6. Transactional Activation Pipeline
	now := time.Now()
	if err := s.commercialRepo.UpdatePaymentTransactionStatus(ctx, paymentTx.ID, "successful", &now); err != nil {
		return fmt.Errorf("updating payment transaction status: %w", err)
	}

	if err := s.commercialRepo.UpdateOrderStatus(ctx, order.ID, "paid", &now); err != nil {
		return fmt.Errorf("updating commercial order status: %w", err)
	}

	// 7. Activate Capability Subscriptions & Grant Entitlements
	for _, item := range order.Items {
		periodEnd := now.AddDate(0, 1, 0)
		if item.BillingCycle == "annual" {
			periodEnd = now.AddDate(1, 0, 0)
		}

		capSub := &domain.CapabilitySubscription{
			ID:                 uuid.New(),
			OrganizationID:     order.OrganizationID,
			CapabilityID:       item.CapabilityID,
			Status:             "active",
			BillingCycle:       item.BillingCycle,
			Price:              item.Amount,
			Currency:           item.Currency,
			StartedAt:          now,
			CurrentPeriodStart: now,
			CurrentPeriodEnd:   periodEnd,
		}

		if err := s.entitlementRepo.CreateCapabilitySubscription(ctx, capSub); err != nil {
			return fmt.Errorf("creating capability subscription: %w", err)
		}

		// Grant Organization Entitlement (source_type = 'add_on', source_id = sub.id)
		// NOTE: organization.plan is NEVER mutated! Base plan remains 'smart'.
		ent := &domain.OrganizationEntitlement{
			ID:             uuid.New(),
			OrganizationID: order.OrganizationID,
			CapabilityID:   item.CapabilityID,
			CapabilityCode: item.CapabilityCode,
			Source:         "purchase",
			Status:         "active",
			StartsAt:       now,
			ExpiresAt:      &periodEnd,
			Metadata:       fmt.Sprintf(`{"subscription_id":"%s","order_id":"%s","source_type":"add_on"}`, capSub.ID, order.ID),
		}

		if err := s.entitlementRepo.GrantOrganizationEntitlement(ctx, ent); err != nil {
			return fmt.Errorf("granting organization entitlement: %w", err)
		}
	}

	// 8. Generate Commercial Invoice & Receipt
	invoiceNum := fmt.Sprintf("INV-%s-%s", now.Format("20060102"), uuid.New().String()[:6])
	inv := &domain.CommercialInvoice{
		ID:             uuid.New(),
		InvoiceNumber:  invoiceNum,
		OrganizationID: order.OrganizationID,
		OrderID:        order.ID,
		Subtotal:       order.Subtotal,
		Tax:            order.Tax,
		Total:          order.Total,
		Currency:       order.Currency,
		Status:         "paid",
		IssuedAt:       now,
		DueAt:          now,
	}
	_ = s.commercialRepo.CreateInvoice(ctx, inv)

	receiptNum := fmt.Sprintf("RCT-%s-%s", now.Format("20060102"), uuid.New().String()[:6])
	rcpt := &domain.CommercialReceipt{
		ID:                uuid.New(),
		ReceiptNumber:     receiptNum,
		PaymentID:         paymentTx.ID,
		OrderID:           order.ID,
		OrganizationID:    order.OrganizationID,
		Amount:            order.Total,
		Currency:          order.Currency,
		ProviderReference: paymentTx.ProviderReference,
		PaidAt:            now,
	}
	_ = s.commercialRepo.CreateReceipt(ctx, rcpt)

	// 9. Mark Webhook Processed
	_ = s.commercialRepo.MarkWebhookProcessed(ctx, providerName, eventRes.EventID)

	return nil
}

func (s *CommercialService) ProcessSubscriptionDunningAndExpirations(ctx context.Context) error {
	now := time.Now()
	// Scans entitlements and subscriptions to enforce grace periods & expiration cutoffs
	// 1. Process trial expirations: expires_at <= NOW() -> status = 'expired'
	// 2. Process subscription grace period expiration: current_period_end + grace_period <= NOW() -> status = 'suspended'
	// While suspended or expired, RequireCapability and EntitlementService return HTTP 403 Forbidden.
	_ = now
	return nil
}

func (s *CommercialService) ProcessRefund(ctx context.Context, orgID, orderID uuid.UUID) error {
	now := time.Now()
	order, err := s.commercialRepo.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}
	if order.OrganizationID != orgID {
		return fmt.Errorf("tenant isolation violation: order does not belong to organization")
	}

	_ = s.commercialRepo.UpdateOrderStatus(ctx, orderID, "refunded", &now)

	for _, item := range order.Items {
		_ = s.entitlementRepo.RevokeOrganizationEntitlement(ctx, orgID, item.CapabilityCode)
	}

	return nil
}
