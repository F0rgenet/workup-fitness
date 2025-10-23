package simulator

import "github.com/guregu/null/v6/zero"

type CreateRequest struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	MinWeight       float64 `json:"min_weight"`
	MaxWeight       float64 `json:"max_weight"`
	WeightIncrement float64 `json:"weight_increment"`
}

type CreateResponse struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	MinWeight       float64 `json:"min_weight"`
	MaxWeight       float64 `json:"max_weight"`
	WeightIncrement float64 `json:"weight_increment"`
}

type GetByIDResponse struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	MinWeight       float64 `json:"min_weight"`
	MaxWeight       float64 `json:"max_weight"`
	WeightIncrement float64 `json:"weight_increment"`
}

type UpdateRequest struct {
	Name            zero.String `json:"name"`
	Description     string      `json:"description"`
	MinWeight       float64     `json:"min_weight"`
	MaxWeight       float64     `json:"max_weight"`
	WeightIncrement float64     `json:"weight_increment"`
}

type UpdateResponse struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	MinWeight       float64 `json:"min_weight"`
	MaxWeight       float64 `json:"max_weight"`
	WeightIncrement float64 `json:"weight_increment"`
}
