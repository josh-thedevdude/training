package main

import "os"

func main() {
	a := App{}
	a.Initialize(os.Getenv("DB_URL"))
	a.Run(":8000")
}
