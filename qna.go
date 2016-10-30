package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"

	yaml "gopkg.in/yaml.v2"
)

type Pair struct {
	Q string
	A []string `yaml:",flow"`
}

type QNA []Pair

func main() {

	// Refactor: getFileArgs
	args := getFileArgs()

	// Transform file arguments to full paths
	// TODO: if starts with /, don't use relative path
	filePaths := []string{}
	for _, arg := range args {
		pwd, err := os.Getwd()
		if err != nil {
			promptErrorAndExit("")
		}
		filePaths = append(filePaths, pwd+"/"+arg)
	}

	// Parse yaml and combine into a single array
	qnas := QNA{}
	for _, file := range filePaths {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			promptErrorAndExit("Can't read file")
		}
		qna := QNA{}
		err = yaml.Unmarshal(data, &qna)
		if err != nil {
			promptErrorAndExit(fmt.Sprintf("%s", err))
		}
		qnas = append(qnas, qna...)
	}
	// fmt.Println(qnas)

	// Histogram array to track # of answered questions. Use index as keys
	//

	reader := bufio.NewReader(os.Stdin)
	num := 1
	for {
		shuffle(qnas)
		for _, qna := range qnas {
			color.Yellow(strconv.Itoa(num) + ". " + qna.Q)

			i := 0
			for i < len(qna.A) {
				answer := strings.TrimSpace(qna.A[i])
				text, _ := reader.ReadString('\n')
				if answer != strings.TrimSpace(text) {
					color.Red("[Try again] " + qna.Q)
					i = 0
					continue
				}
				i++
			}
			num++
		}
	}
}

func shuffle(qnas QNA) {
	for i := range qnas {
		j := rand.Intn(i + 1)
		qnas[i], qnas[j] = qnas[j], qnas[i]
	}
}

func promptErrorAndExit(msg string) {
	os.Stderr.WriteString(msg)
	os.Exit(1)
}

func getFileArgs() []string {
	fileArgs := os.Args[1:]
	if len(fileArgs) <= 0 {
		os.Stderr.WriteString("Please provide qna files as arguments")
		os.Exit(1)
	}
	return fileArgs
}

func getFilePaths(fileArgs []string) []string {
	pwd, err := os.Getwd()
	if err != nil {
		os.Stderr.WriteString("")
		os.Exit(1)
	}
	filePaths := []string{}
	for _, fileArg := range fileArgs {
		filePaths = append(filePaths, pwd+"/"+fileArg)
	}
	return filePaths
}
