package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/subscription/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
)

type CommercialRepository struct {
	server *server.Server

	// In-memory mock store for unit testing when DB is unavailable
	mu                  sync.RWMutex
	orders              map[uuid.UUID]*domain.CommercialOrder
	orderNumbers        map[string]*domain.CommercialOrder
	paymentTransactions map[uuid.UUID]*domain.PaymentTransaction
	providerTxs         map[string]*domain.PaymentTransaction
	webhooks            map[string]*domain.WebhookEventLog
	invoices            map[uuid.UUID][]domain.CommercialInvoice
	receipts            map[uuid.UUID][]domain.CommercialReceipt
}

func NewCommercialRepository(s *server.Server) *CommercialRepository {
	return &CommercialRepository{
		server:              s,
		orders:              make(map[uuid.UUID]*domain.CommercialOrder),
		orderNumbers:        make(map[string]*domain.CommercialOrder),
		paymentTransactions: make(map[uuid.UUID]*domain.PaymentTransaction),
		providerTxs:         make(map[string]*domain.PaymentTransaction),
		webhooks:            make(map[string]*domain.WebhookEventLog),
		invoices:            make(map[uuid.UUID][]domain.CommercialInvoice),
		receipts:            make(map[uuid.UUID][]domain.CommercialReceipt),
	}
}

func (r *CommercialRepository) CreateOrder(ctx context.Context, order *domain.CommercialOrder, items []domain.CommercialOrderItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	for i := range items {
		if items[i].ID == uuid.Nil {
			items[i].ID = uuid.New()
		}
		items[i].OrderID = order.ID
		items[i].CreatedAt = time.Now()
	}
	order.Items = items

	r.orders[order.ID] = order
	r.orderNumbers[order.OrderNumber] = order

	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return nil
	}

	queryOrder := `
		INSERT INTO subscription.commercial_orders 
		(id, order_number, organization_id, status, currency, subtotal, tax, discount, total, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.server.DB.Pool.Exec(ctx, queryOrder,
		order.ID, order.OrderNumber, order.OrganizationID, order.Status, order.Currency,
		order.Subtotal, order.Tax, order.Discount, order.Total, order.CreatedBy, order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting commercial order: %w", err)
	}

	queryItem := `
		INSERT INTO subscription.commercial_order_items
		(id, order_id, capability_id, capability_code, billing_cycle, unit_price, amount, currency, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	for _, item := range items {
		meta := item.Metadata
		if meta == "" {
			meta = "{}"
		}
		_, err := r.server.DB.Pool.Exec(ctx, queryItem,
			item.ID, item.OrderID, item.CapabilityID, item.CapabilityCode, item.BillingCycle,
			item.UnitPrice, item.Amount, item.Currency, meta, item.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("inserting commercial order item: %w", err)
		}
	}

	return nil
}

func (r *CommercialRepository) GetOrder(ctx context.Context, orderID uuid.UUID) (*domain.CommercialOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if order, exists := r.orders[orderID]; exists {
		return order, nil
	}
	return nil, errors.New("commercial order not found")
}

func (r *CommercialRepository) GetOrderByNumber(ctx context.Context, orderNumber string) (*domain.CommercialOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if order, exists := r.orderNumbers[orderNumber]; exists {
		return order, nil
	}
	return nil, errors.New("commercial order not found")
}

func (r *CommercialRepository) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string, paidAt *time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if order, exists := r.orders[orderID]; exists {
		order.Status = status
		order.PaidAt = paidAt
		order.UpdatedAt = time.Now()
	}

	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return nil
	}

	query := `UPDATE subscription.commercial_orders SET status = $1, paid_at = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.server.DB.Pool.Exec(ctx, query, status, paidAt, orderID)
	return err
}

func (r *CommercialRepository) CreatePaymentTransaction(ctx context.Context, tx *domain.PaymentTransaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if tx.ID == uuid.Nil {
		tx.ID = uuid.New()
	}
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = time.Now()

	r.paymentTransactions[tx.ID] = tx
	key := fmt.Sprintf("%s:%s", tx.Provider, tx.ProviderTransactionID)
	r.providerTxs[key] = tx

	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return nil
	}

	meta := tx.Metadata
	if meta == "" {
		meta = "{}"
	}

	query := `
		INSERT INTO subscription.payment_transactions
		(id, order_id, organization_id, provider, provider_transaction_id, provider_reference, payment_url, status, amount, currency, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := r.server.DB.Pool.Exec(ctx, query,
		tx.ID, tx.OrderID, tx.OrganizationID, tx.Provider, tx.ProviderTransactionID, tx.ProviderReference,
		tx.PaymentURL, tx.Status, tx.Amount, tx.Currency, meta, tx.CreatedAt, tx.UpdatedAt,
	)
	return err
}

func (r *CommercialRepository) GetPaymentTransaction(ctx context.Context, txID uuid.UUID) (*domain.PaymentTransaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if tx, exists := r.paymentTransactions[txID]; exists {
		return tx, nil
	}
	return nil, errors.New("payment transaction not found")
}

func (r *CommercialRepository) GetPaymentTransactionByProviderID(ctx context.Context, provider, providerTxID string) (*domain.PaymentTransaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", provider, providerTxID)
	if tx, exists := r.providerTxs[key]; exists {
		return tx, nil
	}

	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return nil, errors.New("payment transaction not found")
	}

	query := `
		SELECT id, order_id, organization_id, provider, provider_transaction_id, provider_reference, payment_url, status, amount, currency, created_at, updated_at
		FROM subscription.payment_transactions
		WHERE provider = $1 AND (provider_transaction_id = $2 OR provider_reference = $2)`
	row := r.server.DB.Pool.QueryRow(ctx, query, provider, providerTxID)

	var tx domain.PaymentTransaction
	err := row.Scan(&tx.ID, &tx.OrderID, &tx.OrganizationID, &tx.Provider, &tx.ProviderTransactionID, &tx.ProviderReference, &tx.PaymentURL, &tx.Status, &tx.Amount, &tx.Currency, &tx.CreatedAt, &tx.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("payment transaction not found")
		}
		return nil, err
	}

	return &tx, nil
}

func (r *CommercialRepository) UpdatePaymentTransactionStatus(ctx context.Context, txID uuid.UUID, status string, paidAt *time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if tx, exists := r.paymentTransactions[txID]; exists {
		tx.Status = status
		tx.PaidAt = paidAt
		tx.UpdatedAt = time.Now()
	}

	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return nil
	}

	query := `UPDATE subscription.payment_transactions SET status = $1, paid_at = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.server.DB.Pool.Exec(ctx, query, status, paidAt, txID)
	return err
}

func (r *CommercialRepository) RecordWebhookEvent(ctx context.Context, event *domain.WebhookEventLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%s:%s", event.Provider, event.EventID)
	if _, exists := r.webhooks[key]; exists {
		return errors.New("webhook event already exists")
	}

	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	event.CreatedAt = time.Now()
	r.webhooks[key] = event

	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return nil
	}

	payload := event.Payload
	if payload == "" {
		payload = "{}"
	}

	query := `
		INSERT INTO subscription.webhook_events
		(id, provider, event_id, event_type, payload, processed, processed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (provider, event_id) DO NOTHING`
	_, err := r.server.DB.Pool.Exec(ctx, query,
		event.ID, event.Provider, event.EventID, event.EventType, payload,
		event.Processed, event.ProcessedAt, event.CreatedAt,
	)
	return err
}

func (r *CommercialRepository) IsWebhookProcessed(ctx context.Context, provider, eventID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", provider, eventID)
	if evt, exists := r.webhooks[key]; exists && evt.Processed {
		return true, nil
	}

	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return false, nil
	}

	var processed bool
	query := `SELECT processed FROM subscription.webhook_events WHERE provider = $1 AND event_id = $2`
	err := r.server.DB.Pool.QueryRow(ctx, query, provider, eventID).Scan(&processed)
	if err != nil {
		return false, nil
	}
	return processed, nil
}

func (r *CommercialRepository) MarkWebhookProcessed(ctx context.Context, provider, eventID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%s:%s", provider, eventID)
	now := time.Now()
	if evt, exists := r.webhooks[key]; exists {
		evt.Processed = true
		evt.ProcessedAt = &now
	}

	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return nil
	}

	query := `UPDATE subscription.webhook_events SET processed = TRUE, processed_at = NOW() WHERE provider = $1 AND event_id = $2`
	_, err := r.server.DB.Pool.Exec(ctx, query, provider, eventID)
	return err
}

func (r *CommercialRepository) CreateInvoice(ctx context.Context, invoice *domain.CommercialInvoice) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if invoice.ID == uuid.Nil {
		invoice.ID = uuid.New()
	}
	invoice.CreatedAt = time.Now()
	r.invoices[invoice.OrganizationID] = append(r.invoices[invoice.OrganizationID], *invoice)

	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return nil
	}

	query := `
		INSERT INTO subscription.commercial_invoices
		(id, invoice_number, organization_id, order_id, subtotal, tax, total, currency, status, issued_at, due_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.server.DB.Pool.Exec(ctx, query,
		invoice.ID, invoice.InvoiceNumber, invoice.OrganizationID, invoice.OrderID,
		invoice.Subtotal, invoice.Tax, invoice.Total, invoice.Currency, invoice.Status,
		invoice.IssuedAt, invoice.DueAt, invoice.CreatedAt,
	)
	return err
}

func (r *CommercialRepository) GetInvoicesByOrg(ctx context.Context, orgID uuid.UUID) ([]domain.CommercialInvoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if invs, exists := r.invoices[orgID]; exists {
		return invs, nil
	}
	return []domain.CommercialInvoice{}, nil
}

func (r *CommercialRepository) CreateReceipt(ctx context.Context, receipt *domain.CommercialReceipt) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if receipt.ID == uuid.Nil {
		receipt.ID = uuid.New()
	}
	receipt.CreatedAt = time.Now()
	r.receipts[receipt.OrganizationID] = append(r.receipts[receipt.OrganizationID], *receipt)

	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return nil
	}

	query := `
		INSERT INTO subscription.commercial_receipts
		(id, receipt_number, payment_id, order_id, organization_id, amount, currency, provider_reference, paid_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.server.DB.Pool.Exec(ctx, query,
		receipt.ID, receipt.ReceiptNumber, receipt.PaymentID, receipt.OrderID, receipt.OrganizationID,
		receipt.Amount, receipt.Currency, receipt.ProviderReference, receipt.PaidAt, receipt.CreatedAt,
	)
	return err
}

func (r *CommercialRepository) GetReceiptsByOrg(ctx context.Context, orgID uuid.UUID) ([]domain.CommercialReceipt, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if rcpts, exists := r.receipts[orgID]; exists {
		return rcpts, nil
	}
	return []domain.CommercialReceipt{}, nil
}

var _ domain.CommercialRepository = (*CommercialRepository)(nil)
