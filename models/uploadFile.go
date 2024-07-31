package models

import "github.com/globalsign/mgo/bson"

type Files struct {
	Id          bson.ObjectId    `bson:"_id"   json:"id"`
	FileName    string           `bson:"filename" json:"filename"`
	FileSize    int64            `bson:"filesize" json:"filesize"`
	Url         string           `bson:"filepath" json:"filepath"`
	PageLink    string           `bson:"pagelink" json:"pagelink"`
	IdPregunta  bson.ObjectId    `bson:"" json:""`
}

type FileUpload struct {
	File string `bson:"file" json:"file`
}

// type FileData struct {
// 	FileURL string `bson:"fileurl" json:"fileurl"`
// }

type FileData struct {
	FileData  Files      `bson:"fileData" json:"fileData"`
	Pregunta  Preguntas  `bson:"pregunta" json:"pregunta"`
}

func (n *Files) InsertFile (file Files) ( error) {
	return Insert(db, "files", file)
}

func (n *Files) FindFileByName (filename string) ( Files, error) {
	var file Files
	err := FindOne(db, "files", bson.M{"filename": filename}, nil, &file)
	return file, err
}

func (n *Files) FindFileByID (id string) ( Files, error) {
	var file Files
	err := FindOne(db, "files", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &file)
	return file, err
}

func (n *Files) FindAllFiles () ( []Files, error) {
	var files []Files
	err := FindAll(db, "files", nil, nil, &files)
	return files, err
}

func (n *Files) FindFilesByIdPregunta (idPregunta string) ( []Files, error) {
	var files []Files
	err := FindAll(db, "files", bson.M{"idpregunta": bson.ObjectIdHex(idPregunta)}, nil, &files)
	return files, err
}


func (n *Files) DeleteFilesByIdPregunta (idPregunta string) error {
	return RemoveAll(db, "files", bson.M{"idpregunta": bson.ObjectIdHex(idPregunta)})
}

func (n *Files) DeleteFileByID (id string) error {
	return Remove(db, "files", bson.M{"_id": bson.ObjectIdHex(id)})
}