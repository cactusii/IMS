package main

func main() {
	server := NewServer("0.0.0.0", 18888)
	server.Start()
}