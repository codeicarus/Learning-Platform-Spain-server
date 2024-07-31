package helper

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"

	"reflect"
	"time"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func GetConfigVars() map[string]string {
	argsWithoutProg := os.Args[1:]
	env := argsWithoutProg[0]
	var config_vars = make(map[string]string)
	// env := "LOCAL"

	// if env == "PROD1" {
	// 	config_vars["ENV"] = "PROD"
	// 	config_vars["BBDD_HOST"] = "localhost:27017"
	// 	config_vars["BBDD_SOURCE"] = "admin"
	// 	config_vars["BBDD_DB"] = "opoprision"
	// 	config_vars["BBDD_USER"] = "opoprision_admin"
	// 	config_vars["BBDD_PASS"] = "MhkI6JP1mOs0s19F"
	// 	config_vars["LISTEN_PORT"] = "8080"
	// 	config_vars["URL"] = "https://www.penitenciarios.com/"
	// 	config_vars["PAYCOMET_PASSWORD"] = "yK3Cj71oAEmOZSMdUigY"
	// 	config_vars["PAYCOMET_MERCHANTCODE"] = "rq8jgpd2"
	// 	config_vars["PAYCOMET_TERMINAL"] = "17605"
	// 	config_vars["PAYCOMET_OPERATION"] = "1"
	// 	config_vars["PAYCOMET_LANGUAGE"] = "ES"
	// 	config_vars["PAYCOMET_CURRENCY"] = "EUR"
	// 	config_vars["PAYCOMET_URL"] = "https://api.paycomet.com/gateway/ifr-bankstore"
	// 	config_vars["PAYCOMET_PRODUCTDESCRIPTION"] = "Curso de preguntas tipo test y ejecicios prácticos"
	// 	config_vars["PAYCOMET_DESCRIPTOR"] = "PENITENCIARIOS.COM"
	// } else
	if env == "DEV" {
		config_vars["ENV"] = "DEV"
		config_vars["BBDD_HOST"] = "localhost:27017"
		config_vars["BBDD_SOURCE"] = "admin"
		config_vars["BBDD_DB"] = "opoprision_dev"
		config_vars["BBDD_USER"] = "opoprision_dev_admin"
		config_vars["BBDD_PASS"] = "MhkI6JP1mOs0s19FDEV"
		config_vars["LISTEN_PORT"] = "3004"
		config_vars["URL"] = "https://dev.api.penitenciarios.com/"
		config_vars["URLF"] = "https://dev.penitenciarios.com/"

		config_vars["PAYCOMET_PASSWORD"] = "yK3Cj71oAEmOZSMdUigY"
		config_vars["PAYCOMET_MERCHANTCODE"] = "rq8jgpd2"
		config_vars["PAYCOMET_TERMINAL"] = "17605"
		config_vars["PAYCOMET_OPERATION"] = "1"
		config_vars["PAYCOMET_LANGUAGE"] = "ES"
		config_vars["PAYCOMET_CURRENCY"] = "EUR"
		config_vars["PAYCOMET_URL"] = "https://api.paycomet.com/gateway/ifr-bankstore"
		config_vars["PAYCOMET_PRODUCTDESCRIPTION"] = "Curso de preguntas tipo test y ejecicios prácticos"
		config_vars["PAYCOMET_DESCRIPTOR"] = "PENITENCIARIOS.COM - DEV"

	} else if env == "LOCAL" {
		config_vars["ENV"] = "LOCAL"
		config_vars["BBDD_HOST"] = "127.0.0.1:27017"
		config_vars["BBDD_SOURCE"] = "admin"
		config_vars["BBDD_DB"] = "opoprision"
		config_vars["BBDD_USER"] = "opoprision_admin"
		config_vars["BBDD_PASS"] = "MhkI6JP1mOs0s19F"
		config_vars["LISTEN_PORT"] = "8080"
		config_vars["URL"] = "http://127.0.0.1:8080/"
		config_vars["URLF"] = "http://127.0.0.1:3000/"

		config_vars["PAYCOMET_PASSWORD"] = "6vAJohBN4w8cnOgKUfTG"
		config_vars["PAYCOMET_MERCHANTCODE"] = "j5a9ptat"
		config_vars["PAYCOMET_TERMINAL"] = "17432"

		config_vars["PAYCOMET_OPERATION"] = "1"
		config_vars["PAYCOMET_LANGUAGE"] = "ES"
		config_vars["PAYCOMET_CURRENCY"] = "EUR"
		config_vars["PAYCOMET_URL"] = "https://api.paycomet.com/gateway/ifr-bankstore"
		config_vars["PAYCOMET_PRODUCTDESCRIPTION"] = "Curso de preguntas tipo test y ejecicios prácticos"
		config_vars["PAYCOMET_DESCRIPTOR"] = "PENITENCIARIOS.COM - LOCAL"
	} else if env == "PROD" {
		config_vars["ENV"] = "PROD"
		config_vars["BBDD_HOST"] = "localhost:27017"
		config_vars["BBDD_SOURCE"] = "admin"
		config_vars["BBDD_DB"] = "opoprision"
		config_vars["BBDD_USER"] = "opoprision_admin"
		config_vars["BBDD_PASS"] = "MhkI6JP1mOs0s19F"
		config_vars["LISTEN_PORT"] = "8080"
		config_vars["URL"] = "https://www.api.penitenciarios.com/"
		config_vars["URLF"] = "https://penitenciarios.com/"

		config_vars["PAYCOMET_PASSWORD"] = "yK3Cj71oAEmOZSMdUigY"
		config_vars["PAYCOMET_MERCHANTCODE"] = "rq8jgpd2"
		config_vars["PAYCOMET_TERMINAL"] = "17605"
		config_vars["PAYCOMET_OPERATION"] = "1"
		config_vars["PAYCOMET_LANGUAGE"] = "ES"
		config_vars["PAYCOMET_CURRENCY"] = "EUR"
		config_vars["PAYCOMET_URL"] = "https://api.paycomet.com/gateway/ifr-bankstore"
		config_vars["PAYCOMET_PRODUCTDESCRIPTION"] = "Curso de preguntas tipo test y ejecicios prácticos"
		config_vars["PAYCOMET_DESCRIPTOR"] = "PENITENCIARIOS.COM"
	}

	return config_vars
}

func CheckUser(email interface{}) bool {
	if email == "djkruske@gmail.com" || email == "avilmor2@gmail.com" || email == "saotomeamelia@gmail.com" || email == "merchanhelianny@gmail.com" || email == "nanitadr18@gmail.com" || email == "zaichemical@gmail.com" || email == "toledof764@gmail.com" {
		return true
	} else {
		return false
	}
}
func IssetIndexArray(arr []string, index int) bool {
	return (len(arr) > index)
}
func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func SendEmail(subject string, msg string, to_email string, from_email string) bool {
	from := mail.Address{"PENITENCIARIOS.COM", "penitenciarios@penitenciarios.com"}
	// from := mail.Address{"PENITENCIARIOS.COM", "info.socialite.app@gmail.com"}

	replyTo := mail.Address{"", from_email}
	to := mail.Address{"", to_email}
	subj := subject
	body := msg

	// Setup headers
	headers := make(map[string]string)
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html"
	headers["charset"] = "UTF-8"
	headers["From"] = from.String()
	headers["Reply-To"] = replyTo.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := "smtp.gmail.com:465"

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", "infopulbong@gmail.com", "sxesimmkfbkxfbfc", host)
	// auth := smtp.PlainAuth("", "info.socialite.app@gmail.com", "xuyrwrievlducusb", host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Println("Ocurrio un error TCP")
		// log.Panic(err)
		return false
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Println("Ocurrio un error newClient")
		// log.Panic(err)
		return false
	}

	// Auth

	if err = c.Auth(auth); err != nil {
		log.Println("Ocurrio un error Auth")
		// log.Panic(err)
		return false
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Println("Ocurrio un error Mail")
		// log.Panic(err)
		return false
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Println("Ocurrio un error Rcpt")
		// log.Panic(err)
		return false
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Println("Ocurrio un error Data")
		// log.Panic(err)
		return false
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Println("Ocurrio un error Write")
		// log.Panic(err)
		return false
	}

	err = w.Close()
	if err != nil {
		log.Println("Ocurrio un error Close")
		// log.Panic(err)
		return false
	}

	c.Quit()
	return true
}

func CreateKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}
func MakeTimestamp() int64 {
	return time.Now().Unix()
}
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
