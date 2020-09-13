package apiserver

import (
	"context"
	"html/template"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodsData ...
type PodsData struct {
	Namespace string
	Pods      []string
}

// HandlerPods ...
func (s *APIServer) HandlerPods() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.RequestURI == "/pods/" {
			io.WriteString(w, "Pods help ")
		} else {
			vars := mux.Vars(r)

			clientset, err := kubernetes.NewForConfig(s.restconf)
			if err != nil {
				s.logger.Error("NewForConfig", err)
				io.WriteString(w, "can't find config file")
				return
			}

			pods, err := clientset.CoreV1().Pods(vars["ns"]).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				s.logger.Error("Pod list", err)
				io.WriteString(w, "can't read pods in namespace: "+vars["ns"])
				return
			}

			a := make([]string, len(pods.Items))
			for i, pod := range pods.Items {
				a[i] = pod.GetName()
			}
			data := PodsData{
				Namespace: vars["ns"],
				Pods:      a,
			}
			templ, _ := template.ParseFiles("templates/pods.html")
			templ.Execute(w, data)
			// io.WriteString(w, body)
		}

	}
}
