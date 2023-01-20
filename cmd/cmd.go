package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
)

const kubecli string = "kubectl"

func ExecuteCommand(arg ...string) []byte {
	command := exec.Command(kubecli, arg...)
	stdout, err := command.Output()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return stdout
}

func ExecuteSessionCommand(arg ...string) {
	cmd := exec.Command(kubecli, arg...)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	fmt.Println(cmd.Args)
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	fmt.Println("2")
	cmd.Wait()
	fmt.Println("3")
}

func findnestedpath(path string, key string, v interface{}) (string, bool) {
	if v == nil {
		return "", false
	}
	switch vv := v.(type) {
	case string:
		if vv == key {
			return path, true
		}
	case map[string]interface{}:
		for k, v := range vv {
			if k == key {
				return path, true
			}
			if found, ok := findnestedpath(fmt.Sprintf("%s.%s", path, k), key, v); ok {
				return found, ok
			}
		}
	case []interface{}:
		for i, v := range vv {
			if found, ok := findnestedpath(fmt.Sprintf("%s[%d]", path, i), key, v); ok {
				return found, ok
			}

		}
	}
	return "", false
}

func GetJsonNodePath(jsonStr []byte, attribute string) string {
	var data map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &data)
	path, ok := findnestedpath("", attribute, data)
	if ok {
		return path
	}
	return ""
}
