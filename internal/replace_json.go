package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReplaceJson(basePath, brandingPath string) map[string]map[string]string {
	baseF, err := os.Open(basePath)
	if err != nil {
		fmt.Println("Error opening base project file: ", basePath, " error mesage: ", err.Error())
		os.Exit(1)
	}

	brandingF, err := os.Open(brandingPath)
	if err != nil {
		fmt.Println("Error opening branding project file: ", basePath, " error mesage: ", err.Error())
		os.Exit(1)
	}
	base := flatten("", getObjects(buildJsonMap(baseF)))
	brand := flatten("", getObjects(buildJsonMap(brandingF)))

	brandingF.Close()
	baseF.Close()

	return overwriteBaseWithBranding(base, brand)

}

func buildJsonMap(f *os.File) map[string]interface{} {
	m := make(map[string]interface{})

	scanner := bufio.NewScanner(f)
	scanner.Scan()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		split := strings.Split(line, ":")

		if len(split) > 1 {
			k := strings.TrimSpace(split[0])
			if strings.TrimSpace(split[1]) == "{" {
				m[k] = buildJsonObject(scanner)

			}
		}
	}
	return m
}

func buildJsonObject(scanner *bufio.Scanner) map[string]interface{} {
	endOfObjectFound := false
	object := make(map[string]interface{})

	for !endOfObjectFound {

		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())
		split := strings.Split(line, ":")

		if strings.Contains(line, "directory:op") || strings.Contains(line, "directory:website") {
			newS := split[0] + ":" + split[1]
			newSplit := []string{newS}
			slice := split[2:]
			newSplit = append(newSplit, slice...)
			split = newSplit
		}

		if len(split) < 2 {
			// End of object
			if split[0] == "}," || split[0] == "}" {
				endOfObjectFound = true
			}
		} else {
			k := strings.TrimSpace(split[0])
			v := strings.TrimSpace(strings.Join(split[1:], ":"))
			if v == "" {
				scanner.Scan()
				v = strings.TrimSpace(scanner.Text())
			}

			if v == "{" {
				object[k] = buildJsonObject(scanner)
			} else {
				if !strings.HasSuffix(v, ",") {
					v = v + ","
				}

				object[k] = v
			}
		}

	}
	return object
}
