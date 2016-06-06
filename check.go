package luhn

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	//"unicode"
)

//this algorithm is used for identifier validation
//to be valid, the value must be of length betweeen min/max + after treatment, the result of the operations
//should be equal to 'rest' modulo 'modulo'
//treatment (algo): digit at pair index should be multiplied by 2, if >10 sum the digits, sum each intermediary result
//starting from the end
func Check(val string, modulo int64, rest int64, minsize int64, maxsize int64) bool {
	log.Println("value sent to the luhn checker:", val)
	l := int64(len(val))
	if l < minsize || l > maxsize {
		log.Println("does not respect min/max length")
		return false
	}
	total := int64(0)
	_, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		log.Println("sequence is not a suite of digits, can not be parsed to an integer")
		return false
	}
	counter := 1
	for i := (l - 1); i > -1; i-- {
		chiffre, err := strconv.ParseInt(string(val[i]), 10, 32)
		if err != nil {
			return false
		}
		if i%2 == 0 {
			c := 2 * chiffre
			if c > 9 {
				total += c - 9
			} else {
				total += c
			}
		} else {
			total += chiffre
		}
		counter++
	}
	if total%modulo == rest {
		return true
	}
	return false
}

//9 digits (SIREN) + 4 digits (NIC) + 1 digit key (conform to luhn)
func siretCheck(siret string) (bool, string) {
	//log.Println("value sent to the siret checker:", siret)
	r := []rune(strings.Replace(strings.TrimSpace(siret), " ", "", -1))
	//log.Println(string(r))
	l := len(r)
	if l != 14 {
		//log.Println("not 14 digits long")
		return false, ""
	}
	checksum := int64(0)
	_, err := strconv.ParseInt(string(r), 10, 64)
	if err != nil {
		//log.Println("sequence is not a suite of digits, can not be parsed to an integer")
		return false, ""
	}
	counter := 1
	for i := (l - 1); i > -1; i-- {
		//log.Println("counter:", counter)
		//log.Println(string(r[i]))
		//chiffre := int64(0)
		chiffre, err := strconv.ParseInt(string(r[i]), 10, 32)
		if err != nil {
			return false, ""
		}
		if counter%2 == 0 {
			c := 2 * chiffre
			//log.Println("pair index:", counter, ", multiplied by 2:", c)
			if c > 9 {
				//log.Println(">9 --> digits summed:", c-9)
				checksum += c - 9
			} else {
				checksum += c
			}
		} else {
			//log.Println("not pair index:", counter, ",value:", chiffre)
			checksum += chiffre
		}
		counter++
	}
	if checksum%10 == 0 {
		return true, string(r)
	}
	return false, ""
}

func extractSirenFromSiret(siret string) string {
	r1 := []rune(siret)
	r2 := []rune("")
	for i, v := range r1 {
		if i < 9 {
			r2 = append(r2, v)
		}
	}
	return string(r2)
}

func GenerateTvaNumber(country string, siret string) (bool, string) {
	valid, valeurOk := LuhnCheck(siret, 10, 14, 14)
	if !valid {
		return false, "invalid siret"
	}
	valeur := ""
	co := ""
	siren := extractSirenFromSiret(valeurOk)
	if strings.ToUpper(country) == "FR" || strings.ToUpper(country) == "FRANCE" {
		co = "FR"
	} else {
		return false, "country not valid"
	}
	for i := int64(1); i < 100; i++ {
		s := ""
		valeur = ""
		if i < 10 {
			a := strconv.Itoa(int(i))
			s = fmt.Sprintf("0%s", a)
		} else {
			s = strconv.Itoa(int(i))
		}
		valeur = fmt.Sprintf("%s%s", s, siren)
		//log.Println(valeur)
		valid, valeurOk = LuhnCheck(valeur, 10, 11, 11)
		if valid {
			return valid, fmt.Sprintf("%s%s", co, valeurOk)
		}
	}
	return false, "not found the tva luhn key"
}

func LuhnCheck(siret string, modulo int64, minsize int64, maxsize int64) (bool, string) {
	//log.Println("value sent to the siret checker:", siret)
	r := []rune(strings.Replace(strings.TrimSpace(siret), " ", "", -1))
	//log.Println(string(r))
	l := int64(len(r))
	if l > maxsize || l < minsize {
		//log.Println("not 14 digits long")
		return false, ""
	}
	checksum := int64(0)
	_, err := strconv.ParseInt(string(r), 10, 64)
	if err != nil {
		//log.Println("sequence is not a suite of digits, can not be parsed to an integer")
		return false, ""
	}
	counter := 1
	for i := (l - 1); i > -1; i-- {
		//log.Println("counter:", counter)
		//log.Println(string(r[i]))
		//chiffre := int64(0)
		chiffre, err := strconv.ParseInt(string(r[i]), 10, 32)
		if err != nil {
			return false, ""
		}
		if counter%2 == 0 {
			c := 2 * chiffre
			//log.Println("pair index:", counter, ", multiplied by 2:", c)
			if c > 9 {
				//log.Println(">9 --> digits summed:", c-9)
				checksum += c - 9
			} else {
				checksum += c
			}
		} else {
			//log.Println("not pair index:", counter, ",value:", chiffre)
			checksum += chiffre
		}
		counter++
	}
	if checksum%10 == 0 {
		return true, string(r)
	}
	return false, ""
}
