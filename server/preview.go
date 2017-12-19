package server

import (
	"github.com/qiangxue/fasthttp-routing"
	"html/template"
	"bytes"
	"github.com/russross/blackfriday"
	"io"
)

var preview *template.Template

func init() {
	preview = template.Must(template.ParseFiles("docs/template.txt"))
}

func (s *Server) handlePreview(c *routing.Context) error {
	promoType := c.Param("type")
	promoId := c.Param("id")

	buf := bytes.Buffer{}
	err := s.RenderPromotion(promoType, promoId, &buf)
	if err != nil {
		return err
	}

	c.Response.Header.Set("Content-Type", "text/html")
	rendered := blackfriday.MarkdownCommon(buf.Bytes())
	c.Response.SetBody(rendered)

	return nil
}

func (s *Server) RenderPromotion(promoType, promoId string, wr io.Writer) (error) {
	pr, err := s.api.GetPromo(promoType, promoId)
	if err != nil {
		return err
	}
	return preview.ExecuteTemplate(wr, pr.GetMeta().Type, &pr)
}