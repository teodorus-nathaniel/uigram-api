package nodejs

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

func ExecScreenshot(url string) string {
	command := "node node/screenshot.js " + url
	parts := strings.Fields(command)
	data, err := exec.Command(parts[0], parts[1:]...).Output()
	if err != nil {
		return ""
	}

	output := string(data)
	output = output[0 : len(output)-1]
	time.AfterFunc(20*time.Minute, func() {
		os.Remove(output)
	})

	return output
}
