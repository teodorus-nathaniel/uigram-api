package nodejs

import (
	"os/exec"
	"strings"
)

func ExecScreenshot(url string) string {
	command := "node screenshot.js " + url
	parts := strings.Fields(command)
	data, err := exec.Command(parts[0], parts[1:]...).Output()
	if err != nil {
		panic(err)
	}

	output := string(data)

	return output
}
