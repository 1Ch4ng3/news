// routes/routes.go
package routes

import (
    "github.com/gorilla/mux"
    "myapp/controllers"
)

func NewRouter() *mux.Router {
    router := mux.NewRouter()

    
    router.HandleFunc("/categories", controllers.HandleCategories).Methods("GET")
   /* router.HandleFunc("/categories", controllers.HandleCreateCategory).Methods("POST")
    router.HandleFunc("/categories/{id}", controllers.HandleGetCategory).Methods("GET")
    router.HandleFunc("/categories/{id}", controllers.HandleUpdateCategory).Methods("PUT")
    router.HandleFunc("/categories/{id}", controllers.HandleDeleteCategory).Methods("DELETE")
*/
    return router
}
