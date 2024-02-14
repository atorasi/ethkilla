package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"ethkilla/src/account"
	"ethkilla/src/constants"
	"ethkilla/src/run"
)

func main() {
	clearTerminal()
	fmt.Printf("%s\n\n", constants.LOGO)
	walletSlice, err := account.SliceOfAccs()
	if err != nil {
		log.Fatal(err)
	}
	for _, wallet := range walletSlice {
		if constants.SETTINGS.NeedNonEth {
			module, err := run.SideActions(wallet.Index, wallet)
			if err != nil {
				log.Printf("Acc.%d | An error was occured with %s: %v", wallet.Index, module, err)
			}
		}
		modules := run.ModulesList()
		for _, module := range modules {
			if _, err := run.EtherActions(wallet.Index, wallet, module); err != nil {
				log.Printf("Acc.%d | An error was occured with %s: %v", wallet.Index, module, err)
			}
		}
		if constants.SETTINGS.NeedDelayAcc {
			acc := account.Account{}
			delay := acc.RandomInt(constants.SETTINGS.DelayAccMin, constants.SETTINGS.DelayAccMax)
			log.Printf("Acc.%d | sleep for %d seconds before the next account", wallet.Index, delay)
			time.Sleep(time.Duration(delay) * time.Second)

		}
	}
	log.Println("The software has shut down. Press Enter to exit.")
	fmt.Scanln()
}

func clearTerminal() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
	default:
		cmd = exec.Command("clear")
	}

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Ошибка при очистке терминала: %v", err)
	}
}
