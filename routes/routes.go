package routes

import (
	"net/http"
	"test/auth"
	"test/controllers"

	"github.com/gorilla/mux"
)

type Route struct {
	Method     string
	Pattern    string
	Handler    http.HandlerFunc
	Middleware mux.MiddlewareFunc
}

var routes []Route

func init() {

	register("GET", "/getCasosPracticos", controllers.GetCasosPracticos, auth.TokenMiddleware)
	register("GET", "/getOficiales", controllers.GetOficiales, auth.TokenMiddleware)
	register("POST", "/delete_preguntas_favoritas", controllers.DeletePreguntasFavoritas, auth.TokenMiddleware)
	register("POST", "/delete_simulacro", controllers.DeleteSimulacro, auth.TokenMiddleware)
	register("POST", "/set_preguntas_favoritas", controllers.SetPreguntasFavoritas, auth.TokenMiddleware)
	register("POST", "/get_preguntas_favoritas", controllers.GetPreguntasFavoritas, auth.TokenMiddleware)
	register("POST", "/delete_preguntas_cp_favoritas", controllers.DeletePreguntasCPFavoritas, auth.TokenMiddleware)
	register("POST", "/set_preguntas_cp_favoritas", controllers.SetPreguntasCPFavoritas, auth.TokenMiddleware)
	register("POST", "/get_preguntas_cp_favoritas", controllers.GetPreguntasCPFavoritas, auth.TokenMiddleware)
	register("POST", "/importadores", controllers.DoImport, auth.TokenMiddleware)
	register("POST", "/importadores_areas", controllers.DoImportAreas, auth.TokenMiddleware)
	register("POST", "/importadores_oficiales", controllers.DoImportOficiales, auth.TokenMiddleware)
	register("POST", "/importadores_casos_practicos", controllers.DoImportCP, auth.TokenMiddleware)
	register("POST", "/importadores_casos_practicos_oficiales", controllers.DoImportCPOficiales, auth.TokenMiddleware)
	register("POST", "/importadores_basico_confundidos", controllers.DoImportBC, auth.TokenMiddleware)
	register("POST", "/importadores_test_nivel", controllers.DoImportTestNivel, auth.TokenMiddleware)
	register("POST", "/guardar_como_simulacro", controllers.GuardarComoSimulacro, auth.TokenMiddleware)
	register("POST", "/hay_error", controllers.HayError, auth.TokenMiddleware)
	register("POST", "/tiene_duda", controllers.TieneDuda, auth.TokenMiddleware)
	register("POST", "/cambiar_tiempo_transcurrido", controllers.CambiarTiempoTranscurrido, auth.TokenMiddleware)
	register("POST", "/corregir_test_ya_contestado", controllers.CorregirTestYaContestado, auth.TokenMiddleware)
	register("POST", "/corregir_test", controllers.CorregirTest, auth.TokenMiddleware)
	register("POST", "/corregir_pregunta_ya_contestada", controllers.CorregirPreguntaYaContestada, auth.TokenMiddleware)
	register("POST", "/corregir_pregunta", controllers.CorregirPregunta, auth.TokenMiddleware)
	register("POST", "/preguntas/{id}", controllers.GetTest, auth.TokenMiddleware)
	register("POST", "/preguntas", controllers.CreateTest, auth.TokenMiddleware)
	register("POST", "/save_pregunta", controllers.SavePregunta, auth.TokenMiddleware)
	register("POST", "/preguntasLegislaciones", controllers.CreateTestLegislaciones, auth.TokenMiddleware)
	register("POST", "/preguntasBasicosConfundidos", controllers.CreateTestBasicosConfundidos, auth.TokenMiddleware)
	register("POST", "/preguntasOficiales", controllers.CreateTestOficiales, auth.TokenMiddleware)
	register("POST", "/preguntasSimulacros", controllers.CreateTestSimulacros, auth.TokenMiddleware)
	register("POST", "/preguntasFalladasBlancasGuardadas", controllers.CreateTestFalladasBlancasGuardadas, auth.TokenMiddleware)
	register("POST", "/getDataFalladasBlancasGuardadas", controllers.GetDataFalladasBlancasGuardadas, auth.TokenMiddleware)
	register("POST", "/repetirTest", controllers.RepetirTest, auth.TokenMiddleware)
	register("POST", "/repetirTestCP", controllers.RepetirTestCP, auth.TokenMiddleware)
	register("POST", "/update-puntaje", controllers.UpdatePuntajetest, auth.TokenMiddleware)
	register("POST", "/update-puntaje-cp", controllers.UpdatePuntajetestCP, auth.TokenMiddleware)
	register("DELETE", "/delete-tema/{id}", controllers.DeleteTema, auth.TokenMiddleware)

	register("GET", "/all_legislacion_data/{id}", controllers.FindDataLegislacion, auth.TokenMiddleware)
	register("PUT", "/chane_name_legislacion", controllers.UpdateNameLegislacion, auth.TokenMiddleware)
	register("DELETE", "/delete_legislacion/{id}", controllers.DeletePreguntaAndRespuesta, auth.TokenMiddleware)
	register("DELETE", "/delete_all_pregunta_legislacion/{id}", controllers.DeleteAllPreguntasLegislacion, auth.TokenMiddleware)
	register("POST", "/upload_legislacion_new/{id}", controllers.UploadQuestionLegislacion, auth.TokenMiddleware)

	register("GET", "/all_tema_data/{id}", controllers.FindDataTema, auth.TokenMiddleware)
	register("PUT", "/chane_name_tema", controllers.UpdateNameTema, auth.TokenMiddleware)
	register("DELETE", "/delete_tema/{id}", controllers.DeletePreguntaAndRespuesta, auth.TokenMiddleware)
	register("DELETE", "/delete_all_pregunta_temas/{id}", controllers.DeleteAllPreguntasTema, auth.TokenMiddleware)
	register("POST", "/upload_tema_new/{id}", controllers.UploadQuestionTema, auth.TokenMiddleware)

	register("POST", "/update-name-simulacro", controllers.UpdateNameSimulacro, auth.TokenMiddleware)

	register("POST", "/save_pregunta_cp", controllers.SavePreguntaCP, auth.TokenMiddleware)
	register("POST", "/preguntasCasosPracticos", controllers.CreateTestCasosPracticos, auth.TokenMiddleware)
	register("POST", "/cambiar_tiempo_transcurrido_cp", controllers.CambiarTiempoTranscurridoCP, auth.TokenMiddleware)
	register("POST", "/preguntas_cp/{id}", controllers.GetTestCP, auth.TokenMiddleware)
	register("POST", "/corregir_pregunta_cp", controllers.CorregirPreguntaCP, auth.TokenMiddleware)
	register("POST", "/corregir_pregunta_cp_ya_contestada", controllers.CorregirPreguntaCPYaContestada, auth.TokenMiddleware)
	register("POST", "/corregir_test_cp", controllers.CorregirTestCP, auth.TokenMiddleware)
	register("POST", "/corregir_test_cp_ya_contestado", controllers.CorregirTestCPYaContestado, auth.TokenMiddleware)
	register("POST", "/hay_error_cp", controllers.HayErrorCP, auth.TokenMiddleware)
	register("POST", "/tiene_duda_cp", controllers.TieneDudaCP, auth.TokenMiddleware)

	register("GET", "/temas", controllers.AllTemas, auth.TokenMiddleware)

	register("GET", "/areas", controllers.AllAreas, auth.TokenMiddleware)
	register("GET", "/areasWLegislaciones", controllers.AllAreasWLegislaciones, auth.TokenMiddleware)
	register("GET", "/areasWBasicosConfundidos", controllers.AllAreasWBasicosConfundidos, auth.TokenMiddleware)
	register("GET", "/areasWTemas", controllers.AllAreasWTemas, auth.TokenMiddleware)
	register("GET", "/areasWBloques", controllers.AllAreasWBloques, auth.TokenMiddleware)

	register("GET", "/niveles", controllers.AllNiveles, auth.TokenMiddleware)
	register("GET", "/niveles/{id}", controllers.FindNivel, auth.TokenMiddleware)
	register("POST", "/niveles", controllers.CreateNivel, auth.TokenMiddleware)
	register("PUT", "/niveles", controllers.UpdateNivel, auth.TokenMiddleware)
	register("DELETE", "/niveles/{id}", controllers.DeleteNivel, auth.TokenMiddleware)

	register("POST", "/usuarios/login", controllers.LoginUsuario, nil)
	register("POST", "/usuarios/register", controllers.RegisterUsuario, nil)
	register("POST", "/usuarios/verify", controllers.VerifyUsuario, nil)
	register("POST", "/usuarios/verify-user-acount", controllers.VerifyUsuarioAcount, auth.TokenMiddleware)
	register("DELETE", "/usuarios/delete-usuario/{id}", controllers.DeleteUserAcount, auth.TokenMiddleware)
	register("POST", "/usuarios/checkMDS", controllers.CheckMDS, auth.TokenMiddleware)
	register("GET", "/usuarios/check", controllers.CheckUsuario, auth.TokenMiddleware)
	register("PUT", "/usuarios/set-password", controllers.SetNewPassword, auth.TokenMiddleware)
	register("POST", "/usuarios/change-day", controllers.ChangeDaySuscripcion, auth.TokenMiddleware)
	register("POST", "/usuarios/delete-historial-cp", controllers.DeleteHistorialCP, auth.TokenMiddleware)
	register("POST", "/usuarios/delete-historial", controllers.DeleteHistorial, auth.TokenMiddleware)

	register("POST", "/usuarios/check_ip", controllers.CheckIpUsuario, auth.TokenMiddleware)
	register("POST", "/usuarios/reset_password", controllers.ResetPassword, nil)
	register("POST", "/usuarios/test_nivel", controllers.CreateTestNivel, auth.TokenMiddleware)
	register("POST", "/usuarios/checksuper", controllers.CheckSuperUsuario, auth.TokenMiddleware)
	register("POST", "/usuarios/find", controllers.FindUsuario, auth.TokenMiddleware)
	register("PUT", "/usuarios/nivel", controllers.SetNivelUsuario, auth.TokenMiddleware)
	register("PUT", "/usuarios", controllers.UpdateUsuario, auth.TokenMiddleware)
	register("POST", "/usuarios/get_stats", controllers.StatsUsuario, auth.TokenMiddleware)
	register("POST", "/usuarios/change_mis_datos", controllers.ChangeMisDatos, auth.TokenMiddleware)
	register("POST", "/usuarios/change_mis_datos", controllers.ChangeMisDatos, auth.TokenMiddleware)
	register("POST", "/usuarios/get_all_user", controllers.GetAllUser, auth.TokenMiddleware)
	register("POST", "/usuarios/get_search_user", controllers.SearchUser, auth.TokenMiddleware)
	register("POST", "/usuarios/get_user_info", controllers.StatsUsuarioInfo, auth.TokenMiddleware)

	register("POST", "/usuarios/createTransacction", controllers.CreateTransacction, auth.TokenMiddleware)
	register("POST", "/usuarios/checkTransacction", controllers.CheckTransacction, nil)
	register("POST", "/contact", controllers.Contact, nil)

	register("POST", "/frases_motivadoras", controllers.GetFraseMotivadora, auth.TokenMiddleware)

	register("POST", "/buscar_pregunta_para_editar", controllers.BuscarPreguntaParaEditar, auth.TokenMiddleware)
	register("POST", "/info_gpt", controllers.GetInfoGPT, auth.TokenMiddleware)
	register("POST", "/continue_gpt", controllers.ContinueChatGPT, auth.TokenMiddleware)

	register("GET", "/get_all_preguntas", controllers.GetAllPreguntas, auth.TokenMiddleware)
	register("POST", "/upload_file/{id}", controllers.UploadFile2, auth.TokenMiddleware)
	register("GET", "/files/{id}", controllers.ViewFile, nil)
	register("GET", "/all_files", controllers.GetViewAllFiles, auth.TokenMiddleware)
	register("DELETE", "/delete_file/{id}", controllers.DeleteFile, auth.TokenMiddleware)
	register("DELETE", "/delete_files_question/{idQuestion}/{idSchema}", controllers.DeleteQuestionAndSchema, auth.TokenMiddleware)

	register("POST", "/start_payment", controllers.StartPayment, auth.TokenMiddleware)
	register("GET", "/view-payments", controllers.ViewAllPayments, auth.TokenMiddleware)
	register("PUT", "/accept-payment", controllers.AcceptPayment, auth.TokenMiddleware)
	register("PUT", "/rechazar-payment", controllers.RechazarPago, auth.TokenMiddleware)
	register("DELETE", "/delete-suscription/{id}", controllers.DeleteSuscripcion, auth.TokenMiddleware)

	register("GET", "/email/validate_email", controllers.SendEmail, nil)
	register("GET", "/email/reset_password", controllers.Resetpassword, nil)
	register("GET", "/email/testSuccess", controllers.TestSuccess, nil)
	register("GET", "/email/testError", controllers.TestError, nil)
	register("GET", "/email/new_suscripcion", controllers.EmailPayStarted, nil)
	register("GET", "/email/accept_suscripcion", controllers.EmailAccepPayStarted, nil)
	register("GET", "/email/rechazar_suscripcion", controllers.EmailRechazarPayStarted, nil)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	for _, route := range routes {
		r := router.Methods(route.Method).
			Path(route.Pattern)
		if route.Middleware != nil {
			r.Handler(route.Middleware(route.Handler))
		} else {
			r.Handler(route.Handler)
		}
	}
	return router
}

func register(method, pattern string, handler http.HandlerFunc, middleware mux.MiddlewareFunc) {
	routes = append(routes, Route{method, pattern, handler, middleware})
}
