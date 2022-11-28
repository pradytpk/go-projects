package main

type apiError struct {
	Err    string
	Status int
}

type User struct {
	ID    int
	Valid bool
}
