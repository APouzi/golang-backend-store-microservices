package database

type Models struct {
	UserCustomer UserCustomer
	UserAdmin    UserAdmin
	Customer     Customer
	Admin        Admin
}

type UserCustomer struct {
	Email      string
	First_Name string
	Last_Name  string
	Customer   *Customer
}

type UserAdmin struct {
	Email      string
	First_Name string
	Last_Name  string
	Admin      *Admin
}

type Customer struct {
	Street_Address string
	Phone_Number   string
	State          string
}

type Profile struct {
}

type Admin struct {
	Privlages []string
	SuperUser bool
}

// --------- Product ---------
