package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Replace(basePath, brandingPath string) map[string]map[string]string {
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

func overwriteBaseWithBranding(base, branding []flatObj) map[string]map[string]string {
	basemap, brandmap := make(map[string]map[string]string), make(map[string]map[string]string)

	for _, obj := range base {
		m := make(map[string]string)
		for _, item := range obj.Items {
			m[item.Key] = item.Value
		}
		basemap[obj.Key] = m
	}

	for _, obj := range branding {
		m := make(map[string]string)
		for _, item := range obj.Items {
			m[item.Key] = item.Value
		}
		brandmap[obj.Key] = m
	}

	for k, items := range brandmap {
		if _, ok := basemap[k]; ok {
			for key, val := range items {
				basemap[k][key] = val
			}
		} else {
			basemap[k] = items
		}
	}

	return basemap
}

func flatten(key string, objects []jsobj) []flatObj {

	flatobjs := []flatObj{}
	flattenedSubObject := []flatObj{}

	for _, object := range objects {
		fobj := flatObj{
			Key:   "",
			Items: nil,
		}
		localKey := ""

		if len(key) == 0 {
			localKey = object.Key
		} else {
			localKey = object.Key
		}

		fobj.Key = localKey

		fobj.Items = append(fobj.Items, object.Items...)

		flatobjs = append(flatobjs, fobj)

		if len(object.SubObject) != 0 {
			flattenedSubObject = flatten(localKey, object.SubObject)
		}
		flatobjs = append(flatobjs, flattenedSubObject...)

	}

	m := make(map[string][]item)

	for _, obj := range flatobjs {
		m[obj.Key] = obj.Items
	}

	return flatobjs
}

func getObjects(m map[string]interface{}) []jsobj {
	objs := []jsobj{}

	for k, v := range m {

		switch val := v.(type) {
		case map[string]interface{}:
			objs = append(objs, mapToFlatObject(k, val))
		default:
			fmt.Println("Error, should be an object: (GetObjects func)", val)
		}
	}
	return objs
}

func mapToFlatObject(key string, m map[string]interface{}) jsobj {

	obj := jsobj{
		Key:       key,
		Items:     []item{},
		SubObject: nil,
	}

	for k, v := range m {
		switch val := v.(type) {
		case map[string]interface{}:
			obj.SubObject = append(obj.SubObject, mapToFlatObject(key+"."+k, val))
		case string:
			i := item{
				Key:   k,
				Value: val,
			}
			obj.Items = append(obj.Items, i)
		}
	}

	return obj
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
