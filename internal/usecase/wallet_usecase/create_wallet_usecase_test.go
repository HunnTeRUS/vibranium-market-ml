package wallet_usecase

import (
	"errors"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateWallet(t *testing.T) {
	t.Run("wallet created successfully", func(t *testing.T) {
		mockRepo := new(MockWalletRepository)

		mockRepo.On("CreateWallet", mock.MatchedBy(func(w *wallet.Wallet) bool {
			// Validate other fields if necessary
			return w.Balance == 0 && w.Vibranium == 0
		})).Return(nil)

		wu := NewWalletUsecase(mockRepo)

		result, err := wu.CreateWallet()
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0.0, result.Balance)
		assert.Equal(t, 0, result.Vibranium)

		mockRepo.AssertExpectations(t)
	})

	t.Run("error trying to create wallet", func(t *testing.T) {
		mockRepo := new(MockWalletRepository)

		mockRepo.On("CreateWallet", mock.Anything).Return(errors.New("create wallet error"))

		wu := NewWalletUsecase(mockRepo)

		result, err := wu.CreateWallet()
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "create wallet error", err.Error())

		mockRepo.AssertExpectations(t)
	})
}

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) CreateWallet(wallet *wallet.Wallet) error {
	args := m.Called(wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) DepositToWallet(userID string, amount float64, vibranium int) error {
	args := m.Called(userID, amount, vibranium)
	return args.Error(0)
}

func (m *MockWalletRepository) GetWallet(userID string) (*wallet.Wallet, error) {
	args := m.Called(userID)
	return args.Get(0).(*wallet.Wallet), args.Error(1)
}

func (m *MockWalletRepository) UpdateWallet(wallet *wallet.Wallet) error {
	args := m.Called(wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) SaveSnapshot() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockWalletRepository) LoadSnapshot() error {
	args := m.Called()
	return args.Error(0)
}
