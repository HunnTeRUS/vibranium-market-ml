package wallet_repository

import (
	"encoding/json"
	"fmt"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"os"
	"sync"
)

type WalletRepository struct {
	sync.RWMutex

	wallets map[string]*wallet.Wallet
}

func NewWalletRepository() *WalletRepository {
	return &WalletRepository{
		wallets: make(map[string]*wallet.Wallet),
	}
}

func (wr *WalletRepository) CreateWallet(wallet *wallet.Wallet) error {
	wr.UpdateLocalWalletReference(wallet)

	return nil
}

func (wr *WalletRepository) LoadSnapshot() error {
	file, err := os.Open(os.Getenv("WALLETS_SNAPSHOT_FILE"))
	if err != nil {
		return err
	}
	defer file.Close()

	wr.Lock()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&wr.wallets)
	if err != nil {
		return err
	}
	wr.Unlock()

	return nil
}

func (wr *WalletRepository) SaveSnapshot() error {
	if len(wr.wallets) == 0 {
		return nil
	}

	walletsSnapshotFile := os.Getenv("WALLETS_SNAPSHOT_FILE")
	if walletsSnapshotFile == "" {
		return fmt.Errorf("environment variable WALLETS_SNAPSHOT_FILE not set")
	}

	wr.RLock()
	defer wr.RUnlock()

	file, err := os.Create(walletsSnapshotFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(wr.wallets)
	if err != nil {
		return err
	}

	return nil
}
