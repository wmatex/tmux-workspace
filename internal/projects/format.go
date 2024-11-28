package projects

import (
	"fmt"
)

var (
	directoryIcon string = ""
	sessionIcon   string = ""
	configIcon    string = ""
	projectIcon   string = ""
)

const (
	blueColor   = 34
	yellowColor = 33
	cyanColor   = 36
	grayColor   = 90
)

func ansiString(code int, s string) string {
	return fmt.Sprintf("\033[%dm%s\033[39m", code, s)
}

func (p *Project) Format() string {
	icon := directoryIcon
	var color int
	switch p.Running {
	case true:
		icon = sessionIcon
		color = cyanColor
	case false:
		icon = directoryIcon
		color = blueColor
	}

	format := fmt.Sprintf("%s %s", ansiString(color, p.Name), ansiString(grayColor, p.Path))

	return fmt.Sprintf("%s %s", ansiString(color, icon), format)
}
