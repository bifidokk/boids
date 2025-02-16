package main

type User struct {
}

func (u *User) Create() {}
func (u *User) Get()    {}
func (u *User) Delete() {}
func (u *User) List()   {}

type Reader interface {
	Get()
	List()
}

type Writer interface {
	Create()
	Delete()
}

func main() {
	var userReader Reader = &User{}
	userWriter := userReader.(Writer)
	//userWriter.Get() // can't call because made casting to Writer interface
	_ = userWriter
}
