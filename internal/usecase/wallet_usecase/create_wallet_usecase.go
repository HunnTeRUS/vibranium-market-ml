package wallet_usecase

import (
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
)

type WalletDepositInputDTO struct {
	UserID string  `json:"userId" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

type WalletOuputDTO struct {
	UserID    string  `json:"user_id"`
	Balance   float64 `json:"balance"`
	Vibranium int     `json:"vibranium"`
}

type WalletUsecaseInterface interface {
	CreateWallet() (*WalletOuputDTO, error)
	DepositWallet(userId string, amount float64) error
	GetWallet(userId string) (*WalletOuputDTO, error)
}

type WalletUsecase struct {
	repositoryInterface wallet.WalletRepositoryInterface
}

func NewWalletUsecase(repositoryInterface wallet.WalletRepositoryInterface) *WalletUsecase {
	return &WalletUsecase{repositoryInterface}
}

func (wu *WalletUsecase) CreateWallet() (*WalletOuputDTO, error) {
	walletEntity := wallet.NewWallet()
	if err := wu.repositoryInterface.CreateWallet(walletEntity); err != nil {
		return nil, err
	}

	return &WalletOuputDTO{
		UserID:    walletEntity.UserID,
		Balance:   walletEntity.Balance,
		Vibranium: walletEntity.Vibranium,
	}, nil
}
