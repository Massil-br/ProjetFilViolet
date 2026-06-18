namespace Mobile;

public partial class LobbyPage : ContentPage
{
    public LobbyPage()
    {
        InitializeComponent();
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

