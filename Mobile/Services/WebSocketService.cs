using System.Net.WebSockets;
using System.Text;
using System.Text.Json;

namespace Mobile.Services;

public class WebSocketService
{
    public event Action<string>? OnMessageReceived;
    private ClientWebSocket _clientWebSocket;
    private readonly string _wsHost = DeviceInfo.Platform == DevicePlatform.Android ? "10.0.2.2" : "localhost";

    public WebSocketService()
    {
        _clientWebSocket = new ClientWebSocket();
    }

    public async Task ConnectAsync(int tableId)
    {
        try
        {
            // On peut envoyer le token JWT dans les headers ou dans le premier message
            var token = App.CurrentAuthToken;
            
            string wsUrl = $"ws://{_wsHost}:8082/ws?token={token}&table_id={tableId}";
            Uri serverUri = new Uri(wsUrl);
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

    public async Task SendActionAsync(string action, ulong amount = 0)
    {
        if (_clientWebSocket.State == WebSocketState.Open)
        {
            var message = new { action = action, amount = amount };
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

