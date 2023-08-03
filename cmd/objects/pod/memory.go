package pod

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/smartaquarius10/kubecli/cmd/objects"
	"github.com/spf13/cobra"
)

type PodDetails struct {
	count     int
	memory    int
	cpu       int
	namespace string
	nodename  string
}

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "This command return pods memory usage per namespace",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		node, _ := cmd.Flags().GetBool("nodes")
		if !node {
			getPodsMemoryInNamespace(namespace)
		} else {
			getPodsMemoryPerNode()
		}
	},
}

func TopCmd() *cobra.Command {
	topCmd.Flags().StringP("namespace", "n", "default", "Pass namespace having pods")
	topCmd.Flags().BoolP("nodes", "d", false, "Get pods memory by nodes")
	return topCmd
}

func getPodsMemoryInNamespace(namespace string) {
	stdout := objects.GetPodMemoryInNamespace(namespace, "pods")
	podNodesDetails := objects.GetPodsWithNodeNameInNamespace(namespace)
	nodescanner := bufio.NewScanner(strings.NewReader(string(podNodesDetails)))
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	m := make(map[string]*PodDetails)
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	for scanner.Scan() {
		if !strings.Contains(scanner.Text(), "NAME") {
			podDetails := strings.Fields(scanner.Text())
			podName := strings.Split(podDetails[0], "-")[0]
			_, ok := m[podName]
			cpu, _ := strconv.Atoi(re.FindString(podDetails[1]))
			memory, _ := strconv.Atoi(re.FindString(podDetails[2]))
			if ok {
				m[podName].count++

				m[podName].cpu = m[podName].cpu + cpu
				m[podName].memory = m[podName].memory + memory
				for nodescanner.Scan() {
					if strings.Contains(nodescanner.Text(), podName) {
						m[podName].nodename = m[podName].nodename + "\n" + strings.Fields(nodescanner.Text())[2]
						break
					}
				}
			} else {
				m[podName] = &PodDetails{count: 1, memory: memory, cpu: cpu, namespace: namespace}
				for nodescanner.Scan() {
					if strings.Contains(nodescanner.Text(), podName) {
						m[podName].nodename = strings.Fields(nodescanner.Text())[2]
						break
					}
				}
			}
		}
	}
	writeTable(m)
}

func getPodsMemoryPerNode() {
	nodename := objects.SelectNodes()
	stdout := objects.GetRunningPodsInNode(nodename)
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	m := make(map[string]*PodDetails)
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	for scanner.Scan() {
		podDetails := strings.Fields(scanner.Text())
		pod := strings.Fields(objects.GetPodMemory(podDetails[0], podDetails[1]))
		val := strings.Split(pod[0], "-")[0]
		_, ok := m[val]
		cpu, _ := strconv.Atoi(re.FindString(pod[1]))
		memory, _ := strconv.Atoi(re.FindString(pod[2]))
		if ok {
			m[val].count++

			m[val].cpu = m[val].cpu + cpu
			m[val].memory = m[val].memory + memory
		} else {
			m[val] = &PodDetails{count: 1, memory: memory, cpu: cpu, namespace: podDetails[0], nodename: nodename}
		}
	}
	writeTable(m)
}

func writeTable(podDetails map[string]*PodDetails) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Namespace", "Pods", "Replicas", "CPU(m)", "Memory(Mi)", "Node"})
	// t.SortBy([]table.SortBy{
	// 	{Name: "Memory(Mi)", Mode: table.DscNumeric},
	// })
	totalcpu := 0
	totalmem := 0
	totalpods := 0
	for k, v := range podDetails {
		totalcpu = totalcpu + v.cpu
		totalmem = totalmem + v.memory
		totalpods = totalpods + v.count
		t.AppendRows([]table.Row{
			{v.namespace, k, v.count, v.cpu, v.memory, v.nodename},
		})
		t.AppendSeparator()
	}

	t.AppendFooter(table.Row{
		"", "Total",
		totalpods,
		totalcpu,
		totalmem,
		"Memory/Pod=" + fmt.Sprint(math.Ceil(float64(totalmem)/float64(totalcpu))),
	})
	t.SetStyle(table.StyleLight)
	t.Render()
}
