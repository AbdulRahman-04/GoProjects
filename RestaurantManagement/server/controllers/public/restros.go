package public

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/AbdulRahman-04/GoProjects/RestaurantManagement/server/config"
	"github.com/AbdulRahman-04/GoProjects/RestaurantManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/RestaurantManagement/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var restroCollection *mongo.Collection

func RestroCollect() {
	restroCollection = utils.MongoClient.Database("Restro_Management").Collection("restros")
}

// get jwt key and url
var JwtRestroKey = []byte(config.AppConfig.JWTKEY)
var URL = config.AppConfig.URL

func GenerateToken(length int) string {
	d := make([]byte, length)
	_, _ = rand.Read(d)
	return hex.EncodeToString(d)
}

// restro signup func
func RestroSignUp(c *gin.Context) {
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// type struct
	type RestroSignUp struct {
		RestroName          string  `json:"restroname" form:"restroname"`
		RestroLocation      string  `json:"restrolocation" form:"restrolocation"` 
		OwnerName           string  `json:"ownername" form:"ownername"`
		OwnerEmail          string  `json:"owneremail" form:"owneremail"`
		Password            string  `json:"password" form:"password"`
		Phone               string  `json:"phone" form:"phone"`
		RestroLicenseNumber string  `json:"licensenumber" form:"licensenumber"`
	}

	// bind into json 
	var inputRestro RestroSignUp 
	if err := c.ShouldBindJSON(&inputRestro); err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid Request",
		})
		return
	}

	// validations 
	if inputRestro.OwnerEmail == "" || inputRestro.Password == "" || inputRestro.OwnerName == "" || inputRestro.RestroName == "" || inputRestro.RestroLocation == "" || inputRestro.Phone == "" || inputRestro.RestroLicenseNumber == "" {
		c.JSON(400, gin.H{
			"msg": "Invalid Request, Fill All Fields",
		})
		return
	}

	if !strings.Contains(inputRestro.OwnerEmail, "@"){
		c.JSON(400, gin.H{
			"msg": "Invalid email",
		})
		return
	}

	if len(inputRestro.Password) < 6 {
		c.JSON(400, gin.H{
			"msg": "Invalid pass length",
		})
		return
	}

	if len(inputRestro.Phone) < 10 {
		c.JSON(400, gin.H{
			"msg": "Invalid Phone number length",
		})
		return
	}

	// duplicate chekc 
	count, err := restroCollection.CountDocuments(ctx, bson.M{"owneremail": inputRestro.OwnerEmail})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid db error",
		})
		return
	}

	if count > 0 {
		c.JSON(400, gin.H{
			"msg": "Restro Already Exists, pls login!❌",
		})
		return
	}

	// hash pass 
	hashPass, err := bcrypt.GenerateFromPassword([]byte(inputRestro.Password), 10)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "couldnt hash password",
		})
		return
	}

	emailToken := GenerateToken(8)
	phoneToken := GenerateToken(8)

	// create new var and push into db 
	var newRestro models.Restro

	newRestro.ID = primitive.NewObjectID()
	newRestro.Role = "owner"
	newRestro.RestroName = inputRestro.RestroName
	newRestro.OwnerEmail = inputRestro.OwnerEmail
	newRestro.Password = string(hashPass)
	newRestro.Phone = inputRestro.Phone
	newRestro.RestroLocation = inputRestro.RestroLocation
	newRestro.OwnerVerified.Email = false 
	newRestro.OwnerVerifyToken.Email = emailToken
	newRestro.OwnerVerifyToken.Phone = phoneToken
    newRestro.CreatedAt = time.Now()
	newRestro.UpdatedAt = time.Now()

	// send email func 
	go func(){

		emailData := utils.EmailData{
			From: "Team RestroManagement",
			To: inputRestro.OwnerEmail,
			Subject: "Email Verification",
			Html: fmt.Sprintf(`<a href="%s/api/public/restro/%s">Verify Email</a>`, URL, emailToken),
		}

		_ = utils.SendEmail(emailData)

	}()

	// push into db 
	_, err = restroCollection.InsertOne(ctx, newRestro)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "db error",
		})
		return
	}

	c.JSON(200, gin.H{
		"Msg": "Restaurant Signed Up🎉✅",
	})
}

// Email verify api 
func RestroEmailVerify(c*gin.Context){
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// take token from param
	token := c.Param("token")

	// check token ind db
	var restro models.Restro
	err := restroCollection.FindOne(ctx, bson.M{"ownerVerifyToken.emailVerifyToken":token}).Decode(&restro)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid token",
		})
		return
	}

	// check if link hasn't been clicked more than once 
	if restro.OwnerVerified.Email {
		c.JSON(200, gin.H{
			"msg": "Restro Owner Email Already Verified",
		})
		return
	}

	// update 
	update := bson.M{
		"$set": bson.M{
			"ownerVerified.emailVerified": true,
			"ownerVerifyToken.emailVerifyToken": nil,
			"updated_at": time.Now(),
		}}
	// update db 
	_, err = restroCollection.UpdateByID(ctx, restro.ID, update)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Couldn't update db",
		})
		return
	}	

	c.JSON(200, gin.H{
		"msg": "Restro Owner Email Verfied🎉, Go login!",
	})
}

// signin 
func RestroOwnerSignin(c*gin.Context){
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// type struct 
	type RestroOwnerSignin struct {
		Email     string     `form:"owneremail" json:"owneremail"`
		Password    string   `form:"password" json:"password"`       
	}

	// bind into json
	var inputRestro RestroOwnerSignin
	if err := c.ShouldBindJSON(&inputRestro); err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid request",
		})
		return
	}

	// validations 
	if inputRestro.Email == "" || inputRestro.Password == "" {
		c.JSON(400, gin.H{
			"msg": "invalid emal and password",
		})
		return
	}

	if !strings.Contains(inputRestro.Email, "@"){
		c.JSON(400, gin.H{
			"msg": "invalid email",
		})
		return
	}

	if len(inputRestro.Password) < 6 {
		c.JSON(400, gin.H{
			"msg": "invalid password length",
		})
		return
	}

	// find user email in db 
	var restro models.Restro
	err := restroCollection.FindOne(ctx, bson.M{"owneremail": inputRestro.Email}).Decode(&restro)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "No email found!",
		})
		return
	}

	// check password 
	err = bcrypt.CompareHashAndPassword([]byte(restro.Password), []byte(inputRestro.Password))
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid Password",
		})
		return
	}

	// token generate 
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": restro.ID,
		"role": restro.Role,
		"owneremail": restro.OwnerEmail,
		"exp": time.Now().Add(6*time.Hour).Unix(),
	}).SignedString(JwtRestroKey)
	
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Token Generation Failed",
		})
		return
	}
	
	c.JSON(200, gin.H{
		"msg": "Logged in successfully!✅", "token": token})
}

//change password 
func RestroOwnerChangePass(c*gin.Context){
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// type struct 
	type RestroOwnerChangePass struct {
		Email string         `json:"owneremail" form:"owneremail"`
		Oldpassword string   `json:"oldpassword" form:"oldpassword"` 
		Newpassword string   `json:"newpassword" form:"newpassword"` 
	}

	// bind into json
	var inputRestro RestroOwnerChangePass
	if err := c.ShouldBindJSON(&inputRestro); err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid request",
		})
		return
	}

	// validations 
	if inputRestro.Email == "" || inputRestro.Oldpassword == "" || inputRestro.Newpassword == "" {
		c.JSON(400, gin.H{
			"msg": "Fill all fields",
		})
		return
	}

	if !strings.Contains(inputRestro.Email, "@"){
		c.JSON(400, gin.H{
			"msg": "Invalid email",
		})
		return
	}

	if len(inputRestro.Newpassword) < 6 {
		c.JSON(400, gin.H{
			"msg": "Invalid new pass length",
		})
		return
	}

	// find email in db 
	var restro models.Restro
	err := restroCollection.FindOne(ctx, bson.M{"owneremail": inputRestro.Email}).Decode(&restro)

	if err != nil {
		c.JSON(400, gin.H{
			"msg": "no email found!",
		})
		return
	}

	// compare old pass 
	err = bcrypt.CompareHashAndPassword([]byte(restro.Password), []byte(inputRestro.Oldpassword))
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid password",
		})
		return
	}

	// hash the new pass 
	hashPass, err := bcrypt.GenerateFromPassword([]byte(inputRestro.Newpassword), 10)
	if err != nil {
      c.JSON(400, gin.H{
			"msg": "couldn't hash pass",
		})
		return
	}

	// db m update 
	update := bson.M{
		"$set": bson.M{
			"password": string(hashPass),
			"updated_at": time.Now(),
		}}

	// update db 
	_ , err = restroCollection.UpdateByID(ctx, restro.ID, update)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "DB ERROR",
		})
		return
	}	

	c.JSON(200, gin.H{
		"msg": "Password Changed successfully!✅",
	})
}

// forgot pass api 
func RestroOwnerForgotPass(c*gin.Context){
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// type struct 
	type RestroOwnerForgotPass struct {
		Email string `json:"owneremail" form:"owneremail"`
	}

	// bind into json
	var inputRestro RestroOwnerForgotPass
	if err := c.ShouldBindJSON(&inputRestro); err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid Request",
		})
		return
	}

	// validations 
	if !strings.Contains(inputRestro.Email, "@") || inputRestro.Email == "" {
		c.JSON(400, gin.H{
			"msg": "Invalid email",
		})
		return
	}

	// find email in db 
	var restro models.Restro
	err := restroCollection.FindOne(ctx, bson.M{"owneremail": inputRestro.Email}).Decode(&restro)

	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid email not found",
		})
		return
	}

	// temp pass generate 
	tempPass := GenerateToken(8)
	hashPass, err := bcrypt.GenerateFromPassword([]byte(tempPass), 10)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "hash generation failed",
		})
		return
	}

	// update the db 
	update := bson.M{
		"$set": bson.M{
			"password": string(hashPass),
			"updated_at": time.Now(),
		}}
	//	update 
	_, err = restroCollection.UpdateByID(ctx, restro.ID, update)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "db error",
		})
		return
	}

	// send email 
	emailData := utils.EmailData {
			From: "Team RestroManagement",
			To: inputRestro.Email,
			Subject: "Email Verification",
			Html: fmt.Sprintf(`<h2>Your Temporary Password is <strong>%s</strong></h2>`, tempPass),
	}
	_ = utils.SendEmail(emailData)

	c.JSON(200, gin.H{
		"msg": "Temporary Password sent to ur email, pls login anc change ur password⚠️",
	})
}