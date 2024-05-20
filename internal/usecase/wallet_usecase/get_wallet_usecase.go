package wallet_usecase

func (wu *WalletUsecase) GetWallet(userId string) (*WalletOuputDTO, error) {
	walletEntity, err := wu.repositoryInterface.GetWallet(userId)
	if err != nil {
		return nil, err
	}
	return &WalletOuputDTO{
		UserID:    walletEntity.UserID,
		Balance:   walletEntity.Balance,
		Vibranium: walletEntity.Vibranium,
	}, nil
}
