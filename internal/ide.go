package internal

import (
	"fmt"
	"strconv"
	"strings"
)

type JetBrainsIDE int

const (
	IntelliJ JetBrainsIDE = iota
	PyCharm
	WebStorm
	PhpStorm
	CLion
	Rider
	DataGrip
	RubyMine
	AppCode
	Studio
)

var IdeNames = []string{
	"idea",
	"pycharm",
	"webstorm",
	"phpstorm",
	"clion",
	"rider",
	"datagrip",
	"rubymine",
	"appcode",
	"studio", // JetBrains android studio
}

func (ide JetBrainsIDE) String() string {
	return IdeNames[ide]
}

func FromString(input string) (JetBrainsIDE, bool) {
	for i, name := range IdeNames {
		if name == input {
			return JetBrainsIDE(i), true
		}
	}
	return 0, false
}

func AllIDEs() []JetBrainsIDE {
	return []JetBrainsIDE{
		IntelliJ, PyCharm, WebStorm, PhpStorm, CLion, Rider, DataGrip, RubyMine, AppCode,
	}
}

func GetSelectedIDEs(allIdesFlag bool) []JetBrainsIDE {
	if allIdesFlag {
		return AllIDEs()
	}

	fmt.Println("Choose the JetBrains IDEs to patch (comma-separated numbers, default is all):")
	allIDEs := AllIDEs()
	for i, ide := range allIDEs {
		fmt.Printf("%d: %s\n", i+1, ide)
	}
	fmt.Print("> ")
	ideInput := ReadLine()

	if ideInput == "" {
		return allIDEs
	}

	var selectedIDEs []JetBrainsIDE
	for _, s := range strings.Split(ideInput, ",") {
		i, err := strconv.Atoi(strings.TrimSpace(s))
		if err == nil && i > 0 && i <= len(allIDEs) {
			selectedIDEs = append(selectedIDEs, allIDEs[i-1])
		}
	}
	return selectedIDEs
}
