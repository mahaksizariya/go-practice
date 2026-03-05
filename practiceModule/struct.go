package practiceModule

import "fmt"

//basic struct

type Student struct {
	Name  string
	age   int
	marks int
}

func Struct() {

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

	college := []struct {
		Name string
		Age  int
	}{
		{"MHK", 21},
		{"Anali", 22},
	}
	fmt.Println(college)

}

func AddStudent() {
	// you are given an array of student data you have to append it to the data struct for Student
	// given student data is

	givenStudentData := []map[string]interface{}{
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
	var students []Student

	for _, data := range givenStudentData {

		studentData := Student{
			Name:  data["Name"].(string),
			age:   data["Age"].(int),
			marks: data["Marks"].(int),
		}

		students = append(students, studentData)
	}

	fmt.Println(students)
}
