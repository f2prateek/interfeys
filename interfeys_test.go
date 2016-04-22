package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestFixtures(t *testing.T) {
	dir, err := ioutil.TempDir("", "interfeys")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Create interfeys in temporary directory.
	interfeys := filepath.Join(dir, "interfeys.exe")
	err = run(".", "go", "build", "-o", interfeys, "interfeys.go")
	if err != nil {
		t.Fatalf("building interfeys: %s", err)
	}

	for _, fixture := range []string{"CoffeeMaker", "Server"} {
		lowercasedName := strings.ToLower(fixture)
		defer os.Remove(fmt.Sprintf("fixtures/%s/in/%s_interface.go", lowercasedName, lowercasedName))
		err = run(fmt.Sprintf("fixtures/%s/in", lowercasedName), interfeys, "-type", fixture)
		if err != nil {
			t.Fatal(err)
		}

		got, err := ioutil.ReadFile(fmt.Sprintf("fixtures/%s/in/%s_interface.go", lowercasedName, lowercasedName))
		if err != nil {
			t.Fatal(err)
		}

		expected, err := ioutil.ReadFile(fmt.Sprintf("fixtures/%s/out/%s_interface.go", lowercasedName, lowercasedName))
		if err != nil {
			t.Fatal(err)
		}

		if string(expected) != string(got) {
			t.Errorf("Fixture: %s\nExpected:\n%s\nGot:\n%s", fixture, expected, got)
		}
	}
}

// run runs a single command and returns an error if it does not succeed.
func run(dir string, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
