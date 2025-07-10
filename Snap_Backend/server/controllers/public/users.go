package public

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/config"
	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/models"
	"github.com/AbdulRahman-04/GoProjects/Snap_Backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

func UserCollect() {
	userCollection = utils.MongoClient.Database("Snap_Backend_Practice").Collection("users")
}

var userKey = []byte(config.AppConfig.JWTKEY)
var URL = config.AppConfig.URL

func generateToken(length int) string {
	d := make([]byte, length)
	_, _ = rand.Read(d)
	return hex.EncodeToString(d)
}

func UserSignup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	type UserSignup struct {
		UserName string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
		Age      int    `json:"age"`
		Location string `json:"location"`
	}

	var inputUser UserSignup
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(400, gin.H{"msg": "invalid request"})
		return
	}

	if inputUser.UserName == "" || inputUser.Email == "" || inputUser.Password == "" || inputUser.Phone == "" || inputUser.Age == 0 || inputUser.Location == ""  {
		c.JSON(400, gin.H{"msg": "fill all fields"})
		return
	}

	if !strings.Contains(inputUser.Email, "@") {
		c.JSON(400, gin.H{"msg": "invalid email"})
		return
	}

	if len(inputUser.Password) < 6 {
		c.JSON(400, gin.H{"msg": "invalid password length"})
		return
	}

	if len(inputUser.Phone) < 10 {
		c.JSON(400, gin.H{"msg": "invalid phone number"})
		return
	}

	count, _ := userCollection.CountDocuments(ctx, bson.M{"email": inputUser.Email})
	if count > 0 {
		c.JSON(400, gin.H{"msg": "user already exists"})
		return
	}

	hashPass, _ := bcrypt.GenerateFromPassword([]byte(inputUser.Password), 10)
	emailToken := generateToken(8)
	phoneToken := generateToken(8)

	var newUser models.User
	newUser.ID = primitive.NewObjectID()
	newUser.Role = "user"
	newUser.UserName = inputUser.UserName
	newUser.Email = inputUser.Email
	newUser.Password = string(hashPass)
	newUser.Phone = inputUser.Phone
	newUser.Age = inputUser.Age
	newUser.Location = inputUser.Location
	newUser.UserVerified.Email = false
	newUser.UserVerifyToken.Email = emailToken
	newUser.UserVerifyToken.Phone = phoneToken
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	go func() {
		emailData := utils.EmailData{
			From:    "Team Snap",
			To:      inputUser.Email,
			Subject: "Email Verification",
			Html:    fmt.Sprintf(`<a href="%s/api/public/user/emailverify/%s">Verify Email</a>`, URL, emailToken),
		}
		_ = utils.SendEmail(emailData)
	}()

	_, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(500, gin.H{"msg": "DB insert error"})
		return
	}

	c.JSON(200, gin.H{"msg": "User signed up successfully, verify email"})
}

func UserEmailVerify(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token := c.Param("token")

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"userVerifyToken.emailVerifyToken": token}).Decode(&user)
	if err != nil {
		c.JSON(400, gin.H{"msg": "invalid email verify token"})
		return
	}

	if user.UserVerified.Email {
		c.JSON(400, gin.H{"msg": "Email already verified"})
		return
	}

	update := bson.M{"$set": bson.M{
		"userVerified.email": true,
		"userVerifyToken.email": nil,
		"updated_at": time.Now(),
	}}

	_, err = userCollection.UpdateByID(ctx, user.ID, update)
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error while verifying"})
		return
	}

	c.JSON(200, gin.H{"msg": "Email verified successfully ✅"})
}

func UserSignIn(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type SignIn struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input SignIn
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"msg": "invalid request"})
		return
	}

	if input.Email == "" || input.Password == "" {
		c.JSON(400, gin.H{"msg": "fill all fields"})
		return
	}

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		c.JSON(400, gin.H{"msg": "user not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(400, gin.H{"msg": "invalid password"})
		return
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID.Hex(),
		"role":  user.Role,
		"email": user.Email,
		"exp":   time.Now().Add(5 * time.Hour).Unix(),
	}).SignedString(userKey)
	if err != nil {
		c.JSON(400, gin.H{"msg": "token error"})
		return
	}

	c.JSON(200, gin.H{"msg": "User logged in successfully", "token": token})
}

func UserChangePass(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	type ChangePass struct {
		Email       string `json:"email"`
		OldPassword string `json:"oldpassword"`
		NewPassword string `json:"newpassword"`
	}

	var input ChangePass
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"msg": "invalid request"})
		return
	}

	if input.Email == "" || input.OldPassword == "" || input.NewPassword == "" {
		c.JSON(400, gin.H{"msg": "fill all fields"})
		return
	}

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		c.JSON(400, gin.H{"msg": "user not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword))
	if err != nil {
		c.JSON(400, gin.H{"msg": "invalid old password"})
		return
	}

	hashNewPass, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 10)
	update := bson.M{"$set": bson.M{"password": string(hashNewPass), "updated_at": time.Now()}}
	_, err = userCollection.UpdateByID(ctx, user.ID, update)
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}

	c.JSON(200, gin.H{"msg": "Password changed successfully ✅"})
}

func UserForgotPass(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	type ForgotPass struct {
		Email string `json:"email"`
	}

	var input ForgotPass
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"msg": "invalid request"})
		return
	}

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		c.JSON(400, gin.H{"msg": "user not found"})
		return
	}

	tempPass := generateToken(8)
	hashTempPass, _ := bcrypt.GenerateFromPassword([]byte(tempPass), 10)
	update := bson.M{"$set": bson.M{"password": string(hashTempPass), "updated_at": time.Now()}}
	_, err = userCollection.UpdateByID(ctx, user.ID, update)
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}

	emailData := utils.EmailData{
		From:    "Team Snap",
		To:      input.Email,
		Subject: "Temporary Password",
		Html:    fmt.Sprintf(`<h2>Password Reset</h2><p>Dear %s, your temporary password is: <strong>%s</strong></p><p>Please log in and change it immediately.</p>`, user.UserName, tempPass),
	}
	_ = utils.SendEmail(emailData)

	c.JSON(200, gin.H{"msg": "Temporary password sent to email ✅"})
}