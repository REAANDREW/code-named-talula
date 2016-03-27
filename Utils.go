package main

import "fmt"

//AdminURL ...
func AdminURL(path string, args ...interface{}) string {
	return fmt.Sprintf("http://%s:%d%s", Host, AdminPort, fmt.Sprintf(path, args...))
}

//TransformURL ...
func TransformURL(path string, args ...interface{}) string {
	return fmt.Sprintf("http://%s:%d%s", Host, ProxyPort, fmt.Sprintf(path, args...))
}
