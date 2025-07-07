CSVReader.DataReader DR = new CSVReader.DataReader();

var (competitors, raceParams, telemetry) = await DR.ReadCsvFilesConcurrentlyAsync();

Console.WriteLine(competitors[1].estimated_speed);
Console.WriteLine(raceParams["degradation_factor"].Description);
Console.WriteLine(telemetry[1].gear);

Functions.Fucntions funcs = new Functions.Fucntions();

// funcs.Function1(double.Parse(raceParams["base_grip"].Value), double.Parse(raceParams["tire_wear_rate"].Value))
