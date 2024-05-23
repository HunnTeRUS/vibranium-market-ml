package wallet_repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

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

	snapshotDir := filepath.Dir(walletsSnapshotFile)
	err := os.MkdirAll(snapshotDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return err
	}

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