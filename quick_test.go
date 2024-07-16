package main

import "testing"

func TestArrSlicing(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	perPage := 2
	currPage := 5
	start := (currPage - 1) * perPage
	end := start + perPage
	if perPage*currPage > len(input) {
		t.Log("sliced: ", input)
	} else {
		t.Log("sliced: ", input[start:end])
	}
}
