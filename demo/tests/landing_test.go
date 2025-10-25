package tests

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Language Selection Landing Page", func() {
	ginkgo.BeforeEach(func() {
		gomega.Expect(page.Navigate(baseURL)).To(gomega.Succeed())
	})

	ginkgo.It("should display language selection", func() {
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("LIVE HOT-RELOAD DEMO"))

		// Check language links - use more specific selectors for landing page
		englishLink := page.All("a[href='/en']")
		gomega.Expect(englishLink.Count()).To(gomega.BeNumerically(">", 0))

		germanLink := page.All("a[href='/de']")
		gomega.Expect(germanLink.Count()).To(gomega.BeNumerically(">", 0))
	})

	ginkgo.It("should navigate to English version", func() {
		englishLink := page.All("a[href='/en']").At(0) // Take first one
		gomega.Expect(englishLink.Click()).To(gomega.Succeed())
		gomega.Eventually(func() string {
			url, _ := page.URL()
			return url
		}, PageLoadTimeout).Should(gomega.ContainSubstring("/en"))
	})

	ginkgo.It("should navigate to German version", func() {
		germanLink := page.All("a[href='/de']").At(0) // Take first one
		gomega.Expect(germanLink.Click()).To(gomega.Succeed())
		gomega.Eventually(func() string {
			url, _ := page.URL()
			return url
		}, PageLoadTimeout).Should(gomega.ContainSubstring("/de"))
	})

	ginkgo.It("should show demo features section", func() {
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("Demo Features"))
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("Multi-Language Routes"))
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("Decentralized i18n"))
	})
})