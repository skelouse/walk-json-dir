package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestCase defines the structure for each test case
type TestCase struct {
	Name           string
	FilesContents  map[string]string // Path to JSON content as string
	ExpectedOutput map[string]interface{}
}

// runParseTest runs a single test case for the parse function
func runParseTest(t *testing.T, tc TestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "runParseTest-*")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %s", err)
		}
		defer os.RemoveAll(tmpDir)

		// Create JSON files based on the test case
		for path, content := range tc.FilesContents {
			fullPath := filepath.Join(tmpDir, path)
			if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
				t.Fatalf("Failed to MkdirAll %s: %s", fullPath, err)
			}

			if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to write JSON to file %s: %s", fullPath, err)
			}
		}

		// Path for the output file
		outputFilePath := filepath.Join(tmpDir, "output.json")

		// Execute the parse function
		if err := parse(tmpDir, outputFilePath); err != nil {
			t.Fatalf("Parse function failed: %s", err)
		}

		// Read and check the output file
		outputFile, err := os.Open(outputFilePath)
		if err != nil {
			t.Fatalf("Failed to open output file: %s", err)
		}
		defer outputFile.Close()

		var outputData map[string]interface{}
		if err := json.NewDecoder(outputFile).Decode(&outputData); err != nil {
			t.Fatalf("Failed to decode output JSON: %s", err)
		}

		// Verify the outputData against the expected output
		if !reflect.DeepEqual(outputData, tc.ExpectedOutput) {
			t.Errorf("Output data does not match expected output.\nExpected: %#v\nGot: %#v", tc.ExpectedOutput, outputData)
		}
	})
}

func TestParseWithDifferentStructures(t *testing.T) {
	testCases := []TestCase{
		{
			Name: "SingleObject",
			FilesContents: map[string]string{
				"test/object.json": `{"key": "value"}`,
			},
			ExpectedOutput: map[string]interface{}{
				"test": map[string]interface{}{
					"object": map[string]interface{}{"key": "value"},
				},
			},
		},
		{
			Name: "SingleList",
			FilesContents: map[string]string{
				"listdir/list.json": `["item1", "item2", "item3"]`,
			},
			ExpectedOutput: map[string]interface{}{
				"listdir": map[string]interface{}{
					"list": []interface{}{"item1", "item2", "item3"},
				},
			},
		},
		{
			Name: "NestedObjects",
			FilesContents: map[string]string{
				"nested/object.json": `{"level1": {"level2": {"key": "value"}}}`,
			},
			ExpectedOutput: map[string]interface{}{
				"nested": map[string]interface{}{
					"object": map[string]interface{}{
						"level1": map[string]interface{}{
							"level2": map[string]interface{}{
								"key": "value",
							},
						},
					},
				},
			},
		},
		{
			Name: "ListAndObject",
			FilesContents: map[string]string{
				"mixed/mixed.json": `[{"objKey": "objValue"}, {"objKey2": "objValue2"}]`,
			},
			ExpectedOutput: map[string]interface{}{
				"mixed": map[string]interface{}{
					"mixed": []interface{}{
						map[string]interface{}{"objKey": "objValue"},
						map[string]interface{}{"objKey2": "objValue2"},
					},
				},
			},
		},
		{
			Name: "MultipleLevels",
			FilesContents: map[string]string{
				"level1/level2/level3/object.json": `{"key": "value"}`,
			},
			ExpectedOutput: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": map[string]interface{}{
							"object": map[string]interface{}{"key": "value"},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		runParseTest(t, tc)
	}
}
