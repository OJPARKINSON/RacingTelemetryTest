namespace Functions
{
    class Fucntions
    {
        public double Function1(int baseGrip, double tireWearRate, int lapsCompleted, double degFactor, double refLaptime, double grip_coefficient)
        {
            double gripLevel = Math.Pow(baseGrip * (1 - tireWearRate * lapsCompleted), degFactor);
            double lapTimeImpact = refLaptime / (1 * grip_coefficient * gripLevel);

            return lapTimeImpact;
        }
    }
}