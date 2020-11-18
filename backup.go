package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/vaughan0/go-ini"
)

const Shell = "bash"

type Config struct {
	date       string
	TargetPath string
	Folders    map[string]string
}

func execute(command string) string {
	out, err := exec.Command(Shell, "-c", command).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	output := string(out[:])

	return output
}

func folders(config Config) {
	fmt.Println("")
	fmt.Println("Folders")
	fmt.Println("------------------")
	var target = config.TargetPath
	for name, source := range config.Folders {
		fmt.Printf(name + " ... ")
		execute("rsync --no-compress -ahWq " + source + " " + target + name + "/")
		fmt.Println("OK")
	}
}

/**
 * Ini file into a validated config
 */
func createConfig() Config {
	var config Config

	configFile, err := ini.LoadFile("config.ini")
	if err != nil {
		fmt.Printf(color.RedString("Fail to read file: %v"), err)
		os.Exit(1)
	}

	targetPath, ok := configFile.Get("setting", "targetPath")
	if !ok {
		fmt.Printf(color.RedString("'targetPath' variable missing from 'setting' section"))
		os.Exit(1)
	}

	config.date = time.Now().Format("2006-01-02")
	config.TargetPath = targetPath + config.date + "/"
	config.Folders = make(map[string]string)
	for folderName, folderPath := range configFile["folder"] {
		config.Folders[folderName] = folderPath
	}

	return config
}

func prepare(path string) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		fmt.Println(color.RedString("prepare: The backup for '" + path + "' already exist"))
		os.Exit(1)
	}

	err := os.Mkdir(path, 0755)
	if err != nil {
		fmt.Printf(color.RedString("prepare: could not create folder [" + path + "]"))
		os.Exit(1)
	}
}

func main() {
	var config = createConfig()

	fmt.Println(color.BlueString("Start backup: " + config.date))

	prepare(config.TargetPath)
	folders(config)
}
