using System.Net.Http.Json;
using Mobile.Models;
using System.Net.Http.Headers;


namespace Mobile.Services;

public class ApiService
{
    private readonly string _baseUrl = "http://127.0.0.1:8081/api";  
    private readonly HttpClient _httpClient;

    public ApiService()
    {
        _httpClient = new HttpClient();
    }

    public async Task<string?> LoginAsync(string email, string password)
    {
        try
        {
            var request = new LoginRequest { Email = email, Password = password };
            var response = await _httpClient.PostAsJsonAsync($"{_baseUrl}/login", request);

            if (response.IsSuccessStatusCode)
            {
                var result = await response.Content.ReadFromJsonAsync<AuthResponse>();
                return result?.Token;
            }
            return null;
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Erreur Login : {ex.Message}");
            return null;
        }
    }

   // Ajoute confirmEmail dans les paramètres
    // Ajoute confirmPassword dans les paramètres
    public async Task<bool> RegisterAsync(string username, string email, string confirmEmail, string password, string confirmPassword)
    {
        var request = new RegisterRequest { 
            Username = username, 
            Email = email, 
            ConfirmEmail = confirmEmail, 
            Password = password,
            ConfirmPassword = confirmPassword // NOUVEAU
        };
        
        var response = await _httpClient.PostAsJsonAsync($"{_baseUrl}/register", request);

        if (response.IsSuccessStatusCode)
        {
            return true;
        }
        else
        {
            string errorContent = await response.Content.ReadAsStringAsync();
            throw new Exception($"Le serveur a répondu avec le code {response.StatusCode} :\n{errorContent}");
        }
    }

    private async Task SetAuthorizationHeader()
{
    var token = await SecureStorage.Default.GetAsync("auth_token");
    if (!string.IsNullOrEmpty(token))
    {
        _httpClient.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", token);
    }
}

// Exemple pour récupérer un salon par son ID (basé sur ton InitSaloonRoutes)
public async Task<Saloon?> GetSaloonAsync(uint saloonId)
{
    await SetAuthorizationHeader();
    try
    {
        var response = await _httpClient.GetAsync($"{_baseUrl}/saloons/{saloonId}");
        if (response.IsSuccessStatusCode)
        {
            return await response.Content.ReadFromJsonAsync<Saloon>();
        }
        return null;
    }
    catch (Exception ex)
    {
        Console.WriteLine($"Erreur GetSaloon : {ex.Message}");
        return null;
    }
}

}
