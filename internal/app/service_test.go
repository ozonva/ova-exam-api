package ovaexamapi_test

import (
	"database/sql"
	"database/sql/driver"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	ovaexamapi "ova-exam-api/internal/app"
	"time"

	"ova-exam-api/internal/repo"
	desc "ova-exam-api/pkg/github.com/ozonva/ova-exam-api/pkg/ova-exam-api"
)
type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

var _ = Describe("Api test", func() {

	var (
		ctrl *gomock.Controller
		dbMock 					*sql.DB
		userRepo repo.Repo
		ovaExamAPI desc.UsersServer
		sqlMock sqlmock.Sqlmock
	)

	BeforeEach(func() {
		db, mock, err := sqlmock.New()
		if err != nil {
			panic(err)
		}
		dbMock = db
		sqlMock = mock
	})

	AfterEach(func() {
		dbMock.Close()
		ctrl.Finish()
	})

	Context("CreateUserV1Request", func() {
		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			userRepo = repo.NewRepo(dbMock)
			ovaExamAPI = ovaexamapi.NewOvaExamAPI(userRepo)
			sqlMock.ExpectExec("INSERT INTO users").WithArgs("testEmail", "testPassword", AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))
		})

		It("ValidCase", func() {
			createUserReq := desc.CreateUserV1Request {
				Email:    "testEmail",
				Password: "testPassword",
			}

			res, err := ovaExamAPI.CreateUserV1(nil, &createUserReq)

			gomega.Ω(err).Should(gomega.BeNil())
			gomega.Ω(*res).Should(gomega.Equal(empty.Empty{}))
		})
	})

	Context("DescribeUserV1Request", func() {
		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			userRepo = repo.NewRepo(dbMock)
			ovaExamAPI = ovaexamapi.NewOvaExamAPI(userRepo)
			rows := sqlmock.NewRows([]string{"id", "Email", "Password", "createdAt", "updatedAt"}).
				AddRow(1, "email", "password", time.Now(), time.Now())
			sqlMock.ExpectQuery("^SELECT (.+) FROM").WillReturnRows(rows)
		})

		It("ValidCase", func() {
			describeUserV1Request := desc.DescribeUserV1Request {
				UserId: 1,
			}

			res, err := ovaExamAPI.DescribeUserV1(nil, &describeUserV1Request)
			userV1 := desc.UserV1{
				UserId:   1,
				Email:    "email",
				Password: "password",
			}

			gomega.Ω(err).Should(gomega.BeNil())
			gomega.Ω(*res).Should(gomega.Equal(userV1))
		})
	})

	Context("ListUsersV1", func() {
		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			userRepo = repo.NewRepo(dbMock)
			ovaExamAPI = ovaexamapi.NewOvaExamAPI(userRepo)
			rows := sqlmock.NewRows([]string{"id", "Email", "Password", "createdAt", "updatedAt"}).
				AddRow(1, "email1", "password1", time.Now(), time.Now()).
				AddRow(1, "email2", "password2", time.Now(), time.Now())
			sqlMock.ExpectQuery("^SELECT (.+) FROM").WillReturnRows(rows)
		})

		It("ValidCase", func() {
			res, err := ovaExamAPI.ListUsersV1(nil, &empty.Empty{})
			gomega.Ω(err).Should(gomega.BeNil())
			gomega.Ω(res.Users).Should(gomega.HaveLen(2))
		})
	})

	Context("RemoveUserV1", func() {
		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			userRepo = repo.NewRepo(dbMock)
			ovaExamAPI = ovaexamapi.NewOvaExamAPI(userRepo)
			sqlMock.ExpectExec("DELETE FROM users").WithArgs(2).WillReturnResult(sqlmock.NewResult(1, 1))
		})

		It("ValidCase", func() {
			req := desc.RemoveUserV1Request{
				UserId: 2,
			}
			res, err := ovaExamAPI.RemoveUserV1(nil, &req)

			gomega.Ω(err).Should(gomega.BeNil())
			gomega.Ω(res).Should(gomega.Equal(&emptypb.Empty{}))
		})
	})
})

