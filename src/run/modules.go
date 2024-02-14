package run

import (
	"ethkilla/src/constants"
	"math/rand"
	"time"
)

func ModulesList() []string {
	var listOfModules []string
	var moduleType = map[string]int{
		"selftrans": constants.SETTINGS.SelfTranTimes,
		"bungee":    constants.SETTINGS.BungeeTimes,
	}
	for moduleName, moduleReps := range moduleType {
		for i := 0; i < moduleReps; i++ {
			listOfModules = append(listOfModules, moduleName)
		}
	}

	// Перемешиваем список
	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(listOfModules), func(i, j int) {
		listOfModules[i], listOfModules[j] = listOfModules[j], listOfModules[i]
	})

	return listOfModules
}
