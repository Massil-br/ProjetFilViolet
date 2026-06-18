namespace Mobile;

using Mobile.Services;

public partial class LobbyPage : ContentPage
{
    public LobbyPage()
    {
        InitializeComponent();
    }

    private WebSocketService _wsService = new WebSocketService();

        private async void OnJoinTableClicked(object sender, EventArgs e)
    {
        try
        {
            // 1. On se connecte à la table 1
            await _wsService.ConnectAsync(1);

            // 2. On envoie une action "buy_in" avec 500 jetons
            await _wsService.SendActionAsync("buy_in", 500);

            // 3. On navigue vers la page de la table de poker (GamePage)
            await Navigation.PushAsync(new GamePage(_wsService));
 
        }
        catch (Exception ex)
        {
            System.Diagnostics.Debug.WriteLine($"=== ERREUR WS : {ex.Message} ===");
            await DisplayAlert("Erreur de connexion", ex.Message, "OK");
        }
    }



    private void OnLogoutClicked(object sender, EventArgs e)
    {
        // On supprime le token pour déconnecter le joueur
        App.CurrentAuthToken = null;
        
        // Nouvelle syntaxe .NET 9 pour revenir au menu de démarrage
        if (Application.Current?.Windows.Count > 0)
        {
            Application.Current.Windows[0].Page = new NavigationPage(new MainPage());
        }
    }
}

