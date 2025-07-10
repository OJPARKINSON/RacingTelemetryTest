package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// TelemetryData holds all telemetry channels
type TelemetryData struct {
	Time              []float64
	Lap               []int
	Distance          []float64
	Speed             []float64
	Throttle          []float64
	BrakePressure     []float64
	TireTempFL        []float64
	TireTempFR        []float64
	TireTempRL        []float64
	TireTempRR        []float64
	FuelFlow          []float64
	EngineRPM         []int
	DRSActive         []int
	BatteryDeployment []float64
	Gear              []int
	SteeringAngle     []float64
}

// RaceParameter represents a single race parameter
type RaceParameter struct {
	Name        string
	Value       interface{}
	Unit        string
	Description string
}

// Competitor represents competitor car data
type Competitor struct {
	CarNumber        int
	Position         int
	GapToLeader      float64
	LastLapTime      float64
	TireCompound     string
	PitStops         int
	EstimatedSpeed   float64 // Now Monaco-realistic top speeds
	FuelLoadEstimate float64
	TireAge          int
	DistanceToOurCar float64
}

// Random number generator with seed
var rng *rand.Rand

func init() {
	// Set seed for reproducible data
	rng = rand.New(rand.NewSource(42))
}

// normalRandom generates a normally distributed random number
func normalRandom(mean, stddev float64) float64 {
	return rng.NormFloat64()*stddev + mean
}

// uniformRandom generates a uniform random number between min and max
func uniformRandom(min, max float64) float64 {
	return min + rng.Float64()*(max-min)
}

// clamp constrains a value between min and max
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// getMonacoSpeedProfile returns realistic Monaco speed based on track position
func getMonacoSpeedProfile(lapProgress float64) float64 {
	// Monaco corner analysis with realistic speeds
	switch {
	case lapProgress < 0.08: // Start/finish straight
		return 140 + 45*math.Sin(lapProgress*25)
	case lapProgress < 0.12: // Turn 1 (Sainte Devote)
		return 85 + 15*math.Sin(lapProgress*30)
	case lapProgress < 0.18: // Beau Rivage climb
		return 95 + 25*lapProgress*10
	case lapProgress < 0.25: // Massenet and Casino Square
		return 70 + 20*math.Sin(lapProgress*15)
	case lapProgress < 0.32: // Mirabeau Haute approach
		return 110 + 30*lapProgress*8
	case lapProgress < 0.38: // Mirabeau (Turn 5)
		return 65 + 20*math.Sin(lapProgress*20)
	case lapProgress < 0.45: // Loews Hairpin approach
		return 80 + 25*lapProgress*6
	case lapProgress < 0.52: // Loews Hairpin (slowest corner)
		return 45 + 15*math.Sin(lapProgress*25)
	case lapProgress < 0.58: // Portier (Turn 8)
		return 75 + 30*lapProgress*5
	case lapProgress < 0.68: // Tunnel entrance to exit
		return 110 + 70*lapProgress*3 // Fastest section
	case lapProgress < 0.75: // Nouvelle Chicane approach
		return 145 + 35*math.Sin(lapProgress*12)
	case lapProgress < 0.82: // Swimming Pool chicane
		return 85 + 25*math.Sin(lapProgress*18)
	case lapProgress < 0.88: // La Rascasse (Turn 17)
		return 70 + 20*lapProgress*4
	case lapProgress < 0.95: // Anthony Noghes (Turn 19)
		return 90 + 35*lapProgress*6
	default: // Back to start/finish
		return 125 + 40*(1-lapProgress)*8
	}
}

// generateTelemetryData creates realistic Monaco F1 telemetry data
func generateTelemetryData() *TelemetryData {
	// Monaco track characteristics
	trackLength := 3.337 // km
	lapTimeBase := 78.5  // seconds base lap time

	// Sampling rate (10Hz for manageable file size)
	sampleRate := 10.0
	samplesPerLap := int(lapTimeBase * sampleRate)
	totalLaps := 10
	totalSamples := samplesPerLap * totalLaps

	// Initialize data structure
	data := &TelemetryData{
		Time:              make([]float64, 0, totalSamples),
		Lap:               make([]int, 0, totalSamples),
		Distance:          make([]float64, 0, totalSamples),
		Speed:             make([]float64, 0, totalSamples),
		Throttle:          make([]float64, 0, totalSamples),
		BrakePressure:     make([]float64, 0, totalSamples),
		TireTempFL:        make([]float64, 0, totalSamples),
		TireTempFR:        make([]float64, 0, totalSamples),
		TireTempRL:        make([]float64, 0, totalSamples),
		TireTempRR:        make([]float64, 0, totalSamples),
		FuelFlow:          make([]float64, 0, totalSamples),
		EngineRPM:         make([]int, 0, totalSamples),
		DRSActive:         make([]int, 0, totalSamples),
		BatteryDeployment: make([]float64, 0, totalSamples),
		Gear:              make([]int, 0, totalSamples),
		SteeringAngle:     make([]float64, 0, totalSamples),
	}

	// Generate data for each lap
	for lap := 1; lap <= totalLaps; lap++ {
		// Tire degradation factor (more aggressive for Monaco)
		tireDeg := 1.0 + float64(lap-1)*0.012 // 1.2% degradation per lap

		// Fuel load effect (lighter car = faster)
		fuelRemaining := 110.0 - float64(lap-1)*2.2     // Starting fuel 110kg, 2.2kg per lap
		fuelEffect := 1.0 - (fuelRemaining-20.0)*0.0003 // Weight penalty

		for sample := 0; sample < samplesPerLap; sample++ {
			// Current time and position
			currentTime := float64(lap-1)*lapTimeBase + float64(sample)/sampleRate
			lapProgress := float64(sample) / float64(samplesPerLap)
			currentDistance := float64(lap-1)*trackLength + lapProgress*trackLength

			// Monaco-specific realistic speed profile
			baseSpeed := getMonacoSpeedProfile(lapProgress)

			// Apply tire degradation and fuel effects
			speed := baseSpeed * fuelEffect / tireDeg

			// Add realistic noise
			speed += normalRandom(0, 1.5)
			speed = clamp(speed, 20, 190) // Monaco realistic speed limits

			// Throttle and brake based on speed and Monaco characteristics
			var throttle, brakePressure float64

			// Major braking zones at Monaco
			brakingZones := []float64{0.12, 0.38, 0.52, 0.75, 0.88} // Sainte Devote, Mirabeau, Hairpin, Chicane, Anthony Noghes
			isBraking := false
			for _, zone := range brakingZones {
				if math.Abs(lapProgress-zone) < 0.025 {
					isBraking = true
					break
				}
			}

			if isBraking { // Heavy braking zones
				throttle = uniformRandom(0, 25)
				brakePressure = uniformRandom(100, 200) // Heavy braking in Monaco
			} else if speed > 140 { // Tunnel section
				throttle = uniformRandom(85, 100)
				brakePressure = uniformRandom(0, 10)
			} else if speed < 70 { // Slow corners
				throttle = uniformRandom(30, 60)
				brakePressure = uniformRandom(20, 50)
			} else { // Medium speed sections
				throttle = uniformRandom(50, 85)
				brakePressure = uniformRandom(0, 25)
			}

			// Tire temperatures (Monaco is demanding on tires due to barriers and track surface)
			baseTireTemp := 90.0 + float64(lap)*2.5 // Increasing with tire degradation
			tireTempFL := baseTireTemp + normalRandom(0, 4) + (throttle * 0.15) + (brakePressure * 0.1)
			tireTempFR := baseTireTemp + normalRandom(0, 4) + (throttle * 0.12) + (brakePressure * 0.08)
			tireTempRL := baseTireTemp + normalRandom(0, 3) + (throttle * 0.18) + (brakePressure * 0.05)
			tireTempRR := baseTireTemp + normalRandom(0, 3) + (throttle * 0.15) + (brakePressure * 0.05)

			// Clamp tire temperatures to realistic ranges
			tireTempFL = clamp(tireTempFL, 80, 140)
			tireTempFR = clamp(tireTempFR, 80, 140)
			tireTempRL = clamp(tireTempRL, 80, 140)
			tireTempRR = clamp(tireTempRR, 80, 140)

			// Fuel flow (higher at high throttle, limited by regulations)
			fuelFlow := 25.0 + (throttle * 0.75) + normalRandom(0, 4)
			fuelFlow = clamp(fuelFlow, 0, 110) // F1 fuel flow limit 110 kg/h

			// Engine RPM based on speed and gear
			var rpm float64
			if speed < 60 {
				rpm = 7000 + speed*35
			} else if speed < 120 {
				rpm = 9000 + (speed-60)*25
			} else {
				rpm = 10500 + (speed-120)*15
			}
			rpm += normalRandom(0, 150)
			rpmInt := int(clamp(rpm, 5000, 15000))

			// DRS (very limited in Monaco - only small section before Sainte Devote)
			var drsActive int
			if lapProgress > 0.95 && lapProgress < 0.08 && speed > 120 && brakePressure < 15 {
				drsActive = 1
			} else {
				drsActive = 0
			}

			// Battery deployment (ERS) - strategic in Monaco due to limited overtaking
			var batteryDeployment float64
			if lapProgress > 0.58 && lapProgress < 0.68 { // Tunnel section
				batteryDeployment = uniformRandom(120, 160) // Maximum deployment
			} else if throttle > 75 {
				batteryDeployment = uniformRandom(60, 120)
			} else {
				batteryDeployment = uniformRandom(0, 40)
			}

			// Gear estimation based on Monaco characteristics
			var gear int
			if speed < 50 {
				gear = int(math.Max(1, math.Min(2, math.Floor(speed/30)+1)))
			} else if speed < 80 {
				gear = int(math.Max(2, math.Min(4, math.Floor(speed/25)+1)))
			} else if speed < 120 {
				gear = int(math.Max(3, math.Min(6, math.Floor(speed/25)+1)))
			} else {
				gear = int(math.Max(5, math.Min(8, math.Floor(speed/30)+2)))
			}

			// Steering angle (Monaco requires constant steering input)
			var steeringAngle float64
			// Monaco corner definitions with realistic steering angles
			majorCorners := []float64{0.12, 0.25, 0.38, 0.52, 0.75, 0.82, 0.88}
			isInCorner := false
			cornerIntensity := 0.0

			for _, corner := range majorCorners {
				if math.Abs(lapProgress-corner) < 0.03 {
					isInCorner = true
					if math.Abs(lapProgress-0.52) < 0.02 { // Hairpin
						cornerIntensity = 1.0
					} else if math.Abs(lapProgress-0.25) < 0.02 || math.Abs(lapProgress-0.75) < 0.02 { // Casino, Swimming Pool
						cornerIntensity = 0.8
					} else {
						cornerIntensity = 0.6
					}
					break
				}
			}

			if isInCorner {
				maxAngle := 35 + cornerIntensity*25 // Up to 60 degrees for hairpin
				steeringAngle = uniformRandom(-maxAngle, maxAngle)
			} else {
				steeringAngle = uniformRandom(-8, 8) // Small corrections on straights
			}

			// Store data with proper rounding
			data.Time = append(data.Time, math.Round(currentTime*10)/10)
			data.Lap = append(data.Lap, lap)
			data.Distance = append(data.Distance, math.Round(currentDistance*1000)/1000)
			data.Speed = append(data.Speed, math.Round(speed*10)/10)
			data.Throttle = append(data.Throttle, math.Round(throttle*10)/10)
			data.BrakePressure = append(data.BrakePressure, math.Round(brakePressure*10)/10)
			data.TireTempFL = append(data.TireTempFL, math.Round(tireTempFL*10)/10)
			data.TireTempFR = append(data.TireTempFR, math.Round(tireTempFR*10)/10)
			data.TireTempRL = append(data.TireTempRL, math.Round(tireTempRL*10)/10)
			data.TireTempRR = append(data.TireTempRR, math.Round(tireTempRR*10)/10)
			data.FuelFlow = append(data.FuelFlow, math.Round(fuelFlow*10)/10)
			data.EngineRPM = append(data.EngineRPM, rpmInt)
			data.DRSActive = append(data.DRSActive, drsActive)
			data.BatteryDeployment = append(data.BatteryDeployment, math.Round(batteryDeployment*10)/10)
			data.Gear = append(data.Gear, gear)
			data.SteeringAngle = append(data.SteeringAngle, math.Round(steeringAngle*10)/10)
		}
	}

	return data
}

// generateRaceParameters creates Monaco-specific race parameters
func generateRaceParameters() []RaceParameter {
	return []RaceParameter{
		{"track_name", "Monaco", "", "Circuit name"},
		{"track_length", 3.337, "km", "Track length"},
		{"total_laps", 78, "laps", "Total race laps"},
		{"base_grip", 0.95, "coefficient", "Base tire grip level"},
		{"tire_wear_rate", 0.015, "per_lap", "Tire degradation rate (higher for Monaco)"},
		{"degradation_factor", 1.9, "factor", "Degradation curve steepness"},
		{"grip_coefficient", 0.82, "coefficient", "Grip to lap time conversion"},
		{"reference_lap_time", 78.5, "seconds", "Reference lap time"},
		{"base_consumption", 2.1, "kg/lap", "Base fuel consumption (lower for Monaco)"},
		{"weight_penalty", 0.0003, "factor", "Fuel weight penalty"},
		{"base_drag", 0.32, "coefficient", "Base drag coefficient (higher downforce setup)"},
		{"damage_factor", 0.25, "factor", "Aero damage impact (higher risk in Monaco)"},
		{"base_downforce", 1200, "N", "Base downforce (high downforce setup)"},
		{"air_density_factor", 1.0, "factor", "Air density correction"},
		{"base_corner_speed", 65, "km/h", "Base cornering speed"},
		{"slipstream_range", 30, "meters", "Slipstream effective range (shorter in Monaco)"},
		{"slipstream_factor", 0.05, "factor", "Slipstream benefit (reduced in Monaco)"},
		{"track_difficulty", 0.95, "factor", "Overtaking difficulty (very high for Monaco)"},
		{"pit_lane_time", 25.2, "seconds", "Pit lane transit time (longer for Monaco)"},
		{"tire_change_time", 2.8, "seconds", "Tire change duration"},
		{"pit_lane_penalty", 0.8, "seconds", "Additional pit penalty"},
		{"average_gap_per_position", 1.2, "seconds", "Time gap per position (larger in Monaco)"},
		{"ambient_temp", 24, "celsius", "Ambient temperature"},
		{"track_temp", 42, "celsius", "Track temperature"},
		{"humidity", 65, "percent", "Relative humidity"},
		{"wind_speed", 8, "km/h", "Wind speed (Monaco can be gusty)"},
		{"tire_compound", "Medium", "", "Current tire compound"},
		{"fuel_capacity", 110, "kg", "Maximum fuel capacity"},
		{"current_fuel", 108.5, "kg", "Current fuel load"},
		{"max_speed", 190, "km/h", "Car maximum speed capability (Monaco limited)"},
		{"aero_damage_percentage", 0.03, "percentage", "Current aerodynamic damage level"},
		{"tire_advantage_per_lap", 1.2, "seconds", "Lap time advantage of fresh tires (higher in Monaco)"},
	}
}

// generateCompetitorData creates Monaco-realistic competitor data
func generateCompetitorData() []Competitor {
	var competitors []Competitor
	tireCompounds := []string{"Soft", "Medium", "Hard"}

	for i := 1; i <= 20; i++ {
		if i == 10 { // Skip our car (car #10)
			continue
		}

		position := i
		if i > 10 {
			position = i - 1
		}

		// Calculate realistic distance to our car based on position
		var distanceToOurCar float64
		ourPosition := 10 // Assume we're in P10
		positionDiff := math.Abs(float64(position - ourPosition))

		if position < ourPosition {
			// Cars ahead of us
			distanceToOurCar = positionDiff * uniformRandom(120, 200) // Larger gaps in Monaco
		} else {
			// Cars behind us
			distanceToOurCar = positionDiff * uniformRandom(150, 250) // Even larger gaps behind
		}

		// Add some variation for cars in close proximity (within 2 positions)
		if positionDiff <= 2 {
			distanceToOurCar = uniformRandom(40, 100) // Still relatively close racing
		}

		// Monaco-realistic top speeds (much lower than high-speed circuits)
		var monacoTopSpeed float64
		// Top speeds vary by car performance and setup
		if position <= 5 { // Top teams
			monacoTopSpeed = uniformRandom(175, 190)
		} else if position <= 10 { // Midfield
			monacoTopSpeed = uniformRandom(165, 180)
		} else { // Back markers
			monacoTopSpeed = uniformRandom(155, 170)
		}

		competitor := Competitor{
			CarNumber:        i,
			Position:         position,
			GapToLeader:      math.Round((float64(position-1)*1.8+uniformRandom(-0.8, 1.2))*100) / 100,
			LastLapTime:      math.Round((78.5+uniformRandom(-2.0, 4.5))*1000) / 1000, // More variation in Monaco
			TireCompound:     tireCompounds[rng.Intn(len(tireCompounds))],
			PitStops:         rng.Intn(2),                        // 0 or 1
			EstimatedSpeed:   math.Round(monacoTopSpeed*10) / 10, // Now Monaco-realistic!
			FuelLoadEstimate: math.Round((uniformRandom(95, 110))*10) / 10,
			TireAge:          rng.Intn(21) + 5, // 5-25 laps
			DistanceToOurCar: math.Round(distanceToOurCar*10) / 10,
		}
		competitors = append(competitors, competitor)
	}

	return competitors
}

// writeTelemetryCSV writes telemetry data to CSV file
func writeTelemetryCSV(data *TelemetryData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"time", "lap", "distance", "speed", "throttle", "brake_pressure",
		"tire_temp_fl", "tire_temp_fr", "tire_temp_rl", "tire_temp_rr",
		"fuel_flow", "engine_rpm", "drs_active", "battery_deployment",
		"gear", "steering_angle",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for i := 0; i < len(data.Time); i++ {
		row := []string{
			fmt.Sprintf("%.1f", data.Time[i]),
			strconv.Itoa(data.Lap[i]),
			fmt.Sprintf("%.3f", data.Distance[i]),
			fmt.Sprintf("%.1f", data.Speed[i]),
			fmt.Sprintf("%.1f", data.Throttle[i]),
			fmt.Sprintf("%.1f", data.BrakePressure[i]),
			fmt.Sprintf("%.1f", data.TireTempFL[i]),
			fmt.Sprintf("%.1f", data.TireTempFR[i]),
			fmt.Sprintf("%.1f", data.TireTempRL[i]),
			fmt.Sprintf("%.1f", data.TireTempRR[i]),
			fmt.Sprintf("%.1f", data.FuelFlow[i]),
			strconv.Itoa(data.EngineRPM[i]),
			strconv.Itoa(data.DRSActive[i]),
			fmt.Sprintf("%.1f", data.BatteryDeployment[i]),
			strconv.Itoa(data.Gear[i]),
			fmt.Sprintf("%.1f", data.SteeringAngle[i]),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// writeRaceParametersCSV writes race parameters to CSV file
func writeRaceParametersCSV(params []RaceParameter, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"parameter", "value", "unit", "description"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write parameter rows
	for _, param := range params {
		row := []string{
			param.Name,
			fmt.Sprintf("%v", param.Value),
			param.Unit,
			param.Description,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// writeCompetitorCSV writes competitor data to CSV file
func writeCompetitorCSV(competitors []Competitor, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"car_number", "position", "gap_to_leader", "last_lap_time",
		"tire_compound", "pit_stops", "estimated_speed", "fuel_load_estimate",
		"tire_age", "distance_to_our_car",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write competitor rows
	for _, comp := range competitors {
		row := []string{
			strconv.Itoa(comp.CarNumber),
			strconv.Itoa(comp.Position),
			fmt.Sprintf("%.2f", comp.GapToLeader),
			fmt.Sprintf("%.3f", comp.LastLapTime),
			comp.TireCompound,
			strconv.Itoa(comp.PitStops),
			fmt.Sprintf("%.1f", comp.EstimatedSpeed),
			fmt.Sprintf("%.1f", comp.FuelLoadEstimate),
			strconv.Itoa(comp.TireAge),
			fmt.Sprintf("%.1f", comp.DistanceToOurCar),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	fmt.Println("Generating Monaco-realistic GP telemetry data...")
	start := time.Now()

	// Generate all data
	telemetryData := generateTelemetryData()
	raceParams := generateRaceParameters()
	competitorData := generateCompetitorData()

	// Write CSV files
	if err := writeTelemetryCSV(telemetryData, "./data/telemetry_data.csv"); err != nil {
		fmt.Printf("Error writing telemetry data: %v\n", err)
		return
	}

	if err := writeRaceParametersCSV(raceParams, "./data/race_parameters.csv"); err != nil {
		fmt.Printf("Error writing race parameters: %v\n", err)
		return
	}

	if err := writeCompetitorCSV(competitorData, "./data/competitor_data.csv"); err != nil {
		fmt.Printf("Error writing competitor data: %v\n", err)
		return
	}

	duration := time.Since(start)

	fmt.Printf("Generated Monaco-realistic files:\n")
	fmt.Printf("- telemetry_data.csv: %d samples\n", len(telemetryData.Time))
	fmt.Printf("- race_parameters.csv: %d parameters\n", len(raceParams))
	fmt.Printf("- competitor_data.csv: %d competitors\n", len(competitorData))
	fmt.Printf("\nKey Monaco improvements:\n")
	fmt.Printf("- Realistic speed ranges: 45-190 km/h (was 45-320 km/h)\n")
	fmt.Printf("- Monaco-specific corner profiles and braking zones\n")
	fmt.Printf("- Competitor top speeds: 155-190 km/h (was 200-250 km/h)\n")
	fmt.Printf("- Higher tire degradation and track difficulty\n")
	fmt.Printf("- Reduced slipstream effectiveness\n")
	fmt.Printf("- Monaco-appropriate race parameters\n")
	fmt.Printf("\nGeneration completed in %v\n", duration)
}
