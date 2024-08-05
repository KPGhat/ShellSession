package plugin

import (
	"fmt"
	"github.com/KPGhat/ShellSession/session"
	"github.com/KPGhat/ShellSession/utils"
	"os"
	"path"
)

func (plugin *commandPlugin) Upload(session *session.Session, args []string) string {
	// TODO Adapt to other platforms
	utils.Warning("Error only support for linux now...")
	if len(args) == 0 || len(args) > 2 {
		plugin.UploadHelp()
	}

	file, err := os.Open(args[0])
	if err != nil {
		utils.Warning(err.Error())
		return ""
	}
	defer file.Close()

	var targetFileName string
	if len(args) == 1 {
		targetFileName = path.Base(args[0])
	} else {
		targetFileName = args[1]
	}

	// touch target file
	session.ExecCmd("echo -ne \"\" >" + targetFileName)
	buffer := make([]byte, 1024)
	var result []byte
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			utils.Warning(err.Error())
			return ""
		}
		//utils.Congrats(fmt.Sprintf("Read %d bytes from %s", n, args[0]))
		result = session.ExecCmd("echo -ne \"" + encodeToHex(buffer[:n]) + "\" >>" + targetFileName)
	}

	// return the last cmd result
	return string(result)
}

func (plugin *commandPlugin) UploadHelp() {
	utils.Congrats("Upload Usage: upload src [dst]")
	utils.Congrats("If dst is not set, will upload to the shell current dir")
}

func encodeToHex(data []byte) string {
	var result string
	for _, b := range data {
		if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') {
			result += string(b)
		} else {
			result += fmt.Sprintf("\\x%02x", b)
		}
	}
	return result
}
