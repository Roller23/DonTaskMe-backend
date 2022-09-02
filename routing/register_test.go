package routing

//
//import (
//	"DonTaskMe-backend/internal/helpers"
//	"testing"
//)
//
//func Testhelpers.IsPasswordValidEmpty(t *testing.T) {
//	res := helpers.IsPasswordValid("")
//	if res != false {
//		t.Errorf("Password should be %v, but got %v", false, res)
//	}
//}
//
//func Testhelpers.IsPasswordValidNotEnoughCharacters(t *testing.T) {
//	res := helpers.IsPasswordValid("Asd1@3e")
//	if res != false {
//		t.Errorf("Password should be %v, but got %v", false, res)
//	}
//}
//
//func Testhelpers.IsPasswordValidNoCapitalLetter(t *testing.T) {
//	res := helpers.IsPasswordValid("@asdqweas")
//	if res != false {
//		t.Errorf("Password should be %v, but got %v", false, res)
//	}
//}
//
//func Testhelpers.IsPasswordValidNoLowercaseLetter(t *testing.T) {
//	res := helpers.IsPasswordValid("@QWFJIHSFIASNDBIF")
//	if res != false {
//		t.Errorf("Password should be %v, but got %v", false, res)
//	}
//}
//
//func Testhelpers.IsPasswordValidNoSpecialCharacter(t *testing.T) {
//	res := helpers.IsPasswordValid("QW123FJI312HSFIASNDBIF")
//	if res != false {
//		t.Errorf("Password should be %v, but got %v", false, res)
//	}
//}
//
//func Testhelpers.IsPasswordValidNoSpecialDigit(t *testing.T) {
//	res := helpers.IsPasswordValid("asddqwe!@#")
//	if res != false {
//		t.Errorf("Password should be %v, but got %v", false, res)
//	}
//}
//
//func Testhelpers.IsPasswordValidValid(t *testing.T) {
//	res := helpers.IsPasswordValid("AaQq@!12")
//	if res != false {
//		t.Errorf("Password should be %v, but got %v", true, res)
//	}
//}
