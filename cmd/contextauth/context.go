package contextauth

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/smartaquarius10/kubecli/cmd"
)

func getCurrentContext() {
	stdout := cmd.ExecuteCommand("config", "current-context")
	fmt.Println("Current context is ", string(stdout))
}
func getContext() string {
	getCurrentContext()
	stdout := cmd.ExecuteCommand("config", "get-contexts", "-o", "name")
	scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
	counter := 1
	cmdmap := make(map[int]string)
	for scanner.Scan() {
		fmt.Println(counter, ")", scanner.Text())
		cmdmap[counter] = scanner.Text()
		counter++
	}
	var context_count int
	fmt.Print("Select:")
	fmt.Scanf("%d", &context_count)
	return cmdmap[context_count]
}
