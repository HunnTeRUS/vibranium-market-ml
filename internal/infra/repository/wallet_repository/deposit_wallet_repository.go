package wallet_repository

func (wr *WalletRepository) DepositToWallet(userID string, amount float64, vibranium int) error {
	if wallet, exists := wr.GetWalletBalance(userID); exists {
		wallet.Balance += amount
		wallet.Vibranium += vibranium

		err := wr.UpdateWallet(wallet)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
