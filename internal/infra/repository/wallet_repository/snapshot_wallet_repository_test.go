package wallet_repository

import (
	"encoding/json"
	"github.com/HunnTeRUS/vibranium-market-ml/internal/entity/wallet"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSnapshot(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "wallets_snapshot_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	os.Setenv("WALLETS_SNAPSHOT_FILE", tmpFile.Name())
	defer os.Unsetenv("WALLETS_SNAPSHOT_FILE")

	wallets := map[string]*wallet.Wallet{
		"user1": {UserID: "user1", Balance: 100.0, Vibranium: 50},
		"user2": {UserID: "user2", Balance: 200.0, Vibranium: 100},
	}
	data, err := json.Marshal(wallets)
	assert.NoError(t, err)

	_, err = tmpFile.Write(data)
	assert.NoError(t, err)
	tmpFile.Close()

	repo := NewWalletRepository()
	err = repo.LoadSnapshot()
	assert.NoError(t, err)

	assert.Equal(t, wallets, repo.wallets)
}

func TestSaveSnapshot(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "wallets_snapshot")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	snapshotFile := filepath.Join(tmpDir, "wallets_snapshot.json")
	os.Setenv("WALLETS_SNAPSHOT_FILE", snapshotFile)
	defer os.Unsetenv("WALLETS_SNAPSHOT_FILE")

	repo := NewWalletRepository()
	repo.wallets["user1"] = &wallet.Wallet{UserID: "user1", Balance: 100.0, Vibranium: 50}
	repo.wallets["user2"] = &wallet.Wallet{UserID: "user2", Balance: 200.0, Vibranium: 100}

	err = repo.SaveSnapshot()
	assert.NoError(t, err)

	file, err := os.Open(snapshotFile)
	assert.NoError(t, err)
	defer file.Close()

	var wallets map[string]*wallet.Wallet
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&wallets)
	assert.NoError(t, err)
	assert.Equal(t, repo.wallets, wallets)
}

func TestSaveSnapshot_EmptyWallets(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "wallets_snapshot")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	snapshotFile := filepath.Join(tmpDir, "wallets_snapshot.json")
	os.Setenv("WALLETS_SNAPSHOT_FILE", snapshotFile)
	defer os.Unsetenv("WALLETS_SNAPSHOT_FILE")

	repo := NewWalletRepository()

	err = repo.SaveSnapshot()
	assert.NoError(t, err)

	_, err = os.Stat(snapshotFile)
	assert.True(t, os.IsNotExist(err))
}
