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
	EstimatedSpeed   float64
	FuelLoadEstimate float64
	TireAge          int
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

// generateTelemetryData creates realistic F1 telemetry data for Monaco GP
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
		// Tire degradation factor
		tireDeg := 1.0 + float64(lap-1)*0.008 // 0.8% degradation per lap

		// Fuel load effect (lighter car = faster)
		fuelRemaining := 110.0 - float64(lap-1)*2.2     // Starting fuel 110kg, 2.2kg per lap
		fuelEffect := 1.0 - (fuelRemaining-20.0)*0.0003 // Weight penalty

		for sample := 0; sample < samplesPerLap; sample++ {
			// Current time and position
			currentTime := float64(lap-1)*lapTimeBase + float64(sample)/sampleRate
			lapProgress := float64(sample) / float64(samplesPerLap)
			currentDistance := float64(lap-1)*trackLength + lapProgress*trackLength

			// Monaco-specific speed profile (street circuit with slow corners)
			var baseSpeed float64
			switch {
			case lapProgress < 0.15: // Casino Square area
				baseSpeed = 45 + 30*math.Sin(lapProgress*10)
			case lapProgress < 0.3: // Uphill to Massenet
				baseSpeed = 60 + 40*lapProgress
			case lapProgress < 0.45: // Casino to Mirabeau
				baseSpeed = 50 + 25*math.Sin(lapProgress*8)
			case lapProgress < 0.55: // Hairpin
				baseSpeed = 25 + 15*math.Sin(lapProgress*20)
			case lapProgress < 0.75: // Portier to Tunnel
				baseSpeed = 80 + 50*lapProgress
			case lapProgress < 0.85: // Swimming Pool
				baseSpeed = 60 + 20*math.Sin(lapProgress*15)
			default: // Back straight
				baseSpeed = 90 + 60*(1-lapProgress)
			}

			// Apply tire degradation and fuel effects
			speed := baseSpeed * fuelEffect / tireDeg

			// Add realistic noise
			speed += normalRandom(0, 2)
			speed = clamp(speed, 20, 320) // Realistic F1 speed limits

			// Throttle and brake based on speed profile
			var throttle, brakePressure float64
			brakingZones := []float64{0.15, 0.45, 0.55, 0.85}
			isBraking := false
			for _, zone := range brakingZones {
				if math.Abs(lapProgress-zone) < 0.02 {
					isBraking = true
					break
				}
			}

			if isBraking { // Braking zones
				throttle = uniformRandom(0, 30)
				brakePressure = uniformRandom(80, 150)
			} else if speed > 200 { // High speed sections
				throttle = uniformRandom(85, 100)
				brakePressure = 0
			} else { // Medium speed corners
				throttle = uniformRandom(40, 80)
				brakePressure = uniformRandom(0, 20)
			}

			// Tire temperatures (Monaco is hard on tires)
			baseTireTemp := 85.0 + float64(lap)*3.0 // Increasing with tire degradation
			tireTempFL := baseTireTemp + normalRandom(0, 5) + (throttle * 0.2)
			tireTempFR := baseTireTemp + normalRandom(0, 5) + (throttle * 0.15)
			tireTempRL := baseTireTemp + normalRandom(0, 4) + (throttle * 0.25)
			tireTempRR := baseTireTemp + normalRandom(0, 4) + (throttle * 0.2)

			// Fuel flow (higher at high throttle)
			fuelFlow := 20.0 + (throttle * 0.8) + normalRandom(0, 3)
			fuelFlow = clamp(fuelFlow, 0, 110) // F1 fuel flow limit

			// Engine RPM
			var rpm float64
			if speed < 50 {
				rpm = 6000 + speed*40
			} else {
				rpm = 8000 + (speed-50)*30
			}
			rpm += normalRandom(0, 100)
			rpmInt := int(clamp(rpm, 4000, 15000))

			// DRS (only on main straight - limited in Monaco)
			var drsActive int
			if lapProgress > 0.75 && speed > 150 && brakePressure < 5 {
				drsActive = 1
			} else {
				drsActive = 0
			}

			// Battery deployment (ERS)
			var batteryDeployment float64
			if throttle > 70 {
				batteryDeployment = uniformRandom(120, 160) // kW
			} else {
				batteryDeployment = uniformRandom(0, 50)
			}

			// Gear estimation
			var gear int
			if speed < 60 {
				gear = int(math.Max(1, math.Min(3, math.Floor(speed/25)+1)))
			} else {
				gear = int(math.Max(3, math.Min(8, math.Floor(speed/40)+2)))
			}

			// Steering angle (Monaco has many turns)
			var steeringAngle float64
			majorCorners := []float64{0.1, 0.2, 0.4, 0.55, 0.8}
			isMajorCorner := false
			for _, corner := range majorCorners {
				if math.Abs(lapProgress-corner) < 0.02 {
					isMajorCorner = true
					break
				}
			}

			if isMajorCorner {
				steeringAngle = uniformRandom(-45, 45)
			} else {
				steeringAngle = uniformRandom(-10, 10)
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

// generateRaceParameters creates race parameters for Monaco GP
func generateRaceParameters() []RaceParameter {
	return []RaceParameter{
		{"track_name", "Monaco", "", "Circuit name"},
		{"track_length", 3.337, "km", "Track length"},
		{"total_laps", 78, "laps", "Total race laps"},
		{"base_grip", 0.95, "coefficient", "Base tire grip level"},
		{"tire_wear_rate", 0.012, "per_lap", "Tire degradation rate"},
		{"degradation_factor", 1.8, "factor", "Degradation curve steepness"},
		{"grip_coefficient", 0.85, "coefficient", "Grip to lap time conversion"},
		{"reference_lap_time", 78.5, "seconds", "Reference lap time"},
		{"base_consumption", 2.2, "kg/lap", "Base fuel consumption"},
		{"weight_penalty", 0.0003, "factor", "Fuel weight penalty"},
		{"base_drag", 0.28, "coefficient", "Base drag coefficient"},
		{"damage_factor", 0.15, "factor", "Aero damage impact"},
		{"base_downforce", 850, "N", "Base downforce"},
		{"air_density_factor", 1.0, "factor", "Air density correction"},
		{"base_corner_speed", 65, "km/h", "Base cornering speed"},
		{"slipstream_range", 50, "meters", "Slipstream effective range"},
		{"slipstream_factor", 0.08, "factor", "Slipstream benefit"},
		{"track_difficulty", 0.7, "factor", "Overtaking difficulty"},
		{"pit_lane_time", 22.5, "seconds", "Pit lane transit time"},
		{"tire_change_time", 2.8, "seconds", "Tire change duration"},
		{"pit_lane_penalty", 0.5, "seconds", "Additional pit penalty"},
		{"average_gap_per_position", 0.8, "seconds", "Time gap per position"},
		{"ambient_temp", 24, "celsius", "Ambient temperature"},
		{"track_temp", 42, "celsius", "Track temperature"},
		{"humidity", 65, "percent", "Relative humidity"},
		{"wind_speed", 5, "km/h", "Wind speed"},
		{"tire_compound", "Medium", "", "Current tire compound"},
		{"fuel_capacity", 110, "kg", "Maximum fuel capacity"},
		{"current_fuel", 108.5, "kg", "Current fuel load"},
	}
}

// generateCompetitorData creates competitor data for overtaking analysis
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

		competitor := Competitor{
			CarNumber:        i,
			Position:         position,
			GapToLeader:      math.Round((float64(i)*1.2+uniformRandom(-0.5, 0.5))*100) / 100,
			LastLapTime:      math.Round((78.5+uniformRandom(-1.5, 3.0))*1000) / 1000,
			TireCompound:     tireCompounds[rng.Intn(len(tireCompounds))],
			PitStops:         rng.Intn(2), // 0 or 1
			EstimatedSpeed:   math.Round((220+uniformRandom(-20, 30))*10) / 10,
			FuelLoadEstimate: math.Round((uniformRandom(95, 110))*10) / 10,
			TireAge:          rng.Intn(21) + 5, // 5-25 laps
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
		"tire_compound", "pit_stops", "estimated_speed", "fuel_load_estimate", "tire_age",
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
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	fmt.Println("Generating Monaco GP telemetry data...")
	start := time.Now()

	// Generate all data
	telemetryData := generateTelemetryData()
	raceParams := generateRaceParameters()
	competitorData := generateCompetitorData()

	// Write CSV files
	if err := writeTelemetryCSV(telemetryData, "../data/telemetry_data.csv"); err != nil {
		fmt.Printf("Error writing telemetry data: %v\n", err)
		return
	}

	if err := writeRaceParametersCSV(raceParams, "../data/race_parameters.csv"); err != nil {
		fmt.Printf("Error writing race parameters: %v\n", err)
		return
	}

	if err := writeCompetitorCSV(competitorData, "../data/competitor_data.csv"); err != nil {
		fmt.Printf("Error writing competitor data: %v\n", err)
		return
	}

	duration := time.Since(start)

	fmt.Printf("Generated files:\n")
	fmt.Printf("- telemetry_data.csv: %d samples\n", len(telemetryData.Time))
	fmt.Printf("- race_parameters.csv: %d parameters\n", len(raceParams))
	fmt.Printf("- competitor_data.csv: %d competitors\n", len(competitorData))
	fmt.Printf("\nData represents 10 laps of Monaco GP telemetry at 10Hz sampling rate\n")
	fmt.Printf("Total telemetry samples: %d\n", len(telemetryData.Time))
	fmt.Printf("Generation completed in %v\n", duration)
}
