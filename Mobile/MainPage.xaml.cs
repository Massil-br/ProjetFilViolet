namespace Mobile;

public partial class MainPage : ContentPage
{
	int count = 0;

	public MainPage()
	{
		InitializeComponent();
	}

	private void OnRegisterClicked(object? sender, EventArgs e)
	{
		count++;

		if (count == 1)
			RegisterBtn.Text = $"Clicked {count} time";
		else
			RegisterBtn.Text = $"Clicked {count} times";

		SemanticScreenReader.Announce(RegisterBtn.Text);
	}

	private void OnLoginClicked(object? sender, EventArgs e)
	{
		count++;

		if (count == 1)
			LoginBtn.Text = $"Clicked {count} time";
		else
			LoginBtn.Text = $"Clicked {count} times";

		SemanticScreenReader.Announce(LoginBtn.Text);
	}
}
