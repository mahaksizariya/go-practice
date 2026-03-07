package practiceModule

import "fmt"

type animal interface {
	speak()
}
type dog struct{}

func (d dog) speak() {
	fmt.Println("dog barks")
}

type cat struct{}

func (c cat) speak() {
	fmt.Println("cat meows")
}

//basic struct

/*type Student struct {
	Name  string
	age   int
	marks int
}

/*func Struct() {

//basic struct
/*	var a1 student
	a1.Name = "mhk"
	a1.age = 12
	a1.marks = 345
	fmt.Println(a1)
*/

//anonymous struct basic

/*college := struct {
	Name  string
	age   int
	marks int
}{
	Name:  "mehak",
	age:   22,
	marks: 323,
}
fmt.Println(college)*/

//anonymous struct with slice

/*college := []struct {
		Name string
		Age  int
	}{
		{"MHK", 21},
		{"Anali", 22},
	}
	fmt.Println(college)

}
*/
func AddStudent() {

	var a animal
	a = dog{}
	a.speak()
	a = cat{}
	a.speak()

	// you are given an array of student data you have to append it to the data struct for Student
	// given student data is

	/*givenStudentData := []map[string]interface{}{
		{
			"Name":  "Ram",
			"Age":   23,
			"Marks": 23,
		},
		{
			"Name":  "Rohan",
			"Age":   13,
			"Marks": 35,
		},
		{
			"Name":  "Ravi",
			"Age":   21,
			"Marks": 33,
		},
		{
			"Name":  "Don",
			"Age":   43,
			"Marks": 53,
		},
	}
	/*
		[
			{
				Name: "Ram"
				Age: 23,
				Marks: 23,
			},
			{
				Name: "Rohan"
				Age: 13,
				Marks: 35,
			},
			{
				Name: "Ravi"
				Age: 21,
				Marks: 33,
			},
			{
				Name: "Don"
				Age: 43,
				Marks: 53,
			},

		]
	*/
	/*var students []Student

	for _, data := range givenStudentData {

		studentData := Student{
			Name:  data["Name"].(string),
			age:   data["Age"].(int),
			marks: data["Marks"].(int),
		}

		students = append(students, studentData)
	}

	fmt.Println(students)*/

	//map create basic
	/*students := map[string]  map[string]int{
	/*student["Math"] = 90
	student["english"] = 90

	/*fmt.Println(student)
	fmt.Println(student["Math"])*/
	/*	"Math":    90,
			"English": 80,
		}
		/*
			for key, value := range student {
				fmt.Println(value, key)*/
	/*delete(student, "Math")*/
	/*fmt.Println(len(student))*/

	/*"Gupta":{
		"math":90,
		"eng":80,
	},
	"Mhk":{
		"math":50,
		"eng":80,
	},
	for name, subjects := range students {
		fmt.Println("student:",name)
	}
	for subjects,marks:=range subjects{
		fmt.Println(subject,marks)
	}*/
	//BAISC POINTER EXAMPLE
	/* x := 10
	p := &x
	fmt.Println(x)
	fmt.Println(*p)
	fmt.Println(p)*/

}
