package routes

import (
	"fmt"
	"net/http"
	"package/controller"
)

func StartServer() {
	fmt.Println("server started at 8000")
	http.ListenAndServe(":8000", nil)

}
func InitializeRoutes() {
	http.HandleFunc("/", controller.Contact)
	http.HandleFunc("/success", controller.Success)
	http.HandleFunc("/show", controller.Show)
	http.HandleFunc("/edit", controller.Edit)
	http.HandleFunc("/update", controller.Update)
	http.HandleFunc("/delete", controller.Delete)
}
