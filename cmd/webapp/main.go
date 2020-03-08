package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"todoapp/internal/storage"
	"todoapp/pkg/linkedlist"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	err := storage.Localstore.Load()
	if err != nil {
		fmt.Printf("Failed to load data!")
		os.Exit(1)
	}
	go func() {
		// Save TODO: add sigterm handling to save on exit
		for {
			time.Sleep(5 * time.Second)
			err := storage.Localstore.Save()
			if err != nil {
				fmt.Printf("saving file failed: %v", err)
			}
		}
	}()

	r := chi.NewRouter()
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Get("/list/all", allLists)
	r.Post("/list/new", newList)
	r.Get("/list/id/{uuid}", getList)
	r.Delete("/list/id/{uuid}", deleteList)
	r.Post("/list/id/{uuid}/item/add", addItemToList)
	r.Post("/list/id/{list_uuid}/item/id/{item_uuid}", editItemFromList)
	r.Delete("/list/id/{list_uuid}/item/id/{item_uuid}", deleteItemFromList)
	port := ":8080"
	fmt.Println("starting on port", port)
	http.ListenAndServe(port, r)
}

func allLists(w http.ResponseWriter, r *http.Request) {
	l := storage.Localstore.AllLists()
	// we don't want to return the actual lists, just summary
	var l2 []*storage.List
	for _, v := range l {
		v2 := *v
		v2.List = nil
		l2 = append(l2, &v2)
	}
	b, err := json.Marshal(l2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func newList(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["name"]
	if len(name) == 0 || name[0] == "" {
		http.Error(w, "please use 'name' parameter", http.StatusUnprocessableEntity)
		return
	}
	l, err := storage.Localstore.NewList(name[0], 1000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func getList(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	list := storage.Localstore.GetList(uuid)
	if list == nil {
		http.Error(w, fmt.Sprintf("uuid '%v' does not exist", uuid), http.StatusUnprocessableEntity)
		return
	}
	b, err := json.Marshal(list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func deleteList(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	storage.Localstore.DeleteList(uuid)
	return
}

func addItemToList(w http.ResponseWriter, r *http.Request) {
	listUUID := chi.URLParam(r, "uuid")
	list := storage.Localstore.GetList(listUUID)
	if list == nil {
		http.Error(w, fmt.Sprintf("list uuid '%v' does not exist", listUUID), http.StatusUnprocessableEntity)
		return
	}
	r.ParseForm()
	text := r.Form["text"]
	if len(text) == 0 || text[0] == "" {
		http.Error(w, "please use 'text' parameter", http.StatusUnprocessableEntity)
		return
	}
	itemUUIDList := r.Form["uuid"]
	var itemUUID string
	if len(itemUUIDList) > 0 {
		itemUUID = itemUUIDList[0]
	}
	entry := linkedlist.Item{
		Text: text[0],
		UUID: itemUUID,
	}
	id, err := list.List.AddItem(entry)
	if err != nil {
		// TODO: different status code if error is because of max
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	entry.UUID = id
	var b []byte
	if len(r.Form["return_list"]) > 0 && r.Form["return_list"][0] == "true" {
		b, err = json.Marshal(list)
	} else {
		b, err = json.Marshal(entry)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
	return
}

func editItemFromList(w http.ResponseWriter, r *http.Request) {
	listUUID := chi.URLParam(r, "list_uuid")
	itemUUID := chi.URLParam(r, "item_uuid")
	list := storage.Localstore.GetList(listUUID)
	if list == nil {
		http.Error(w, fmt.Sprintf("list '%v' does not exist", listUUID), http.StatusNotFound)
		return
	}
	r.ParseForm()
	text := r.Form["text"]
	if len(text) == 0 {
		http.Error(w, fmt.Sprintf("parameter 'text' must be used"), http.StatusUnprocessableEntity)
		return
	}
	list.List.EditItem(itemUUID, text[0])
	if len(r.Form["return_list"]) > 0 && r.Form["return_list"][0] == "true" {
		b, err := json.Marshal(list)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(b)
	}
}

func deleteItemFromList(w http.ResponseWriter, r *http.Request) {
	listUUID := chi.URLParam(r, "list_uuid")
	itemUUID := chi.URLParam(r, "item_uuid")
	list := storage.Localstore.GetList(listUUID)
	if list == nil {
		http.Error(w, fmt.Sprintf("list '%v' does not exist", listUUID), http.StatusNotFound)
		return
	}
	list.List.DeleteItem(itemUUID)
}
