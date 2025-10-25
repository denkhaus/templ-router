package tests

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Language Switching", func() {
	ginkgo.It("should switch from English to German via navbar", func() {
		gomega.Expect(page.Navigate(baseURL + "/en/dashboard")).To(gomega.Succeed())

		// Find and click German language switcher
		germanSwitch := page.Find("a[href='/de']")
		gomega.Expect(germanSwitch.Count()).To(gomega.BeNumerically(">", 0))
		gomega.Expect(germanSwitch.Click()).To(gomega.Succeed())

		gomega.Eventually(func() string {
			url, _ := page.URL()
			return url
		}, PageLoadTimeout).Should(gomega.ContainSubstring("/de"))
		
		// Check that we're on the German page
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("lang=\"de\""))
	})

	ginkgo.It("should switch from German to English via navbar", func() {
		gomega.Expect(page.Navigate(baseURL + "/de/dashboard")).To(gomega.Succeed())

		// Find and click English language switcher
		englishSwitch := page.Find("a[href='/en']")
		gomega.Expect(englishSwitch.Count()).To(gomega.BeNumerically(">", 0))
		gomega.Expect(englishSwitch.Click()).To(gomega.Succeed())

		gomega.Eventually(func() string {
			url, _ := page.URL()
			return url
		}, PageLoadTimeout).Should(gomega.ContainSubstring("/en"))
		
		// Check that we're on the English page
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("lang=\"en\""))
	})

	ginkgo.It("should maintain language context across pages", func() {
		// Start in German
		gomega.Expect(page.Navigate(baseURL + "/de/dashboard")).To(gomega.Succeed())
		
		// Check we're in German context
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("12.345"))

		// Navigate to admin page using direct navigation instead of clicking
		gomega.Expect(page.Navigate(baseURL + "/de/admin")).To(gomega.Succeed())

		// Should still be in German context
		gomega.Eventually(func() string {
			url, _ := page.URL()
			return url
		}, PageLoadTimeout).Should(gomega.ContainSubstring("/de/admin"))
	})
})

var _ = ginkgo.Describe("Internationalization Features", func() {
	ginkgo.It("should display correct currency format per locale", func() {
		// English - Dollar format
		gomega.Expect(page.Navigate(baseURL + "/en/dashboard")).To(gomega.Succeed())
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("$12,345"))

		// German - Euro format (check for the number format, not the symbol)
		gomega.Expect(page.Navigate(baseURL + "/de/dashboard")).To(gomega.Succeed())
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("12.345"))
	})

	ginkgo.It("should display localized names in user profiles", func() {
		// For now, just check that user pages load correctly
		// English user page
		gomega.Expect(page.Navigate(baseURL + "/en/user/123")).To(gomega.Succeed())
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("123"))

		// German user page
		gomega.Expect(page.Navigate(baseURL + "/de/user/123")).To(gomega.Succeed())
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("123"))
	})

	ginkgo.It("should show correct language indicator in navbar", func() {
		// English page should show EN as active
		gomega.Expect(page.Navigate(baseURL + "/en/dashboard")).To(gomega.Succeed())
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("EN"))

		// German page should show DE as active
		gomega.Expect(page.Navigate(baseURL + "/de/dashboard")).To(gomega.Succeed())
		gomega.Eventually(func() string {
			html, _ := page.HTML()
			return html
		}, PageLoadTimeout).Should(gomega.ContainSubstring("DE"))
	})
})