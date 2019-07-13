package main

import "testing"

func TestCheckSubmission(t *testing.T) {
	c := NewCheck()

	veredict, err := c.CheckSubmission("P71701_ca", 2)
	if err != nil {
		t.Error(err)
	}

	t.Log(veredict)

	veredict, err = c.CheckSubmission("P68688_ca", 6)
	if err != nil {
		t.Error(err)
	}

	t.Log(veredict)

	veredict, err = c.CheckSubmission("P68688_ca", 40)
	if err != nil {
		t.Error(err)
	}

	t.Log(veredict)
}

func TestGetNumSubmissions(t *testing.T) {
	c := NewCheck()
	n, err := c.GetNumSubmissions("P68688_ca")
	if err != nil {
		t.Error(err)
	}

	t.Log(n)
}
