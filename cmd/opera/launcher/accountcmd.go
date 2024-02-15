// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package launcher

import (
	"fmt"
	"github.com/Fantom-foundation/go-opera/cmd/opera/launcher/utils"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/console/prompt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"gopkg.in/urfave/cli.v1"
	"strings"
)

// unlockAccounts unlocks any account specifically requested.
func unlockAccounts(ctx *cli.Context, stack *node.Node) error {
	var unlocks []string
	inputs := strings.Split(ctx.GlobalString(utils.UnlockedAccountFlag.Name), ",")
	for _, input := range inputs {
		if trimmed := strings.TrimSpace(input); trimmed != "" {
			unlocks = append(unlocks, trimmed)
		}
	}
	// Short circuit if there is no account to unlock.
	if len(unlocks) == 0 {
		return nil
	}
	// If insecure account unlocking is not allowed if node's APIs are exposed to external.
	// Print warning log to user and skip unlocking.
	if !stack.Config().InsecureUnlockAllowed && stack.Config().ExtRPCEnabled() {
		return fmt.Errorf("account unlock with HTTP access is forbidden")
	}
	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	passwords, err := utils.MakePasswordList(ctx)
	if err != nil {
		return err
	}
	for i, account := range unlocks {
		if _, _, err := UnlockAccount(ks, account, i, passwords); err != nil {
			return err
		}
	}
	return nil
}

// UnlockAccount tries unlocking the specified account a few times.
func UnlockAccount(ks *keystore.KeyStore, address string, i int, passwords []string) (accounts.Account, string, error) {
	if !common.IsHexAddress(address) {
		return accounts.Account{}, "", fmt.Errorf("could not unlock account - '%s' is not an address", address)
	}
	account := accounts.Account{Address: common.HexToAddress(address)}
	var err error
	for trials := 0; trials < 3; trials++ {
		prompt := fmt.Sprintf("Unlocking account %s | Attempt %d/%d", address, trials+1, 3)
		password, errPass := GetPassPhrase(prompt, false, i, passwords)
		if errPass != nil {
			return accounts.Account{}, "", errPass
		}
		err = ks.Unlock(account, password)
		if err == nil {
			log.Info("Unlocked account", "address", account.Address.Hex())
			return account, password, nil
		}
		if err, ok := err.(*keystore.AmbiguousAddrError); ok {
			log.Info("Unlocked account", "address", account.Address.Hex())
			accountRecovered, errRecovery := ambiguousAddrRecovery(ks, err, password)
			if errRecovery != nil {
				return accounts.Account{}, "", errRecovery
			}
			return accountRecovered, password, nil
		}
		if err != keystore.ErrDecrypt {
			// No need to prompt again if the error is not decryption-related.
			break
		}
	}
	// All trials expended to unlock account, bail out
	return accounts.Account{}, "", fmt.Errorf("failed to unlock account %s (%w)", address, err)
}

// GetPassPhrase retrieves the password associated with an account, either fetched
// from a list of preloaded passphrases, or requested interactively from the user.
func GetPassPhrase(msg string, confirmation bool, i int, passwords []string) (string, error) {
	// If a list of passwords was supplied, retrieve from them
	if len(passwords) > 0 {
		if i < len(passwords) {
			return passwords[i], nil
		}
		return passwords[len(passwords)-1], nil
	}
	// Otherwise prompt the user for the password
	if msg != "" {
		fmt.Println(msg)
	}
	password, err := prompt.Stdin.PromptPassword("Passphrase: ")
	if err != nil {
		return "", fmt.Errorf("failed to read passphrase: %v", err)
	}
	if confirmation {
		confirm, err := prompt.Stdin.PromptPassword("Repeat passphrase: ")
		if err != nil {
			return "", fmt.Errorf("failed to read passphrase confirmation: %v", err)
		}
		if password != confirm {
			return "", fmt.Errorf("passphrases do not match")
		}
	}
	return password, nil
}

func ambiguousAddrRecovery(ks *keystore.KeyStore, err *keystore.AmbiguousAddrError, auth string) (accounts.Account, error) {
	fmt.Printf("Multiple key files exist for address %x:\n", err.Addr)
	for _, a := range err.Matches {
		fmt.Println("  ", a.URL)
	}
	fmt.Println("Testing your passphrase against all of them...")
	var match *accounts.Account
	for _, a := range err.Matches {
		if err := ks.Unlock(a, auth); err == nil {
			match = &a
			break
		}
	}
	if match == nil {
		return accounts.Account{}, fmt.Errorf("none of the listed files could be unlocked")
	}
	fmt.Printf("Your passphrase unlocked %s\n", match.URL)
	fmt.Println("In order to avoid this warning, you need to remove the following duplicate key files:")
	for _, a := range err.Matches {
		if a != *match {
			fmt.Println("  ", a.URL)
		}
	}
	return *match, nil
}
