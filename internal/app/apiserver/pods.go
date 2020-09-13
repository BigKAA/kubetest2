package apiserver

import (
	"context"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

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
			body := "Pods </br>"
			for _, pod := range pods.Items {
				body += " " + pod.GetName()
			}

			io.WriteString(w, body)
		}

	}
}
