package routing

import "testing"

func TestIsPasswordValidEmpty(t *testing.T) {
	res := isPasswordValid("")
	if res != false {
		t.Errorf("Password should be %v, but got %v", false, res)
	}
}

func TestIsPasswordValidNotEnoughCharacters(t *testing.T) {
	res := isPasswordValid("Asd1@3e")
	if res != false {
		t.Errorf("Password should be %v, but got %v", false, res)
	}
}

func TestIsPasswordValidNoCapitalLetter(t *testing.T) {
	res := isPasswordValid("@asdqweas")
	if res != false {
		t.Errorf("Password should be %v, but got %v", false, res)
	}
}

func TestIsPasswordValidNoLowercaseLetter(t *testing.T) {
	res := isPasswordValid("@QWFJIHSFIASNDBIF")
	if res != false {
		t.Errorf("Password should be %v, but got %v", false, res)
	}
}

func TestIsPasswordValidNoSpecialCharacter(t *testing.T) {
	res := isPasswordValid("QW123FJI312HSFIASNDBIF")
	if res != false {
		t.Errorf("Password should be %v, but got %v", false, res)
	}
}

func TestIsPasswordValidNoSpecialDigit(t *testing.T) {
	res := isPasswordValid("asddqwe!@#")
	if res != false {
		t.Errorf("Password should be %v, but got %v", false, res)
	}
}

func TestIsPasswordValidValid(t *testing.T) {
	res := isPasswordValid("AaQq@!12")
	if res != false {
		t.Errorf("Password should be %v, but got %v", true, res)
	}
}
