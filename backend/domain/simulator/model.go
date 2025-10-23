package simulator

import (
	"github.com/guregu/null/v6"
	"github.com/guregu/null/v6/zero"
)

type Simulator struct {
	ID              int         `json:"id"`
	Name            zero.String `json:"name"`
	Description     string      `json:"description"`
	MinWeight       float64     `json:"min_weight"`
	MaxWeight       float64     `json:"max_weight"`
	WeightIncrement float64     `json:"weight_increment"`
	CreatedAt       null.Time   `json:"created_at"`
}
