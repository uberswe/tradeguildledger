package server

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"
	"time"
)

type AddonData struct {
	Version string
}

func buildWindowsClient() {
	// Build latest client version
	// TODO inject version specific for each client
	log.Println("Building windows client")
	executable := "go"
	if runtime.GOOS == "linux" {
		executable = "/usr/local/go/bin/go"
	}
	cmd := exec.Command(
		executable,
		"build",
		"-ldflags",
		fmt.Sprintf("-H=windowsgui -X 'main.Version=%s' -X 'main.BuildTime=%s' -X 'main.APIKey=%s'", serverVersion, time.Now().Format(time.RFC3339), "DEV"),
		"-o",
		"./downloads/tgl.exe",
		"./cmd/client/main.go")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if runtime.GOOS == "linux" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOPATH=%s", os.Getenv("GOPATH")))
		cmd.Env = append(cmd.Env, fmt.Sprintf("HOME=%s", os.Getenv("HOME")))
		cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
		cmd.Env = append(cmd.Env, "GOOS=windows")
		cmd.Env = append(cmd.Env, "GO386=softfloat")
		cmd.Env = append(cmd.Env, "GOARCH=386")
		cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
		cmd.Env = append(cmd.Env, "CXX=i686-w64-mingw32-g++")
		cmd.Env = append(cmd.Env, "CC=i686-w64-mingw32-gcc")
	} else if runtime.GOOS == "darwin" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("GOPATH=%s", os.Getenv("GOPATH")))
		cmd.Env = append(cmd.Env, fmt.Sprintf("HOME=%s", os.Getenv("HOME")))
		cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
		cmd.Env = append(cmd.Env, "GOOS=windows")
		cmd.Env = append(cmd.Env, "GOARCH=amd64")
		cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
		cmd.Env = append(cmd.Env, "CC=x86_64-w64-mingw32-gcc")
	}
	if err := cmd.Run(); err != nil {
		log.Println(out.String())
		log.Println(stderr.String())
		log.Println(err)
		return
	}
	log.Println("Completed building windows client")
}

type ItemInfo struct {
	ID                    string
	Name                  string
	TwentyFourHourAverage string
	SevenDayAverage       string
	ThirtyDayAverage      string
	LowBuy                string
	HighBuy               string
}

func buildAddonZip() {
	// Build latest addon version
	log.Println("Building zip addon files")
	buf := new(bytes.Buffer)

	w := zip.NewWriter(buf)

	var files = map[string]interface{}{
		"./TradeGuildLedger.iml": AddonData{Version: serverVersion},
		"./TradeGuildLedger.lua": AddonData{Version: serverVersion},
		"./TradeGuildLedger.txt": AddonData{Version: serverVersion},
		"./TradeGuildLedgerItems.lua": []ItemInfo{
			{
				ID:                    "-",
				Name:                  "-",
				TwentyFourHourAverage: "-",
				SevenDayAverage:       "-",
				ThirtyDayAverage:      "-",
				LowBuy:                "-",
				HighBuy:               "-",
			},
		},
	}
	for file, data := range files {
		f, err := w.Create(strings.TrimPrefix(file, "./"))
		if err != nil {
			log.Println(err)
			return
		}
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Println(err)
			return
		}
		tmpl, err := template.New("lua").Parse(string(content))
		if err != nil {
			log.Println(err)
			return
		}
		buf := &bytes.Buffer{}
		err = tmpl.Execute(buf, data)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = f.Write(buf.Bytes())
		if err != nil {
			log.Println(err)
			return
		}
	}

	err := w.Close()
	if err != nil {
		log.Println(err)
		return
	}
	f, err := os.Create("./downloads/tgl.zip")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = f.Write(buf.Bytes())
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Completed building zip addon files")
}
