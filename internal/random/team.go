package random

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
)

// Generator produces random team compositions based on configuration.
type Generator struct {
	composition     []string
	allowDuplicates bool
	heroes          map[string][]string
}

// Member represents a hero pick in the generated team.
type Member struct {
	Role string `json:"role"`
	Hero string `json:"hero"`
}

// Team holds the resulting team members.
type Team struct {
	Members []Member `json:"members"`
}

// NewGenerator constructs a team generator with validation.
func NewGenerator(composition []string, allowDuplicates bool, heroes map[string][]string) (*Generator, error) {
	if len(composition) == 0 {
		return nil, fmt.Errorf("random: composition must not be empty")
	}

	for _, role := range composition {
		if len(heroes[role]) == 0 {
			return nil, fmt.Errorf("random: no heroes configured for role %s", role)
		}
	}

	return &Generator{
		composition:     append([]string(nil), composition...),
		allowDuplicates: allowDuplicates,
		heroes:          heroes,
	}, nil
}

// Generate returns a random team that satisfies the configured composition.
func (g *Generator) Generate() (Team, error) {
	picked := make(map[string]struct{})
	team := Team{Members: make([]Member, 0, len(g.composition))}

	for _, role := range g.composition {
		pool := g.heroes[role]
		if len(pool) == 0 {
			return Team{}, fmt.Errorf("random: empty hero pool for role %s", role)
		}

		choice, err := g.pick(pool, picked)
		if err != nil {
			return Team{}, err
		}

		team.Members = append(team.Members, Member{Role: role, Hero: choice})

		if !g.allowDuplicates {
			picked[choice] = struct{}{}
		}
	}

	// Provide deterministic ordering for easier testing/consumers.
	sort.Slice(team.Members, func(i, j int) bool {
		if team.Members[i].Role == team.Members[j].Role {
			return team.Members[i].Hero < team.Members[j].Hero
		}
		return team.Members[i].Role < team.Members[j].Role
	})

	return team, nil
}

func (g *Generator) pick(pool []string, picked map[string]struct{}) (string, error) {
	available := make([]string, 0, len(pool))
	for _, hero := range pool {
		if g.allowDuplicates {
			available = append(available, hero)
			continue
		}
		if _, exists := picked[hero]; !exists {
			available = append(available, hero)
		}
	}

	if len(available) == 0 {
		return "", fmt.Errorf("random: insufficient unique heroes to satisfy composition")
	}

	idx, err := cryptoRand(len(available))
	if err != nil {
		return "", fmt.Errorf("random: failed to pick hero %w", err)
	}

	return available[idx], nil
}

func cryptoRand(n int) (int, error) {
	max := big.NewInt(int64(n))
	v, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return int(v.Int64()), nil
}
