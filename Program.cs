using System.Globalization;
using Microsoft.Extensions.FileSystemGlobbing.Internal.PathSegments;

CSVReader.DataReader DR = new CSVReader.DataReader();

var (competitors, raceParams) = await DR.ReadCsvFilesConcurrentlyAsync();

Functions.Fucntions funcs = new Functions.Fucntions();

using var reader = new StreamReader("./data/telemetry_data.csv");
using (var csv = new CsvHelper.CsvReader(reader, CultureInfo.InvariantCulture))
{
    csv.Read();
    csv.ReadHeader();

    List<double> fiveLapTopSpeed = new List<double>();
    int currentLap = 1;
    double fastestSpeed = 0.0;

    while (csv.Read())
    {
        if (csv.GetField<int>(1) != currentLap)
        {
            currentLap = csv.GetField<int>(1);

            fiveLapTopSpeed.Add(fastestSpeed);
            fastestSpeed = 0.0;
            if (fiveLapTopSpeed.Count > 5)
            {
                fiveLapTopSpeed.RemoveAt(0);
            }



        }

        var currentSpeed = csv.GetField<double>(3);
        if (currentSpeed > fastestSpeed)
        {
            fastestSpeed = currentSpeed;
        }


        // double func1 = funcs.Function1(double.Parse(raceParams["base_grip"].Value), double.Parse(raceParams["tire_wear_rate"].Value), csv.GetField<int>(1) - 1, double.Parse(raceParams["degradation_factor"].Value), double.Parse(raceParams["reference_lap_time"].Value), double.Parse(raceParams["grip_coefficient"].Value));
        // Console.WriteLine(func1);

        // double func2 = funcs.Function2(double.Parse(raceParams["base_consumption"].Value), double.Parse(raceParams["weight_penalty"].Value), double.Parse(raceParams["current_fuel"].Value), csv.GetField<int>(1) - 1);
        // Console.WriteLine(func2);

        // var (straightLineSpeed, cornerSpeed) = funcs.Function3(double.Parse(raceParams["base_drag"].Value), double.Parse(raceParams["damage_factor"].Value), double.Parse(raceParams["base_downforce"].Value), double.Parse(raceParams["air_density_factor"].Value), double.Parse(raceParams["max_speed"].Value), double.Parse(raceParams["base_corner_speed"].Value));

        // Console.WriteLine(straightLineSpeed);
        // Console.WriteLine(cornerSpeed);

        var yourAvgSpeed = funcs.calculateAvgPace(fiveLapTopSpeed);

        if (csv.GetField<double>(0) == 784.3)
        {

            foreach (CompetitorData comp in competitors)
            {

                var overtakingProbability = funcs.Function4(
                    yourAvgSpeed,
                    comp.estimated_speed,
                    comp.distance_to_our_car,
                    int.Parse(raceParams["slipstream_range"].Value),
                    double.Parse(raceParams["slipstream_factor"].Value),
                    double.Parse(raceParams["track_difficulty"].Value));

                Console.WriteLine(overtakingProbability);
            }
            Console.WriteLine(yourAvgSpeed);
        }


    }
}
;

