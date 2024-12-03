package randutil

import "testing"

func TestGenerateID(t *testing.T) {
	got := GenerateID(10)
	if len(got) != 10 {
		t.Errorf("GenerateID(10) = %s; want length 10", got)
	}

	got2 := GenerateID(10)
	if got == got2 {
		t.Errorf("GenerateID(10) = %s; GenerateID(10) = %s; want different", got, got2)
	}

	got3 := GenerateID(5)
	if len(got3) != 5 {
		t.Errorf("GenerateID(5) = %s; want length 5", got3)
	}

	got4 := GenerateID(50)
	if len(got4) != 50 {
		t.Errorf("GenerateID(50) = %s; want length 50", got4)
	}
}

func TestGenerateIDForPocketBase(t *testing.T) {
	got := GenerateIDForPocketBase()
	if len(got) != idPocketBaseIDLength {
		t.Errorf("GenerateIDForPocketBase() = %s; want length %d", got, idPocketBaseIDLength)
	}

	got2 := GenerateIDForPocketBase()
	if got == got2 {
		t.Errorf("GenerateIDForPocketBase() = %s; GenerateIDForPocketBase() = %s; want different", got, got2)
	}
}
