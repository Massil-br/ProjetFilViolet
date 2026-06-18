using Mobile.Services;

namespace Mobile;

public partial class RegisterPage : ContentPage
{
    private readonly ApiService _apiService = new ApiService();

    public RegisterPage()
    {
        InitializeComponent();
    }

   private async void OnSubmitRegisterClicked(object sender, EventArgs e)
    {
        string username = UsernameEntry.Text;
        string email = EmailEntry.Text;
        string confirmEmail = ConfirmEmailEntry.Text;
        string password = PasswordEntry.Text;
        string confirmPassword = ConfirmPasswordEntry.Text; // NOUVEAU

        // On vérifie que tout est rempli
        if (string.IsNullOrWhiteSpace(username) || string.IsNullOrWhiteSpace(email) || 
            string.IsNullOrWhiteSpace(confirmEmail) || string.IsNullOrWhiteSpace(password) || 
            string.IsNullOrWhiteSpace(confirmPassword))
        {
            await DisplayAlert("Erreur", "Veuillez remplir tous les champs", "OK");
            return;
        }

        try
        {
            // N'oublie pas de passer confirmPassword à la fin
            bool isSuccess = await _apiService.RegisterAsync(username, email, confirmEmail, password, confirmPassword);

            if (isSuccess)
            {
                await DisplayAlert("Succès", "Compte créé ! Vous pouvez maintenant vous connecter.", "OK");
                await Navigation.PopAsync();
            }
        }
        catch (Exception ex)
        {
            await DisplayAlert("Erreur Réelle", ex.ToString(), "OK");
        }
    }


}
