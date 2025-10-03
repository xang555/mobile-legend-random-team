package random

import "testing"

func TestGeneratorGenerate(t *testing.T) {
	composition := []string{"Tank", "Mage", "Marksman"}
	heroes := map[string][]string{
		"Tank":     {"Khufra", "Atlas"},
		"Mage":     {"Lunox", "Yve"},
		"Marksman": {"Beatrix", "Brody"},
	}

	gen, err := NewGenerator(composition, false, heroes)
	if err != nil {
		t.Fatalf("expected generator, got error %v", err)
	}

	team, err := gen.Generate()
	if err != nil {
		t.Fatalf("expected team, got error %v", err)
	}

	if len(team.Members) != len(composition) {
		t.Fatalf("expected %d members, got %d", len(composition), len(team.Members))
	}
}

func TestGeneratorGenerateDisallowDuplicates(t *testing.T) {
	composition := []string{"Tank", "Tank"}
	heroes := map[string][]string{
		"Tank": {"Khufra"},
	}

	gen, err := NewGenerator(composition, false, heroes)
	if err != nil {
		t.Fatalf("expected generator, got error %v", err)
	}

	if _, err := gen.Generate(); err == nil {
		t.Fatalf("expected error due to insufficient heroes")
	}
}

func TestGeneratorErrorsOnEmptyComposition(t *testing.T) {
	if _, err := NewGenerator(nil, false, map[string][]string{}); err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestGeneratorErrorsOnMissingRole(t *testing.T) {
	composition := []string{"Tank"}
	if _, err := NewGenerator(composition, false, map[string][]string{}); err == nil {
		t.Fatalf("expected error, got none")
	}
}
