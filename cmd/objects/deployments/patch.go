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
		object, _ := cmd.Flags().GetString("object")
		backupAttribute(namespace, attribute, object, filter)
	},
}

func BackupCmd() *cobra.Command {
	backupCmd.Flags().StringP("namespace", "n", "", "Pass namespace having deployments")
	backupCmd.Flags().StringP("attrib", "a", "", "Pass attribute name")
	backupCmd.Flags().StringP("object", "o", "", "Pass object name")
	backupCmd.Flags().StringArrayP("filter", "f", []string{"*"}, "Pass search characters or full name of deployment")
	backupCmd.MarkFlagRequired("namespace")
	backupCmd.MarkFlagRequired("attrib")
	backupCmd.MarkFlagRequired("object")
	return backupCmd
}

func backupAttribute(namespace string, attribute string, object string, filter []string) {
	stdout := objects.GetKubernetesObjects(namespace, object, "name")
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	f, _ := os.Create(namespace)
	var query string = ""
	defer f.Close()
	for scanner.Scan() {
		for _, fi := range filter {
			if fi == "*" || strings.Contains(scanner.Text(), fi) {
				name := objects.RemoveExtraChars(scanner.Text(), "deployment.apps/")
				if query == "" {
					stdoutj := objects.GetObjectJson(name, namespace, object)
					query = cmd.GetJsonNodePath(stdoutj, attribute)
				}
				if query != "" {
					stdout = objects.GetKubernetesObject(namespace, object, name, "jsonpath='{"+query+"}'")
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
		object, _ := cmd.Flags().GetString("object")
		removeAttribute(namespace, attrib, object, filter)
	},
}

func RemoveCmd() *cobra.Command {
	removeCmd.Flags().StringP("namespace", "n", "", "Pass namespace")
	removeCmd.Flags().StringP("attrib", "a", "", "Pass attribute name")
	removeCmd.Flags().StringP("object", "o", "", "Pass object name")
	removeCmd.Flags().StringArrayP("filter", "f", []string{"*"}, "Pass search characters of deployment")
	removeCmd.MarkFlagRequired("namespace")
	removeCmd.MarkFlagRequired("attrib")
	removeCmd.MarkFlagRequired("object")
	return removeCmd
}
func removeAttribute(namespace string, attribute string, object string, filter []string) {
	if !checkfileexist(namespace) {
		backupAttribute(namespace, attribute, object, filter)
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
	Use:   "update",
	Short: "This command updates the singular attributes of object definition",
	Long:  "For eg. update label or annotation etc. However, it won't work when in the section with multiple containers",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		query, _ := cmd.Flags().GetString("query")
		all, _ := cmd.Flags().GetString("all")
		value, _ := cmd.Flags().GetString("value")
		applyAttribute(namespace, query, all, value)
	},
}

func ApplyCmd() *cobra.Command {
	applyCmd.Flags().StringP("namespace", "n", "", "Pass namespace having deployments")
	applyCmd.Flags().StringP("query", "q", "", "Pass json query path split with dot(.) of the attibute. For eg. spec.replicas ")
	applyCmd.Flags().StringP("all", "a", "", "Patch attribute in all deployments")
	applyCmd.Flags().StringP("value", "v", "", "Pass value of the attribute specific deployment. For eg. 2")
	applyCmd.MarkFlagRequired("query")
	applyCmd.MarkFlagRequired("namespace")
	applyCmd.MarkFlagRequired("value")
	return applyCmd
}
func applyAttribute(namespace string, query string, all string, value string) {
	jsonquery := getQuery(query, `"`+value+`"`)
	if all != "" {
		stdout := objects.GetKubernetesObjects(namespace, "deployments", "name")
		scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
		for scanner.Scan() {
			stdout := cmd.ExecuteCommand("patch", "deployments", scanner.Text(), "-n", namespace, "-p", jsonquery)
			fmt.Println(string(stdout))
		}
	} else {
		deployment := objects.SelectObject(namespace, "deployments", "deployment.apps/", "name")
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
