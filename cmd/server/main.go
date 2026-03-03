/*for i := 0; i <= 10; i++ {
	fmt.Print(i)
}
*/

/*i := 0
for i <= 50 {
	fmt.Println(i)
	i++
}*/

/*for {
	fmt.Println("Hello")
}
*/

/*for i := 1; i <= 10; i++ {
	if i == 8 {
		break
	}
	fmt.Println(i)
}*/

/*nums := []int{10, 23, 34}
for index, value := range nums {
	fmt.Println("index:", index, "value:", value)
}
*/
/*nums := []int{10, 20, 30}
for _, value := range nums {
	fmt.Println("value:", value)
}
*/

/*nums := []int{10, 20, 30}
for i := range nums {
	nums[i] = nums[i]*10 + 2
}
fmt.Println(nums)
*/

/*

	 nums := [4]int{10, 20, 30, 40}
	for i := range nums {
		nums[i] = nums[i]*10 + 2
	}
	fmt.Println(nums)
*/

/*nums := [4]int{10, 20, 30, 40}
for i := 0; i <= 3; i++ {
	fmt.Println(nums[i])
}*/

/*nums := []string{"a", "n", "j", "l", "i"}
for i := 0; i <= 4; i++ {
	if nums[i] == "j" {
		fmt.Println(nums[i])
	}
}*/

/*for i := 1; i <= 20; i++ {
	if i%2 == 0 {
		fmt.Println(i)
	}
}*/

/*for i := 1; i <= 20; i++ {
	if i%2 != 0 {
		fmt.Println(i)
	}
}*/

/*nums := 7
for i := 1; i <= 10; i++ {
	fmt.Println(nums, "x", i, "=", nums*i)
}*/

/*arr := []int{10, 20, 30, 40}

for i := 0; i < len(arr); i++ {
	fmt.Print(arr[i], " ")
}*/

/*arr := []int{10, 20, 30, 40}
for i := range arr {
	fmt.Print(arr[i], " ")
}
*/

/*arr := []int{10, 20, 30, 40}
sum := 0

for _, v := range arr {
	sum += v
}

fmt.Println("Sum:", sum)
*/

/*arr := []int{10, 20, 30, 40}
sum := 0
for i := 0; i < len(arr); i++ {
	sum += arr[i]
}
fmt.Println("Sum:", sum)
*/

/*str := "Hello"

for i := 0; i < len(str); i++ {
	fmt.Println(string(str[i]))
}*/

/*str := "Hello"
for _, v := range str {
	fmt.Println(string(v))
}*/
/*str := "Hello"
result := ""
for i := len(str) - 1; i >= 0; i-- {
	result += string(str[i])
}
fmt.Println(result)
*/
/*str := "Mangooo"
count := 0

for _, ch := range str {
	if ch == 'a' || ch == 'e' || ch == 'i' || ch == 'o' || ch == 'u' {
		count++
	}
}

fmt.Println("Vowels:", count)
*/

/*
	func main() {
		arr := [][]int{{1, 2}, {3, 4}, {5, 6}}
		var result []int
		for i := 0; i < len(arr); i++ {
			for j := 0; j < len(arr[i]); j++ {
				result = append(result, arr[i][j])
			}
		}
		fmt.Println(result)

		arr1 := []int{1, 2, 3, 4, 5, 6}
		fmt.Println(len(arr1), " got my array with length")
		// for i := 0; i < len(arr1); i++ {
		// 	arr1[i] = arr1[i] * 2
		// }

		// fmt.Println(arr1)

		// var arr2 []int
		// num := 2
		// for i := 1; i <= 10; i++ {
		// 	arr2 = append(arr2, num*i)
		// }
		// fmt.Println(arr2, " got my array with multiple of 2")
	}

// functions in Golang
*/

/*func greet() {
	fmt.Println("Hello Mhk")
}
func main() {
	greet()
}*/

/*func greet(name string) {
	fmt.Println("Hello", name)
}
func main() {
	greet("Mhk Gupta")
}
*/
/*func add(a int, b int) int {
	return a - b
}
func main() {
	result := add(88, 78)
	fmt.Println("The sub is:", result)
}*/
package main

import "fmt"

/*func divide(a int, b int) (int, int) {
	quotient := a / b
	remainder := a % b
	return quotient, remainder
}
func main() {
	quotient, remainder := divide(10, 3)
	fmt.Println("Quotient:", quotient)
	fmt.Println("Remainder:", remainder)
}
*/
func multiplyby2(arr []int) []int {
	var multiple []int
	for i := 0; i < len(arr); i++ {
		multiple = append(multiple, arr[i]*2)
	}
	return multiple
}
func emptyarr(arr []int) []int {
	var empty []int
	num := 2
	for i := 0; i < len(arr); i++ {
		empty = append(empty, num*arr[i])
	}
	fmt.Println("hello world")
	return empty
}
func multiplearr(arr2 [][]int) []int {
	var multiplearr []int
	for i := 0; i < len(arr2); i++ {
		for j := 0; j < len(arr2[i]); j++ {
			multiplearr = append(multiplearr, arr2[i][j])
		}
	}
	return multiplearr
}
func sum(a int, b int) int {
	fmt.Println("hello")
	return a + b
}
func divi2(arr3 []int) []int {
	var divi2 []int
	for i := 0; i < len(arr3); i++ {
		if arr3[i]%2 == 0 && arr3[i] != 8 {
			divi2 = append(divi2, arr3[i])
		}
	}
	return divi2
}
func reverse(arr4 []int) []int {
	i := 0
	j := len(arr4) - 1
	for i < j {
		arr4[i], arr4[j] = arr4[j], arr4[i]
		i++
		j--
	}
	return arr4
}

func main() {
	arr := []int{1, 2, 3, 4, 5}
	fmt.Println(multiplyby2(arr))
	fmt.Println(emptyarr(arr))
	arr2 := [][]int{{1, 2}, {3, 4}, {5, 6}}
	fmt.Println(multiplearr(arr2))
	fmt.Println(sum(10, 20))
	arr3 := []int{2, 7, 8, 9, 4}
	fmt.Println(divi2(arr3))
	arr4 := []int{1, 2, 3, 4, 5, 6}
	fmt.Println(reverse(arr4))
}
