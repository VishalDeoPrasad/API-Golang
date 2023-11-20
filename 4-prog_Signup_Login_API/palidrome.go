package main

import "fmt"

var pallindrom []string
func main() {
	s := "apabobbbodid"
	for i := 0; i<len(s); i++{
		for j := i+1; j<len(s); j++{
			checkPalindrom(s[i:j+1])
		}
	}
	fmt.Println(pallindrom)
}

func checkPalindrom(s string){
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			return	
		}
	}
	pallindrom = append(pallindrom, s)

}
