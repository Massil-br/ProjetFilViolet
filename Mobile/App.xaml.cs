namespace Mobile;

public partial class App : Application
{
    public static string? CurrentAuthToken { get; set; }

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
        // On désactive l'auto-login via SecureStorage pour permettre de tester avec plusieurs instances locales
        // var token = await SecureStorage.Default.GetAsync("auth_token");
        // if (!string.IsNullOrEmpty(token))
        // {
        //     window.Page = new NavigationPage(new LobbyPage());
        // }
    }
}

