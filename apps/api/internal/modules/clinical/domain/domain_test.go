package domain_test

import (
	"testing"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/clinical/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCatalogItemEntity(t *testing.T) {
	itemID := uuid.New()
	item := domain.CatalogItem{
		ID:        itemID,
		Code:      "FBC",
		Name:      "Full Blood Count",
		Type:      "test",
		BasePrice: 7500.00,
		TatHours:  24,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.Equal(t, itemID, item.ID)
	assert.Equal(t, "FBC", item.Code)
	assert.Equal(t, "Full Blood Count", item.Name)
	assert.Equal(t, 7500.00, item.BasePrice)
	assert.Equal(t, 24, item.TatHours)
}
