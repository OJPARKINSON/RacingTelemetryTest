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

        public (double, double) Function3(double baseDrag, double damageFactor, double baseDownforce, double airDensity, double maxSpeed, double baseCornerSpeed)
        {
            var dragCoefficient = baseDrag + (damageFactor * 0.1);
            var downForceLoss = baseDownforce * (1 - 0.1);
            var straightLineSpeed = maxSpeed * (1 - dragCoefficient * airDensity);
            var corneringSpeed = baseCornerSpeed * Math.Sqrt(downForceLoss / baseDownforce);

            return (straightLineSpeed, corneringSpeed);
        }
    }
}