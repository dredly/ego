package types

const startingELO = 1000

type Person struct {
	Name string
	ELO float64
}

func NewPerson(name string) *Person {
	return &Person{
		Name: name,
		ELO: startingELO, 
	}
}

func (p *Person) UpdateELO(change float64) {	
	p.ELO += change
}