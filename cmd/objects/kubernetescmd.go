package objects

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/smartaquarius10/kubecli/cmd"
)

func GetKubernetesObjects(namespace string, objectType string, selector string) []byte {
	return cmd.ExecuteCommand("get", objectType, "-n", namespace, "-o", selector)
}

func GetKubernetesObject(namespace string, objectType string, objectName string, selector string) []byte {
	return cmd.ExecuteCommand("get", objectType, objectName, "-n", namespace, "-o", selector)
}

func GetObjectJson(name string, namespace string, object string) []byte {
	return cmd.ExecuteCommand("get", object, name, "-n", namespace, "-o", "json")
}

func SelectObject(namespace string, object string, extrachars string, selector string) string {
	stdout := GetKubernetesObjects(namespace, object, selector)
	return selectApp(stdout, extrachars, true)
}

func SelectContainer(namespace string, podName string, selector string, extrachars string) string {
	stdout := GetKubernetesObject(namespace, "pods", podName, selector)
	return selectApp(stdout, extrachars, false)
}

func selectApp(stdout []byte, extrachars string, isPod bool) string {
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	counter := 1
	cmdmap := make(map[int]string)
	for scanner.Scan() {
		name := RemoveExtraChars(scanner.Text(), extrachars)
		if len(name) != 0 {
			fmt.Println(counter, ")", name)
			cmdmap[counter] = name
			counter++
		}
	}
	if counter > 2 {
		var deployment_count int
		if isPod {
			fmt.Print("Select Pod:")
		} else {
			fmt.Print("Select Container:")
		}
		fmt.Scanf("%d", &deployment_count)
		return cmdmap[deployment_count]
	} else if counter > 1 {
		return cmdmap[counter]
	}
	return ""
}

func RemoveExtraChars(text string, extrachars string) string {
	return strings.ReplaceAll(text, extrachars, "")
}
