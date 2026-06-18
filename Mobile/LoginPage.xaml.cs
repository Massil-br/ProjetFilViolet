using Mobile.Services;

namespace Mobile;

public partial class LoginPage : ContentPage
{
    private readonly ApiService _apiService = new ApiService();

    public LoginPage()
    {
        InitializeComponent();
    }

    private async void OnSubmitLoginClicked(object sender, EventArgs e)
    {
        string email = EmailEntry.Text;
        string password = PasswordEntry.Text;

        if (string.IsNullOrWhiteSpace(email) || string.IsNullOrWhiteSpace(password))
        {
            await DisplayAlert("Erreur", "Veuillez remplir tous les champs", "OK");
            return;
        }

        // Appel au vrai backend
        string? token = await _apiService.LoginAsync(email, password);

        if (!string.IsNullOrEmpty(token))
        {
            App.CurrentAuthToken = token;
            
            if (Application.Current?.Windows.Count > 0)
            {
                Application.Current.Windows[0].Page = new NavigationPage(new LobbyPage());
            }
        }
        else
        {
            await DisplayAlert("Erreur", "Identifiants incorrects ou compte inexistant.", "OK");
        }
    }
}

