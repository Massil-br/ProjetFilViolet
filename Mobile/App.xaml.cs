namespace Mobile;

public partial class App : Application
{
    public App()
    {
        InitializeComponent();
    }

    protected override Window CreateWindow(IActivationState? activationState)
    {
        // On crée la fenêtre de base avec la MainPage
        var window = new Window(new NavigationPage(new MainPage()));
        
        // On vérifie en arrière-plan si le joueur est déjà connecté
        CheckUserLoginAsync(window);
        
        return window;
    }

    private async void CheckUserLoginAsync(Window window)
    {
        var token = await SecureStorage.Default.GetAsync("auth_token");

        if (!string.IsNullOrEmpty(token))
        {
            // Si le joueur a un token, on remplace la page par le Lobby
            window.Page = new NavigationPage(new LobbyPage());
        }
    }
}

