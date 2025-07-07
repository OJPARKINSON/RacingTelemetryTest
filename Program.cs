using System.Globalization;
using CSVReader;

DataReader DR = new DataReader();

var (competitors, raceParams, telemetry) = await DR.ReadCsvFilesConcurrentlyAsync();

Console.WriteLine(competitors[1].estimated_speed);
Console.WriteLine(raceParams[1].value);
Console.WriteLine(telemetry[1].gear);

