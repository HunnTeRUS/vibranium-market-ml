package wallet_repository

func (wr *walletRepository) DepositToWallet(userID string, amount float64, vibranium int) error {
	if wallet, exists := wr.GetWalletBalance(userID); exists {
		wallet.Balance += amount
		wallet.Vibranium += vibranium

		err := wr.UpdateWallet(wallet)
		if err != nil {
			return err
		}

		return nil
	}

	wallet, err := wr.GetWallet(userID)
	if err != nil {
		return err
	}

	wallet.Balance += amount
	wallet.Vibranium += vibranium

	err = wr.UpdateWallet(wallet)
	if err != nil {
		return err
	}

	return nil
}
