package internal

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func makeJsObjects(m map[string]map[string]string) []jsobj {
	n := make(map[string][]flatObj)

	for k, v := range m {
		root := strings.Split(k, ".")[0]
		items := []item{}
		for key, value := range v {
			i := item{
				Key:   key,
				Value: value,
			}
			items = append(items, i)
		}
		obj := flatObj{
			Key:   k,
			Items: items,
		}
		n[root] = append(n[root], obj)
	}

	objects := []jsobj{}
	for _, o := range n {
		objects = append(objects, flatObjsToJsObj(o))
	}
	return objects
}

func writeObjectToFile(o jsobj, f *os.File, indent string) {

	key := indent + o.Key + ": {"
	wtof(key, f)

	for _, item := range o.Items {
		line := indent + "\t" + item.Key + ": " + item.Value
		wtof(line, f)
	}
	for _, item := range o.SubObject {
		writeObjectToFile(item, f, indent+"\t")
	}

	endOfObject := indent + "},"
	wtof(endOfObject, f)
}

func flatObjsToJsObj(fa []flatObj) jsobj {

	ob := jsobj{
		Key:       strings.Split(fa[0].Key, ".")[0],
		Items:     []item{},
		SubObject: nil,
	}

	nonRootFa := []flatObj{}
	for _, obj := range fa {
		split := strings.Split(obj.Key, ".")
		if len(split) == 1 {
			ob.Items = obj.Items
		} else {
			obj.Key = strings.Join(split[1:], ".")
			nonRootFa = append(nonRootFa, obj)
		}
	}

	m := make(map[string][]flatObj)
	for _, obj := range nonRootFa {
		split := strings.Split(obj.Key, ".")
		key := split[0]
		m[key] = append(m[key], obj)
	}

	for _, v := range m {
		ob.SubObject = append(ob.SubObject, flatObjsToJsObj(v))
	}

	return ob
}

func wtof(line string, f *os.File) {
	_, err := f.WriteString(line + "\n")
	if err != nil {
		fmt.Println("error writing to file: ", err.Error())
	}
}

func CopyFile(src, dest string) {
	f, err := os.Open(src)

	if err != nil {
		fmt.Println("error opening file: ", src)
		fmt.Println(err.Error())
	}
	defer f.Close()

	d, err := os.Create(dest)

	if err != nil {
		fmt.Println("error creating a copied file: ", dest)
		fmt.Println(err.Error())
	}
	defer d.Close()

	_, err = io.Copy(d, f)
	if err != nil {
		fmt.Println("error copying to file: ", dest)
		fmt.Println(err.Error())
	}

}

func WriteToFile(m map[string]map[string]string, outputDir, filename string) {

	objects := makeJsObjects(m)
	f, err := os.Create(path.Join(outputDir, filename))
	if err != nil {
		fmt.Println("Error creating file: ", err.Error())
	}

	wtof("export default {", f)

	for _, o := range objects {
		writeObjectToFile(o, f, "\t")
	}

	wtof("}", f)
}
