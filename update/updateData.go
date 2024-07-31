package updateData

import (
	"log"
	"test/models"
)
var (
	daoUsuarios = models.Usuarios{}
	daoHistorial = models.Historiales{}
	daoHistorialesCP = models.HistorialesCP{}
)
func UpdateData() {
	users,err := daoUsuarios.FindAllUsuarios()
	if err != nil{
		log.Println(err)
		log.Println("Error xd")
		return
	}

	//Actualizar fechas de suscripcion
	var fechaNueva int64 = 1689476399 // 15 de julio de 2023
	for _, user := range users {
		user.MaximoDiaSuscripcion = fechaNueva
		if err := daoUsuarios.UpdateDatauser(user); err != nil {
			log.Println("Error al actuaizar el usuario '"+user.Name+"'")
		}else{
			log.Println("Se actualizo la fecha de suscripcion al usuario '"+user.Name+"'")
		}
	}

	//Actualizar notas test
	historiales, err := daoHistorial.FindAllHistoriales()
	if err != nil{
		log.Println("Ocurrio un error al buscar historiales")
		return
	}

	for _, historial := range historiales {
		if historial.Puntuacion < 0 {
			historial.Puntuacion = 0
			if err:=daoHistorial.UpdateHistorial(historial); err != nil {
				log.Println("Error al actualizar historial")
			}
		}
	}

	//Actualizar nota de test-cp

	historialesCP, err := daoHistorialesCP.FindAllHistorialesCP()
	if err != nil{
		log.Println("Ocurrio un error al buscar historiales")
		return
	}

	for _, historialcp := range historialesCP {
		if historialcp.Puntuacion < 0 {
			historialcp.Puntuacion = 0
			if err:=daoHistorialesCP.UpdateHistorial(historialcp); err != nil {
				log.Println("Error al actualizar historial-cp")
			}
		}
	}

}	