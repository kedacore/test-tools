package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	authv1 "k8s.io/api/authentication/v1"
	authzv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	PORT = 8080
)

var value = 0

type response struct {
	Value   int  `json:"value"`
	Success bool `json:"success"`
}

type contextKey string

const userContextKey contextKey = "user"

func authAndAuthz(clientset *kubernetes.Clientset, next http.Handler) http.Handler {
	return authenticate(clientset, authorize(clientset, next))
}

func authenticate(clientset *kubernetes.Clientset, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		tokenReview := &authv1.TokenReview{
			Spec: authv1.TokenReviewSpec{
				Token: tokenString,
			},
		}

		// make sure there's rbac to allow this
		response, err := clientset.AuthenticationV1().TokenReviews().Create(context.TODO(), tokenReview, metav1.CreateOptions{})
		if err != nil || !response.Status.Authenticated {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, response.Status.User.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func authorize(clientset *kubernetes.Clientset, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the user from the context
		user, ok := r.Context().Value(userContextKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// check some arbitrary permission that we want to enforce
		// here we will check if pods can view the /api/value requestPath
		subjectAccessReview := &authzv1.SubjectAccessReview{
			Spec: authzv1.SubjectAccessReviewSpec{
				NonResourceAttributes: &authzv1.NonResourceAttributes{
					Path: "/api/value",
					Verb: "get",
				},
				User: user,
			},
		}

		response, err := clientset.AuthorizationV1().SubjectAccessReviews().Create(context.TODO(), subjectAccessReview, metav1.CreateOptions{})
		if err != nil || !response.Status.Allowed {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getValue(w http.ResponseWriter, r *http.Request) {
	rsp := response{
		Value:   value,
		Success: true,
	}

	json.NewEncoder(w).Encode(rsp)
}

func setValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	number := vars["number"]

	value, _ = strconv.Atoi(number)
	json.NewEncoder(w).Encode(value)
	fmt.Printf("new value: %d\n", value)
}

func main() {
	fmt.Printf("Running server on port: %d\n", PORT)

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to create in-cluster config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	r := mux.NewRouter()

	// authenticate only requests that are verified by the k8s auth apiserver
	r.Handle("/api/value", authAndAuthz(clientset, http.HandlerFunc(getValue))).Methods("GET")
	r.HandleFunc("/api/value/{number:[0-9]+}", setValue).Methods("POST")
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}).Methods("GET")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil))
}
