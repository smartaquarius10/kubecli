package objects

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/smartaquarius10/kubecli/cmd"
)

func GetKubernetesNodes(objectType string, selector string) []byte {
	return cmd.ExecuteCommand("get", objectType, "-o", selector)
}

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
	return selectApp(stdout, extrachars, true, false)
}

func SelectContainer(namespace string, podName string, selector string, extrachars string) string {
	stdout := GetKubernetesObject(namespace, "pods", podName, selector)
	return selectApp(stdout, extrachars, false, false)
}

func selectApp(stdout []byte, extrachars string, isPod bool, isNode bool) string {
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
		} else if !isNode {
			fmt.Print("Select Container:")
		} else {
			fmt.Print("Select Node:")
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

func GetPodMemoryInNamespace(namespace string, objectType string) []byte {
	return cmd.ExecuteCommand("top", objectType, "-n", namespace, "--no-headers")
}
func GetPodsWithNodeNameInNamespace(namespace string) []byte {
	return cmd.ExecuteCommand("get", "pods", "-o", "custom-columns=NAME:.metadata.name,Namespace:.metadata.namespace,NodeName:.spec.nodeName", "-n", namespace, "--no-headers", "--field-selector=status.phase==Running")
}
func GetPodMemory(namespace string, podname string) string {
	return string(cmd.ExecuteCommand("top", "pods", podname, "-n", namespace, "--no-headers"))
}
func SelectNodes() string {
	stdout := GetKubernetesNodes("nodes", "name")
	return selectApp(stdout, "node/", false, true)
}
func GetRunningPodsInNode(node_name string) []byte {
	return cmd.ExecuteCommand("get", "pods", "-A", "--no-headers", "--field-selector=status.phase==Running,spec.nodeName=="+node_name)
}

func GetObjectLastUpdatedTimeStamp(objectType string, objectName string, namespace string) string {
	return string(cmd.ExecuteCommand("get", objectType, objectName, "--namespace", namespace, `--show-managed-fields -o jsonpath='{range .metadata.managedFields[*]}{.manager}{" did "}{.operation}{" at "}{.time}{"\n"}{end}'`))
}
