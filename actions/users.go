package actions

import "github.com/gobuffalo/buffalo"

type UsersResource struct {
	buffalo.Resource
}

// List default implementation.
func (v UsersResource) List(c buffalo.Context) error {
	return c.Render(200, r.String("Users#List"))
}

// Show default implementation.
func (v UsersResource) Show(c buffalo.Context) error {
	return c.Render(200, r.String("Users#Show"))
}

// New default implementation.
func (v UsersResource) New(c buffalo.Context) error {
	return c.Render(200, r.String("Users#New"))
}

// Create default implementation.
func (v UsersResource) Create(c buffalo.Context) error {
	return c.Render(200, r.String("Users#Create"))
}

// Edit default implementation.
func (v UsersResource) Edit(c buffalo.Context) error {
	return c.Render(200, r.String("Users#Edit"))
}

// Update default implementation.
func (v UsersResource) Update(c buffalo.Context) error {
	return c.Render(200, r.String("Users#Update"))
}

// Destroy default implementation.
func (v UsersResource) Destroy(c buffalo.Context) error {
	return c.Render(200, r.String("Users#Destroy"))
}

// UsersShow default implementation.
func UsersShow(c buffalo.Context) error {
	return c.Render(200, r.HTML("users/show.html"))
}

// UsersIndex default implementation.
func UsersIndex(c buffalo.Context) error {
	return c.Render(200, r.HTML("users/index.html"))
}

// UsersCreate default implementation.
func UsersCreate(c buffalo.Context) error {
	return c.Render(200, r.HTML("users/create.html"))
}
