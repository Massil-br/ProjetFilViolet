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
            // 1. On se connecte
            await _wsService.ConnectAsync();

            // 2. On envoie une action "JOIN_TABLE" au Go
            await _wsService.SendActionAsync("JOIN_TABLE", new { table_id = 1 });

            // 3. On navigue vers la page de la table de poker (GamePage)
            // Au lieu de new GamePage(), on lui passe le service :
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
        SecureStorage.Default.Remove("auth_token");
        
        // Nouvelle syntaxe .NET 9 pour revenir au menu de démarrage
        if (Application.Current?.Windows.Count > 0)
        {
            Application.Current.Windows[0].Page = new NavigationPage(new MainPage());
        }
    }
}

