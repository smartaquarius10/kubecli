package deployments

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/smartaquarius10/kubecli/cmd"
	"github.com/smartaquarius10/kubecli/cmd/objects"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup attribute value of the objects in a file",
	Long:  "If filter is passed then filtered objects shall be selected else all objects",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		attribute, _ := cmd.Flags().GetString("attrib")
		filter, _ := cmd.Flags().GetStringArray("filter")
		backupAttribute(namespace, attribute, filter)
	},
}

func BackupCmd() *cobra.Command {
	backupCmd.Flags().StringP("namespace", "n", "", "Pass namespace having deployments")
	backupCmd.Flags().StringP("attrib", "a", "", "Pass attribute name")
	backupCmd.Flags().StringArrayP("filter", "f", []string{"*"}, "Pass search characters or full name of deployment")
	backupCmd.MarkFlagRequired("namespace")
	backupCmd.MarkFlagRequired("attrib")
	return backupCmd
}

func backupAttribute(namespace string, attribute string, filter []string) {
	stdout := objects.GetKubernetesObjects(namespace, "deployments", "name")
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	f, _ := os.Create(namespace)
	var query string = ""
	defer f.Close()
	for scanner.Scan() {
		for _, fi := range filter {
			if fi == "*" || strings.Contains(scanner.Text(), fi) {
				name := objects.RemoveExtraChars(scanner.Text(), "deployment.apps/")
				if query == "" {
					stdoutj := objects.GetObjectJson(name, namespace, "deployments")
					query = cmd.GetJsonNodePath(stdoutj, attribute)
				}
				if query != "" {
					stdout = objects.GetKubernetesObject(namespace, "deployments", name, "jsonpath='{"+query+"}'")
					if len(stdout) > 2 {
						fmt.Println(name + "|" + getQuery(query, strings.Trim(string(stdout), "'")))
						f.WriteString(name + "|" + getQuery(query, strings.Trim(string(stdout), "'")) + "\n")
					} else {
						fmt.Println("object not found in", name)
					}
				} else {
					fmt.Println("object not found in", name)
				}
			}
		}
	}
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove labels specified in yaml",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		attrib, _ := cmd.Flags().GetString("attrib")
		filter, _ := cmd.Flags().GetStringArray("filter")
		removeAttribute(namespace, attrib, filter)
	},
}

func RemoveCmd() *cobra.Command {
	removeCmd.Flags().StringP("namespace", "n", "", "Pass namespace")
	removeCmd.Flags().StringP("attrib", "a", "", "Pass attribute name")
	removeCmd.Flags().StringArrayP("filter", "f", []string{"*"}, "Pass search characters of deployment")
	removeCmd.MarkFlagRequired("namespace")
	removeCmd.MarkFlagRequired("attrib")
	return removeCmd
}
func removeAttribute(namespace string, attribute string, filter []string) {
	if !checkfileexist(namespace) {
		backupAttribute(namespace, attribute, filter)
	}
	stdout := objects.GetKubernetesObjects(namespace, "deployments", "name")
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	query := ""
	jsonquery := ""
	for scanner.Scan() {

		for _, f := range filter {
			if f == "*" || strings.Contains(scanner.Text(), f) {
				name := objects.RemoveExtraChars(scanner.Text(), "deployment.apps/")
				if query == "" {
					stdoutj := objects.GetObjectJson(name, namespace, "deployments")
					query = cmd.GetJsonNodePath(stdoutj, attribute)
					jsonquery = getQuery(query, `{"`+attribute+`":null}`)
				}
				if query != "" {
					stdout = cmd.ExecuteCommand("patch", "deployments", name, "-n", namespace, "-p", jsonquery)
					if len(stdout) > 2 {
						fmt.Println(string(stdout))
					}
				}
			}
		}
	}
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply objects specified in yaml",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		query, _ := cmd.Flags().GetString("query")
		file, _ := cmd.Flags().GetString("file")
		value, _ := cmd.Flags().GetString("value")
		applyAttribute(namespace, query, file, value)
	},
}

func ApplyCmd() *cobra.Command {
	applyCmd.Flags().StringP("namespace", "n", "", "Pass namespace having deployments")
	applyCmd.Flags().StringP("query", "q", "", "Pass query path of the object. If file is not passed")
	applyCmd.Flags().StringP("file", "f", "", "Pass file name if selectors are in backup file")
	applyCmd.Flags().StringP("value", "v", "", "Pass node value for specific deployment. If file is not passed")
	applyCmd.MarkFlagRequired("namespace")
	return applyCmd
}
func applyAttribute(namespace string, query string, file string, value string) {
	if file != "" {
		readFile, _ := os.Open(file)
		defer readFile.Close()
		scanner := bufio.NewScanner(readFile)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			parts := strings.Split(scanner.Text(), "|")
			stdout := cmd.ExecuteCommand("patch", "deployments", parts[0], "-n", namespace, "-p", parts[1])
			fmt.Println(string(stdout))
		}
	} else {
		deployment := objects.SelectObject(namespace, "deployments", "deployment.apps/", "name")
		jsonquery := getQuery(query, `"`+value+`"`)
		stdout := cmd.ExecuteCommand("patch", "deployments", deployment, "-n", namespace, "-p", jsonquery)
		fmt.Println(string(stdout))
	}
}

func getQuery(query string, attribute string) string {
	parts := strings.Split(query, ".")
	var jsonquery string
	var curlybrackets string
	for _, v := range parts {
		if len(v) > 0 {
			jsonquery = jsonquery + `{"` + v + `":`
			curlybrackets = curlybrackets + "}"
		}
	}
	jsonquery = jsonquery + attribute + curlybrackets
	return jsonquery
}

func checkfileexist(namespace string) bool {
	if _, err := os.Stat(namespace); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
