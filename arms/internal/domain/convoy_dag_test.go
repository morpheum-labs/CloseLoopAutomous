package domain

import "testing"

func TestValidateConvoySubtasks_Acyclic(t *testing.T) {
	a := SubtaskID("a")
	b := SubtaskID("b")
	err := ValidateConvoySubtasks([]Subtask{
		{ID: a, AgentRole: "x"},
		{ID: b, AgentRole: "y", DependsOn: []SubtaskID{a}},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestValidateConvoySubtasks_Cycle(t *testing.T) {
	a := SubtaskID("a")
	b := SubtaskID("b")
	err := ValidateConvoySubtasks([]Subtask{
		{ID: a, AgentRole: "x", DependsOn: []SubtaskID{b}},
		{ID: b, AgentRole: "y", DependsOn: []SubtaskID{a}},
	})
	if err == nil {
		t.Fatal("want cycle error")
	}
}

func TestValidateConvoySubtasks_UnknownDep(t *testing.T) {
	err := ValidateConvoySubtasks([]Subtask{
		{ID: "a", AgentRole: "x", DependsOn: []SubtaskID{"nope"}},
	})
	if err == nil {
		t.Fatal("want unknown dep error")
	}
}
