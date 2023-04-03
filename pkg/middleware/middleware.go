package middleware

//
//import (
//	"bytes"
//	"encoding/json"
//	"github.com/gin-gonic/gin"
//	"github.com/sirupsen/logrus"
//	"gitlab.com/jfcore/common/logger"
//	"io/ioutil"
//	"iot/pkg/service"
//	"net/http"
//	"runtime/debug"
//)
//
//type MiddlewareHandler struct {
//	service service.MiddlewareInterface
//}
//
//type MiddlewareInterface interface {
//	HandleMiddleware() gin.HandlerFunc
//}
//
//func NewMiddlewareHandlers(service service.MiddlewareInterface) *MiddlewareHandler {
//	return &MiddlewareHandler{service: service}
//}
//func (h *MiddlewareHandler) HandleMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		log := logger.WithCtx(c, "Request Detail")
//		defer func() {
//			if r := recover(); r != nil {
//				log.Error(r)
//				debug.PrintStack()
//				panic(r)
//			}
//		}()
//
//		LoggingRequest(c)
//		c.Next()
//	}
//}
//
//func LoggingRequest(c *gin.Context) {
//	log := logger.WithCtx(c, "LoggingRequest")
//
//	r := c.Request
//	header := c.Request.Header
//	params := c.Request.URL.Query()
//	buf, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		log.WithError(err).Errorf("Error reading request body: %v", err.Error())
//		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	var obj map[string]interface{}
//	json.Unmarshal(buf, &obj)
//
//	log.WithFields(logrus.Fields{
//		"Request params": params,
//		"Request body":   obj,
//		"header":         header,
//	}).Info("uri: ", c.Request.RequestURI)
//
//	reader := ioutil.NopCloser(bytes.NewBuffer(buf))
//	c.Request.Body = reader
//}
