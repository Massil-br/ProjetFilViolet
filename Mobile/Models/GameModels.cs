using System.Text.Json.Serialization;

namespace Mobile.Models;

public class Saloon
{
    [JsonPropertyName("id")]
    public uint Id { get; set; }

    [JsonPropertyName("name")]
    public string Name { get; set; } = string.Empty;

    // Ajoute d'autres propriétés si ton modèle Go en a (ex: MinBet, MaxBet, etc.)
}

public class Table
{
    [JsonPropertyName("id")]
    public uint Id { get; set; }

    [JsonPropertyName("saloon_id")]
    public uint SaloonId { get; set; }

    [JsonPropertyName("slots_available")]
    public ulong SlotsAvailable { get; set; }
}
