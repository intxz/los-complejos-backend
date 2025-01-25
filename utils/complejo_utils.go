package utils

import (
	"fmt"
	"strconv"
)

// CalcIMC calculates the Body Mass Index (BMI) based on weight and height.
// If weight or height is empty, it returns "N/A" to indicate that the IMC cannot be calculated.
func CalcIMC(weight, height string) string {
	// Check if weight or height is empty
	if weight == "" || height == "" {
		return "N/A" // Return "N/A" if the values are not provided
	}

	// Convert weight to float
	weightF, err := strconv.ParseFloat(weight, 32)
	if err != nil {
		fmt.Println("Error converting weight:", err)
		return "Invalid input"
	}

	// Convert height to float
	heightF, err := strconv.ParseFloat(height, 32)
	if err != nil {
		fmt.Println("Error converting height:", err)
		return "Invalid input"
	}

	// Calculate IMC
	calcIMC := weightF / (heightF * heightF)

	// Return IMC category
	if calcIMC < 18.5 {
		return "Soldado del Burgo De Los No Muertos"
	} else if calcIMC >= 18.5 && calcIMC < 25 {
		return "NPC"
	} else if calcIMC >= 25 && calcIMC < 30 {
		return "Susi Slayer"
	} else {
		return "Burger King Slayer"
	}
}
