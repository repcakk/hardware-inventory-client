package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var config = loadConfig("config/config.json")

// Config structure for storing server properties
type Config struct {
	ServerIP string `json:"serverIP"`
	Port     string `json:"port"`
	Username string `json:"username"`
}

func loadConfig(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func getGpuName() string {
	cmd, err := exec.Command("cmd", "/C", "wmic path win32_VideoController get name").CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}
	gpuName := string(cmd)
	gpuName = strings.Replace(gpuName, "Name", "", 1)
	gpuName = strings.Replace(gpuName, "\n", "", 1)
	gpuName = strings.TrimSpace(gpuName)
	return gpuName
}

func getPnpDeviceID() string {
	cmd, err := exec.Command("cmd", "/C", "wmic path win32_videocontroller get pnpdeviceid").CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}
	gpuName := string(cmd)
	gpuName = strings.Replace(gpuName, "PNPDeviceID", "", 1)
	gpuName = strings.Replace(gpuName, "\n", "", 1)
	gpuName = strings.TrimSpace(gpuName)
	return gpuName
}

func getHostname() string {
	cmd, err := exec.Command("cmd", "/C", "hostname").CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}
	hostname := string(cmd)
	hostname = strings.TrimSpace(hostname)
	return hostname
}

func updateDataOnServer() {

	hostname := getHostname()
	gpuSN := getPnpDeviceID()
	gpuName := getGpuName()
	username := config.Username

	var serverAddress string = "http://" + config.ServerIP + ":" + config.Port
	var requestString string = serverAddress + "/update"
	requestString = strings.TrimSpace(requestString)

	formData := url.Values{
		"hostname": {hostname},
		"username": {username},
		"gpuSN":    {gpuSN},
		"gpuName":  {gpuName},
	}
	_, err := http.PostForm(requestString, formData)

	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	updateDataOnServer()
}
