package main

// return error without including packages

type CustomError struct {
	message string
}

func (ce *CustomError) Error() string {
	return ce.message
}

func main() {
	println(handle().Error())
}

func handle() error {
	return &CustomError{"An error occurred"}
}
