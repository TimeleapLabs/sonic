package launcher

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-opera/inter/validatorpk"
	"github.com/Fantom-foundation/go-opera/valkeystore"
)

func addFakeValidatorKey(ctx *cli.Context, key *ecdsa.PrivateKey, pubkey validatorpk.PubKey, valKeystore valkeystore.RawKeystoreI) error {
	// add fake validator key
	if key != nil && !valKeystore.Has(pubkey) {
		err := valKeystore.Add(pubkey, crypto.FromECDSA(key), validatorpk.FakePassword)
		if err != nil {
			return fmt.Errorf("failed to add fake validator key: %v", err)
		}
	}
	return nil
}

// makeValidatorPasswordList reads password lines from the file specified by the global --validator.password flag.
func makeValidatorPasswordList(ctx *cli.Context) ([]string, error) {
	if path := ctx.GlobalString(validatorPasswordFlag.Name); path != "" {
		text, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read validator password file: %w", err)
		}
		lines := strings.Split(string(text), "\n")
		// Sanitise DOS line endings.
		for i := range lines {
			lines[i] = strings.TrimRight(lines[i], "\r")
		}
		return lines, nil
	}
	if ctx.GlobalIsSet(FakeNetFlag.Name) {
		return []string{validatorpk.FakePassword}, nil
	}
	return nil, nil
}

func unlockValidatorKey(ctx *cli.Context, pubKey validatorpk.PubKey, valKeystore valkeystore.KeystoreI) error {
	if !valKeystore.Has(pubKey) {
		return valkeystore.ErrNotFound
	}
	var err error
	for trials := 0; trials < 3; trials++ {
		prompt := fmt.Sprintf("Unlocking validator key %s | Attempt %d/%d", pubKey.String(), trials+1, 3)
		passwordList, err := makeValidatorPasswordList(ctx)
		if err != nil {
			return err
		}
		password, err := GetPassPhrase(prompt, false, 0, passwordList)
		if err != nil {
			return err
		}
		err = valKeystore.Unlock(pubKey, password)
		if err == nil {
			log.Info("Unlocked validator key", "pubkey", pubKey.String())
			return nil
		}
		if err.Error() != "could not decrypt key with given password" {
			return err
		}
	}
	// All trials expended to unlock account, bail out
	return err
}
