package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth/gothic"
	"github.com/zaimnazif974/budgeting-BE/pkg/config"
	"github.com/zaimnazif974/budgeting-BE/pkg/models"
	"github.com/zaimnazif974/budgeting-BE/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	// Getting user data from google
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		result := fmt.Sprintf("ComplateUserAuthL: %s", err)
		utils.WriteError(w, http.StatusUnauthorized, result)
		return
	}

	// search or add user to database
	var appUser models.User

	db := config.GetDB()

	searchQuery := models.User{
		Provider: user.Provider,
		Email:    user.Email,
	}

	newUser := models.User{
		Provider:  user.Provider,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	//Search or create
	db.Where(searchQuery).FirstOrCreate(&appUser, newUser)

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &utils.JWTClaims{
		UserID:     appUser.ID,
		AcessToken: user.AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET_KEY", "nil")))

	response := map[string]any{
		"user":  appUser,
		"token": tokenString,
	}
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	log.Printf("acess Token: %v", user.AccessToken)
	w.Header().Set("Content-Type", "application/json")
	utils.ResponseJSON(w, http.StatusAccepted, response, "Loggin google sucessfully")
}

func GoogleLogout(w http.ResponseWriter, r *http.Request) {

	err := gothic.Logout(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Failed to login using google")
		return
	}

	// utils.ResponseJSON(w, http.StatusAccepted, userLoggedOut, "Logout successfully")
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	utils.ParseBody(r, &user)
	db := config.GetDB()

	//validator
	if user.Email == "" {
		utils.WriteError(w, http.StatusBadRequest, "Email must be filled")
		return
	}
	if user.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "Password must be filled")
		return
	}

	var existingUser models.User

	//Check if user is exist
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		response := fmt.Sprintf("%s email has been registered", user.Email)
		utils.WriteError(w, http.StatusBadRequest, response)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to hash password")
		return
	}

	newUser := models.User{Email: user.Email, Password: string(hash)}
	db.Create(&newUser)

	utils.ResponseJSON(w, http.StatusCreated, newUser, "Sucessfully registered user")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var existingUser models.User

	utils.ParseBody(r, &user)

	db := config.GetDB()

	if user.Email == "" {
		utils.WriteError(w, http.StatusBadRequest, "Email must be filled")
		return
	}
	if user.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "Password must be filled")
		return
	}

	//Check if user is exist
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		utils.WriteError(w, http.StatusBadRequest, "User not Found")
		return
	}

	//Checking if the password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid email or password")
		return
	}

	//Creating token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &utils.JWTClaims{
		UserID: existingUser.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET_KEY", "nil")))

	response := map[string]any{
		"user":  existingUser,
		"token": tokenString,
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.ResponseJSON(w, http.StatusAccepted, response, "Loggin google sucessfully")
}

func BeginGoogleAuth(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func Logout(w http.ResponseWriter, r *http.Request) {

	gothic.Logout(w, r)

	utils.ResponseJSON(w, http.StatusAccepted, nil, "Logout sucessfully")
}
