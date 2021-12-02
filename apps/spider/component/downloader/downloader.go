package downloader

import (
	component2 "github.com/shipengqi/example.v1/apps/spider/component"
	"github.com/shipengqi/example.v1/apps/spider/component/stub"
	"net/http"
)

type DownloaderImpl struct {
	stub.ComponentInternal

	client *http.Client
}

func New(
	cid component2.CID,
	client *http.Client,
	scoreCalculator component2.CalculateScore) (component2.Downloader, error) {

	ci, _ := stub.NewComponentInternal(cid, scoreCalculator)
	return &DownloaderImpl{
		ComponentInternal: ci,
		client:            client,
	}, nil
}

func (d *DownloaderImpl) Download(req *component2.Request) (*component2.Response, error) {
	return nil, nil
}
