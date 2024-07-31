package models

import (
	"errors"
	"strings"
	"test/helper"

	"github.com/globalsign/mgo/bson"
)

type Usuarios struct {
	Id                   bson.ObjectId `bson:"_id" json:"id"`
	Name                 string        `bson:"name" json:"name"`
	Email                string        `bson:"email" json:"email"`
	Password             string        `bson:"password" json:"password"`
	Estado               string        `bson:"estado" json:"estado"` /*PENDIENTE | VERIFICADO | PREMIUM | BANEADO*/
	FirstLogin           bool          `bson:"first_login" json:"first_login"`
	Token                string        `bson:"token" json:"token"`
	IdNivel              bson.ObjectId `bson:"id_nivel" json:"id_nivel"`
	LastFrase            int64         `bson:"last_frase" json:"last_frase"`
	CodigoPromocional    string        `bson:"codigo_promocional" json:"codigo_promocional"`
	MaximoDiaSuscripcion int64         `bson:"maximo_dia_suscripcion" json:"maximo_dia_suscripcion"`
	LastHeartBeat        int64         `bson:"last_heart_beat,omitempty" json:"last_heart_beat,omitempty"`
	LastIp               string        `bson:"last_ip,omitempty" json:"last_ip,omitempty"`
	Newsletter           bool          `bson:"newsletter" json:"newsletter"`
	Detalle              []Question    `bson:"detalle" json:"detalle"`
	Socket_id            string        `bson:"socket_id" json:"socket_id"`
	Connected            string          `bson:"connected" json:"connected"`
}

type Info struct {
	Preguntas          string `bson:"Preguntas"  jaon:"Preguntas"`
	Historiales        string `bson:"Historiales"  jaon:"Historiales"`
	AreasWTemas        string `bson:"AreasWTemas"  jaon:"AreasWTemas"`
	Examenes_oficiales string `bson:"Examenes_oficiales"  jaon:"Examenes_oficiales"`
	Simulacros         string `bson:"Simulacros"  jaon:"Simulacros"`
	Niveles            string `bson:"Niveles"  jaon:"Niveles"`
	Usuario            string `bson:"Usuario"  jaon:"Usuario"`
	Historiales_cp     string `bson:"Historiales_cp"  jaon:"Historiales_cp"`
	Casos_practicos    string `bson:"Casos_practicos"  jaon:"Casos_practicos"`
	Basicosconfundidos string `bson:"Basicosconfundidos"  jaon:"Basicosconfundidos"`
	Legislaciones      string `bson:"Legislaciones"  jaon:"Legislaciones"`
}

type Login struct {
	Ip       string `bson:"ip" json:"ip"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}
type Verify struct {
	Verificador string `bson:"verificador" json:"verificador"`
}
type VerifyAcountUser struct {
	Id string `bson:"id" json:"id"`
}
type JwtToken struct {
	Token string `json:"token"`
}
type FindUsuario struct {
	Email string `bson:"email" json:"email"`
}

type Question struct {
	Question string `bson:"question" json:"question"`
	Answer   string `bson:"answer" json:"answer"`
}

type SetNivelUsuario struct {
	Email     string        `bson:"email" json:"email"`
	IdNivel   bson.ObjectId `bson:"id_nivel" json:"id_nivel"`
	Preguntas []Question    `bson:"preguntas" json:"preguntas"`
}

type NewDayUserUpdate struct {
	Id string `bson:"id"   json:"id"`
	Maximo_dia_suscripcion int64 `bson:"maximo_dia_suscripcion" json:"maximo_dia_suscripcion"`
}

type DeleteHistorial struct {
	Id string `bson:"id"   json:"id"`
}

func (n *Usuarios) InsertUsuario(nivel Usuarios) error {
	return Insert(db, "usuarios", nivel)
}

func (n *Usuarios) FindAllUsuarios() ([]Usuarios, error) {
	var result []Usuarios
	err := FindAll(db, "usuarios", nil, nil, &result)
	return result, err
}

func (n *Usuarios) FindUsuarioByEmail(email string) (Usuarios, error) {
	var result Usuarios
	err := FindOne(db, "usuarios", bson.M{"email": email}, nil, &result)
	return result, err
}

func (n *Usuarios) FindUsuarioById(id string) (Usuarios, error) {
	var result Usuarios
	err := FindOne(db, "usuarios", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}


func (n *Usuarios) FindUsuarioByName(name string) ([]Usuarios, error) {
	var result []Usuarios
	query := bson.M{"name": bson.M{"$regex": name, "$options": "i"}}
	err := FindAll(db, "usuarios", query, nil, &result)
	return result, err
}

func (n *Usuarios) FindUsuariosByEmail(name string) ([]Usuarios, error) {
	var result []Usuarios
	query := bson.M{"email": bson.M{"$regex": name, "$options": "i"}}
	err := FindAll(db, "usuarios", query, nil, &result)
	return result, err
}

func (n *Usuarios) FindUsuarioByVerificador(verificador string) (Usuarios, error) {
	var result []Usuarios
	var usuario Usuarios
	var encontrado = false

	err := FindAll(db, "usuarios", nil, nil, &result)
	for _, data := range result {
		if helper.GetMD5Hash("VERIFICAREMAIL"+strings.ToLower(strings.TrimSpace(data.Email))) == verificador {
			if data.Estado == "PENDIENTE" || data.Estado == "VERIFICADO" {
				usuario = data
				encontrado = true
			}
		}
	}
	if !encontrado {
		err = errors.New("No encontrado")
	}
	return usuario, err
}

func (n *Usuarios) UpdateUsuario(nivel Usuarios) error {
	return Update(db, "usuarios", bson.M{"_id": nivel.Id}, nivel)
}

func (n* Usuarios) UpdateUsuarioPassword(usuario Usuarios) error {
	return Update(db, "usuarios", bson.M{"_id": usuario.Id}, bson.M{"$set": bson.M{"password": usuario.Password}})
}


func (n *Usuarios) RemoveUsuario(id string) error {
	return Remove(db, "usuarios", bson.M{"_id": bson.ObjectIdHex(id)})
}

func (n *Usuarios) UpdateDatauser(usuario Usuarios) error {
	return Update(db, "usuarios", bson.M{"_id": usuario.Id}, usuario)
}