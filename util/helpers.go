package util

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

type SignedDetails struct {
	FirstName string
	LastName  string
	Email     string
	UserType  string
	Uuid      string
	jwt.StandardClaims
}

func GenerateAllToken(firstName, lastName, email, userType, uuid string) (token string, refreshToken string, err error) {
	claims := &SignedDetails{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		UserType:  userType,
		Uuid:      uuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, err
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, err
	}

	return claims, err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func VerifyPassword(hashedPass string, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))
	if err != nil {
		return err
	}
	return nil
}

func CheckUserType(r *http.Request, role string) error {
	userType := r.Header.Get("user_type")

	if userType != role {
		return fmt.Errorf("Unauthorized to access this resource")
	}

	return nil
}

func MatchUserTypeToUUID(r *http.Request, claimsUUID string) error {
	userType := r.Header.Get("user_type")
	uuid := r.Header.Get("uuid")

	if userType != "ADMIN" && uuid != claimsUUID {
		return fmt.Errorf("Unauthorized to access this resource")
	}
	err := CheckUserType(r, userType)
	return err
}
