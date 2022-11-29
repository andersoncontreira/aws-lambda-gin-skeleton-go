package config

import (
	"os"
	"testing"
)

func TestLoadWithNoArgsLoadsDotEnv(t *testing.T) {
	err := LoadVariables()
	pathError := err.(*os.PathError)
	if pathError == nil || pathError.Op != "open" || pathError.Path != ".env" {
		t.Errorf("Didn't try and open .env by default")
	}
}

func TestOverloadWithPathArgOverloadsDotEnv(t *testing.T) {
	err := LoadVariables()
	pathError := err.(*os.PathError)
	if pathError == nil || pathError.Op != "open" || pathError.Path != ".env" {
		t.Errorf("Didn't try and open .env by default")
	}
}

func TestLoadFileNotFound(t *testing.T) {
	err := LoadVariables("somefilethatwillneverexistever.env")
	if err == nil {
		t.Error("File wasn't found but Load didn't return an error")
	}
}

//func TestOverloadFileNotFound(t *testing.T) {
//	err := Overload("somefilethatwillneverexistever.env")
//	if err == nil {
//		t.Error("File wasn't found but Overload didn't return an error")
//	}
//}
