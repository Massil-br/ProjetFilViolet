using System.Text.Json; 
using Mobile.Services;
using Mobile.Models;
using System.Text;

namespace Mobile;

public partial class GamePage : ContentPage
{
    private WebSocketService _wsService;
    private uint _myUserId = 0;

    // État courant pour la validation de la mise
    private ulong _currentMinRaise = 20;
    private ulong _currentBigBlind = 20;
    private ulong _myCurrentChips = 0;
    private ulong _currentTableBet = 0;
    private bool _raiseVisible = false;

    public GamePage(WebSocketService wsService)
    {
        InitializeComponent();
        _wsService = wsService;
        ExtractMyUserId();
    }

    private void ExtractMyUserId()
    {
        try
        {
            if (!string.IsNullOrEmpty(App.CurrentAuthToken))
            {
                var parts = App.CurrentAuthToken.Split('.');
                if (parts.Length >= 2)
                {
                    string payload = parts[1];
                    int mod4 = payload.Length % 4;
                    if (mod4 > 0) payload += new string('=', 4 - mod4);
                    payload = payload.Replace('-', '+').Replace('_', '/');

                    var jsonBytes = Convert.FromBase64String(payload);
                    var jsonStr = System.Text.Encoding.UTF8.GetString(jsonBytes);
                    var claims = JsonSerializer.Deserialize<Dictionary<string, JsonElement>>(jsonStr);
                    if (claims != null && claims.TryGetValue("user_id", out var uidElement))
                    {
                        _myUserId = (uint)uidElement.GetDouble();
                    }
                }
            }
        }
        catch (Exception ex)
        {
            System.Diagnostics.Debug.WriteLine($"Error parsing JWT: {ex.Message}");
        }
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
            try
            {
                var options = new JsonSerializerOptions {
                    PropertyNameCaseInsensitive = true,
                    NumberHandling = System.Text.Json.Serialization.JsonNumberHandling.AllowReadingFromString
                };
                var message = JsonSerializer.Deserialize<GameStateMessage>(jsonMessage, options);

                if (message?.Type == "state" && message.State != null)
                    UpdateTableUI(message.State);
                else if (message?.Type == "error")
                    DisplayAlert("Erreur", message.Error, "OK");
            }
            catch (Exception ex)
            {
                System.Diagnostics.Debug.WriteLine($"Erreur parsing JSON: {ex.Message}");
            }
        });
    }

    private void UpdateTableUI(SecureTableState state)
    {
        // Sauvegarder les infos pour la validation du raise
        _currentMinRaise = state.MinRaise > 0 ? state.MinRaise : state.BigBlindAmount;
        _currentBigBlind = state.BigBlindAmount;
        _currentTableBet = state.CurrentBet;

        // Message d'action
        ActionMessageLabel.Text = string.IsNullOrEmpty(state.LastActionMessage)
            ? "En attente..." : state.LastActionMessage;

        // Pot total = jetons collectés (pots) + mises en cours (current bets des joueurs)
        ulong totalPot = 0;
        foreach (var pot in state.Pots ?? new List<Pot>()) totalPot += pot.Amount;
        // Ajouter les mises du tour en cours qui ne sont pas encore dans le pot
        foreach (var pl in state.Players ?? new List<SecurePlayerState>()) totalPot += pl.CurrentBet;
        PotLabel.Text = $"POT: {totalPot} ♦";
        StageLabel.Text = state.Stage;

        // Cartes communes
        var boardCards = (state.CommunityCards ?? new List<Card>())
            .Select(MapCard).ToList();
        BindableLayout.SetItemsSource(CommunityCardsLayout, boardCards);

        // Séparer : MOI vs ADVERSAIRES
        SecurePlayerState? myPlayer = null;
        int myIndex = -1;
        bool isMyTurn = false;

        var opponents = new List<PlayerUIModel>();

        for (int i = 0; i < state.Players.Count; i++)
        {
            var p = state.Players[i];

            if (p.Id == _myUserId)
            {
                myPlayer = p;
                myIndex = i;
                if (state.GameInProgress && i == state.CurrentTurnIdx)
                    isMyTurn = true;
            }
            else
            {
                string role = GetRole(state, i);
                bool isCurrentTurn = state.GameInProgress && i == state.CurrentTurnIdx;
                bool isShowdown = state.Stage == "Showdown";

                var cards = new List<CardUIModel>();
                if (p.HoleCards != null && p.HoleCards.Count > 0)
                {
                    // Cartes réelles : disponibles si c'est le Showdown ou si envoyées par le serveur
                    cards = p.HoleCards.Select(MapCard).ToList();
                }
                else if (p.IsActive && !isShowdown)
                {
                    // Dos de carte : adversaire actif hors Showdown
                    cards = new List<CardUIModel> {
                        new() { Display = "🂠", Color = Color.FromArgb("#3A6EA5") },
                        new() { Display = "🂠", Color = Color.FromArgb("#3A6EA5") }
                    };
                }
                // Au Showdown sans cartes => joueur couché, on n'affiche rien

                // Couleur de fond spéciale au Showdown pour indiquer qui a les cartes
                Color bgColor;
                if (isShowdown && cards.Count > 0)
                    bgColor = Color.FromArgb("#1A3D1A"); // Vert foncé au Showdown
                else if (p.IsActive)
                    bgColor = Color.FromArgb("#1E3A5F");
                else
                    bgColor = Color.FromArgb("#141414");

                opponents.Add(new PlayerUIModel
                {
                    NameAndRole = role + p.NickName,
                    ChipsText = $"♦ {p.Chips}",
                    CurrentBetText = $"Mise: {p.CurrentBet}",
                    HasBet = p.CurrentBet > 0,
                    BackgroundColor = bgColor,
                    Opacity = p.IsActive || isShowdown ? 1.0 : 0.35,
                    BorderColor = isCurrentTurn ? Colors.Orange : (isShowdown && cards.Count > 0 ? Colors.LimeGreen : Color.FromArgb("#2A2A3A")),
                    HoleCards = cards
                });
            }
        }

        BindableLayout.SetItemsSource(OpponentsLayout, opponents);

        // === Mettre à jour ma zone ===
        if (myPlayer != null)
        {
            string myRole = myIndex >= 0 ? GetRole(state, myIndex) : "";
            MyNameLabel.Text = myRole + myPlayer.NickName + " (Moi)";
            _myCurrentChips = myPlayer.Chips;
            MyChipsLabel.Text = $"♦ {myPlayer.Chips}";

            if (myPlayer.CurrentBet > 0)
            {
                MyBetLabel.Text = $"Mise: {myPlayer.CurrentBet}";
                MyBetLabel.IsVisible = true;
            }
            else
            {
                MyBetLabel.IsVisible = false;
            }

            // Mes cartes
            var myCards = (myPlayer.HoleCards != null && myPlayer.HoleCards.Count > 0)
                ? myPlayer.HoleCards.Select(MapCard).ToList()
                : new List<CardUIModel>();
            BindableLayout.SetItemsSource(MyCardsLayout, myCards);

            // Bordure de mon badge si c'est mon tour
            MyPlayerBadge.Stroke = isMyTurn ? Colors.Orange : Color.FromArgb("#2A4A6A");
            MyPlayerBadge.StrokeThickness = isMyTurn ? 2.5 : 1;
        }
        else
        {
            BindableLayout.SetItemsSource(MyCardsLayout, new List<CardUIModel>());
        }

        // Indicateurs de tour
        MyTurnLabel.IsVisible = isMyTurn;

        if (state.MinRaise > 0)
            MinRaiseLabel.Text = $"Min: {state.MinRaise} | Mise actuelle: {state.CurrentBet}";
        else
            MinRaiseLabel.Text = "";

        // Afficher le timer si la partie est en cours
        if (state.GameInProgress && state.SecondsRemaining > 0)
        {
            int totalSecs = 30; // ActionTimeoutSecs côté serveur
            TimerLabel.Text = $"⏱ {state.SecondsRemaining}s";
            TimerLabel.IsVisible = true;
            TimerBar.Progress = (double)state.SecondsRemaining / totalSecs;
        }
        else
        {
            TimerLabel.IsVisible = false;
            TimerBar.Progress = 1.0;
        }

        // Pré-remplir le champ de mise avec le min raise si c'est mon tour
        if (isMyTurn && !_raiseVisible)
        {
            RaiseEntry.Text = _currentMinRaise.ToString();
        }
    }

    private string GetRole(SecureTableState state, int idx)
    {
        if (!state.GameInProgress) return "";
        if (idx == state.DealerIdx) return "[D] ";
        if (idx == state.SmallBlindIdx) return "[SB] ";
        if (idx == state.BigBlindIdx) return "[BB] ";
        return "";
    }

    private CardUIModel MapCard(Card c)
    {
        string rankStr = c.Rank switch {
            11 => "J", 12 => "Q", 13 => "K", 14 => "A", _ => c.Rank.ToString()
        };
        string suitStr = c.Suit switch {
            0 => "♠", 1 => "♥", 2 => "♦", 3 => "♣", _ => "?"
        };
        Color color = (c.Suit == 1 || c.Suit == 2) ? Colors.Red : Colors.Black;
        return new CardUIModel { Display = $"{rankStr}{suitStr}", Color = color };
    }

    // === Bouton RELANCER : affiche/masque le panneau de saisie ===
    private async void OnRaiseClicked(object sender, EventArgs e)
    {
        if (!_raiseVisible)
        {
            // Afficher le panneau
            _raiseVisible = true;
            RaisePanel.IsVisible = true;
            RaiseEntry.Text = _currentMinRaise.ToString();
            RaiseButton.Text = "✓ VALIDER";
            RaiseButton.BackgroundColor = Color.FromArgb("#E65100");
        }
        else
        {
            // Valider et envoyer
            if (!ulong.TryParse(RaiseEntry.Text?.Trim(), out ulong raiseAmount))
            {
                await DisplayAlert("Montant invalide", "Veuillez entrer un nombre valide.", "OK");
                return;
            }

            // Règle poker : le raise doit être au moins le min raise
            // (table_bet actuel + min_raise, ou all-in)
            ulong minRequired = _currentTableBet + _currentMinRaise;
            if (raiseAmount > _myCurrentChips)
            {
                // All-in autorisé
                raiseAmount = _myCurrentChips;
            }
            else if (raiseAmount < minRequired && raiseAmount < _myCurrentChips)
            {
                await DisplayAlert("Mise trop basse",
                    $"Le raise minimum est {minRequired} (mise actuelle {_currentTableBet} + relance min {_currentMinRaise}).\n" +
                    $"Tu peux aussi faire all-in pour {_myCurrentChips}.", "OK");
                return;
            }

            await _wsService.SendActionAsync("raise", raiseAmount);

            // Réinitialiser le panneau
            _raiseVisible = false;
            RaisePanel.IsVisible = false;
            RaiseButton.Text = "↑ RELANCER";
            RaiseButton.BackgroundColor = Color.FromArgb("#1565C0");
        }
    }

    // Incrément / Décrément rapide du montant
    private void OnRaiseIncrement(object sender, EventArgs e)
    {
        if (ulong.TryParse(RaiseEntry.Text?.Trim(), out ulong current))
        {
            ulong step = _currentBigBlind > 0 ? _currentBigBlind : 20;
            ulong newVal = current + step;
            if (newVal > _myCurrentChips) newVal = _myCurrentChips;
            RaiseEntry.Text = newVal.ToString();
        }
    }

    private void OnRaiseDecrement(object sender, EventArgs e)
    {
        if (ulong.TryParse(RaiseEntry.Text?.Trim(), out ulong current))
        {
            ulong step = _currentBigBlind > 0 ? _currentBigBlind : 20;
            ulong minRequired = _currentTableBet + _currentMinRaise;
            ulong newVal = current > step ? current - step : minRequired;
            if (newVal < minRequired) newVal = minRequired;
            RaiseEntry.Text = newVal.ToString();
        }
    }

    private async void OnLeaveClicked(object sender, EventArgs e)
    {
        await _wsService.SendActionAsync("leave", 0);
        await Navigation.PopAsync();
    }

    private async void OnStartGameClicked(object sender, EventArgs e)
    {
        await _wsService.SendActionAsync("start_game", 0);
    }

    private async void OnFoldClicked(object sender, EventArgs e)
    {
        // Annuler le raise si ouvert
        if (_raiseVisible)
        {
            _raiseVisible = false;
            RaisePanel.IsVisible = false;
            RaiseButton.Text = "↑ RELANCER";
            RaiseButton.BackgroundColor = Color.FromArgb("#1565C0");
        }
        await _wsService.SendActionAsync("fold", 0);
    }

    private async void OnCallClicked(object sender, EventArgs e)
    {
        await _wsService.SendActionAsync("call", 0);
    }
}

public class PlayerUIModel
{
    public string NameAndRole { get; set; } = string.Empty;
    public string ChipsText { get; set; } = string.Empty;
    public string CurrentBetText { get; set; } = string.Empty;
    public bool HasBet { get; set; }
    public Color BackgroundColor { get; set; } = Colors.Transparent;
    public double Opacity { get; set; } = 1.0;
    public Color BorderColor { get; set; } = Colors.Transparent;
    public List<CardUIModel> HoleCards { get; set; } = new();
}

public class CardUIModel
{
    public string Display { get; set; } = string.Empty;
    public Color Color { get; set; } = Colors.Black;
}
