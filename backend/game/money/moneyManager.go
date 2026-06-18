package money

import (
	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/api/models"

	"gorm.io/gorm"
)

// GetUserMoney récupère l'argent possédé par l'utilisateur en base de données.
func GetUserMoney(userID uint) (uint64, error) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return 0, err
	}
	return user.Money, nil
}

// AccountMoneyManager modifie l'argent de l'utilisateur de manière sécurisée.
// Un amount positif ajoute de l'argent. Un amount négatif en retire.
// Si le retrait dépasse le solde disponible, le solde est ramené à 0.
func AccountMoneyManager(userID uint, amount int64) (uint64, error) {
	var newBalance uint64
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		// Utilise SELECT FOR UPDATE pour éviter les race conditions
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&user, userID).Error; err != nil {
			return err
		}

		if amount >= 0 {
			user.Money += uint64(amount)
		} else {
			deduction := uint64(-amount)
			if user.Money <= deduction {
				user.Money = 0
			} else {
				user.Money -= deduction
			}
		}

		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		newBalance = user.Money
		return nil
	})

	return newBalance, err
}
