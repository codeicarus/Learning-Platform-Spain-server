package models

import "github.com/globalsign/mgo/bson"

type Transacciones struct {
	Id            bson.ObjectId `bson:"_id" json:"id"`
	IdTransaccion int64         `bson:"id_transaccion" json:"id_transaccion"`
	Fecha         int64         `bson:"fecha" json:"fecha"`
	Estado        int           `bson:"estado" json:"estado"`
	Producto      string        `bson:"producto" json:"producto"`
	Precio        float64       `bson:"precio" json:"precio"`
	FechaCobro    int64         `bson:"fecha_cobro" json:"fecha_cobro"`
	Token         string        `bson:"token" json:"token"`
	IdUsuario     bson.ObjectId `bson:"id_usuario" json:"id_usuario"`
}

func (n *Transacciones) InsertTransaccion(transaccion Transacciones) error {
	return Insert(db, "transacciones", transaccion)
}

func (n *Transacciones) FindAllTransacciones() ([]Transacciones, error) {
	var result []Transacciones
	err := FindAll(db, "transacciones", nil, nil, &result)
	return result, err
}

func (n *Transacciones) FindTransaccionById(id string) (Transacciones, error) {
	var result Transacciones
	err := FindOne(db, "transacciones", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}
func (n *Transacciones) FindTransaccionLastTransaction() (Transacciones, error) {
	var result Transacciones
	err := FindOneOrder(db, "transacciones", nil, nil, &result, "-id_transaccion")
	return result, err
}
func (n *Transacciones) FindTransaccionByIdTransaccion(id_transaccion int64) (Transacciones, error) {
	var result Transacciones
	err := FindOne(db, "transacciones", bson.M{"id_transaccion": id_transaccion}, nil, &result)
	return result, err
}

func (n *Transacciones) UpdateTransaccion(transaccion Transacciones) error {
	return Update(db, "transacciones", bson.M{"_id": transaccion.Id}, transaccion)
}

func (n *Transacciones) RemoveTransaccion(id string) error {
	return Remove(db, "transacciones", bson.M{"_id": bson.ObjectIdHex(id)})
}
