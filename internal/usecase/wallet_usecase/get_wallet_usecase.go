package wallet_usecase

import "github.com/HunnTeRUS/vibranium-market-ml/config/logger"

func (wu *WalletUsecase) GetWallet(userId string) (*WalletOuputDTO, error) {
	walletEntity, err := wu.repositoryInterface.GetWallet(userId)
	if err != nil {
		logger.Error("action=GetWallet, message=error calling GetWallet repository", err)
		return nil, err
	}
	return &WalletOuputDTO{
		UserID:    walletEntity.UserID,
		Balance:   walletEntity.Balance,
		Vibranium: walletEntity.Vibranium,
	}, nil
}
