package main

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	judgemap = judgemapStore()

	ss()
}
