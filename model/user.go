package model
type User struct {
	ID   int
	Name string
	Role string
}

type Users []User

func (u Users) Exists(id int) bool {
	//...
	return  false
}

func (u Users) FindByName(name string) (User, error) {
	//...
	return  User{10,"u","是"},nil
}
