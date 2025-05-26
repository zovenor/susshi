package app

import (
	_ "embed"

	"github.com/rivo/tview"
)

//go:embed assets/icon.txt
var icon string
var iconTview = tview.TranslateANSI(icon)
