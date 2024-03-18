package GRPC_test

import "fmt"

type PersonAction interface {
	SayHello(anybody string)
	SayGoodBye(anybody string)
}

type Person struct {
	name string
	gender string
	age int
}

func NewPerson(name, gender string, age int) *Person {
	return &Person{
		name: name,
		gender: gender,
		age: age,
	}
}

func (p *Person) SayHello(anybody string) {
	fmt.Printf("%s say \"Hello\" to %s", p.name, anybody)
}

func (p *Person) SayGoodBye(anybody string) {
	fmt.Printf("%s say \"Good bye\" to %s", p.name, anybody)
}
