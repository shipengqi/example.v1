package downloader

import (
	"net/http"

	"github.com/shipengqi/example.v1/spider/component"
	"github.com/shipengqi/example.v1/spider/component/stub"
)

type DownloaderImpl struct {
	stub.ComponentInternal

	client *http.Client
}

func New(
	cid component.CID,
	client *http.Client,
	scoreCalculator component.CalculateScore) (component.Downloader, error) {

	ci, _ := stub.NewComponentInternal(cid, scoreCalculator)
	return &DownloaderImpl{
		ComponentInternal: ci,
		client:            client,
	}, nil
}

func (d *DownloaderImpl) Download(req *component.Request) (*component.Response, error) {
	return nil, nil
}
