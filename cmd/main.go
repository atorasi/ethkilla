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
	"ethkilla/utils"
)

func main() {
	clearTerminal()
	fmt.Printf("%s\n\n", constants.LOGO)
	log.Println("t.me/tripleshizu t.me/tripleshizu t.me/tripleshizu t.me/tripleshizu")
	log.Println("Donate - 0x4163dfa9eE4A25e950ce1a0A2221FafA29fe2df6 - Any EVM")
	fmt.Println()

	walletSlice, err := account.SliceOfAccs()
	if err != nil {
		log.Fatal(err)
	}
	proxyList, err := utils.ReadFile(`..\data\proxy.txt`)
	if err != nil {
		log.Fatal(err)
	}
	for index, wallet := range walletSlice {
		sideClient, err := account.NewClient(index, proxyList, constants.SETTINGS.SideChain)
		if err != nil {
			log.Printf("An error with RPC connection with %s, check your node %s",
				constants.SETTINGS.SideChain, err)
		}
		defer sideClient.Close()

		ethClient, err := account.NewClient(index, proxyList, "ETHEREUM")
		if err != nil {
			log.Printf("An error with RPC connection with %s, check your node %s",
				"ETHEREUM", err)
		}
		defer ethClient.Close()

		if constants.SETTINGS.NeedNonEth {
			module, err := run.SideActions(wallet.Index, wallet, sideClient, ethClient)
			if err != nil {
				log.Printf("Acc.%d | An error was occured with %s: %v", wallet.Index, module, err)
			}
		}
		modules := run.ModulesList()
		for _, module := range modules {
			if _, err := run.EtherActions(wallet.Index, wallet, module, ethClient); err != nil {
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

	log.Println("t.me/tripleshizu t.me/tripleshizu t.me/tripleshizu t.me/tripleshizu")
	log.Println("Donate - 0x4163dfa9eE4A25e950ce1a0A2221FafA29fe2df6 - Any EVM")
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
