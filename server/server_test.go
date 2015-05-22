package server_test

import (
	"net/http"

	"github.com/nkcraddock/numzero"
	"github.com/nkcraddock/numzero/game"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("server integration tests", func() {
	var s *ServerHarness

	req_coffee := map[string]interface{}{
		"code":   "coffee",
		"desc":   "made coffee",
		"points": 5,
	}

	BeforeEach(func() {
		authStore := numzero.NewMemoryStore()
		store := game.NewMemoryStore()
		s = NewServerHarness(authStore, store)
		s.Authenticate("username", "password")
	})

	Context("/rules", func() {
		It("gets a list of rules", func() {
			res := s.PUT("/rules", &req_coffee)

			var results []game.Rule
			res = s.GET("/rules", &results)
			Ω(res.Code).Should(Equal(http.StatusOK))
			Ω(results).Should(HaveLen(1))
		})

		It("gets an existing rule by code", func() {
			res := s.PUT("/rules", &req_coffee)

			var rule game.Rule
			res = s.GET("/rules/coffee", &rule)
			Ω(res.Code).Should(Equal(http.StatusOK))
			Ω(rule.Points).Should(Equal(5))
		})

		It("adds a new rule", func() {
			res := s.PUT("/rules", &req_coffee)
			Ω(res.Code).Should(Equal(http.StatusCreated))
		})

		It("updates an existing rule", func() {
			s.PUT("/rules", &req_coffee)

			req_modified := map[string]interface{}{
				"code":   "coffee",
				"desc":   "talked about coffee",
				"points": 1,
			}

			res := s.PUT("/rules", &req_modified)
			Ω(res.Code).Should(Equal(http.StatusCreated))

			var rule game.Rule
			res = s.GET("/rules/coffee", &rule)
			Ω(res.Code).Should(Equal(http.StatusOK))

			Ω(rule.Description).Should(Equal("talked about coffee"))
			Ω(rule.Points).Should(Equal(1))
		})
	})
})
