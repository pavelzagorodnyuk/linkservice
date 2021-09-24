package linkservice

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pavelzagorodnyuk/linkservice/internal/api"
)

// параметры подключения к базе данных для проведения тестов
// комбинация параметров для подключения к БД, запущенной скриптом run_test_db.sh
var DBConnParamsForTests = "user=postgres password=passw0rd host=0.0.0.0 port=5433 dbname=linkservice sslmode=disable"

// var connectionParamsForTests = os.ExpandEnv("user=$POSTGRES_USER password=$POSTGRES_PASSWORD host=$DB_HOST port=$DB_PORT dbname=$POSTGRES_DB sslmode=disable")
// var connectionParamsForTests = "user=postgres password=passw0rd host=127.0.0.1 port=5432 dbname=linkservice sslmode=disable"

var TestCreateCases = []struct {
	name     string
	req      *api.URL
	expError error
}{
	{
		name: "test_1",
		req: &api.URL{
			Url: "https://golang.org/",
		},
		expError: nil,
	},
	{
		name: "test_2",
		req: &api.URL{
			Url: "this is not a URL",
		},
		expError: ErrInvalidURL,
	},
}

var TestCreateCasesTwo = []struct {
	req *api.URL
}{
	{
		req: &api.URL{
			Url: "http://abc.abc/",
		},
	},
	{
		req: &api.URL{
			Url: "http://archive.org/filename.txt",
		},
	},
	{
		req: &api.URL{
			Url: "http://qwerty.bca/",
		},
	},
}

func TestCreate(t *testing.T) {
	// устанавливаем подключение к базе данных
	db, err := sql.Open("postgres", DBConnParamsForTests)
	if err != nil {
		t.Fatalf("failed connecting to the database: %v", err)
	}

	defer db.Close()

	// проверка метода Create на тест-кейсах
	for _, testCase := range TestCreateCases {
		t.Run(testCase.name, func(t *testing.T) {

			service := GRPCServer{Database: db}
			res, err := service.Create(context.Background(), testCase.req)

			switch {
			case err == nil && testCase.expError == nil:
				// если ошибки нет и не ожидалось, то проверяем сокращенную
				// ссылку на корректный формат
				if res == nil {
					t.Errorf("it was expected that the response would not be equal to nil")
					return
				}

				if !linkTemplate.MatchString(res.GetLink()) {
					t.Errorf("the abbreviated link has an incorrect format")
				}

			case err != testCase.expError:
				// получена ошибка err, отличная от той, что ожидалась
				t.Errorf("an error with a value of \"%v\" was expected, but \"%v\" was received",
					testCase.expError, err)

			case err == testCase.expError:
				// если произошла ожидаемая ошибка, то все OK, но res должен
				// быть равен nil
				if res != nil {
					t.Errorf("it was expected that the response would be equal to nil")
				}
			}
		})
	}

	// проверка метода Create для случаев, когда в нескольких запросах
	// содержится один и тот же URL
	for _, testCase := range TestCreateCasesTwo {
		t.Run(testCase.req.GetUrl(), func(t *testing.T) {

			service := GRPCServer{Database: db}

			res, err := service.Create(context.Background(), testCase.req)

			if err != nil {
				t.Errorf("Create method reported an error: %v", err)
				return
			}

			link1 := res.GetLink()

			res, err = service.Create(context.Background(), testCase.req)

			if err != nil {
				t.Errorf("Create method reported an error: %v", err)
				return
			}

			link2 := res.GetLink()

			if link1 != link2 {
				t.Errorf("different abbreviated links were generated for the same URL")
			}
		})
	}
}

var TestGetCases = []struct {
	name     string
	req      *api.Link
	url      string
	expError error
}{
	{
		name:     "test_1",
		req:      nil,
		expError: nil,
	},
	{
		name: "test_2",
		req: &api.Link{
			Link: "123_abcABC",
		},
		expError: ErrURLNotFound,
	},
	{
		name: "test_3",
		req: &api.Link{
			Link: "@5gfh35^Gdfh&EWR",
		},
		expError: ErrInvalidLink,
	},
}

func TestGet(t *testing.T) {
	// устанавливаем подключение к базе данных
	db, err := sql.Open("postgres", DBConnParamsForTests)
	if err != nil {
		t.Fatalf("failed connecting to the database: %v", err)
	}

	defer db.Close()

	// дополняем тест-кейсы необходимыми корректными короткими ссылками
	service := GRPCServer{Database: db}

	for i := range TestGetCases {
		if TestGetCases[i].req == nil {
			url := "http://abc.abc/"

			res, err := service.Create(context.Background(), &api.URL{
				Url: url,
			})

			if err != nil {
				t.Fatalf("failed to add values to test cases")
			}

			TestGetCases[i].req = &api.Link{Link: res.GetLink()}
			TestGetCases[i].url = url
		}
	}

	// проверка метода Get на тест-кейсах
	for _, testCase := range TestGetCases {
		t.Run(testCase.name, func(t *testing.T) {

			service := GRPCServer{Database: db}
			res, err := service.Get(context.Background(), testCase.req)

			switch {
			case err == nil && testCase.expError == nil:
				// если ошибки нет и не ожидалось, то проверяем URL
				if res == nil {
					t.Errorf("it was expected that the response would not be equal to nil")
					return
				}

				if res.GetUrl() != testCase.url {
					t.Errorf("URL contained in the response does not match the expected one")
				}

			case err != testCase.expError:
				// получена ошибка err, отличная от той, что ожидалась
				t.Errorf("an error with a value of \"%v\" was expected, but \"%v\" was received",
					testCase.expError, err)

			case err == testCase.expError:
				// если произошла ожидаемая ошибка, то все OK, но res должен
				// быть равен nil
				if res != nil {
					t.Errorf("it was expected that the response would be equal to nil")
				}
			}
		})
	}
}

func TestGenerateRandomСharacters(t *testing.T) {
	var n = 1000

	for i := 0; i < n; i++ {
		link := generateRandomСharacters(lengthLink)
		if !linkTemplate.MatchString(link) {
			t.Errorf("link \"%s\" is incorrect", link)
			t.FailNow()
		}
	}
}
