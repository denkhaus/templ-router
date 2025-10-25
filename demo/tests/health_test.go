package tests

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Health Check", func() {
	ginkgo.It("should respond to health endpoint", func() {
		gomega.Expect(page.Navigate(baseURL + "/api/health")).To(gomega.Succeed())
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("healthy"))
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("clean"))
	})

	ginkgo.It("should return correct JSON structure", func() {
		gomega.Expect(page.Navigate(baseURL + "/api/health")).To(gomega.Succeed())
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("architecture"))
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("dependency_injection"))
	})
})