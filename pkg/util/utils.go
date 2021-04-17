package util

import "reflect"

// items := []string{"A", "1", "B", "2", "C", "3"}
// // Missing Example
//  _, found := Find(items, "golangcode.com")
//  if !found {
//    fmt.Println("Value not found in slice")
//  }
//  // Found example
//  k, found := Find(items, "B")
//  if !found {
//    fmt.Println("Value not found in slice")
//  }
//  fmt.Printf("B found at key: %d\n", k)

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (string, bool) {
	for _, item := range slice {
		if item == val {
			return item, true
		}
	}
	return "", false
}

// ItemExists - check if array element exists
// strArray := [5]string{"India", "Canada", "Japan", "Germany", "Italy"}
// fmt.Println(itemExists(strArray, "Canada"))
func ItemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func GetNodeRole(val string) string {
	serverRoles := []string{"master", "server"}
	// myRoles := "master"
	// myRoles = "server"
	// myRoles = "single"
	// myRoles = "worker"
	// myRoles = "noname"

	if val == "single" {
		return "single"
	}
	_, found := Find(serverRoles, val)
	if !found {
		return "agent"
	}
	return "server"
}
