package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"time"
)

func main() {
	// setUpST()
	// getListOfPids()
	// startUpST()
	// stopST("28288")
	cleanST()
}

func getListOfPids() []string {
	//go to the supertokens root directory
	files, err := ioutil.ReadDir("../supertokens-root/.started/")
	if err != nil {
		return []string{}
	}
	//then find out all the all the files in the .started directory
	//iterate through all those files
	var result []string
	for _, file := range files {
		//read the contents of those files and push them to an array
		pathOfFileToBeRead := "../supertokens-root/.started/" + file.Name()
		data, err := ioutil.ReadFile(pathOfFileToBeRead)
		if err != nil {
			log.Fatalf(err.Error())
		}
		result = append(result, string(data))
	}
	fmt.Println(result)
	return result
}

func setUpST() {
	installationPath := "../supertokens-root/"
	cmd := exec.Command("cp", "temp/config.yaml", "./config.yaml")
	cmd.Dir = installationPath
	err := cmd.Run()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func startUpST() {
	host := "localhost"
	port := "8000"
	installationPath := "../supertokens-root/"
	pidsBefore := getListOfPids()
	returned := false

	command := fmt.Sprintf(`java -Djava.security.egd=file:/dev/urandom -classpath "./core/*:./plugin-interface/*" io.supertokens.Main ./ DEV host=%s port=%s test_mode`, host, port)

	cmd := exec.Command("bash", "-c", command)

	cmd.Dir = installationPath
	err := cmd.Run()
	if err != nil {
		if !returned {
			returned = true
			log.Fatalf(err.Error(), "could not initiate a supertokens instance")
		}
	}

	startTime := time.Now().Unix()
	for (time.Now().Unix() - startTime) < 30000 {
		pidsAfter := getListOfPids()
		if len(pidsAfter) <= len(pidsBefore) {
			//ask joel about the the process that just stops here why is that required
			continue
		}
		var nonIntersection = []string{""}
		for i := 0; i < len(pidsAfter); i++ {
			for j := 0; j < len(pidsBefore); j++ {
				if pidsAfter[i] != pidsBefore[j] {
					nonIntersection = append(nonIntersection, pidsAfter[i])
				}
			}
		}
		if len(nonIntersection) != 1 {
			if !returned {
				returned = true
				log.Fatalf("something went wrong while starting up the core")
			}
		} else {
			if !returned {
				returned = true
				return nonIntersection[0]
			}
		}
	}
	if !returned {
		returned = true
		log.Fatalf("something went wrong while starting up the core")
	}
}

func stopST(pid string) {
	installationPath := "../supertokens-root/"
	pidsBefore := getListOfPids()
	if len(pidsBefore) == 0 {
		return
	}
	cmd := exec.Command("kill", pid)
	cmd.Dir = installationPath
	err := cmd.Run()
	if err != nil {
		log.Fatalf(err.Error(), "could not close the supertokens instance")
	}
	startTime := time.Now().Unix()
	for (time.Now().Unix() - startTime) < 30000 {
		pidsAfter := getListOfPids()
		includes := false
		for i := 0; i < len(pidsAfter); i++ {
			if pidsAfter[i] == pid {
				includes = true
			}
		}
		if includes {
			//*stopping login for 10mils
			continue
		} else {
			return
		}
	}
	log.Fatalf(err.Error(), "error while stopping st")
}

func killAllSTCoresOnly() {
	pids := getListOfPids()
	for i := 0; i < len(pids); i++ {
		stopST(pids[i])
	}
}

func cleanST() {
	installationPath := "../supertokens-root/"
	cmd := exec.Command("rm", "config.yaml")
	cmd.Dir = installationPath
	err := cmd.Run()
	if err != nil {
		log.Fatalf(err.Error(), "could not delete the config yaml file")
	}

	cmd = exec.Command("rm", "-rf", ".webserver-temp-*")
	cmd.Dir = installationPath
	err = cmd.Run()
	if err != nil {
		log.Fatalf(err.Error(), "could not delete the webserver-temp files")
	}

	cmd = exec.Command("rm", "-rf", ".started")
	cmd.Dir = installationPath
	err = cmd.Run()
	if err != nil {
		log.Fatalf(err.Error(), "could not delete the .started file")
	}
}
