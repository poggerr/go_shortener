package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/poggerr/go_shortener/internal/authorization"
	"github.com/poggerr/go_shortener/internal/models"
	pb "github.com/poggerr/go_shortener/internal/proto"
	"github.com/poggerr/go_shortener/internal/service"
	"github.com/poggerr/go_shortener/internal/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"net/http"
	"time"
)

var (
	ErrLinkIsDeleted = errors.New("ссылка удалена")
	ErrLinkNotFound  = errors.New("ссылка не найдена")
)

type URLShortener struct {
	linkRepo service.URLShortenerService
	baseURL  string
	pb.UnimplementedShortenerServer
}

var _ pb.ShortenerServer = (*URLShortener)(nil)

// NewURLShortener создает URLShortener и инициализирует его адресом, по которому будут доступны методы,
// и репозиторием хранения ссылок.
func NewURLShortener(base string, repo service.URLShortenerService) *URLShortener {
	hand := URLShortener{}
	hand.linkRepo = repo
	hand.baseURL = base

	return &hand
}

// CreateShortURL хендлер создания короткой ссылки
func (a *URLShortener) CreateShortURL(res http.ResponseWriter, req *http.Request) {
	userID := authorization.FromContext(req.Context())

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	if !utils.IsURL(string(body)) {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Некорректный URL"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shortURL, err := a.linkRepo.Store(ctx, userID, string(body))
	if err != nil {
		switch {
		case shortURL != "":
			res.WriteHeader(http.StatusConflict)
			res.Write([]byte(a.baseURL + "/" + shortURL))
			return
		default:
			log.Debug().Msg(fmt.Sprintf("store error: %s", err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	res.Header().Set("content-type", "text/plain; charset=utf-8")

	res.WriteHeader(http.StatusCreated)

	res.Write([]byte(a.baseURL + "/" + shortURL))

}

// DBConnect хендлер проверки подключения к БД
func (a *URLShortener) DBConnect(res http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := a.linkRepo.Ping(ctx); err != nil {
		log.Debug().Msg(fmt.Sprintf("database connection error: %s", err))
		res.WriteHeader(http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
}

// CreateBatch Хендлер создания нескольких ссылок
func (a *URLShortener) CreateBatch(res http.ResponseWriter, req *http.Request) {
	userID := authorization.FromContext(req.Context())
	var list models.BatchList
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &list)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("unmarshal error: %s", err))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	list, err = a.linkRepo.StoreBatch(ctx, userID, list, a.baseURL)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	marshal, err := json.Marshal(list)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("marshal error: %s", err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json ")
	res.WriteHeader(http.StatusCreated)
	res.Write(marshal)

}

// CreateJSONShorten хендлер создания ссылки из json
func (a *URLShortener) CreateJSONShorten(res http.ResponseWriter, req *http.Request) {
	userID := authorization.FromContext(req.Context())

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}

	var url models.URL

	err = json.Unmarshal(body, &url)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("unmarshal error: %s", err))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if !utils.IsURL(url.LongURL) {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Некорректный URL"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res.Header().Set("content-type", "application/json ")
	shortURL, err := a.linkRepo.Store(ctx, userID, url.LongURL)
	if err != nil {
		switch {
		case shortURL != "":
			res.WriteHeader(http.StatusConflict)
			res.Write([]byte(a.baseURL + "/" + shortURL))
			return
		default:
			log.Debug().Msg(fmt.Sprintf("store error: %s", err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	shortURL = a.baseURL + "/" + shortURL
	shortenMap := make(map[string]string)
	shortenMap["result"] = shortURL

	marshal, err := json.Marshal(shortenMap)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("marshal error: %s", err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write(marshal)

}

// GetUrlsByUser хендлер получения списка ссылок пользователя
func (a *URLShortener) GetUrlsByUser(res http.ResponseWriter, req *http.Request) {
	userID := authorization.FromContext(req.Context())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	strg, err := a.linkRepo.GetUserStorage(ctx, userID, a.baseURL)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("error: %s", err))
		res.WriteHeader(http.StatusNoContent)
		return
	}

	bucket := MapToBucket(a.baseURL, strg)
	marshal, err := json.Marshal(bucket)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("marshal error: %s", err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json ")
	res.WriteHeader(http.StatusOK)
	res.Write(marshal)

}

// ReadOriginalURL хендлер получения оригинальной ссылки
func (a *URLShortener) ReadOriginalURL(res http.ResponseWriter, req *http.Request) {
	shortURL := chi.URLParam(req, "id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ans, err := a.linkRepo.Restore(ctx, shortURL)
	switch {
	case errors.Is(err, ErrLinkNotFound):
		log.Debug().Err(err)
		res.WriteHeader(http.StatusNoContent)
		return
	case errors.Is(err, ErrLinkIsDeleted):
		log.Debug().Err(err)
		res.WriteHeader(http.StatusGone)
		return
	case err != nil:
		log.Debug().Msg(fmt.Sprintf("error: %s", err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", "text/plain ")
	res.Header().Set("Location", ans)
	res.WriteHeader(http.StatusTemporaryRedirect)

}

// DeleteUrls хендлер удаления ссылок
func (a *URLShortener) DeleteUrls(res http.ResponseWriter, req *http.Request) {
	userID := authorization.FromContext(req.Context())

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("error: %s", err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var keys []string

	err = json.Unmarshal(body, &keys)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("unmarshal error: %s", err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.linkRepo.Delete(ctx, userID, keys)

	res.WriteHeader(http.StatusAccepted)
}

func (a *URLShortener) GetStats(res http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stat, err := a.linkRepo.Statistics(ctx)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("error: %s", err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	marshal, err := json.Marshal(stat)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("marshal error: %s", err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(marshal)

}

func (a *URLShortener) Shorten(ctx context.Context, request *pb.ShortRequest) (*pb.ShortResponse, error) {
	var resp pb.ShortResponse
	userID := uuid.MustParse(request.User)

	shortURL, err := a.linkRepo.Store(ctx, &userID, request.Url)
	if err != nil {
		switch {
		case shortURL != "":
			return nil, status.Error(codes.AlreadyExists, "link is already shortened")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	resp.ShortUrl = fmt.Sprintf("%s/%s", a.baseURL, shortURL)
	return &resp, nil
}

func (a *URLShortener) Expand(ctx context.Context, request *pb.ExpandRequest) (*pb.ExpandResponse, error) {
	var resp pb.ExpandResponse
	url, err := a.linkRepo.Restore(ctx, request.ShortUrl)
	if err != nil {
		switch {
		case errors.Is(err, ErrLinkNotFound):
			return nil, status.Error(codes.NotFound, "link not found")
		case errors.Is(err, ErrLinkIsDeleted):
			return nil, status.Error(codes.NotFound, "link is deleted")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	resp.Url = url
	return &resp, nil
}

func (a *URLShortener) Batch(ctx context.Context, request *pb.ShortBatchRequest) (*pb.ShortBatchResponse, error) {
	var resp pb.ShortBatchResponse
	userID := uuid.MustParse(request.UserId)

	for _, url := range request.Original {
		shortURL, err := a.linkRepo.Store(ctx, &userID, url.OriginalUrl)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		ans := &pb.CorrelationShortURL{
			CorrelationId: url.CorrelationId,
			ShortUrl:      shortURL,
			OriginalUrl:   url.OriginalUrl,
		}
		resp.Original = append(resp.Original, ans)
	}
	return &resp, nil
}

func (a *URLShortener) GetUserBucket(ctx context.Context, request *pb.GetUserBucketRequest) (*pb.GetUserBucketResponse, error) {
	var resp pb.GetUserBucketResponse
	userID := uuid.MustParse(request.UserId)
	urlsMap, err := a.linkRepo.GetUserStorage(ctx, &userID, a.baseURL)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp.Pair = make([]*pb.Pair, 0, len(urlsMap))
	for k, v := range urlsMap {
		resp.Pair = append(resp.Pair, &pb.Pair{
			ShortUrl:    k,
			OriginalUrl: fmt.Sprintf("%s%s", a.baseURL, v),
		})
	}
	return &resp, nil
}

func (a *URLShortener) Stats(ctx context.Context, empty *emptypb.Empty) (*pb.StatResponse, error) {
	var resp pb.StatResponse
	ans, err := a.linkRepo.Statistics(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp.Urls = int32(ans.LinksCount)
	resp.User = int32(ans.UsersCount)
	return &resp, nil
}

func (a *URLShortener) Delete(ctx context.Context, request *pb.DeleteRequest) (*emptypb.Empty, error) {
	userID := uuid.MustParse(request.User)
	a.linkRepo.Delete(ctx, &userID, request.Ids)
	return &emptypb.Empty{}, nil
}
