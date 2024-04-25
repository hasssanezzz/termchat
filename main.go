package main

func main() {
	server := NewServer(":8080")
	server.start()
}
