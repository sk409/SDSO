package main

import "fmt"

func route(path string) string {
	return fmt.Sprintf("%s/%s", serverOrigin, path)
}
