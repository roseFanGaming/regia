package regia

import (
	"fmt"
	"strings"
)

const (
	colorRed     = iota + 91 // red
	colorGreen               // green
	colorYellow              // yellow
	colorBlue                // blue
	colorMagenta             // magenta
)

var Banner = `
██████╗ ███████╗ ██████╗ ██╗ █████╗ 
██╔══██╗██╔════╝██╔════╝ ██║██╔══██╗
██████╔╝█████╗  ██║  ███╗██║███████║
██╔══██╗██╔══╝  ██║   ██║██║██╔══██║
██║  ██║███████╗╚██████╔╝██║██║  ██║
╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝╚═╝  ╚═╝
`

func formatColor(text string, color int) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, text)
}

func getCleanedRequestMapping(mapping map[string]string) map[string]string {
	cleanedMapping := make(map[string]string)
	for handleName, requestMethod := range mapping {
		requestMethodUpper := strings.ToUpper(requestMethod)
		for index, method := range httpMethods {
			if requestMethodUpper == method {
				break
			} else if index == (len(httpMethods)-1) && requestMethodUpper != method {
				panic("invalid method" + requestMethod)
			}
		}
		cleanedMapping[handleName] = requestMethodUpper
	}
	return cleanedMapping
}
