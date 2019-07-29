package main

import "fmt"

func Hello(name, language string) string {
	if name == "" {
		name = "World"
	}

	return greetingPrefix(language) + name
}

const (
	spanish = "in Spanish"
	french  = "in French"
)

func greetingPrefix(language string) (prefix string) {
	switch language {
	case spanish:
		prefix = "Hola, "
	case french:
		prefix = "Bonjour, "
	default:
		prefix = "Hello, "
	}
	return
}

func HelloWorld() {
	fmt.Println(Hello("", ""))
}
