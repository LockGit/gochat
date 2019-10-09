/**
 * Created by lock
 * Date: 2019-08-12
 * Time: 11:36
 */
package site

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gochat/config"
	"net/http"
)

type Site struct {
}

func New() *Site {
	return &Site{}
}

func (s *Site) Run() {
	siteConfig := config.Conf.Site
	port := siteConfig.SiteBase.ListenPort
	addr := fmt.Sprintf(":%d", port)
	logrus.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir("./site/"))))
}
