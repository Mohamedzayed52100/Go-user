package utils

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"github.com/goplaceapp/goplace-user/internal/role/adapters/convertors"
	roleDomain "github.com/goplaceapp/goplace-user/pkg/roleservice/domain"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm/logger"
)

func GetLogLevel() logger.LogLevel {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		return logger.Error
	}

	switch logLevel {
	case "DEBUG":
		return logger.Info
	case "WARN":
		return logger.Warn
	case "INFO":
		return logger.Info
	case "SILENT":
		return logger.Silent
	default:
		return logger.Error
	}
}

func GeneratePassword() string {
	n := rand.Intn(8) + 10
	password := GenerateRandomString(n)
	return password
}

func GenerateRandomString(n int) string {
	result := make([]rune, n)
	for i := range result {
		x := rand.Int()
		if x%2 == 0 {
			result[i] = rune(rand.Intn(26) + 97)
		} else {
			result[i] = rune(rand.Intn(26) + 65)
		}
	}
	return string(result)
}

func GenerateRandomDigits(n int) string {
	result := make([]rune, n)
	for i := range result {
		result[i] = rune(rand.Intn(10) + 48)
	}
	return string(result)
}

func HashPassword(password string) string {
	resPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(resPassword)
}

func ConvertToKebabCase(role string) string {
	// Capitalize first letter in each word
	role = strings.ToLower(role)

	// Replace spaces with hyphens
	role = ReplaceSpacesWithHyphens(role)

	return role
}

func ReplaceSpacesWithHyphens(role string) string {
	return strings.ReplaceAll(role, " ", "-")
}

func GeneratePinCode() string {
	generator := rand.Intn(100000000)

	return fmt.Sprintf("%08d", generator)
}

func CategorizeAndArrangePermissions(permissions []*roleDomain.Permission) []*userProto.Permission {
	categorized := make(map[string][]*userProto.Permission)
	var uncategorized []*userProto.Permission

	for _, perm := range permissions {
		if perm.Category == "" {
			uncategorized = append(uncategorized, convertors.BuildPermissionResponse(perm))
		} else {
			if _, exists := categorized[perm.Category]; !exists {
				categorized[perm.Category] = []*userProto.Permission{}
			}
			categorized[perm.Category] = append(categorized[perm.Category], convertors.BuildPermissionResponse(perm))
		}
	}

	var response []*userProto.Permission
	for category, perms := range categorized {
		response = append(response, &userProto.Permission{
			Name:        category,
			DisplayName: category,
			Description: "This category has permissions related to " + category,
			Permissions: perms,
			CreatedAt:   timestamppb.Now(),
			UpdatedAt:   timestamppb.Now(),
		})
	}

	for _, perm := range uncategorized {
		response = append(response, perm)
	}

	return response
}
