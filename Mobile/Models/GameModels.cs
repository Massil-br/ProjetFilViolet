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

public class Card
{
    [JsonPropertyName("rank")]
    public int Rank { get; set; }

    [JsonPropertyName("suit")]
    public int Suit { get; set; }
}

public class SecurePlayerState
{
    [JsonPropertyName("id")]
    public uint Id { get; set; }

    [JsonPropertyName("nick_name")]
    public string NickName { get; set; } = string.Empty;

    [JsonPropertyName("chips")]
    public ulong Chips { get; set; }

    [JsonPropertyName("hole_cards")]
    public List<Card>? HoleCards { get; set; }

    [JsonPropertyName("current_bet")]
    public ulong CurrentBet { get; set; }

    [JsonPropertyName("total_contribution")]
    public ulong TotalContribution { get; set; }

    [JsonPropertyName("is_active")]
    public bool IsActive { get; set; }

    [JsonPropertyName("is_all_in")]
    public bool IsAllIn { get; set; }

    [JsonPropertyName("has_acted")]
    public bool HasActed { get; set; }
}

public class Pot
{
    [JsonPropertyName("amount")]
    public ulong Amount { get; set; }

    [JsonPropertyName("eligible")]
    public List<uint>? Eligible { get; set; }
}

public class SecureTableState
{
    [JsonPropertyName("id")]
    public uint Id { get; set; }

    [JsonPropertyName("players")]
    public List<SecurePlayerState> Players { get; set; } = new();

    [JsonPropertyName("community_cards")]
    public List<Card> CommunityCards { get; set; } = new();

    [JsonPropertyName("stage")]
    public string Stage { get; set; } = string.Empty;

    [JsonPropertyName("dealer_idx")]
    public int DealerIdx { get; set; }

    [JsonPropertyName("small_blind_idx")]
    public int SmallBlindIdx { get; set; }

    [JsonPropertyName("big_blind_idx")]
    public int BigBlindIdx { get; set; }

    [JsonPropertyName("current_turn_idx")]
    public int CurrentTurnIdx { get; set; }

    [JsonPropertyName("current_bet")]
    public ulong CurrentBet { get; set; }

    [JsonPropertyName("min_raise")]
    public ulong MinRaise { get; set; }

    [JsonPropertyName("small_blind_amount")]
    public ulong SmallBlindAmount { get; set; }

    [JsonPropertyName("big_blind_amount")]
    public ulong BigBlindAmount { get; set; }

    [JsonPropertyName("game_in_progress")]
    public bool GameInProgress { get; set; }

    [JsonPropertyName("pots")]
    public List<Pot> Pots { get; set; } = new();

    [JsonPropertyName("seconds_remaining")]
    public int SecondsRemaining { get; set; }

    [JsonPropertyName("last_action_message")]
    public string LastActionMessage { get; set; } = string.Empty;
}

public class GameStateMessage
{
    [JsonPropertyName("type")]
    public string Type { get; set; } = string.Empty;

    [JsonPropertyName("state")]
    public SecureTableState? State { get; set; }

    [JsonPropertyName("error")]
    public string? Error { get; set; }
}
