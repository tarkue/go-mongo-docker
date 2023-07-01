package endpoint

import (
	"log"
	"minly-backend/internal/app/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
)

type Service interface {
	MatchLink(ctx echo.Context, link string) error
}

type DataBase interface {
	UpdateCounter(initial string, counter int64) error
	GetByInitial(initial string) (bson.M, error)
	GetByResult(result string) (bson.M, error)
}

type Endpoint struct {
	s  Service
	db DataBase
}

func New(s Service, db DataBase) *Endpoint {
	return &Endpoint{
		s:  s,
		db: db,
	}
}

func (e *Endpoint) GetResult(ctx echo.Context) error {

	if ctx.QueryParams().Has(ParamInitial) {
		link := &models.Link{}
		initial := ctx.QueryParams().Get(ParamInitial)

		err := e.s.MatchLink(ctx, initial)

		if err != nil {

			er := &models.HandlerError{
				Error: ParamBroken + ParamInitial,
			}

			ctx.JSON(http.StatusBadRequest, &er)
			return nil
		}

		res, err := e.db.GetByInitial(initial)
		if err != nil {
			return err
		}

		mapstructure.Decode(
			res,
			&link,
		)

		err = e.db.UpdateCounter(link.Initial, link.Counter+1)
		if err != nil {
			return err
		}

		ctx.JSON(http.StatusOK, &link)
		return nil
	}

	er := &models.HandlerError{
		Error: ParamNotFound + ParamInitial,
	}
	ctx.JSON(http.StatusBadRequest, &er)
	return nil
}

func (e *Endpoint) GetLink(ctx echo.Context) error {
	if ctx.QueryParams().Has(ParamResult) {
		link := &models.Link{}
		result := ctx.QueryParams().Get(ParamResult)

		res, err := e.db.GetByResult(result)

		log.Println(res, err)
		if err != nil {
			er := &models.HandlerError{
				Error: ParamBroken + ParamResult,
			}
			ctx.JSON(http.StatusBadRequest, &er)
			return nil
		}

		mapstructure.Decode(
			res,
			&link,
		)

		ctx.JSON(http.StatusOK, &link.Initial)
		return nil
	}

	er := &models.HandlerError{
		Error: ParamNotFound + ParamResult,
	}

	ctx.JSON(http.StatusBadRequest, &er)
	return nil
}
