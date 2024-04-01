
package main

import ("fmt")

func main(){

	arrayString := []string{"I","am","stupid","and","weak"}
	fmt.Printf("arrayString is %+v\n", arrayString)

	for index, str1 := range arrayString {
		/*if str1 == arrayString[2] {
			arrayString[index] = "smart"
		} else if str1 == arrayString[4] {
			arrayString[index] = "strong"
		}*/

		switch index {
		case 2:
			arrayString[index] = "smart"
		case 4:
			arrayString[index] = "strong"
		default:

		}
	}
	fmt.Printf("arrayString is %+v\n", arrayString)

}