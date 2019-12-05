package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var config = loadConfig("config/config.json")

// Config structure for storing server properties
type Config struct {
	ServerIP string `json:"serverIP"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
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
	pnpDeviceID := string(cmd)
	pnpDeviceID = strings.Replace(pnpDeviceID, "PNPDeviceID", "", 1)
	pnpDeviceID = strings.Replace(pnpDeviceID, "\n", "", 1)
	pnpDeviceID = strings.TrimSpace(pnpDeviceID)

	splittedID := strings.Split(pnpDeviceID, "&")

	gpuSN, err := strconv.ParseUint(splittedID[4], 16, 64)
	if err != nil {
		fmt.Println(err.Error())
	}

	return strconv.FormatUint(gpuSN, 10)
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

func getMac() string {
	cmd, err := exec.Command("cmd", "/C", "getmac").CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}

	resultString := strings.Replace(string(cmd), "Physical Address", "", 1)
	resultString = strings.Replace(resultString, "Transport Name", "", 1)
	resultString = strings.Replace(resultString, "=", "", -1)
	resultString = strings.TrimSpace(resultString)

	mac := strings.Split(resultString, " ")[0]
	mac = strings.TrimSpace(mac)

	return mac
}

func updateDataOnServer() {

	macAddress := getMac()
	hostname := getHostname()
	gpuSN := getPnpDeviceID()
	gpuName := getGpuName()
	username := config.Username
	surname := config.Surname
	email := config.Email

	var serverAddress string = "http://" + config.ServerIP + ":" + config.Port
	var requestString string = serverAddress + "/update"
	requestString = strings.TrimSpace(requestString)

	formData := url.Values{
		"macAddress": {macAddress},
		"hostname":   {hostname},
		"username":   {username},
		"surname":    {surname},
		"email":      {email},
		"gpuSN":      {gpuSN},
		"gpuName":    {gpuName},
	}
	_, err := http.PostForm(requestString, formData)

	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	updateDataOnServer()
}
