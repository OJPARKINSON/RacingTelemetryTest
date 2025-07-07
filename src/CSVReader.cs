using System.Globalization;
using CsvHelper.Configuration.Attributes;

namespace CSVReader
{
    class DataReader
    {

        private async Task<List<CompetitorData>> ReadCompetitorData()
        {
            return await Task.Run(() =>
            {
                using var reader = new StreamReader("./data/competitor_data.csv");
                using var csv = new CsvHelper.CsvReader(reader, CultureInfo.InvariantCulture);
                return csv.GetRecords<CompetitorData>().ToList();
            });
        }

        public async Task<Dictionary<string, RaceParamValue>> ReadParamData()
        {
            using var reader = new StreamReader("./data/race_parameters.csv");
            using var csv = new CsvHelper.CsvReader(reader, CultureInfo.InvariantCulture);

            return await Task.Run(() =>
                csv.GetRecords<RaceParams>()
                .ToDictionary(
                    r => r.Parameter,
                    r => new RaceParamValue
                    {
                        Value = r.Value,
                        Unit = r.Unit,
                        Description = r.Description
                    }
                ));
        }


        public class RaceParamValue
        {
            public string Value { get; set; }
            public string Unit { get; set; }
            public string Description { get; set; }
        }



        public async Task<(List<CompetitorData>, Dictionary<string, RaceParamValue>)> ReadCsvFilesConcurrentlyAsync()
        {
            var competitorTask = ReadCompetitorData();
            var paramTask = ReadParamData();

            return (await competitorTask, await paramTask);
        }
    }
}


public class CompetitorData
{
    public int car_number { get; set; }
    public int position { get; set; }
    public double gap_to_leader { get; set; }
    public double last_lap_time { get; set; }
    public string tire_compound { get; set; }
    public int pit_stops { get; set; }
    public double estimated_speed { get; set; }
    public double fuel_load_estimate { get; set; }
    public int tire_age { get; set; }
}
public class RaceParams
{
    [Name("parameter")]
    public string Parameter { get; set; }

    [Name("value")]
    public string Value { get; set; }

    [Name("unit")]
    public string Unit { get; set; }

    [Name("description")]
    public string Description { get; set; }
}
public class TelemetryData
{
    public double time { get; set; }
    public double lap { get; set; }
    public double distance { get; set; }
    public double speed { get; set; }
    public double throttle { get; set; }
    public double brake_pressure { get; set; }
    public double tire_temp_fl { get; set; }
    public double tire_temp_fr { get; set; }
    public double tire_temp_rl { get; set; }
    public double tire_temp_rr { get; set; }
    public double fuel_flow { get; set; }
    public double engine_rpm { get; set; }
    public double drs_active { get; set; }
    public double battery_deployment { get; set; }
    public double gear { get; set; }
    public double steering_angle { get; set; }

}
