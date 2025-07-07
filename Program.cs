using System.Globalization;

CSVReader.DataReader DR = new CSVReader.DataReader();

var (competitors, raceParams) = await DR.ReadCsvFilesConcurrentlyAsync();

Functions.Fucntions funcs = new Functions.Fucntions();



using var reader = new StreamReader("./data/telemetry_data.csv");
using (var csv = new CsvHelper.CsvReader(reader, CultureInfo.InvariantCulture))
{
    csv.Read();
    csv.ReadHeader();
    while (csv.Read())
    {

        // double func1 = funcs.Function1(double.Parse(raceParams["base_grip"].Value), double.Parse(raceParams["tire_wear_rate"].Value), csv.GetField<int>(1) - 1, double.Parse(raceParams["degradation_factor"].Value), double.Parse(raceParams["reference_lap_time"].Value), double.Parse(raceParams["grip_coefficient"].Value));
        // Console.WriteLine(func1);

    }
}
;

