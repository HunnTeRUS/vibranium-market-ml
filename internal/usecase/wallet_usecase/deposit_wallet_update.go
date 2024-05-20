package wallet_usecase

import "github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"

func (wu *WalletUsecase) DepositWallet(userId string, amount float64) error {
	if err := wallet.ValidateDeposit(userId, amount); err != nil {
		return err
	}

	return wu.repositoryInterface.DepositToWallet(userId, amount)
}
