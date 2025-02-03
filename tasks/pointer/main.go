package main

import "fmt"

type User struct {
	Name string
}

func main() {
	user := User{Name: "John"}

	fmt.Println("user:", user.Name) // John

	updateUser(user)

	fmt.Println("user after update:", user.Name) // John

}

// send by copy
func updateUser(u User) {
	u.Name = "Bob"

	fmt.Println("in updateUser ", u.Name) // Bob

	resetUser(&u)

	fmt.Println("after resetUser ", u.Name) // Bob
}

func resetUser(u *User) {
	u = &User{Name: "Noname"} // rewrite address, so now it points to another memory address and doesn't affect scope of updateUser method

	fmt.Println("in resetUser ", u.Name) // Noname
}
