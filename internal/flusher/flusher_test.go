package flusher_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"ova-exam-api/internal/domain/entity/user"
	"ova-exam-api/internal/flusher"
	"ova-exam-api/internal/mocks"
)

var _ = Describe("Flusher", func() {

	var (
		ctrl *gomock.Controller
		mockRepo      *mocks.MockRepo
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepo(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Valid case", func() {
		var users []user.User
		BeforeEach(func() {
			mockRepo.EXPECT().AddEntities(gomock.Any()).Return(nil).Times(5)
			users = getUsers(10)
		})

		It("should all save all users", func() {
			flusher := flusher.NewFlusher(2, mockRepo)
			notSaved := flusher.Flush(users)

			gomega.Ω(notSaved).Should(gomega.BeNil())
		})
	})

	Context("Valid case", func() {
		var users []user.User
		BeforeEach(func() {
			mockRepo.EXPECT().AddEntities(gomock.Any()).Return(nil).Times(6)
			users = getUsers(11)
		})

		It("should not saved one user", func() {
			flusher := flusher.NewFlusher(2, mockRepo)
			notSaved := flusher.Flush(users)

			gomega.Ω(notSaved).Should(gomega.HaveLen(0))
		})
	})
})

func getUsers(userCount int) []user.User{
	var users = make([]user.User, userCount)
	for i := 1; i < userCount; i++ {
		users[i] = user.User{
			UserId: uint64(i),
		}
	}

	return users
}
