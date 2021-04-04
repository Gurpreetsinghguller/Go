package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type RegisterDetails struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type SigninDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type SigninResponse struct {
	Role  string `json:"role"`
	Token string `json:"token"`
}
type Error struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}
type Role struct {
	Role string `json:"role"`
}

var mySigningKey = []byte("generatetokenusingthiskey")

func AutoMigrate() {
	db := DbConn()
	db.AutoMigrate(&RegisterDetails{})
	defer db.Close()
}
func SetError(err Error, message string) Error {
	err.IsError = true
	err.Message = message
	return err
}
func DbConn() *gorm.DB {
	connection, err := gorm.Open("postgres", "postgres://admin:1234@localhost/login?sslmode=disable")
	if err != nil {
		log.Fatalln("wrong database url")
	}
	sqldb := connection.DB()
	err = sqldb.Ping()
	if err != nil {
		log.Fatal("database is disconnected")
	}
	return connection
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(email, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func IsAuthorized(handle func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There is an error")
				}
				return mySigningKey, nil
			})
			if err != nil {
				var err Error
				// err.IsError = true
				err = SetError(err, "Your Token has been expired")
				json.NewEncoder(w).Encode(err)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if token.Valid && claims["role"] == "user" && ok {
				r.Header.Set("Role", "User")
				handle(w, r)
			} else if token.Valid && claims["role"] == "admin" && ok {
				r.Header.Set("Role", "Admin")
				handle(w, r)
			}
		} else {

			w.Write([]byte("Error in Token"))
		}
	})
}
func Registeration(w http.ResponseWriter, r *http.Request) {
	bodydata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("Error in Reading body")
	}
	var registerationDetails RegisterDetails
	err = json.Unmarshal(bodydata, &registerationDetails)
	if err != nil {
		log.Fatalln("Error in Unmarshaling")
	}
	DB := DbConn()
	defer DB.Close()
	var checkuser RegisterDetails
	DB.Where("email = 	?", registerationDetails.Email).First(&checkuser)
	//checks email is alredy in use
	if checkuser.Email != "" {
		var err Error
		err = SetError(err, "Email already in use")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	registerationDetails.Password, err = HashPassword(registerationDetails.Password)
	if err != nil {
		log.Fatalln("Error in Password Hashing")
	}
	DB.Create(&registerationDetails)
	bytedata, err := json.MarshalIndent(registerationDetails, "", "  ")
	if err != nil {
		var err Error
		err = SetError(err, "Error in Marshaling")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytedata)
}

func Signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	DB := DbConn()
	defer DB.Close()
	var signinDetails SigninDetails
	var registerationDetails RegisterDetails
	json.NewDecoder(r.Body).Decode(&signinDetails)
	DB.Where("email =?", signinDetails.Email).First(&registerationDetails)
	if registerationDetails.Email == "" {
		var err Error
		err = SetError(err, "Invalid Email or Password")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	chkpass := CheckPasswordHash(signinDetails.Password, registerationDetails.Password)
	if !chkpass {
		var err Error
		err = SetError(err, "Invalid Email or Password")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	token, err := GenerateJWT(registerationDetails.Email, registerationDetails.Role)
	if err != nil {
		var err Error
		err = SetError(err, "Failed to generate token")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}
	var signinResponse SigninResponse
	signinResponse.Token = token
	signinResponse.Role = registerationDetails.Role
	pBytes, err := json.MarshalIndent(signinResponse, " ", " ")
	w.Write(pBytes)

}
func AdminDashboard(w http.ResponseWriter, r *http.Request) {

	fmt.Println("admin dashboad called")

	if r.Header.Get("Role") != "Admin" {
		fmt.Println("Unauthorized accesss")
		w.Write([]byte("Unauthorized accesss"))
		return
	}
	fmt.Println("admin dashboad")
	w.Write([]byte("This is admin Dashboard"))
}
func UserDashboard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("user dashboad called")

	if r.Header.Get("Role") != "User" {
		fmt.Println("Unauthorized accesss")
		w.Write([]byte("Unauthorized accesss"))
		return
	}
	fmt.Println("user dashboad")
	w.Write([]byte("This is User Dashboard"))
}

func main() {
	AutoMigrate()
	router := mux.NewRouter()
	router.HandleFunc("/api/Register", Registeration).Methods("POST")
	router.HandleFunc("/api/Signin", Signin).Methods("POST")
	router.Handle("/api/admin", IsAuthorized(AdminDashboard)).Methods("GET")
	router.Handle("/api/user", IsAuthorized(UserDashboard)).Methods("GET")
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})
	fmt.Println("Server started at http://localhost:8000")
	http.ListenAndServe(":8000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
}
