package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli/v3"
)

// Function that determines whether the JSON content is an object or an array and processes it accordingly.
func jsonToMap(filePath string) ([]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode JSON without assuming its structure
	var rawData interface{}
	err = json.NewDecoder(file).Decode(&rawData)
	if err != nil {
		return nil, err
	}

	// Handle arrays as-is, wrap single objects for consistency
	switch v := rawData.(type) {
	case []interface{}:
		return v, nil // Return the array directly
	case map[string]interface{}:
		return []interface{}{v}, nil // Wrap the object in a slice
	default:
		return nil, fmt.Errorf("unknown JSON structure")
	}
}

func insertData(nestedMap map[string]interface{}, segments []string, data []interface{}) {
	for i, segment := range segments {
		if i == len(segments)-1 {
			// Check if data contains only one item and insert it directly, else insert the entire slice.
			if len(data) == 1 {
				nestedMap[segment] = data[0]
			} else {
				nestedMap[segment] = data
			}
		} else {
			if _, exists := nestedMap[segment]; !exists {
				nestedMap[segment] = make(map[string]interface{})
			}
			nestedMap = nestedMap[segment].(map[string]interface{})
		}
	}
}

func walkDir(root string, finalData map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".json") {
			jsonData, err := jsonToMap(path)
			if err != nil {
				fmt.Printf("Error: processing JSON file %s: %v\n", path, err)
				return nil
			}

			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			relPath = strings.TrimSuffix(relPath, ".json")

			pathSegments := strings.Split(relPath, string(os.PathSeparator))
			insertData(finalData, pathSegments, jsonData)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error: walking the path %q: %v\n", root, err)
	}
}

func parse(root, outputFilePath string) error {
	var wg sync.WaitGroup
	finalData := make(map[string]interface{})

	wg.Add(1)
	go walkDir(root, finalData, &wg)

	// Progress indicator for large directory processing, maybe this is a bad indicator
	go func() {
		for range time.Tick(10 * time.Second) {
			size := len(finalData)
			fmt.Printf("Current size of finalData: %d items\n", size)
		}
	}()
	wg.Wait()

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed creating output file: %s", err)
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)
	if err := encoder.Encode(finalData); err != nil {
		return fmt.Errorf("failed encoding final data to JSON: %s", err)
	}

	fmt.Printf("JSON data written to `%s`\n", outputFilePath)
	return nil
}

func valueOrDefault(c *cli.Command, flag *cli.StringFlag) string {
	if len(flag.DefaultText) == 0 {
		log.Fatalf("attempting to get valueOrDefault from flag with no DefaultText: `%s`", flag.Name)
	}

	value := c.String(flag.Name)
	if len(value) == 0 {
		return flag.DefaultText
	}

	return value
}

func main() {
	dirFlag := &cli.StringFlag{
		Name:        "dir",
		Aliases:     []string{"d"},
		DefaultText: "./",
		Usage:       "directory to search for JSON files",
	}
	outputFilePathFlag := &cli.StringFlag{
		Name:        "output-file-path",
		Aliases:     []string{"o"},
		DefaultText: "output.json",
		Usage:       "path and filename for the output JSON file",
	}

	cmd := &cli.Command{
		Description: "recursively scans the specified directory for JSON files, merging their contents into a single JSON file",
		Flags: []cli.Flag{
			dirFlag,
			outputFilePathFlag,
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			dir := valueOrDefault(c, dirFlag)
			outputFilePath := valueOrDefault(c, outputFilePathFlag)

			return parse(dir, outputFilePath)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
