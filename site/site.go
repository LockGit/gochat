/**
 * Created by lock
 * Date: 2019-08-12
 * Time: 11:36
 */
package site

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type Site struct {
}

func New() *Site {
	return &Site{}
}

func (s *Site) Run() {
	logrus.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("./"))))
}
