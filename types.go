package main

type Clinic struct{
	Name 	string `json:"name"`
	State	string `json:"state"`
	Availability Availability `json:"availability"`
}

type Availability struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Clinics struct{
	Clinics [] Clinic `json:"clinics"`
}