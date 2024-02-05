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

type DockerfileAnalysisResult struct {
    DockerfilePath string
    BaseImages     map[string]string // Maps base images to their versions
}

func main() {
    dockerfileCmd := flag.NewFlagSet("dockerfile", flag.ExitOnError)
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

    file, err := os.Open("./gh_crawler/src/gh_enumerator/gh_data/dockerfiles_paths.csv")
    if err != nil {
        fmt.Printf("Failed to open CSV file: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    // Skip the first header line by reading it and not using it
    if _, err := reader.Read(); err != nil {
        fmt.Printf("Failed to read the header line from CSV file: %v\n", err)
        os.Exit(1)
    }

    // Read the rest of the data after the header
    dockerfiles, err := reader.ReadAll()
    if err != nil {
        fmt.Printf("Failed to read CSV file: %v\n", err)
        os.Exit(1)
    }

    // Initialize slice to store results for all Dockerfiles
    var results [][]string
    // Optionally, if you want to include a header in the output CSV, uncomment the next line
    // results = append(results, []string{"DockerfilePath", "BaseImage", "Version"}) // Header for output CSV

    // Analyze Dockerfiles
    for _, record := range dockerfiles {
		subdir_path := "./gh_crawler/src/gh_enumerator/gh_data"
		simple_path := record[0] +"/"+ record[1]
        path := subdir_path + "/" + record[0] +"/"+ record[1]
        if _, err := os.Stat(path); os.IsNotExist(err) {
            fmt.Printf("File does not exist: %s\n", path)
            continue
        }

        opts := common.Options{
            DisableRules:    &disableListPlaceholder, // Assuming this is defined somewhere
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
            results = append(results, []string{simple_path, baseImage, version})
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
