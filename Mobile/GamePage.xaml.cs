using System.Text.Json; 
using Mobile.Services;

namespace Mobile;

public partial class GamePage : ContentPage
{
    private WebSocketService _wsService;

    public GamePage(WebSocketService wsService)
    {
        InitializeComponent();
        _wsService = wsService;
    }

    protected override void OnAppearing()
    {
        base.OnAppearing();
        _wsService.OnMessageReceived += HandleServerMessage;
    }

    protected override void OnDisappearing()
    {
        base.OnDisappearing();
        _wsService.OnMessageReceived -= HandleServerMessage;
    }

    private void HandleServerMessage(string jsonMessage)
    {
        MainThread.BeginInvokeOnMainThread(() =>
        {
            System.Diagnostics.Debug.WriteLine($"GamePage a reçu : {jsonMessage}");
            
            // Met à jour le texte sur le tapis vert
            StatusLabel.Text = "Message reçu : " + jsonMessage;
        });
    }

    private async void OnLeaveClicked(object sender, EventArgs e)
    {
        await _wsService.SendActionAsync("LEAVE_TABLE", new { table_id = 1 });
        await Navigation.PopAsync();
    }
}

