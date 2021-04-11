package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"deploy/tools"
)

const ActionDeploy = "rollout"
const ActionDelete = "delete"

var (
	action    string
	namespace string
	release   string
	chart     string
	host      string
	ip        string
)

func init() {
	flag.Parse()
	action = flag.Arg(0)
	namespace = flag.Arg(1)
	release = flag.Arg(2)
	chart = flag.Arg(3)
	host = flag.Arg(4)
	ip = flag.Arg(5)
}

func main() {
	ans := tools.NewAnalytics()
	dns, err := tools.NewDnsProvider()
	if err != nil {
		log.Fatal(err)
	}

	var cmd *exec.Cmd
	if action == ActionDeploy {
		cmd = exec.Command("helm",
			"upgrade",
			"-n", namespace,
			"--install",
			"--wait",
			release,
			chart,
			"--set", fmt.Sprintf("image.tag=%s", os.Getenv("DEPLOY_TAG")),
			"--set", fmt.Sprintf("ingress.baseHost=%s", host),
		)
	} else {
		cmd = exec.Command("helm", "uninstall", release)
	}
	if kubeConfig := os.Getenv("KUBE_CONFIG"); kubeConfig != "" {
		file, err := ioutil.TempFile("kube", "config")
		if err != nil {
			log.Fatal(err)
		}
		path, err := filepath.Abs(filepath.Dir(file.Name()))
		fName := filepath.Base(file.Name())
		if err != nil {
			log.Fatal(err)
		}
		_, err = file.Write([]byte(kubeConfig))
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(file.Name())
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("KUBECONFIG=%s", fmt.Sprintf("%s/%s", path, fName)),
		)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Command started with error: %v", err)
	}
	err = cmd.Wait()
	if err != nil {
		_ = ans.TrackRollout(getDataSource(), getRepoOwner(), getActor())
		log.Fatalf("Command finished with error: %v", err)
	}

	if action == ActionDeploy {
		err = dns.AddRecord(host, ip)
		if err != nil {
			log.Print(err)
		}
	}
	if action == ActionDelete {
		err = dns.RemoveRecord(host, ip)
		if err != nil {
			log.Print(err)
		}
	}

	_ = ans.TrackDeploy(getDataSource(), getRepoOwner(), getActor())
}

func getDataSource() tools.DataSource {
	if test := os.Getenv("GITHUB_ACTIONS"); test != "" {
		return tools.Github
	}
	if test := os.Getenv("GITLAB_CI"); test != "" {
		return tools.Gitlab
	}
	return ""
}

func getRepoOwner() string {
	if test := os.Getenv("GITHUB_REPOSITORY_OWNER"); test != "" {
		return test
	}
	if test := os.Getenv("CI_PROJECT_NAMESPACE"); test != "" {
		return test
	}
	return ""
}

func getActor() string {
	if test := os.Getenv("GITHUB_ACTOR"); test != "" {
		return test
	}
	if test := os.Getenv("GITLAB_USER_LOGIN"); test != "" {
		return test
	}
	return ""
}
