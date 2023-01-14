package directories

import (
	"os"

	"create-go-app.com/code"
)

// Create function creates a new directory with the name passed as argument
func Create(dir string) error {
	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// Exists function checks if a directory already exists or not
func Exists(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// Change function changes the current working directory
func Change(dir string) error {
	if err := os.Chdir(dir); err != nil {
		return err
	}
	return nil
}

func GetWorkingDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

func CreateFile(dirname, flagName string) error {
	var fileContent []byte
	if flagName == "cli" {
		fileContent = []byte(code.CLI)
	} else if flagName == "http" {
		fileContent = []byte(code.HTTP)
	} else if flagName == "lib" {
		fileContent = []byte(code.LIB)
	}
	err := os.WriteFile(dirname+"/main.go", fileContent, os.FileMode(0777))
	if err != nil {
		return err
	}
	return nil
}
