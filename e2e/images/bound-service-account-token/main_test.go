package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	authv1 "k8s.io/api/authentication/v1"
	authzv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func TestAuthenticationAndAuthorization(t *testing.T) {
	for _, tc := range []struct {
		name              string
		audience          string
		acceptedAudiences []string
		authenticated     bool
		allowed           bool
		reviewError       bool
		wantStatus        int
		wantAuthz         bool
	}{
		{name: "legacy audience omitted", authenticated: true, allowed: true, wantStatus: http.StatusOK, wantAuthz: true},
		{name: "dedicated audience accepted", audience: "keda-metrics-e2e", acceptedAudiences: []string{"keda-metrics-e2e"}, authenticated: true, allowed: true, wantStatus: http.StatusOK, wantAuthz: true},
		{name: "audience missing from response", audience: "keda-metrics-e2e", authenticated: true, allowed: true, wantStatus: http.StatusUnauthorized},
		{name: "different audience returned", audience: "keda-metrics-e2e", acceptedAudiences: []string{"other"}, authenticated: true, allowed: true, wantStatus: http.StatusUnauthorized},
		{name: "token rejected", audience: "keda-metrics-e2e", acceptedAudiences: []string{"keda-metrics-e2e"}, allowed: true, wantStatus: http.StatusUnauthorized},
		{name: "token review fails", audience: "keda-metrics-e2e", reviewError: true, wantStatus: http.StatusUnauthorized},
		{name: "authorization still required", audience: "keda-metrics-e2e", acceptedAudiences: []string{"keda-metrics-e2e"}, authenticated: true, wantStatus: http.StatusForbidden, wantAuthz: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("TOKEN_AUDIENCE", tc.audience)
			reviewed, authorized := false, false
			api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/apis/authentication.k8s.io/v1/tokenreviews":
					reviewed = true
					var review authv1.TokenReview
					if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
						t.Error(err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					var wantAudiences []string
					if tc.audience != "" {
						wantAudiences = []string{tc.audience}
					}
					if review.Spec.Token != "test-token" || !slices.Equal(review.Spec.Audiences, wantAudiences) {
						t.Errorf("unexpected TokenReview: %#v", review.Spec)
					}
					if tc.reviewError {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					_ = json.NewEncoder(w).Encode(authv1.TokenReview{
						TypeMeta: metav1.TypeMeta{APIVersion: "authentication.k8s.io/v1", Kind: "TokenReview"},
						Status: authv1.TokenReviewStatus{Authenticated: tc.authenticated, Audiences: tc.acceptedAudiences,
							User: authv1.UserInfo{Username: "system:serviceaccount:apps:metrics-reader"}},
					})
				case "/apis/authorization.k8s.io/v1/subjectaccessreviews":
					authorized = true
					var review authzv1.SubjectAccessReview
					if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
						t.Error(err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					if review.Spec.User != "system:serviceaccount:apps:metrics-reader" || review.Spec.NonResourceAttributes == nil ||
						review.Spec.NonResourceAttributes.Path != "/api/value" || review.Spec.NonResourceAttributes.Verb != "get" {
						t.Errorf("unexpected SubjectAccessReview: %#v", review.Spec)
					}
					_ = json.NewEncoder(w).Encode(authzv1.SubjectAccessReview{
						TypeMeta: metav1.TypeMeta{APIVersion: "authorization.k8s.io/v1", Kind: "SubjectAccessReview"},
						Status:   authzv1.SubjectAccessReviewStatus{Allowed: tc.allowed},
					})
				default:
					t.Errorf("unexpected API path: %s", r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer api.Close()
			clientset, err := kubernetes.NewForConfig(&rest.Config{Host: api.URL})
			if err != nil {
				t.Fatal(err)
			}
			nextCalled := false
			handler := authAndAuthz(clientset, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			}))
			request := httptest.NewRequest(http.MethodGet, "/api/value", nil)
			request.Header.Set("Authorization", "Bearer test-token")
			result := httptest.NewRecorder()
			handler.ServeHTTP(result, request)
			if result.Code != tc.wantStatus || !reviewed || authorized != tc.wantAuthz || nextCalled != (tc.wantStatus == http.StatusOK) {
				t.Errorf("status=%d, reviewed=%t, authorized=%t, nextCalled=%t", result.Code, reviewed, authorized, nextCalled)
			}
		})
	}
}
