package termio

import (
	"github.com/fatih/color"
)

var (
	Bold = color.New(color.Bold).SprintFunc()
	Gray = color.HiBlackString
)
