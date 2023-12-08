package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tdewolff/minify/v2"
	minifyhtml "github.com/tdewolff/minify/v2/html"
	"net/http"
	"strings"
	"time"
)

var minifier *minify.M

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = newTemplates()

	e.Static("/images", "./views/static/images")
	e.Static("/stylesheets", "/views/static/stylesheets")

	minifier = minify.New()
	// GOTCHA: minify considers echo.MIMETextHTMLCharsetUTF8 as a media type, not a mime type
	minifier.Add(echo.MIMETextHTML, &minifyhtml.Minifier{
		KeepQuotes:     true,
		TemplateDelims: minifyhtml.GoTemplateDelims,
	})

	vm := newIndexViewModel()

	e.GET("/", func(ctx echo.Context) error {
		return ctx.Render(200, "index", vm)
	})

	e.POST("/count", func(ctx echo.Context) error {
		vm.Count++
		return ctx.Render(http.StatusOK, "index", vm)
	})

	e.POST("/contacts", func(ctx echo.Context) error {
		nameVal := strings.TrimSpace(ctx.FormValue("name"))
		emailVal := strings.TrimSpace(ctx.FormValue("email"))
		vm.ContactForm.Values["name"] = nameVal
		vm.ContactForm.Values["email"] = emailVal

		clear(vm.ContactForm.Errors)
		if len(nameVal) == 0 {
			vm.ContactForm.Errors["name"] = "Name must not be empty."
		}
		if len(emailVal) == 0 {
			vm.ContactForm.Errors["email"] = "Email must not be empty."
		} else if vm.hasEmail(emailVal) {
			vm.ContactForm.Errors["email"] = "Email already existed."
		}
		if len(vm.ContactForm.Errors) > 0 {
			return ctx.Render(http.StatusUnprocessableEntity, "contact-form", vm)
		}

		newContact := newContactDto(nameVal, emailVal)
		vm.Contacts = append(vm.Contacts, *newContact)

		buf0 := new(bytes.Buffer)
		if err := e.Renderer.Render(buf0, "contact-div-oob", *newContact, ctx); err != nil {
			return err
		}
		buf1 := new(bytes.Buffer)
		if err := e.Renderer.Render(buf1, "contact-form", vm, ctx); err != nil {
			return err
		}
		buf0.Write(buf1.Bytes())

		return ctx.HTMLBlob(http.StatusOK, buf0.Bytes())
	})

	e.DELETE("/contacts/:id", func(ctx echo.Context) error {
		// simulate long-running request
		time.Sleep(time.Second)

		respHeaders := ctx.Response().Header()
		vm.ContactsErrors = []string{}

		idVal := ctx.Param("id")
		contactId, err := uuid.Parse(idVal)

		var existingContact *ContactDto
		var contactIdx int
		if err != nil {
			vm.ContactsErrors = append(vm.ContactsErrors, fmt.Sprintf("Not a valid contact ID: '%s'.", idVal))

			respHeaders.Set("HX-Retarget", "section#contacts-display-container")
			return ctx.Render(http.StatusBadRequest, "contacts-display", vm)
		}
		if existingContact, contactIdx = vm.findContactIndex(contactId); existingContact == nil {
			vm.ContactsErrors = append(vm.ContactsErrors, fmt.Sprintf("Contact '%s' does not exist.", idVal))

			respHeaders.Set("HX-Retarget", "section#contacts-display-container")
			return ctx.Render(http.StatusUnprocessableEntity, "contacts-display", vm)
		}

		vm.Contacts = append(vm.Contacts[:contactIdx], vm.Contacts[contactIdx+1:]...)
		return ctx.Render(http.StatusOK, "contacts-display", vm)
	})

	err := e.Start("3201")
	e.Logger.Fatal(err)
}
