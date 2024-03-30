package merkledag

import (
	"encoding/json"
	"strings"
)

const (
	// TypeTree 表示目录类型
	TypeTree = "tree"
	// TypeBlob 表示文件类型
	TypeBlob = "blob"
	// TypeList 表示列表类型
	TypeList = "list"
	// Step 表示每个对象的步长
	Step = 4
)

func Hash2File(store KVStore, hash []byte, path string) []byte {
	if !hasHash(store, hash) {
		return nil
	}

	objBinary, _ := store.Get(hash)
	obj := binaryToObj(objBinary)
	pathArr := strings.Split(path, "/")
	return getFileByPath(obj, pathArr, store)
}

func getFileByPath(obj *Object, pathArr []string, store KVStore) []byte {
	for _, part := range pathArr {
		if obj == nil {
			return nil
		}
		obj = findLinkByName(obj, part, store)
	}
	if obj == nil {
		return nil
	}
	return getObjectData(obj, store)
}

func hasHash(store KVStore, hash []byte) bool {
	flag, _ := store.Has(hash)
	return flag
}

func findLinkByName(obj *Object, name string, store KVStore) *Object {
	for _, link := range obj.Links {
		if link.Name == name {
			objBinary, _ := store.Get(link.Hash)
			return binaryToObj(objBinary)
		}
	}
	return nil
}

func getObjectData(obj *Object, store KVStore) []byte {
	if obj.Type == TypeBlob {
		data, _ := store.Get(obj.Hash)
		return data
	}
	if obj.Type == TypeList {
		var result []byte
		for _, link := range obj.Links {
			linkObj := findLinkByHash(link.Hash, store)
			result = append(result, getObjectData(linkObj, store)...)
		}
		return result
	}
	return nil
}

func findLinkByHash(hash []byte, store KVStore) *Object {
	objBinary, _ := store.Get(hash)
	return binaryToObj(objBinary)
}

func binaryToObj(objBinary []byte) *Object {
	var obj Object
	json.Unmarshal(objBinary, &obj)
	return &obj
}
