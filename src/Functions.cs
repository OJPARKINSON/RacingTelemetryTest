namespace Functions
{
    class Fucntions
    {
        public double Function1(double baseGrip, double tireWearRate, int lapsCompleted, double degFactor, double refLaptime, double gripCoefficient)
        {
            double gripLevel = baseGrip * Math.Pow(1 - tireWearRate * lapsCompleted, degFactor);

            double lapTimeImpact = refLaptime / (1 + gripCoefficient * gripLevel);

            return lapTimeImpact;
        }

        // This function meets critera but doesn't caculate the value for the end of the race
        public double Function2(double baseConsumption, double weightPenalty, double currentFuelLoad, int remainingLaps)
        {
            var fuelPerLap = baseConsumption + (weightPenalty * currentFuelLoad);
            var remainingFuel = currentFuelLoad - (fuelPerLap * remainingLaps);
            var fuelSaveRequired = remainingFuel < 0.0 ? Math.Abs(remainingFuel) / remainingLaps : 0;
            return fuelSaveRequired;

        }
    }
}