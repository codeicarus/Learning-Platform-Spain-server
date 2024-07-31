package models

import (
	"math/rand"
	"time"

	"github.com/globalsign/mgo/bson"
)

type Preguntas struct {
	Id                   bson.ObjectId      `bson:"_id,omitempty" json:"id"`
	Pregunta             string             `bson:"pregunta" json:"pregunta"`
	Explicacion          string             `bson:"explicacion" json:"explicacion"`
	IdTema               bson.ObjectId      `bson:"id_tema,omitempty" json:"id_tema,omitempty"`
	IdNivel              bson.ObjectId      `bson:"id_nivel,omitempty" json:"id_nivel,omitempty"`
	IdBasicosConfundidos bson.ObjectId      `bson:"id_basicos_confundidos,omitempty" json:"id_basicos_confundidos,omitempty"`
	IdLegislacion        bson.ObjectId      `bson:"id_legislacion,omitempty" json:"id_legislacion,omitempty"`
	Oficial              bool               `bson:"oficial" json:"oficial"`
	AnioOficial          string             `bson:"anio_oficial" json:"anio_oficial"`
	SchemaId               []bson.ObjectId    `bson:"schemaid" json:"schemaid"`
	Schema               []Files    `bson:"schema" json:"schema"`
}

type PreguntasWRespuestas struct {
	Id                   bson.ObjectId `bson:"_id" json:"id"`
	Pregunta             string        `bson:"pregunta" json:"pregunta"`
	Explicacion          string        `bson:"explicacion" json:"explicacion"`
	IdTema               bson.ObjectId `bson:"id_tema" json:"id_tema"`
	IdNivel              bson.ObjectId `bson:"id_nivel" json:"id_nivel"`
	IdBasicosConfundidos bson.ObjectId `bson:"id_basicos_confundidos,omitempty" json:"id_basicos_confundidos,omitempty"`
	IdLegislacion        bson.ObjectId `bson:"id_legislacion,omitempty" json:"id_legislacion,omitempty"`
	Oficial              bool          `bson:"oficial" json:"oficial"`
	AnioOficial          string        `bson:"anio_oficial" json:"anio_oficial"`
	Respuestas           []Respuestas  `bson:"respuestas" json:"respuestas"`
	Schema               []Files      `bson:"schema" json:"schema"`
}
type PreguntasADevolver struct {
	Id          bson.ObjectId         `bson:"_id" json:"id"`
	Pregunta    string                `bson:"pregunta" json:"pregunta"`
	IdTema      bson.ObjectId         `bson:"id_tema" json:"id_tema"`
	Explicacion string                `bson:"explicacion" json:"explicacion"`
	Respuestas  []RespuestasADevolver `bson:"respuestas" json:"respuestas"`
	Oficial     bool                  `bson:"oficial" json:"oficial"`
	AnioOficial string                `bson:"anio_oficial" json:"anio_oficial"`
	Schema               []Files      `bson:"schema" json:"schema"`
}

type PreguntasRespuestas struct {
	Id         bson.ObjectId `bson:"_id" json:"id"`
	Pregunta   string        `bson:"pregunta" json:"pregunta"`
	Respuestas []string      `bson:"respuestas" json:"respuestas"`
}
type CreateTest struct {
	Temas               []string `bson:"temas" json:"temas"`
	NumeroPreguntas     int      `bson:"numero_preguntas" json:"numero_preguntas"`
	RespuestaAutomatica bool     `bson:"respuesta_automatica" json:"respuesta_automatica"`
	Tipo                string   `bson:"tipo" json:"tipo"`
	Email               string   `bson:"email" json:"email"`
}
type SavePregunta struct {
	Email    string                 `bson:"email" json:"email"`
	Pregunta PreguntaRecibidaEditar `bson:"pregunta" json:"pregunta"`
}
type PreguntaRecibidaEditar struct {
	Id          bson.ObjectId         `bson:"_id" json:"id"`
	Pregunta    string                `bson:"pregunta" json:"pregunta"`
	Explicacion string                `bson:"explicacion" json:"explicacion"`
	Nivel       string                `bson:"nivel" json:"nivel"`
	Respuestas  []RespuestasADevolver `bson:"respuestas" json:"respuestas"`
}
type CreateTestLegislacion struct {
	Legislaciones       []string `bson:"legislaciones" json:"legislaciones"`
	NumeroPreguntas     int      `bson:"numero_preguntas" json:"numero_preguntas"`
	RespuestaAutomatica bool     `bson:"respuesta_automatica" json:"respuesta_automatica"`
	Tipo                string   `bson:"tipo" json:"tipo"`
	Email               string   `bson:"email" json:"email"`
}
type CreateTestBasicoConfundido struct {
	BasicosConfundidos  []string `bson:"basicosconfundidos" json:"basicosconfundidos"`
	NumeroPreguntas     int      `bson:"numero_preguntas" json:"numero_preguntas"`
	RespuestaAutomatica bool     `bson:"respuesta_automatica" json:"respuesta_automatica"`
	Tipo                string   `bson:"tipo" json:"tipo"`
	Email               string   `bson:"email" json:"email"`
}
type CreateTestOficiales struct {
	Oficiales           []string `bson:"oficiales" json:"oficiales"`
	NumeroPreguntas     int      `bson:"numero_preguntas" json:"numero_preguntas"`
	RespuestaAutomatica bool     `bson:"respuesta_automatica" json:"respuesta_automatica"`
	Tipo                string   `bson:"tipo" json:"tipo"`
	Email               string   `bson:"email" json:"email"`
}
type CreateTestSimulacros struct {
	Simulacros          []string `bson:"simulacros" json:"simulacros"`
	NumeroPreguntas     int      `bson:"numero_preguntas" json:"numero_preguntas"`
	RespuestaAutomatica bool     `bson:"respuesta_automatica" json:"respuesta_automatica"`
	Tipo                string   `bson:"tipo" json:"tipo"`
	Email               string   `bson:"email" json:"email"`
}
type CorregirTest struct {
	IdTest             string `bson:"id_test" json:"id_test""`
	RespuestasMarcadas string `bson:"respuestas_marcadas" json:"respuestas_marcadas"`
	Tiempo             int    `bson:"tiempo" json:"tiempo"`
}

func (n *Preguntas) InsertPregunta(pregunta Preguntas) error {
	return Insert(db, "preguntas", pregunta)
}

func (n *Preguntas) FindAllPreguntas() ([]Preguntas, error) {
	var result []Preguntas
	err := FindAll(db, "preguntas", nil, nil, &result)
	return result, err
}

func (n *Preguntas) FindPreguntaById(id string) (Preguntas, error) {
	var result Preguntas
	err := FindOne(db, "preguntas", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}
func (n *PreguntasWRespuestas) FindPreguntaByIdWRespuestas(id string) (PreguntasWRespuestas, error) {
	var result PreguntasWRespuestas
	err := FindOneWRespuestas(db, "preguntas", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}
func (n *PreguntasWRespuestas) FindPreguntasByIdsWRespuestas(ids []bson.ObjectId) ([]PreguntasWRespuestas, error) {
	var result []PreguntasWRespuestas
	err := FindAllWRespuestas(db, "preguntas", bson.M{"_id": bson.M{"$in": ids}}, nil, &result)
	return result, err
}
func (n *Preguntas) FindPreguntasByIds(ids []bson.ObjectId, selector interface{}) ([]Preguntas, error) {
	var result []Preguntas
	err := FindAll(db, "preguntas", bson.M{"_id": bson.M{"$in": ids}}, selector, &result)
	return result, err
}
func (n *Preguntas) FindPreguntaByIdsTemas(ids_temas []bson.ObjectId) ([]Preguntas, error) {
	var result []Preguntas
	err := FindAll(db, "preguntas", bson.M{"id_tema": bson.M{"$in": ids_temas}}, nil, &result)
	return result, err
}

/* Get preguntas por temas  */
func (n *Preguntas) FindPreguntaByIdsTemasIdNivelLimitAnt(ids_temas []bson.ObjectId, id_nivel bson.ObjectId, numero int) ([]Preguntas, error) {

	var result []Preguntas
	err := FindLimit(db, "preguntas", bson.M{"id_tema": bson.M{"$in": ids_temas}}, nil, &result, numero)
	// err := FindLimit(db, "preguntas", bson.M{"id_tema": bson.M{"$in": ids_temas}, "id_nivel": id_nivel}, nil, &result, numero)
	return result, err
}
func (n *Preguntas) FindPreguntaByIdsTemasIdNivelLimit(ids_temas []bson.ObjectId, id_nivel bson.ObjectId, numero int) ([]Preguntas, error) {
	var result []Preguntas
	err := FindAll(db, "preguntas", bson.M{"id_tema": bson.M{"$in": ids_temas}}, nil, &result)
	// err := FindAll(db, "preguntas", bson.M{"id_tema": bson.M{"$in": ids_temas}, "id_nivel": id_nivel}, nil, &result)
	rand.Seed(time.Now().UnixNano())
	for i := len(result) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		result[i], result[j] = result[j], result[i]
	}

	if len(result) < numero {
		numero = len(result)
	}
	resultSlice := result[0:numero]

	return resultSlice, err
}
func (n *Preguntas) FindPreguntaByIdsTemasIdNivelLimitNotInPreguntasAnt(ids_temas []bson.ObjectId, id_nivel bson.ObjectId, numero int, ids_preguntas []bson.ObjectId) ([]Preguntas, error) {
	var result []Preguntas
	err := FindLimit(db, "preguntas", bson.M{"id_tema": bson.M{"$in": ids_temas}, "id_nivel": id_nivel, "_id": bson.M{"$nin": ids_preguntas}}, nil, &result, numero)
	return result, err
}

func (n *Preguntas) FindPreguntaByIdsTemasIdNivelLimitNotInPreguntas(ids_temas []bson.ObjectId, id_nivel bson.ObjectId, numero int, ids_preguntas []bson.ObjectId) ([]Preguntas, error) {
	var result []Preguntas
	err := FindAll(db, "preguntas", bson.M{"id_tema": bson.M{"$in": ids_temas}, "id_nivel": id_nivel, "_id": bson.M{"$nin": ids_preguntas}}, nil, &result)

	rand.Seed(time.Now().UnixNano())
	for i := len(result) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		result[i], result[j] = result[j], result[i]
	}

	if len(result) < numero {
		numero = len(result)
	}
	resultSlice := result[0:numero]

	return resultSlice, err
}

/*
func (n *Preguntas) FindPreguntaByIdsTemasIdNivelLimitNormales(ids_temas []bson.ObjectId, id_nivel bson.ObjectId, numero int) ([]Preguntas, error) {
	var result []Preguntas
	err := FindLimit(db, "preguntas", bson.M{"id_tema": bson.M{"$in": ids_temas}, "id_nivel": id_nivel, "id_basicos_confundidos": bson.M{"$exists": false}}, nil, &result, numero)
	return result, err
}
func (n *Preguntas) FindPreguntaByIdsTemasIdNivelLimitBasicosConfundidos(ids_temas []bson.ObjectId, id_nivel bson.ObjectId, numero int) ([]Preguntas, error) {
	var result []Preguntas
	err := FindLimit(db, "preguntas", bson.M{"id_tema": bson.M{"$in": ids_temas}, "id_nivel": id_nivel, "id_basicos_confundidos": bson.M{"$exists": true}}, nil, &result, numero)
	return result, err
}
*/

func (n *Preguntas) FindPreguntaByIdsLegislacionesIdNivelLimitAnt(ids_legislaciones []bson.ObjectId, id_nivel bson.ObjectId, numero int) ([]Preguntas, error) {
	var result []Preguntas
	err := FindLimit(db, "preguntas", bson.M{"id_legislacion": bson.M{"$in": ids_legislaciones}, "id_nivel": id_nivel}, nil, &result, numero)
	return result, err
}
func (n *Preguntas) FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones []bson.ObjectId, id_nivel bson.ObjectId, numero int) ([]Preguntas, error) {
	var result []Preguntas
	err := FindAll(db, "preguntas", bson.M{"id_legislacion": bson.M{"$in": ids_legislaciones}, "id_nivel": id_nivel}, nil, &result)
	rand.Seed(time.Now().UnixNano())
	for i := len(result) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		result[i], result[j] = result[j], result[i]
	}

	if len(result) < numero {
		numero = len(result)
	}
	resultSlice := result[0:numero]

	return resultSlice, err
}

//#################################################
//#################################################
//#################################################
func (n *Preguntas) FindRandomPreguntasByLegislacion(numPreguntas int, idLegislacion string) ([]Preguntas, error) {
	var result []Preguntas

	err := FindAll(db, "preguntas", bson.M{"id_legislacion": bson.ObjectIdHex(idLegislacion)}, nil, &result)
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	if len(result) < numPreguntas {
		numPreguntas = len(result)
	}

	return result[:numPreguntas], nil
}

func (n *Preguntas) FindRandomPreguntasByTemas(numPreguntas int, idTema string) ([]Preguntas, error) {
	var result []Preguntas

	err := FindAll(db, "preguntas", bson.M{"id_tema": bson.ObjectIdHex(idTema)}, nil, &result)
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	if len(result) < numPreguntas {
		numPreguntas = len(result)
	}

	return result[:numPreguntas], nil
}

//################################################333
//################################################333
//################################################333

func (n *Preguntas) FindPreguntaByIdsBasicosConfundidosIdNivelLimitAnt(ids_basicosconfundidos []bson.ObjectId, id_nivel bson.ObjectId, numero int) ([]Preguntas, error) {
	var result []Preguntas
	err := FindLimit(db, "preguntas", bson.M{"id_basicos_confundidos": bson.M{"$in": ids_basicosconfundidos}, "id_nivel": id_nivel}, nil, &result, numero)
	return result, err
}
func (n *Preguntas) FindPreguntaByIdsBasicosConfundidosIdNivelLimit(ids_basicosconfundidos []bson.ObjectId, id_nivel bson.ObjectId, numero int) ([]Preguntas, error) {
	var result []Preguntas
	err := FindAll(db, "preguntas", bson.M{"id_basicos_confundidos": bson.M{"$in": ids_basicosconfundidos}, "id_nivel": id_nivel}, nil, &result)

	rand.Seed(time.Now().UnixNano())
	for i := len(result) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		result[i], result[j] = result[j], result[i]
	}
	if len(result) < numero {
		numero = len(result)
	}

	resultSlice := result[0:numero]

	return resultSlice, err
}

func (n *Preguntas) UpdatePregunta(pregunta Preguntas) error {
	return Update(db, "preguntas", bson.M{"_id": pregunta.Id}, pregunta)
}

func (n *Preguntas) RemovePregunta(id string) error {
	return Remove(db, "preguntas", bson.M{"_id": bson.ObjectIdHex(id)})
}


//##################################
//##################################
//##################################

func (n *Preguntas) FindAllPreguntasByLegislacion(idLegislacin string) ([]Preguntas, error) {
	var preguntas []Preguntas
	err := FindAll(db, "preguntas", bson.M{"id_legislacion": bson.ObjectIdHex(idLegislacin)}, nil, &preguntas)
	return preguntas, err
}

func (n *Preguntas) FindAllPreguntasByTema(idTema string) ([]Preguntas, error) {
	var preguntas []Preguntas
	err := FindAll(db, "preguntas", bson.M{"id_tema": bson.ObjectIdHex(idTema)}, nil, &preguntas)
	return preguntas, err
}

func (n *Preguntas) FindFileByPreguntaId(idPregunta string) ([]Files, error) {
	var files []Files
	err := FindAll(db, "files", bson.M{"idpregunta": bson.ObjectIdHex(idPregunta)}, nil, &files)
	return files, err
}