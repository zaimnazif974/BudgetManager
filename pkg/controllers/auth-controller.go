package controllers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth/gothic"
	"github.com/zaimnazif974/budgeting-BE/pkg/config"
	"github.com/zaimnazif974/budgeting-BE/pkg/models"
	"github.com/zaimnazif974/budgeting-BE/pkg/utils"
)

func GoogleLogin(w http.ResponseWriter, r *http.Request) {

	// Getting user data from google
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Failed to login using google")
		return
	}

	// search or add user to database
	var appUser models.User

	db := config.GetDB()

	query := map[any]string{
		"Name":      user.Name,
		"Provider":  user.Provider,
		"Email":     user.Email,
		"FirstName": user.FirstName,
		"LastName":  user.LastName,
	}

	//Search or create
	db.Find(&appUser, query).FirstOrCreate(&appUser)

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &utils.JWTClaims{
		UserID: appUser.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString("github.com/golang-jwt/jwt/v5")

	utils.ResponseJSON(w, http.StatusAccepted, appUser, "Login successfully")
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.ResponseJSON(w, http.StatusAccepted, tokenString, "Loggin google sucessfully")
}

func GoogleLogout(w http.ResponseWriter, r *http.Request) {

	err := gothic.Logout(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Failed to login using google")
		return
	}

	// utils.ResponseJSON(w, http.StatusAccepted, userLoggedOut, "Logout successfully")
}
