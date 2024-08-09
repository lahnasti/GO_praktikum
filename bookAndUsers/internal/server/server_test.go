package server

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/lahnasti/GO_praktikum/internal/models"
	mock_server "github.com/lahnasti/GO_praktikum/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsersHandler(t *testing.T) {
	gin.SetMode(gin.TestMode) // Устанавливаем тестовый режим для Gin
	var srv Server
	r := gin.Default()
	r.GET("/users", srv.GetAllUsersHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close() // Закрываем сервер в конце теста

	type want struct {
		code  int
		users string
	}
	type test struct {
		name    string
		request string
		method  string
		users   []models.User
		err     error
		want    want
	}
	tests := []test{
		{
			name:    "Test 'GetAllUsersHandler' #1; Default call",
			request: "/users",
			method:  http.MethodGet,
			users: []models.User{
				{
					UID:      1,
					Name:     "Nastya",
					Login:    "lahnasti",
					Password: "123",
				},
				{
					UID:      2,
					Name:     "Karina",
					Login:    "dmit",
					Password: "kar",
				},
			},
			want: want{
				code:  http.StatusOK,
				users: `{"message":"List users","users":[{"uId":1,"name":"Nastya","login":"lahnasti","password":"123"},{"uId":2,"name":"Karina","login":"dmit","password":"kar"}]}`,
			},
		},
		{
			name:    "Test 'GetAllUsersHandler' #2; Error call",
			request: "/users",
			method:  http.MethodGet,
			users:   nil,
			err:     errors.New("test error"),
			want: want{
				code:  http.StatusInternalServerError,
				users: `{"error":"test error"}`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			//делаем мок контроллер
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // Убеждаемся, что все ожидаемые вызовы выполнены
			// Создаем mock-объект
			m := mock_server.NewMockRepository(ctrl)
			m.EXPECT().GetAllUsers().Return(tc.users, tc.err)
			// Присваиваем mock-объект полю Db сервера
			srv.Db = m
			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpSrv.URL + tc.request
			resp, err := req.Send()
			// Проверяем результаты
			assert.NoError(t, err)
			assert.Equal(t, tc.want.users, string(resp.Body()))
			assert.Equal(t, tc.want.code, resp.StatusCode())
		})
	}

}

func TestRegisterUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode) // Устанавливаем тестовый режим для Gin
	var srv Server
	r := gin.Default()
	r.POST("/users/register", srv.RegisterUserHandler)
	httpSrv := httptest.NewServer(r)

	type want struct {
		code   int
		answer string
	}
	type test struct {
		name    string
		request string
		method  string
		user    string
		err     any
		errCall bool
		dbFlag  bool
		want    want
	}
	tests := []test{
		{
			name:    "Test 'RegisterHandler' #1; Default call",
			request: "/users/register",
			method:  http.MethodPost,
			user:    `{"uId":1,"name":"Vitya","login":"login1","password":"pass1"}`,
			err:     nil,
			dbFlag:  true,
			want: want{
				code:   http.StatusOK,
				answer: `{"message":"User successfully registered","user_id":1}`,
			},
		},
		{
			name:    "Test 'RegisterHandler' #2; BadRequest call",
			request: "/users/register",
			method:  http.MethodPost,
			user:    "",
			err:     nil,
			errCall: true,
			dbFlag:  false,
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:    "Test 'RegisterHandler' #3; Conflict call",
			request: "/users/register",
			method:  http.MethodPost,
			user:    `{"uId":1,"name":"Vitya","login":"login1","password":"pass1"}`,
			err:     errors.New(`ERROR`),
			errCall: true,
			dbFlag:  true,
			want: want{
				code:   http.StatusInternalServerError,
				answer: `{"error":"ERROR"}`,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.dbFlag {
				ctrl := gomock.NewController(t)
				m := mock_server.NewMockRepository(ctrl)
				defer ctrl.Finish()
				m.EXPECT().AddUser(gomock.Any()).Return(1, tc.err)
				srv.Db = m
			}
			req := resty.New().R()
			req.Method = tc.method
			req.Body = tc.user
			req.URL = httpSrv.URL + tc.request
			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, tc.want.code, resp.StatusCode())
			if tc.errCall {
				if tc.want.answer != "" {
					assert.Equal(t, tc.want.answer, string(resp.Body()))
					return
				}
			} else {
				assert.Contains(t, string(resp.Body()), tc.want.answer)
				assert.NotEmpty(t, resp.Header().Get("Authorization"))
				_, err = ValidateJWT(string(resp.Header().Get("Authorization")))
				assert.NoError(t, err)
			}
		})
	}
	httpSrv.Close()
}

func TestGetUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode) // Устанавливаем тестовый режим для Gin

	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // Убеждаемся, что все ожидаемые вызовы выполнены

	// Создаем mock-объект
	m := mock_server.NewMockRepository(ctrl)

	// Инициализация сервера с моком
	srv := &Server{Db: m}
	r := gin.Default()
	r.GET("/users", srv.GetUserHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()

	type want struct {
		code int
		user string
	}
	type test struct {
		name    string
		request string
		method  string
		user    models.User
		err     error
		want    want
	}
	tests := []test{
		{
			name:    "Test 'GetUserHandler' #1; Default call",
			request: "/users?uid=1",
			method:  http.MethodGet,
			user: models.User{
				UID:      1,
				Name:     "Liza",
				Login:    "login1",
				Password: "hashedpass",
			},
			want: want{
				code: http.StatusOK,
				user: `{"uId":1,"name":"Liza","login":"login1","password":"hashedpass"}`,
			},
		},
		{
			name:    "Test 'GetUserHandler' #2; BadRequest - invalid UID",
			request: "/users?uid=invalid",
			method:  http.MethodGet,
			user:    models.User{}, // Не используется
			err:     nil,
			want: want{
				code: http.StatusBadRequest,
				user: `{"error":"invalid argument"}`,
			},
		},
		{
			name:    "Test 'GetUserHandler' #3; InternalServerError - GetUser error",
			request: "/users?uid=2",
			method:  http.MethodGet,
			user:    models.User{}, // Не используется
			err:     errors.New("error getting user"),
			want: want{
				code: http.StatusInternalServerError,
				user: `{"error":"error getting user"}`,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Настройка мока в зависимости от теста
			if tc.want.code == http.StatusOK {
				uid, _ := strconv.Atoi(tc.request[len("/users?uid="):])
				m.EXPECT().GetUser(uid).Return(tc.user, tc.err).Times(1)
			} else if tc.want.code == http.StatusInternalServerError {
				uid, _ := strconv.Atoi(tc.request[len("/users?uid="):])
				m.EXPECT().GetUser(uid).Return(tc.user, tc.err).Times(1)
			}

			resp, err := http.Get(httpSrv.URL + tc.request)
			assert.NoError(t, err)
			assert.Equal(t, tc.want.code, resp.StatusCode)

			body := ""
			if resp.Body != nil {
				defer resp.Body.Close()
				bodyBytes, _ := io.ReadAll(resp.Body)
				body = string(bodyBytes)
			}
			assert.JSONEq(t, tc.want.user, body)
		})
	}
}

//TODO: доделать хэндлер, разобраться с токеном
/*func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode) // Устанавливаем тестовый режим для Gin
	var srv Server
	r := gin.Default()
	r.POST("/users/login", srv.LoginHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close() // Закрываем сервер в конце теста

	type want struct {
		code int
		body string
	}
	type test struct {
		name    string
		request string
		method  string
		user    string
		err     any
		errCall bool
		dbFlag  bool
		want    want
	}

	tests := []test{
		{
			name:    "Test 'LoginHandler' #1; Default call",
			request: "/users/login",
			method:  http.MethodPost,
			user:    `{"login":"login1","password":"pass1"}`,
			err:     nil,
			dbFlag:  true,
			want: want{
				code: http.StatusOK,
				body: "User Liza was logined",
			},
		},
		{
			name:    "Test 'LoginHandler' #2; BadRequest call",
			request: "/users/login",
			method:  http.MethodPost,
			user:    "",
			err:     nil,
			errCall: true,
			dbFlag:  false,
			want: want{
				code: http.StatusBadRequest,
				body: `{"error": err.Error()}`,
			},
		},
		{
			name:    "Test 'LoginHandler' #3; Unauthorized user not found",
			request: "/users/login",
			method:  http.MethodPost,
			user:    `{"login":"login1","password":"pass1"}`,
			err:     errors.New("user not found"),
			want: want{
				code: http.StatusUnauthorized,
				body: `{"error":"Invalid credentials: user not found"}`,
			},
		},
		{
			name:    "Test 'LoginHandler' #4; Unauthorized incorrect password",
			request: "users/login",
			method:  http.MethodPost,
			user:    `{"login":"login1","password":"wrongpass"}`,
			err:     bcrypt.ErrMismatchedHashAndPassword,
			want: want{
				code: http.StatusUnauthorized,
				body: `{"error": "Invalid credentials: incorrect password"}`,
			},
		},
		/*{
			name: "Test 'LoginHandler' #5; StatusInternalServerError on token generation",
			request: "users/login",
			method: http.MethodPost,
			user: `{"login":"login1","password":"pass1"}`,
			err: nil,
			want: want{
				code: http.StatusInternalServerError,
				body: `{"error": "Failed to generate token"}`,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_server.NewMockRepository(ctrl)
			m.EXPECT().GetUserByLogin(gomock.Any()).Return(tc.user, tc.err)
			srv.Db = m
			req := resty.New().R()
			req.Method = tc.method
			req.Body = tc.user
			req.URL = httpSrv.URL + tc.request
			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, tc.want.code, resp.StatusCode())
			if tc.errCall {
				if tc.want.body != "" {
					assert.Equal(t, tc.want.body, string(resp.Body()))
					return
				}
			} else {
				assert.Contains(t, string(resp.Body()), tc.want.body)
				assert.NotEmpty(t, resp.Header().Get("Authorization"))
				_, err = ValidateJWT(string(resp.Header().Get("Authorization")))
				assert.NoError(t, err)
			}
		})
	}
} */

func TestUpdateUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode) // Устанавливаем тестовый режим для Gin
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // Убеждаемся, что все ожидаемые вызовы выполнены
	// Создаем mock-объект
	m := mock_server.NewMockRepository(ctrl)
	// Инициализация сервера с моком
	srv := &Server{Db: m}
	r := gin.Default()
	r.PUT("/users/:id", srv.UpdateUserHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()

	type want struct {
		code int
		user string
	}
	type test struct {
		name    string
		request string
		method  string
		user    string
		err     error
		dbFlag  bool
		want    want
	}
	tests := []test{
		{
			name:    "Test 'UpdateUserHandler' #1; Successful update",
			request: "/users/1",
			method:  http.MethodPut,
			user:    `{"uId":1,"name":"Liza","login":"login1","password":"newhashedpass"}`,
			err:     nil,
			dbFlag:  true,
			want: want{
				code: http.StatusOK,
				user: `{"message":"User updated","user":{"uId":1,"name":"Liza","login":"login1","password":"newhashedpass"}}`,
			},
		},
		{
			name:    "Test 'UpdateUserHandler' #2; BadRequest - Invalid JSON",
			request: "/users/1",
			method:  http.MethodPut,
			user:    `{"uId":1,"name":}`, // invalid JSON
			err:     nil,
			dbFlag:  false,
			want: want{
				code: http.StatusBadRequest,
				user: `{"error":"invalid character '}' looking for beginning of value"}`,
			},
		},
		{
			name:    "Test 'UpdateUserHandler' #3; InternalServerError - Invalid ID",
			request: "/users/1",
			method:  http.MethodPut,
			user:    `{"uId":1,"name":"Liza","login":"login1","password":"newhashedpass"}`,
			err:     errors.New("user not found"),
			dbFlag:  true,
			want: want{
				code: http.StatusNotFound,
				user: `{"error":"user not found"}`,
			},
		},
		{
			name:    "Test 'UpdateUserHandler' #4; NotFound - UpdateUser error",
			request: "/users/invalid",
			method:  http.MethodPut,
			user:    `{"uId":1,"name":"Liza","login":"login1","password":"newhashedpass"}`,
			err:     nil,
			dbFlag:  false,
			want: want{
				code: http.StatusInternalServerError,
				user: `{"error":"strconv.Atoi: parsing \"invalid\": invalid syntax"}`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.dbFlag {
				m.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(tc.err)
			}

			req := resty.New().R()
			req.Method = tc.method
			req.SetBody(tc.user)
			req.URL = httpSrv.URL + tc.request

			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, tc.want.code, resp.StatusCode())
			assert.JSONEq(t, tc.want.user, resp.String())
		})
	}
}

//TODO: Доделать хэндлер
/*func TestDeleteUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock_server.NewMockRepository(ctrl)
	srv := &Server {
			Db: m,
			deleteChan: make(chan int, 1),
	}
	r := gin.Default()
	r.DELETE("/users/:id", srv.DeleteUserHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()

	type want struct {
		code int
		answer string
	}
	type test struct {
		name string
		request string
		method string
		token string
		userID string
		err error
		dbFlag bool
		want want
	}
	tests := []test {
		{
			name: "Test 'DeleteUserHandler' #1; StatusOK",
			request: "/users/1",
			method: http.MethodDelete,
			token: "validToken",
			userID: "1",
			err: nil,
			dbFlag: true,
			want: want {
				code: http.StatusOK,
				answer: `{"message":"User deleted","user_id":1}`,
			},
		},
		{
			name:    "Test 'DeleteUserHandler' #2; Invalid token",
			request: "/user/1",
			method:  http.MethodDelete,
			token:   "invalid-token",
			userID:  "1",
			err:     errors.New("invalid token"),
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"message":"Bad auth token","error":"invalid token"}`,
			},
		},
		{
			name:    "Test 'DeleteUserHandler' #3; Invalid user ID",
			request: "/user/invalid",
			method:  http.MethodDelete,
			token:   "valid-token",
			userID:  "invalid",
			err:     nil,
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"message":"Bad user_id","error":"strconv.Atoi: parsing \"invalid\": invalid syntax"}`,
			},
		},
		{
			name:    "Test 'DeleteUserHandler' #4; Database error",
			request: "/user/1",
			method:  http.MethodDelete,
			token:   "valid-token",
			userID:  "1",
			err:     errors.New("database error"),
			dbFlag:  true,
			want: want{
				code:   http.StatusInternalServerError,
				answer: `{"error":"database error"}`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.dbFlag {
				userIDInt, _ := strconv.Atoi(tc.userID)
				m.EXPECT().SetDeleteStatus(userIDInt).Return(tc.err)
			}

			req := resty.New().R().
				SetHeader("Authorization", tc.token)

			resp, err := req.Execute(tc.method, httpSrv.URL+tc.request)

			assert.NoError(t, err)
			assert.Equal(t, tc.want.code, resp.StatusCode())
			assert.JSONEq(t, tc.want.answer, resp.String())
		})
	}
}*/

func TestGetAllBooksHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var srv Server
	r := gin.Default()
	r.GET("/books", srv.GetAllBooksHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()

	type want struct {
		code   int
		answer string
	}
	type test struct {
		name    string
		request string
		method  string
		books   []models.Book
		err     error
		want    want
	}
	tests := []test{
		{
			name:    "Test 'GetAllBooksHandler' #1; StatusOK",
			request: "/books",
			method:  http.MethodGet,
			books: []models.Book{
				{
					BID:    1,
					Title:  "War and Peace",
					Author: "Lev Tolstoy",
					UID:    1,
				},
				{
					BID:    2,
					Title:  "Diary of a monkey",
					Author: "Jain Birkin",
					UID:    2,
				},
			},
			want: want{
				code:   http.StatusOK,
				answer: `{"message":"Storage books","books":[{"bId":1,"title":"War and Peace","author":"Lev Tolstoy","uId":1},{"bId":2,"title":"Diary of a monkey","author":"Jain Birkin","uId":1}]}`,
			},
		},
		{
			name:    "Test 'GetAllBooksHandler' #2; Error call",
			request: "/books",
			method:  http.MethodGet,
			books:   nil,
			err:     errors.New("test error"),
			want: want{
				code:   http.StatusInternalServerError,
				answer: `{"error":"test error"}`,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_server.NewMockRepository(ctrl)
			m.EXPECT().GetAllBooks().Return(tc.books, tc.err)
			srv.Db = m
			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpSrv.URL + tc.request
			resp, err := req.Send()
			assert.NoError(t, err)
			assert.JSONEq(t, tc.want.answer, string(resp.Body()))
			assert.Equal(t, tc.want.code, resp.StatusCode())
		})
	}
}

func generateValidToken() string {
	// Используйте вашу функцию для генерации токена
	token, _ := GenerateJWT("1") // Параметры могут отличаться в зависимости от реализации
	return token
}

func TestSaveBookHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var srv Server
	r := gin.Default()
	r.POST("/books/add", srv.SaveBookHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()

	type want struct {
		code   int
		answer string
	}
	type test struct {
		name    string
		request string
		method  string
		token   string
		book    string
		err     error
		dbFlag  bool
		want    want
	}
	tests := []test{
		{
			name:    "Test 'SaveBookHandler' #1; Successful save",
			request: "/books/add",
			method:  http.MethodPost,
			token:   generateValidToken(),
			book:    `{"bId":1,"title":"War and Peace","author":"Lev Tolstoy","uId":1}`,
			err:     nil,
			dbFlag:  true,
			want: want{
				code:   http.StatusOK,
				answer: `{"message":"Book saved"}`,
			},
		},
		{
			name:    "Test 'SaveBookHandler' #2; BadRequest call - Invalid JSON",
			request: "/books/add",
			method:  http.MethodPost,
			token:   generateValidToken(),
			book:    "",
			err:     nil,
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"error":"EOF"}`,
			},
		},
		{
			name:    "Test 'SaveBookHandler' #3; Unauthorized call - Invalid token",
			request: "/books/add",
			method:  http.MethodPost,
			token:   "invalid.token",
			book:    `{"bId":1,"title":"War and Peace","author":"Lev Tolstoy","uId":1}`,
			err:     nil,
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"message":"Bad auth token","error":"token contains an invalid number of segments"}`,
			},
		},
		{
			name:    "Test 'SaveBookHandler' #4; InternalServerError call - DB error",
			request: "/books/add",
			method:  http.MethodPost,
			token:   generateValidToken(),
			book:    `{"bId":1,"title":"War and Peace","author":"Lev Tolstoy","uId":1}`,
			err:     errors.New("DB error"),
			dbFlag:  true,
			want: want{
				code:   http.StatusInternalServerError,
				answer: `{"error":"DB error"}`,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_server.NewMockRepository(ctrl)
			if tc.dbFlag {
				m.EXPECT().SaveBook(gomock.Any()).Return(tc.err)
			}
			srv.Db = m

			req := resty.New().R()
			req.Method = tc.method
			req.Body = tc.book
			req.Header.Add("Authorization", tc.token)
			req.URL = httpSrv.URL + tc.request
			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, tc.want.code, resp.StatusCode())
			if tc.want.answer != "" {
				assert.JSONEq(t, tc.want.answer, string(resp.Body()))
			}
		})
	}
}

func TestSaveBooksHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var srv Server
	r := gin.Default()
	r.POST("/books/adds", srv.SaveBooksHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()

	type want struct {
		code   int
		answer string
	}
	type test struct {
		name    string
		request string
		method  string
		token   string
		dbFlag  bool
		books   string
		err     error
		want    want
	}
	tests := []test{
		{
			name:    "Test 'SaveBooksHandler' #1; StatusOK",
			request: "/books/adds",
			method:  http.MethodPost,
			token:   generateValidToken(),
			dbFlag:  true,
			books:   `[{"bId":1,"title":"Book1","author":"Author1"},{"bId":2,"title":"Book2","author":"Author2"}]`,
			want: want{
				code:   http.StatusOK,
				answer: `{"message": "All books saved"}`,
			},
		},
		{
			name:    "Test 'SaveBooksHandler' #2; Missing Authorization header",
			request: "/books/adds",
			method:  http.MethodPost,
			token:   "",
			books:   `[{"bId":1,"title":"Book1","author":"Author1"},{"bId":2,"title":"Book2","author":"Author2"}]`,
			dbFlag:  false,
			want: want{
				code:   http.StatusUnauthorized,
				answer: `{"message":"Authorization header missing"}`,
			},
		},
		{
			name:    "Test 'SaveBooksHandler' #3; Invalid token",
			request: "/books/adds",
			method:  http.MethodPost,
			token:   "invalid token",
			books:   `[{"bId":1,"title":"Book1","author":"Author1"},{"bId":2,"title":"Book2","author":"Author2"}]`,
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"message":"Bad auth token","error":"token contains an invalid number of segments"}`,
			},
		},
		{
			name:    "Test 'SaveBooksHandler' #4; Invalid JSON",
			request: "/books/adds",
			method:  http.MethodPost,
			token:   generateValidToken(),
			books:   `{"invalid":`, // Некорректный JSON
			dbFlag:  false,
			want: want{
				code:   http.StatusBadRequest,
				answer: `{"error":"unexpected EOF"}`,
			},
		},
		{
			name:    "Test 'SaveBooksHandler' #5; Database error",
			request: "/books/adds",
			method:  http.MethodPost,
			token:   generateValidToken(),
			books:   `[{"bId":1,"title":"Book1","author":"Author1"},{"bId":2,"title":"Book2","author":"Author2"}]`,
			err:     errors.New("db error"),
			dbFlag:  true,
			want: want{
				code:   http.StatusInternalServerError,
				answer: `{"error":"db error"}`,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_server.NewMockRepository(ctrl)
			if tc.dbFlag {
				m.EXPECT().SaveBooks(gomock.Any(), gomock.Any()).Return(tc.err)
				srv.Db = m
			}

			req := resty.New().R()
			req.Method = tc.method
			req.Body = tc.books
			req.Header.Add("Authorization", tc.token)
			req.URL = httpSrv.URL + tc.request
			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, tc.want.code, resp.StatusCode())
			if tc.want.answer != "" {
				assert.JSONEq(t, tc.want.answer, string(resp.Body()))
			}
		})
	}
}

func TestGetBooksByUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_server.NewMockRepository(ctrl)
	srv := &Server{Db: m}
	r := gin.Default()
	r.GET("/books", srv.GetBooksByUserHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()

	type want struct {
		code   int
		answer string
	}
	type test struct {
		name    string
		request string
		method  string
		token   string
		uId     string
		books   []models.Book
		dbFlag bool
		err     error
		want    want
	}
	tests := []test{
		{
			name:    "Test 'GetBooksByUserHandler' #1; Success",
			request: "/books",
			method:  http.MethodGet,
			token:   generateValidToken(),
			uId:     "1",
			dbFlag:  true,
			err:     nil,
			books: []models.Book{
				{
					BID:    1,
					Title:  "War and Peace",
					Author: "Lev Tolstoy",
					UID:    1,
				},
				{
					BID:    2,
					Title:  "Diary of a monkey",
					Author: "Jain Birkin",
					UID:    1,
				},
			},
			want: want{
				code: http.StatusOK,
				answer: `{"message":"Books retrieved","books":[{"bId":1,"title":"War and Peace","author":"Lev Tolstoy","uId":1},{"bId":2,"title":"Diary of a monkey","author":"Jain Birkin","uId":1}]}`,
			},
		},
		{
			name:    "Test 'GetBooksByUserHandler' #2; Missing Authorization header",
			request: "/books",
			method:  http.MethodGet,
			token:   "",
			uId:     "1",
			want: want{
				code:   http.StatusUnauthorized,
				answer: `{"message":"Authorization header missing"}`,
			},
		},
		{
			name:    "Test 'GetBooksByUserHandler' #3; Invalid token",
			request: "/books",
			method:  http.MethodGet,
			token:   "invalid.token",
			uId:     "1",
			want: want{
				code:   http.StatusUnauthorized,
				answer: `{"message":"Invalid token","error":"token contains an invalid number of segments"}`,
			},
		},
		{
			name:    "Test 'GetBooksByUserHandler' #4; Invalid user ID",
			request: "/books",
			method:  http.MethodGet,
			token:   generateValidToken(),
			uId:     "invalid",
			want: want{
				code:   http.StatusInternalServerError,
				answer: `{"error":"strconv.Atoi: parsing \"invalid\": invalid syntax"}`,
			},
		},
		{
			name:    "Test 'GetBooksByUserHandler' #5; Database error",
			request: "/books",
			method:  http.MethodGet,
			token:   generateValidToken(),
			uId:     "1",
			err:     errors.New("db error"),
			want: want{
				code:   http.StatusNotFound,
				answer: `{"message": "Books not found", "error": "db error"}`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.method == http.MethodGet && tc.token != "" {
				if tc.uId != "invalid" {
					uId, _ := strconv.Atoi(tc.uId)
					if tc.err != nil {
						m.EXPECT().GetBooksByUser(uId).Return(nil, tc.err).Times(1)
					} else {
						m.EXPECT().GetBooksByUser(uId).Return(tc.books, nil).Times(1)
					}
				}
			}

			req := httptest.NewRequest(tc.method, httpSrv.URL+tc.request+"?uid="+tc.uId, nil)
			if tc.token != "" {
				req.Header.Set("Authorization", tc.token)
			}

			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			assert.Equal(t, tc.want.code, resp.Code)
			assert.JSONEq(t, tc.want.answer, resp.Body.String())
		})
	}
}