package replacelayer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReplaceJs(basePath, brandingPath string) map[string]map[string]string {
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

	base := flatten("", getObjects(buildMap(baseF)))
	brand := flatten("", getObjects(buildMap(brandingF)))

	brandingF.Close()
	baseF.Close()

	return overwriteBaseWithBranding(base, brand)

}

func buildObject(scanner *bufio.Scanner) map[string]interface{} {
	endOfObjectFound := false
	object := make(map[string]interface{})

	for !endOfObjectFound {

		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())
		split := strings.Split(line, ":")

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
				object[k] = buildObject(scanner)
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

func buildMap(f *os.File) map[string]interface{} {
	m := make(map[string]interface{})

	scanner := bufio.NewScanner(f)
	scanner.Scan()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		split := strings.Split(line, ":")

		if len(split) > 1 {
			k := strings.TrimSpace(split[0])
			if strings.TrimSpace(split[1]) == "{" {
				m[k] = buildObject(scanner)
			}
		}
	}
	return m
}
