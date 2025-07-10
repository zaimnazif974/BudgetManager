package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/markbates/goth/gothic"
	"github.com/zaimnazif974/budgeting-BE/pkg/config"
	"github.com/zaimnazif974/budgeting-BE/pkg/models"
	"github.com/zaimnazif974/budgeting-BE/pkg/utils"
	"golang.org/x/crypto/bcrypt"
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
