using System.Net.WebSockets;
using System.Text;
using System.Text.Json;

namespace Mobile.Services;

public class WebSocketService
{
    public event Action<string>? OnMessageReceived;
    private ClientWebSocket _clientWebSocket;
    private readonly string _wsUrl = "ws://10.0.2.2:8081/ws/game"; // Remplace par ton IP/Port (10.0.2.2 pour l'émulateur Android)

    public WebSocketService()
    {
        _clientWebSocket = new ClientWebSocket();
    }

    public async Task ConnectAsync()
    {
        try
        {
            // On peut envoyer le token JWT dans les headers ou dans le premier message
            var token = await SecureStorage.Default.GetAsync("auth_token");
            
            Uri serverUri = new Uri(_wsUrl);
            await _clientWebSocket.ConnectAsync(serverUri, CancellationToken.None);
            
            Console.WriteLine("Connecté au serveur de jeu !");

            // On lance l'écoute des messages du serveur en tâche de fond
            _ = ReceiveMessagesAsync();
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Erreur de connexion WS : {ex.Message}");
        }
    }

    public async Task SendActionAsync(string action, object data)
    {
        if (_clientWebSocket.State == WebSocketState.Open)
        {
            var message = new { Action = action, Data = data };
            string jsonMessage = JsonSerializer.Serialize(message);
            var bytes = Encoding.UTF8.GetBytes(jsonMessage);

            await _clientWebSocket.SendAsync(new ArraySegment<byte>(bytes), WebSocketMessageType.Text, true, CancellationToken.None);
        }
    }

    private async Task ReceiveMessagesAsync()
    {
        var buffer = new byte[1024 * 4];

        while (_clientWebSocket.State == WebSocketState.Open)
        {
            var result = await _clientWebSocket.ReceiveAsync(new ArraySegment<byte>(buffer), CancellationToken.None);

            if (result.MessageType == WebSocketMessageType.Close)
            {
                await _clientWebSocket.CloseAsync(WebSocketCloseStatus.NormalClosure, string.Empty, CancellationToken.None);
            }
            else
            {
                // On décode le message reçu
                string message = Encoding.UTF8.GetString(buffer, 0, result.Count);
                Console.WriteLine($"Message du serveur : {message}");

                // 3. ON DÉCLENCHE L'ÉVÉNEMENT ICI
                // S'il y a des pages qui écoutent (comme GamePage), on leur envoie le message
                OnMessageReceived?.Invoke(message);
            }
        }
    }
}

