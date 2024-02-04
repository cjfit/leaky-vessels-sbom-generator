/*
 * SBOM Generator
 * by cjfit
 * 
 * Forked from: https://github.com/snyk/leaky-vessels-static-detector/
 * © 2024 Snyk Limited
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

 package main

 import (
	 "encoding/csv"
	 "flag"
	 "fmt"
	 "os"
	 "static-detector/internal/cmd"
	 "static-detector/internal/common"

 )
 
 var (

	logger = common.GetLogger()

	disableListPlaceholder []string

 )
/*  // Placeholder for the actual analysis logic
 func AnalyzeDockerfile(dockerfilePath string, analyzeBase bool) (map[string]string, error) {
	 // Placeholder response
	 return map[string]string{"Path": dockerfilePath, "Analysis": "Success", "BaseImageAnalyzed": fmt.Sprintf("%t", analyzeBase)}, nil
 } */

 type DockerfileAnalysisResult struct {
    DockerfilePath string
    BaseImages     map[string]string // Maps base images to their versions
}

 func main() {
	
	dockerfileCmd := flag.NewFlagSet("dockerfile", flag.ExitOnError)
    csvFilePath := dockerfileCmd.String("csv", "", "Path to the CSV file containing Dockerfile paths")
    analyzeBase := dockerfileCmd.Bool("base", false, "Whether to analyze base image")

    if len(os.Args) < 2 {
        fmt.Println("dockerfile subcommand is required")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "dockerfile":
        dockerfileCmd.Parse(os.Args[2:])
    default:
        fmt.Println("Only 'dockerfile' mode is supported at this time")
        os.Exit(1)
    }

    if *csvFilePath == "" {
        fmt.Println("CSV file path is required")
        os.Exit(1)
    }

    file, err := os.Open(*csvFilePath)
    if err != nil {
        fmt.Printf("Failed to open CSV file: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    dockerfiles, err := reader.ReadAll()
    if err != nil {
        fmt.Printf("Failed to read CSV file: %v\n", err)
        os.Exit(1)
    }

    // Initialize slice to store results for all Dockerfiles
    var results [][]string
    results = append(results, []string{"DockerfilePath", "BaseImage", "Version"}) // Header

    // Analyze Dockerfiles
    for _, record := range dockerfiles {
        path := record[0]
        if _, err := os.Stat(path); os.IsNotExist(err) {
            fmt.Printf("File does not exist: %s\n", path)
            continue
        }

        opts := common.Options{
            DisableRules: &disableListPlaceholder, // Assuming this is defined somewhere
            AnalyzeBaseImage: *analyzeBase,
        }

        code, analysisResult, err := cmd.AnalyzeDockerfile(path, opts)
        if err != nil || code != 0 {
            fmt.Printf("Error analyzing Dockerfile '%s': %v\n", path, err)
            // Optionally, log the error into the CSV or handle it as needed
            continue
        }

        // Append the analysis result for each Dockerfile to results slice
        for baseImage, version := range analysisResult.BaseImages {
            results = append(results, []string{analysisResult.DockerfilePath, baseImage, version})
        }
    }

    // Now, write the collected results to the CSV file
    outputFile, err := os.Create("analysis_results.csv")
    if err != nil {
        fmt.Printf("Failed to create output CSV file: %v\n", err)
        return
    }
    defer outputFile.Close()

    writer := csv.NewWriter(outputFile)
    if err := writer.WriteAll(results); err != nil {
        fmt.Printf("Failed to write to output CSV file: %v\n", err)
    } else {
        fmt.Println("Analysis completed successfully. Results saved to analysis_results.csv.")
    }
}