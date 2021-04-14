package util

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
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
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
