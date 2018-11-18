package main
type person1 struct {

	 name string
	 age int
	 account float64
}

func Newperson(name string)*person1  {
	return &person1{
		name :name,
	}
}
func (p *person1)PersonGetAge(age int)  {
	if age>20 {
		p.age=age
	}
}
func main(){
	
}
