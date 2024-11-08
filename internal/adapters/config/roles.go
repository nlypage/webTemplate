package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func RoleHasRights(role string, requiredRights []string) bool {
	userRights := viper.GetStringSlice(fmt.Sprintf("roles.%s", role))

	rightSet := make(map[string]struct{}, len(userRights))
	for _, right := range userRights {
		rightSet[right] = struct{}{}
	}

	for _, right := range requiredRights {
		if _, exists := rightSet[right]; !exists {
			return false
		}
	}
	return true
}
