package public

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/config"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var adminCollection *mongo.Collection

func AdminCollect() {
	adminCollection = utils.MongoClient.Database("Event_Booking").Collection("admin")
}

// jwt key and url
var adminJwtKey = []byte(config.AppConfig.JWTKEY)
var adminUrl = config.AppConfig.URL

// generate token func
func GenerateToken(length int) string {
	d := make([]byte, length)
	_, _ = rand.Read(d)
	return hex.EncodeToString(d)
}

// admin signup api
func AdminSignUp(c *gin.Context) {
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// type struct
	type AdminSignUp struct {
		AdminName string `json:"name" form:"name"`
		Email     string `json:"email" form:"email"`
		Password  string `json:"password" form:"password"`
		Phone     string `json:"phone" form:"phone"`
		Language  string `json:"language" form:"language"`
		Location  string `json:"location" form:"location"`
	}

	// bind into json
	var inputAdmin AdminSignUp
	if err := c.ShouldBindJSON(&inputAdmin); err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid Request",
		})
		return
	}

	// validations
	if inputAdmin.AdminName == "" || inputAdmin.Email == "" || inputAdmin.Password == "" || inputAdmin.Phone == "" || inputAdmin.Language == "" || inputAdmin.Location == "" {
		c.JSON(400, gin.H{
			"msg": "Invalid Request, please fill all fields⚠️",
		})
		return
	}

	if !strings.Contains(inputAdmin.Email, "@") {
		c.JSON(400, gin.H{
			"msg": "Invalid email",
		})
		return
	}

	if len(inputAdmin.Password) < 6 {
		c.JSON(400, gin.H{
			"msg": "Invalid password length",
		})
		return
	}

	if len(inputAdmin.Phone) < 10 {
		c.JSON(400, gin.H{
			"msg": "Invalid phone number length",
		})
		return
	}

	// duplicate chekc
	count, err := adminCollection.CountDocuments(ctx, bson.M{"email": inputAdmin.Email})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid db error",
		})
		return
	}

	if count > 0 {
		c.JSON(400, gin.H{
			"msg": "Admin Already Exists, Please Go Login⚠️",
		})
		return
	}

	// hash the pass
	hashPass, err := bcrypt.GenerateFromPassword([]byte(inputAdmin.Password), 10)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "couldn't hash pass",
		})
		return
	}

	emailToken := GenerateToken(8)
	phoneToken := GenerateToken(8)

	// create new var and send email
	var newAdmin models.Admin

	newAdmin.ID = primitive.NewObjectID()
	newAdmin.Role = "admin"
	newAdmin.AdminName = inputAdmin.AdminName
	newAdmin.Email = inputAdmin.Email
	newAdmin.Password = string(hashPass)
	newAdmin.Location = inputAdmin.Location
	newAdmin.Language = inputAdmin.Language
	newAdmin.Phone = inputAdmin.Phone
	newAdmin.AdminVerified.Email = false
	newAdmin.AdminVerifyToken.Email = emailToken
	newAdmin.AdminVerifyToken.Phone = phoneToken
	newAdmin.CreatedAt = time.Now()
	newAdmin.UpdatedAt = time.Now()

	// send email
	go func() {

		emailData := utils.EmailData{
			From:    "Team Ivents Plannerz🎉",
			To:      inputAdmin.Email,
			Subject: "Email Verification",
			Html:    fmt.Sprintf(`<a href="%s/api/public/user/emailverify/%s">Verify email</a>`, adminUrl, emailToken),
		}

		_ = utils.SendEmail(emailData)

	}()

	// push into db
	_, err = adminCollection.InsertOne(ctx, newAdmin)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Db error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Admin Signed Up🎉, Verify Your Email and then login✅",
	})

}

// email verify api
func EmailVerifyAdmin(c *gin.Context) {
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token := c.Param("token")

	// compare token
	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"adminverifytoken.emailVerifyToken": token}).Decode(&admin)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Db error",
		})
		return
	}

	if admin.AdminVerified.Email {
		c.JSON(200, gin.H{
			"msg": "Email Verified already, u can login now!",
		})
		return
	}

	// update
	update := bson.M{
		"$set": bson.M{
			"adminverified.emailVerified": true,
			"adminVerifyToken.emailVerifyToken": nil,
			"updated_at":                        time.Now(),
		}}

	// update db
	_, err = adminCollection.UpdateByID(ctx, admin.ID, update)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "db error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "email Verified✨🙌",
	})
}

// admin signin api
func AdminSignIn(c *gin.Context) {
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// type struct
	type AdminSignIn struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}

	var inputAdmin AdminSignIn
	if err := c.ShouldBindJSON(&inputAdmin); err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid request",
		})
		return
	}

	// validations
	if inputAdmin.Email == "" || inputAdmin.Password == "" {
		c.JSON(400, gin.H{
			"msg": "fill all fields!❌",
		})
		return
	}
	if !strings.Contains(inputAdmin.Email, "@") {
		c.JSON(400, gin.H{
			"msg": "invalid email",
		})
		return
	}
	if len(inputAdmin.Password) < 6 {
		c.JSON(400, gin.H{
			"msg": "invalid password length",
		})
		return
	}

	// find admin ind db
	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"email": inputAdmin.Email}).Decode(&admin)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "no email found!",
		})
		return
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(inputAdmin.Password))
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid password",
		})
		return
	}

	// token generation
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    admin.ID,
		"role":  admin.Role,
		"email": admin.Email,
		"exp":   time.Now().Add(5 * time.Hour).Unix(),
	}).SignedString(adminJwtKey)

	if err != nil {
		c.JSON(400, gin.H{
			"msg": "token generation failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Logged in successfull!✨", "token": token})
}

// change pass api
func AdminChangePass(c *gin.Context) {
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// type struct
	type AdminSignIn struct {
		Email       string `json:"email" form:"email"`
		Oldpassword string `json:"oldpassword" form:"oldpassword"`
		Newpassword string `json:"newpassword" form:"newpassword"`
	}

	// bind into json
	var inputAdmin AdminSignIn
	if err := c.ShouldBindJSON(&inputAdmin); err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid request",
		})
		return
	}

	// validations
	if inputAdmin.Email == "" || inputAdmin.Oldpassword == "" || inputAdmin.Newpassword == "" {
		c.JSON(400, gin.H{
			"msg": "invalid request, fill all fields",
		})
		return
	}

	if !strings.Contains(inputAdmin.Email, "@") {
		c.JSON(400, gin.H{
			"msg": "invalid email",
		})
		return
	}

	if len(inputAdmin.Newpassword) < 6 {
		c.JSON(400, gin.H{
			"msg": "invalid new pass length",
		})
		return
	}

	// find email in db
	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"email": inputAdmin.Email}).Decode(&admin)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid request, no email found",
		})
		return
	}

	// compare old pass
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(inputAdmin.Oldpassword))
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid old password",
		})
		return
	}

	// hash new pass
	hashPass, err := bcrypt.GenerateFromPassword([]byte(inputAdmin.Newpassword), 10)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "hashing failed",
		})
		return
	}

	// update db
	update := bson.M{
		"$set": bson.M{
			"password":   string(hashPass),
			"updated_at": time.Now(),
		}}
	// update db
	_, err = adminCollection.UpdateByID(ctx, admin.ID, update)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid db error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Password Changed Successfully!✅",
	})

}

// forgot passs api 
func AdminForgotPass(c*gin.Context){
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// type struct
	type AdminForgotPass struct {
		Email       string `json:"email" form:"email"`
	}

	// bind 
	var inputAdmin AdminForgotPass
	if err := c.ShouldBindJSON(&inputAdmin); err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid request",
		})
		return
	}

	// validations 
	if !strings.Contains(inputAdmin.Email, "@") || inputAdmin.Email == "" {
		c.JSON(400, gin.H{
			"msg": "invalid Email",
		})
		return
	}

	// find user in db 
	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"email": inputAdmin.Email}).Decode(&admin)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid db error",
		})
		return
	}

	// generate temp pass 
	tempPass := GenerateToken(8)
	hashNewPass, err := bcrypt.GenerateFromPassword([]byte(tempPass), 10)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid request, couldnt hash pass",
		})
		return
	}

	// update the db 
    update := bson.M{
		"$set": bson.M{
			"password": string(hashNewPass),
			"updated_at": time.Now(),
		}}

	// update db
	 _, err  =adminCollection.UpdateByID(ctx, admin.ID, update)
	 if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid db error",
		})
		return
	 }	

	// email 
	emailData := utils.EmailData {
		From:    "Team Ivents Plannerz🎉",
			To:      inputAdmin.Email,
			Subject: "Email Verification",
			Html:    fmt.Sprintf(`<h2>Your Temporary password is <strong>%s</strong></h2>`, tempPass),
	} 

	_ = utils.SendEmail(emailData)

	c.JSON(200, gin.H{
		"msg": "Temporary email sent to ur mail✅✨",
	})
}