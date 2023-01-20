package objects

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/smartaquarius10/kubecli/cmd"
)

func GetKubernetesObject(namespace string, objectType string) []byte {
	return cmd.ExecuteCommand("get", objectType, "-n", namespace, "-o", "name")
}

func GetObjectJson(name string, namespace string, object string) []byte {
	return cmd.ExecuteCommand("get", object, name, "-n", namespace, "-o", "json")
}

func SelectObject(namespace string, object string, extrachars string) string {
	stdout := GetKubernetesObject(namespace, object)
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	counter := 1
	cmdmap := make(map[int]string)
	for scanner.Scan() {
		name := RemoveExtraChars(scanner.Text(), extrachars)
		fmt.Println(counter, ")", name)
		cmdmap[counter] = name
		counter++
	}
	var deploment_count int
	fmt.Print("Select:")
	fmt.Scanf("%d", &deploment_count)
	return cmdmap[deploment_count]
}

func RemoveExtraChars(text string, extrachars string) string {
	return strings.ReplaceAll(text, extrachars, "")
}
