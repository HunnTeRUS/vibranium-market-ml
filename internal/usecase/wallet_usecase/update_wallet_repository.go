package wallet_usecase

import "github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"

func (wu *WalletUsecase) UpdateWallet(wallet *wallet.Wallet) error {
	return wu.repositoryInterface.UpdateWallet(wallet)
}
