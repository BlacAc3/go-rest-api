package config 

import (
    "fmt"
    "os"
    "path/filepath"
)

func GetProjectDir() (string,error){
    workingDir, err := os.Getwd()
    if err != nil {
        return "", fmt.Errorf("Error: %v",err)
    }
    projectDir := filepath.Dir(workingDir)
    return projectDir, nil
}
