package replacelayer

import "fmt"

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

		if len(object.SubObjects) != 0 {
			flattenedSubObject = flatten(localKey, object.SubObjects)
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
	a := []jsobj{}
	for k, v := range m {

		switch val := v.(type) {
		case map[string]interface{}:
			a = append(a, mapToFlatObject(k, val))
		default:
			fmt.Println("Error, should be an object: (GetObjects func)", val)
		}
	}
	return a
}

func mapToFlatObject(key string, m map[string]interface{}) jsobj {

	obj := jsobj{
		Key:        key,
		Items:      []item{},
		SubObjects: nil,
	}

	for k, v := range m {
		switch val := v.(type) {
		case map[string]interface{}:
			obj.SubObjects = append(obj.SubObjects, mapToFlatObject(key+"."+k, val))
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
