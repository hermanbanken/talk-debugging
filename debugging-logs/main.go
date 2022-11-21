package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	http.ListenAndServe(addr(), &shop{})
}

type shop struct {
	carts map[string][]string
}

func (s *shop) get(user string) []string {
	if s.carts == nil {
		s.carts = make(map[string][]string)
	}
	return s.carts[user]
}
func (s *shop) add(user string, item string) {
	if s.carts == nil {
		s.carts = make(map[string][]string)
	}
	s.carts[user] = append(s.carts[user], item)
}

func (s *shop) getUser(w http.ResponseWriter, r *http.Request) (user string) {
	// Anonymous users
	c, err := r.Cookie("user")
	if err != nil || c.Value == "" {
		user = strconv.FormatInt(rand.Int63(), 10)
		c = &http.Cookie{Name: "user", Value: user}
		http.SetCookie(w, c)
	} else {
		user = c.Value
	}
	return
}

func (s *shop) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := s.getUser(w, r)

	// Routing
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) <= 1 {
		// Print cart
		cart := s.get(user)
		if len(cart) == 0 {
			w.Write([]byte("Cart empty"))
		} else {
			w.Write([]byte(fmt.Sprintf("Cart: %v", cart)))
		}
		w.Write([]byte("<br>"))

		w.Write([]byte(`Our products:
<ul>
<li><a href="/product/1">Light bulb</a></li>
<li><a href="/product/2">Light bulb 2</a></li>
<li><a href="/product/3">LED strip</a></li>
</ul>
	`))
	}

	if parts[0] == "products" {
		switch parts[1] {
		case "1", "2", "3":
			if r.Method == http.MethodPost {
				s.add(user, parts[1])
			}
			w.Write([]byte(fmt.Sprintf(`
Product %s; <br>
<form><input name='product' value='%s' type='hidden' /><input type='submit' value='add to cart' /></form>`, parts[1], parts[1])))

		default:
			w.WriteHeader(404)
			w.Write([]byte("Not found"))
		}
		w.Write([]byte(parts[1]))
	}
}

func addr() string {
	if port, hasPort := os.LookupEnv("PORT"); hasPort {
		return ":" + port
	}
	return ":8080"
}
