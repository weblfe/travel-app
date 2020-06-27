package repositories

import "github.com/weblfe/travel-app/models"

func isForbid(data *models.User) bool {
		return data.DeletedAt!=0 || data.Status!=1
}
