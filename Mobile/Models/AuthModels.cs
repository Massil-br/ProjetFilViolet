using System.Text.Json.Serialization;
namespace Mobile.Models;

public class LoginRequest
{
    public string Email { get; set; } = string.Empty;
    public string Password { get; set; } = string.Empty;
}

public class RegisterRequest
{
    // Côté Go : json:"nick_name"
    [JsonPropertyName("nick_name")]
    public string Username { get; set; } = string.Empty;

    // Côté Go : json:"email"
    [JsonPropertyName("email")]
    public string Email { get; set; } = string.Empty;

    // Côté Go : json:"confirmEmail" (basé sur ton contrôleur Go d'avant)
    [JsonPropertyName("confirmEmail")]
    public string ConfirmEmail { get; set; } = string.Empty;

    // Côté Go : json:"password"
    [JsonPropertyName("password")]
    public string Password { get; set; } = string.Empty;

    // Côté Go : json:"confirmPassword"
    [JsonPropertyName("confirmPassword")]
    public string ConfirmPassword { get; set; } = string.Empty;
}


public class AuthResponse
{
    public string Token { get; set; } = string.Empty;
    public string Message { get; set; } = string.Empty;
}
