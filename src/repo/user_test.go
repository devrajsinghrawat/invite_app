package repo_test

import (
	"github.com/devrajsinghrawat/invite_app/src/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("db operations", func() {
	Context("Check user", func() {
		It("is admin user present", func() {
			u, err := r.GetUser(&model.User{
				Username: "admin",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(u.Role).Should(Equal("ADMIN"))
		})
	})
})
