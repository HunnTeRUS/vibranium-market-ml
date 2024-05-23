package order

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestNewOrder(t *testing.T) {
	t.Run("Create valid buy order", func(t *testing.T) {
		order, err := NewOrder("user1", OrderTypeBuy, 100, 10.0)
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, OrderTypeBuy, order.Type)
		assert.Equal(t, 100, order.Amount)
		assert.Equal(t, 10.0, order.Price)
		assert.Equal(t, OrderStatusPending, order.Status)
		assert.Equal(t, 0, order.SellValueRemaining)
	})

	t.Run("Create valid sell order", func(t *testing.T) {
		order, err := NewOrder("user1", OrderTypeSell, 100, 10.0)
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, OrderTypeSell, order.Type)
		assert.Equal(t, 100, order.Amount)
		assert.Equal(t, 10.0, order.Price)
		assert.Equal(t, OrderStatusPending, order.Status)
		assert.Equal(t, 100, order.SellValueRemaining)
	})

	t.Run("Invalid order type", func(t *testing.T) {
		order, err := NewOrder("user1", 3, 100, 10.0)
		assert.Error(t, err)
		assert.Nil(t, order)
	})

	t.Run("Invalid amount", func(t *testing.T) {
		order, err := NewOrder("user1", OrderTypeBuy, -100, 10.0)
		assert.Error(t, err)
		assert.Nil(t, order)
	})

	t.Run("Invalid price", func(t *testing.T) {
		order, err := NewOrder("user1", OrderTypeBuy, 100, -10.0)
		assert.Error(t, err)
		assert.Nil(t, order)
	})
}

func TestOrder_CompleteOrder(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	order := &Order{
		ID:     "order1",
		UserID: "user1",
		Type:   OrderTypeBuy,
		Amount: 100,
		Price:  10.0,
		Status: OrderStatusPending,
	}

	mockRepo.On("UpsertOrder", order).Return(nil)

	order.CompleteOrder(mockRepo)

	assert.Equal(t, OrderStatusCompleted, order.Status)
	mockRepo.AssertExpectations(t)
}

func TestOrder_CancelOrder(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	order := &Order{
		ID:     "order1",
		UserID: "user1",
		Type:   OrderTypeBuy,
		Amount: 100,
		Price:  10.0,
		Status: OrderStatusPending,
	}

	mockRepo.On("UpsertOrder", order).Return(nil)

	err := order.CancelOrder(mockRepo, "test cancel")

	assert.Equal(t, OrderStatusCanceled, order.Status)
	assert.EqualError(t, err, "test cancel")
	mockRepo.AssertExpectations(t)
}

// MockOrderRepository to mock OrderRepositoryInterface
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) UpsertOrder(order *Order) {
	m.Called(order)
}

func (m *MockOrderRepository) GetMemOrder(orderId string) (*Order, bool) {
	args := m.Called(orderId)
	return args.Get(0).(*Order), args.Bool(1)
}

func (m *MockOrderRepository) GetBuyingMatchingOrder(orderEntity *Order) (*Order, error) {
	args := m.Called(orderEntity)
	return args.Get(0).(*Order), args.Error(1)
}

func (m *MockOrderRepository) GetSellingMatchingOrder(orderEntity *Order) (*Order, error) {
	args := m.Called(orderEntity)
	return args.Get(0).(*Order), args.Error(1)
}

func (m *MockOrderRepository) LoadSnapshot() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockOrderRepository) SaveSnapshot() error {
	args := m.Called()
	return args.Error(0)
}
