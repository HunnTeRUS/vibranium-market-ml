package wallet_usecase

import (
	"github.com/HunnTeRUS/vibranium-market-ml/config/logger"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
)

func (wu *WalletUsecase) DepositWallet(userId string, amount float64, vibranium int) error {
	if err := wallet.ValidateDeposit(userId, amount, vibranium); err != nil {
		logger.Error("action=DepositWallet, message=error calling ValidateDeposit validation", err)
		return err
	}

	return wu.repositoryInterface.DepositToWallet(userId, amount, vibranium)
}
