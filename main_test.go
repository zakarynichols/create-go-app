package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/zakarynichols/create-go-app/dir"
)

const testSrc = "tests_src"
const testDst = "tests_dst"

func TestCopyDir(t *testing.T) {
	var expectedPaths = []string{
		testDst,
		testDst + "/nested",
		testDst + "/nested/deep",
		testDst + "/nested/deep/deep_file.txt",
		testDst + "/nested/deep/deeper",
		testDst + "/nested/deep/deeper/deeper_file.txt",
	}

	if runtime.GOOS == "windows" {
		expectedPaths = []string{
			testDst,
			testDst + "\\nested",
			testDst + "\\nested\\deep",
			testDst + "\\nested\\deep\\deep_file.txt",
			testDst + "\\nested\\deep\\deeper",
			testDst + "\\nested\\deep\\deeper\\deeper_file.txt",
		}
	}

	err := os.Mkdir(testDst, 0750)
	if err != nil {
		log.Fatal(err)
	}

	err = dir.CpdirAll(
		testSrc, testDst,
	)

	if err != nil {
		t.Errorf("copying error, %s", err)
	}

	i := 0 // Increment a counter each walk.

	// Walk the paths. They should match the expected array in order.
	err = filepath.WalkDir(testDst, func(path string, d fs.DirEntry, err error) error {
		if path != expectedPaths[i] {
			t.Fatalf("%s does not equal %s", path, expectedPaths[i])
		}
		i++
		return nil
	})
	if err != nil {
		t.Error("walk dir err:", err.Error())
	}

	// Cleanup directories created by tests
	err = os.RemoveAll(testDst)
	if err != nil {
		t.Fatalf("failed to remove all: %s", err.Error())
	}
}
