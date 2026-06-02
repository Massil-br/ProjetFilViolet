package main

// À implémenter quand le CRUD prêts sera disponible.
// Récupère l'argent possédé par l'utilisateur en base.
func moneyTaker(userID int) int {
	// TODO: récupération en DB
	return 0
}

func acountMoneyManager(userID int, amount int) int {
	userMoney := moneyTaker(userID)

	if amount >= 0 {
		userMoney += amount
		return userMoney
	} else if userMoney + amount <= 0 {
		userMoney = 0
		return userMoney
	} else {
		userMoney -= amount
		return userMoney
	}
}
