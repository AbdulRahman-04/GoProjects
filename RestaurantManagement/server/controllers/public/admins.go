package public

import (
	"context"
	// "crypto/rand"
	// "encoding/hex"
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

var adminCollection *mongo.Collection

func AdminCollect() {
	adminCollection = utils.MongoClient.Database("Restro_Management").Collection("admin")
}

// JWT KEY AND URL
var adminKey = []byte(config.AppConfig.JWTKEY)

var adminUrl = config.AppConfig.URL


// func generateAdminToken(length int) string {
// 	d := make([]byte, length)
// 	_, _ = rand.Read(d)
// 	return hex.EncodeToString(d)
// }

// Admin Signup Function
func AdminSignup(c *gin.Context) {
	// ctx, cancel
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// type struct
	type AdminSignup struct {
		AdminName string `json:"adminname" form:"adminname"`
		Email     string `json:"email" form:"email"`
		Password  string `json:"password" form:"password"`
		Phone     string `json:"phone" form:"phone"`
		Location  string `json:"location" form:"location"`
		Language  string `json:"language" form:"language"`
	}

	// bind into json and validate
	var inputAdmin AdminSignup

	if err := c.ShouldBindJSON(&inputAdmin); err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid request",
		})
		return
	}

	// validations
	if inputAdmin.AdminName == "" || inputAdmin.Email == "" || inputAdmin.Password == "" || inputAdmin.Phone == "" || inputAdmin.Location == "" || inputAdmin.Language == "" {
		c.JSON(400, gin.H{
			"msg": "invalid request, pls fill all fields",
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
			"msg": "invalid pass length, must be 6 minimum",
		})
		return
	}

	if len(inputAdmin.Phone) < 10 {
		c.JSON(400, gin.H{
			"msg": "invalid phone length, must be 10 or more",
		})
		return
	}

	// if inputAdmin.Age < 10 {
	// 	c.JSON(400, gin.H{
	// 		"msg": "invalid age, u need to be atleast 10 years old or more to use this app!🙌",
	// 	})
	// 	return
	// }

	// duplicate check
	count, err := adminCollection.CountDocuments(ctx, bson.M{"email": inputAdmin.Email})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "DB ERROR",
		})
		return
	}

	if count > 0 {
		c.JSON(400, gin.H{
			"msg": "admin already exists, pls login",
		})
		return
	}

	// hash the pass
	hashPass, err := bcrypt.GenerateFromPassword([]byte(inputAdmin.Password), 10)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "couldn't hhash pass",
		})
		return
	}

	emailToken := generateToken(8)
	phoneToken := generateToken(8)

	// create var and push values inside it
	var newAdmin models.Admin

	newAdmin.ID = primitive.NewObjectID()
	newAdmin.Role = "admin"
	newAdmin.AdminName = inputAdmin.AdminName
	newAdmin.Email = inputAdmin.Email
	newAdmin.Password = string(hashPass)
	newAdmin.Phone = inputAdmin.Phone
	newAdmin.Location = inputAdmin.Location
	newAdmin.Language = inputAdmin.Language
	newAdmin.AdminVerified.Email = false
	newAdmin.AdminVerifyToken.Email = emailToken
	newAdmin.AdminVerifyToken.Phone = phoneToken
	newAdmin.CreatedAt = time.Now()
	newAdmin.UpdatedAt = time.Now()

	// send email func
	go func() {
		emailData := utils.EmailData{
			From:    "Team RestroManagement",
			To:      inputAdmin.Email,
			Subject: "Email Verifiction",
			Html: fmt.Sprintf(`<a href="%s/api/public/admin/emailverify/%s">Verify ur email</a>`, adminUrl, emailToken),
		}
		_ = utils.SendEmail(emailData)
	}()

	// insert the admin into db
	_, err = adminCollection.InsertOne(ctx, newAdmin)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid db error while signing up admin",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Admin Signedup✨, Please verify ur email and then login",
	})
}

// email verify api
func EmailVerify(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token := c.Param("token")

	// check token indb
	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"adminVerifyToken.emailVerifyToken": token}).Decode(&admin)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid email verify token",
		})
		return
	}

	// check if link hasnt been clicked more than once
	if admin.AdminVerified.Email {
		c.JSON(400, gin.H{
			"msg": "Admin Email Already verified. you can go and login now✨",
		})
		return
	}

	// db update
	update := bson.M{
		"$set": bson.M{
			"adminVerified.emailVerified":       true,
			"adminVerifyToken.emailVerifyToken": nil,
			"updated_at":                        time.Now(),
		},
	}

	_, err = adminCollection.UpdateByID(ctx, admin.ID, update)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid db error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Admin Email Verified✨✅",
	})
}

// admin signin api
func AdminSignIn(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// type
	type AdminSignIn struct {
		Email    string `bson:"email" form:"email"`
		Password string `bson:"password" form:"password"`
	}

	// bindo into json
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
			"msg": "Please fill all fields",
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
			"msg": "password length must be 6",
		})
		return
	}

	// check if email exists in db
	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"email": inputAdmin.Email}).Decode(&admin)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid email❌, no admin found",
		})
		return
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(inputAdmin.Password))
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid Password!❌",
		})
		return
	}

	// jwt token generate
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    admin.ID.Hex(),
		"role":  admin.Role,
		"email": admin.Email,
		"exp":   time.Now().Add(5 * time.Hour).Unix(),
	}).SignedString(adminKey)

	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid jwt tokwn generation",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Admin Loggedin Successfully!✨", "token": token})
}

// change password
func AdminChangePass(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// type
	type AdminChangePass struct {
		Email       string `json:"email" form:"email"`
		Olpassword  string `json:"oldpassword" form:"oldpassword"`
		Newpassword string `json:"newpassword" form:"newpassword"`
	}

	// bind into json and validate
	var inputAdmin AdminChangePass
	if err := c.ShouldBindJSON(&inputAdmin); err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid request",
		})
		return
	}

	// validations
	if inputAdmin.Email == "" || inputAdmin.Newpassword == "" || inputAdmin.Olpassword == "" {
		c.JSON(400, gin.H{
			"msg": "fill all fields",
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

	// find user in db
	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"email": inputAdmin.Email}).Decode(&admin)

	if err != nil {
		c.JSON(400, gin.H{
			"msg": "no email found!❌",
		})
		return
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(inputAdmin.Olpassword))
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid Old Password",
		})
		return
	}

	hashNewPass, err := bcrypt.GenerateFromPassword([]byte(inputAdmin.Newpassword), 10)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "couldn't hash pass",
		})
		return
	}

	// update in db
	update := bson.M{
		"$set": bson.M{
			"password":   string(hashNewPass),
			"updated_at": time.Now(),
		}}

	_, err = adminCollection.UpdateByID(ctx, admin.ID, update)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid db error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Password changes successfully!✅",
	})
}

// forgot admin pass api
func AdminForgotPass(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// type
	type AdminForgotPass struct {
		Email string `json:"email" form:"email"`
	}

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
			"msg": "invalid email",
		})
		return
	}

	// find email in db
	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"email": inputAdmin.Email}).Decode(&admin)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid email, no admin found❌",
		})
		return
	}

	// generate tempass
	tempPass := generateToken(8)
	hashTempPass, err := bcrypt.GenerateFromPassword([]byte(tempPass), 10)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "couldn't hash pass",
		})
		return
	}

	// db update krdo
	update := bson.M{
		"$set": bson.M{
			"password":   string(hashTempPass),
			"updated_at": time.Now(),
		}}

	_, err = adminCollection.UpdateByID(ctx, admin.ID, update)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid db error",
		})
		return
	}

	// send pass on email
	emailData := utils.EmailData{
		From:    "Team RestroManagement",
		To:      inputAdmin.Email,
		Subject: "Temporary Password",
		Html:    fmt.Sprintf(`<h2>Password Reset</h2><p>Dear %s, your temporary password is: <strong>%s</strong></p><p>Please log in and change it immediately.</p>`, admin.AdminName, tempPass),
	}

	_ = utils.SendEmail(emailData)

	c.JSON(200, gin.H{
		"msg": "Temporary password sent to email ✅",
	})
}
