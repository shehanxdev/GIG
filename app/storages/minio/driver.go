package minio

import (
	"GIG/app/storages/interfaces"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/revel/revel"
	"log"
)

type Handler struct {
	interfaces.StorageHandlerInterface
	Client         *minio.Client
	CacheDirectory string
}

func (h Handler) GetCacheDirectory() string {
	return h.CacheDirectory
}

/*
NewHandler - Always use the NewHandler method to create an instance.
Otherwise, the handler will not be configured
*/
func NewHandler(cacheDirectory string) *Handler {
	var err error
	handler := new(Handler)
	endpoint, _ := revel.Config.String("minio.endpoint")
	accessKeyID, _ := revel.Config.String("minio.accessKeyID")
	secretAccessKey, _ := revel.Config.String("minio.secretAccessKey")
	secureUrl, _ := revel.Config.Bool("minio.secureUrl")
	handler.CacheDirectory = cacheDirectory

	// Initialize minio client object.
	handler.Client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secureUrl,
	})
	if err != nil {
		log.Println("error connecting to Minio file server")
		panic(err)
	}
	return handler
}
