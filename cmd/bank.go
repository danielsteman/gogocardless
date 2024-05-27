package main

import "github.com/go-chi/chi/v5"

type bankResource struct{}

func (rs bankResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Get("/", rs.List)    // GET /users - read a list of users
	// r.Post("/", rs.Create) // POST /users - create a new user and persist it
	// r.Put("/", rs.Delete)

	// r.Route("/{id}", func(r chi.Router) {
	// 	// r.Use(rs.TodoCtx) // lets have a users map, and lets actually load/manipulate
	// 	r.Get("/", rs.Get)       // GET /users/{id} - read a single user by :id
	// 	r.Put("/", rs.Update)    // PUT /users/{id} - update a single user by :id
	// 	r.Delete("/", rs.Delete) // DELETE /users/{id} - delete a single user by :id
	// })

	return r
}
