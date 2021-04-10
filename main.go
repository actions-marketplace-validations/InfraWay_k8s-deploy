package main

import (
	"log"
	"os"
	"os/exec"

	"deploy/tools"
)

func main() {
	ans := tools.NewAnalytics()
	var err error
	//dns, err := tools.NewDnsProvider()
	//if err != nil {
	//	log.Fatal(err)
	//}

	cmd := exec.Command("helm", "help")
	cmd.Stdout = os.Stdout

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Command started with error: %v", err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Command finished with error: %v", err)
	}

	err = ans.TrackDeploy(tools.Github, "plyo", "mazahaca")
	if err != nil {
		log.Print(err)
	}

	//out, err := exec.Command("helm", "upgrade").Output()
	//if err != nil {
	//	log.Fatalf("Command finished with error: %v", err)
	//}
	//log.Print(string(out))

	//err = ans.TrackDeploy("plyo")
	//if err != nil {
	//	log.Printf("Analytics finished with error: %v", err)
	//}
}
