package linkservice

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"regexp"

	"github.com/pavelzagorodnyuk/linkservice/internal/api"
)

var (
	// длина коротких ссылок
	lengthLink = 10

	// URLTemplate представляет собой скомпилированное регулярное выражение для
	// проверки строки на соответствие требованиям URL
	URLTemplate = regexp.MustCompile(`^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`)

	// linkTemplate представляет собой скомпилированное регулярное выражение
	// для проверки строки на соответствие требованиям короткой ссылки
	linkTemplate = regexp.MustCompile(`^[0-9a-zA-Z_]{10}$`)
)

var (
	// ErrReqProc возвращается в случаях, когда не удается обработать
	// gRPC-запрос
	ErrReqProc = errors.New("linkservice: the request could not be processed")

	// ErrInvalidURL возвращается в случаях, когда gRPC-запрос содержит
	// некорректный URL
	ErrInvalidURL = errors.New("linkservice: the request contains an invalid URL")

	// ErrInvalidLink возвращается в случаях, когда gRPC-запрос содержит
	// некорректную короткую ссылку
	ErrInvalidLink = errors.New("linkservice: the request contains an invalid link")

	// ErrURLNotFound возвращается в случаях, когда для указанной короткой
	// ссылки не существует оригинальной ссылки URL
	ErrURLNotFound = errors.New("linkservice: unknown abbreviated link — the original URL was not found")
)

type GRPCServer struct {
	Database *sql.DB
	api.UnimplementedLinkServiceServer
}

func (s *GRPCServer) Create(ctx context.Context, req *api.URL) (*api.Link, error) {
	// проверка переданной в запросе строки на соответствие требованиям URL
	if !URLTemplate.MatchString(req.GetUrl()) {
		return nil, ErrInvalidURL
	}

	// проверяем, сгенерирована ли короткая ссылка для указанного URL
	row := s.Database.QueryRow("SELECT link FROM links WHERE original_url = $1;", req.GetUrl())

	var link string
	err := row.Scan(&link)

	// если во время запроса произошла ошибка и она не является sql.ErrNoRows,
	// то отправляем сообщение с невозможностью обработать запрос
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Create method: %v\n", err)
		return nil, ErrReqProc
	}

	// если ошибок (включая sql.ErrNoRows) во время запроса не произошло, то
	// для указанного в запросе URL найдена сокращенная ссылка
	if err == nil {
		return &api.Link{Link: link}, nil
	}

	// генерируем для указанного URL короткую ссылку и добавляем новую запись
	// в базу данных. Если подобная короткая ссылка уже существует, то
	// генерируем новую и повторяем попытку добавления записи. Повторяем до
	// тех пор, пока не добавится новая запись или не произойдет иная ошибка

	// UCViolation представляет собой текстовое описание ошибки, возникающей
	// при нарушении ограничения уникальности в PostgreSQL
	UCViolation := "pq: duplicate key value violates unique constraint \"link_pk\""

	for {
		// генерируем для указанного URL короткую ссылку
		link = generateRandomСharacters(lengthLink)

		_, err := s.Database.Exec("INSERT INTO links (link, original_url) VALUES ($1, $2);", link, req.GetUrl())

		// если произошла ошибка, которая не является шибкой UCViolation, то
		// завершаем работу метода и сообщаем о ситуации
		if err != nil && err.Error() != UCViolation {
			log.Printf("Create method: %v\n", err)
			return nil, ErrReqProc
		}

		if err == nil {
			break
		}
	}

	return &api.Link{Link: link}, nil
}

func (s *GRPCServer) Get(ctx context.Context, req *api.Link) (*api.URL, error) {
	// проверка переданной в запросе строки на соответствие требованиям
	// короткой ссылки
	if !linkTemplate.MatchString(req.GetLink()) {
		return nil, ErrInvalidLink
	}

	// запрашиваем исходный URL по сокращенной ссылке
	row := s.Database.QueryRow("SELECT original_url FROM links WHERE link = $1;", req.GetLink())

	var url string
	err := row.Scan(&url)

	// если во время запроса произошла ошибка и она не является sql.ErrNoRows,
	// то отправляем сообщение с невозможностью обработать запрос
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Get method: %v\n", err)
		return nil, ErrReqProc
	}

	// если записей в базе данных для данной сокращенной ссылки не найдено, то
	// возвращаем соответствующую ошибку
	if err == sql.ErrNoRows {
		return nil, ErrURLNotFound
	}

	return &api.URL{Url: url}, nil
}

// generateRandomCharacters генерирует строки длиной length случайных символов.
// При генерации используются символы латинского алфавита в нижнем и верхнем
// регистре, цифры и символ подчеркивания (_).
func generateRandomСharacters(length int) string {
	// задаем исходный алфавит символов
	alphabet := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz")

	rc := make([]rune, length)

	// заполняем срез rc случайными символами алфавита
	for i := 0; i < length; i++ {
		rc[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(rc)
}
