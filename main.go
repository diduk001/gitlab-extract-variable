package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func printExample() {
	fmt.Println("Examples:")
	fmt.Println("gitlab-extract-variable -token=TOKEN -project=ProjectOwner/ProjectName -output=output.txt")
	fmt.Println("gitlab-extract-variable -token=TOKEN -project=ProjectOwner/ProjectName -compact")
}

type ApiResponse []struct {
	VariableType     string `json:"variable_type"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	Protected        bool   `json:"protected"`
	Masked           bool   `json:"masked"`
	Raw              bool   `json:"raw"`
	EnvironmentScope string `json:"environment_scope"`
}

func getApiResponse(token string, projectName string) ApiResponse {
	projectNameEscaped := url.QueryEscape(projectName)
	apiUrl := fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/variables", projectNameEscaped)

	client := http.Client{}
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		fmt.Println("Error while creating request")
		panic(err)
	}
	req.Header.Add("PRIVATE-TOKEN", token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while performing request")
		panic(err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Status code is %d\n", resp.StatusCode)
		panic(resp.Body)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error while reading response")
		panic(err)
	}

	var result ApiResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Println("Error while deserializing JSON")
		panic(err)
	}
	return result
}

func writeCSV(response ApiResponse, filename string, compactFlag bool) {
	file, err := os.Create(filename)
	defer file.Close()

	if err != nil {
		fmt.Println("Error while creating output file")
		panic(err)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var header []string
	if compactFlag {
		header = []string{"key", "value"}
	} else {
		header = []string{"variable_type", "key", "value", "protected", "masked", "raw", "environment_scope"}
	}

	err = writer.Write(header)
	if err != nil {
		fmt.Println("Error while writing csv header")
		panic(err)
	}
	for i, variable := range response {
		var row []string

		if compactFlag {
			row = []string{variable.Key, variable.Value}
		} else {
			row = []string{variable.VariableType, variable.Key, variable.Value,
				strconv.FormatBool(variable.Protected), strconv.FormatBool(variable.Masked),
				strconv.FormatBool(variable.Raw), variable.EnvironmentScope}
		}
		if err = writer.Write(row); err != nil {
			fmt.Printf("Error while writing %d-th row\n", i)
			panic(err)
		}
	}
	writer.Flush()
}

func main() {
	privateTokenPtr := flag.String("token", "", "GitLab user's private token")
	projectNamePtr := flag.String("project", "", "GitLab project name")
	outputFilePtr := flag.String("output", "output.csv", "Output file")
	compactPtr := flag.Bool("compact", false, "Compact output (only key and value)")

	flag.Parse()

	if *privateTokenPtr == "" {
		fmt.Println("Private token is not specified!")
		printExample()
		return
	}
	if *projectNamePtr == "" {
		fmt.Println("Project name is not specified!")
		printExample()
		return
	}

	response := getApiResponse(*privateTokenPtr, *projectNamePtr)
	writeCSV(response, *outputFilePtr, *compactPtr)
}
