package wallet_repository

func (wr *walletRepository) DepositToWallet(userID string, amount float64) error {
	if wallet, exists := wr.GetWalletBalance(userID); exists {
		wallet.Balance += amount

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

	err = wr.UpdateWallet(wallet)
	if err != nil {
		return err
	}

	return nil
}
