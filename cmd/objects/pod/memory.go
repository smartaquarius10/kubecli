package pod

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/smartaquarius10/kubecli/cmd/objects"
	"github.com/spf13/cobra"
)

type PodDetails struct {
	count  int
	memory int
	cpu    int
}

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "This command return pods memory usage per namespace",
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		getPodsMemoryInNamespace(namespace)
	},
}

func TopCmd() *cobra.Command {
	topCmd.Flags().StringP("namespace", "n", "", "Pass namespace having pods")
	return topCmd
}

func getPodsMemoryInNamespace(namespace string) {
	stdout := objects.GetPodMemory(namespace, "pods")
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	m := make(map[string]*PodDetails)
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	for scanner.Scan() {
		if !strings.Contains(scanner.Text(), "NAME") {
			podDetails := strings.Fields(scanner.Text())
			val := strings.Split(podDetails[0], "-")[0]
			_, ok := m[val]
			cpu, _ := strconv.Atoi(re.FindString(podDetails[1]))
			memory, _ := strconv.Atoi(re.FindString(podDetails[2]))
			if ok {
				m[val].count++

				m[val].cpu = m[val].cpu + cpu
				m[val].memory = m[val].memory + memory
			} else {
				m[val] = &PodDetails{count: 1, memory: memory, cpu: cpu}
			}
		}
	}
	writeTable(m)
}

func writeTable(podDetails map[string]*PodDetails) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Pods", "ReplicaCount", "CPU(m)", "Memory(Mi)"})
	totalcpu := 0
	totalmem := 0
	for k, v := range podDetails {
		totalcpu = totalcpu + v.cpu
		totalmem = totalmem + v.memory
		t.AppendRows([]table.Row{
			{k, v.count, v.cpu, v.memory},
		})
		t.AppendSeparator()
	}

	t.AppendFooter(table.Row{"Total", totalcpu, totalmem})
	t.SetStyle(table.StyleLight)
	t.Render()
}
