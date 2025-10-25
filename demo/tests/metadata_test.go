package tests

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Metadata and Layout System", func() {
	ginkgo.Context("Page Metadata", func() {
		ginkgo.It("should have correct page titles", func() {
			// Test dashboard title
			gomega.Expect(page.Navigate(baseURL + "/en/dashboard")).To(gomega.Succeed())
			gomega.Eventually(func() string {
				title, _ := page.Title()
				return title
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Dashboard"))
			
			// Test German dashboard title
			gomega.Expect(page.Navigate(baseURL + "/de/dashboard")).To(gomega.Succeed())
			gomega.Eventually(func() string {
				title, _ := page.Title()
				return title
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Dashboard"))
		})
		
		ginkgo.It("should have proper meta tags", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/dashboard")).To(gomega.Succeed())
			
			// Check for viewport meta tag
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("name=\"viewport\""))
			
			// Check for charset
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("charset=\"UTF-8\""))
		})
		
		ginkgo.It("should have correct lang attribute", func() {
			// English pages should have lang="en"
			gomega.Expect(page.Navigate(baseURL + "/en/dashboard")).To(gomega.Succeed())
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("lang=\"en\""))
			
			// German pages should have lang="de"
			gomega.Expect(page.Navigate(baseURL + "/de/dashboard")).To(gomega.Succeed())
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("lang=\"de\""))
		})
		
		ginkgo.It("should include CSS and JS assets", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/dashboard")).To(gomega.Succeed())
			
			// Check for CSS
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("href=\"/assets/css/output.css\""))
			
			// Check for HTMX
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("src=\"/assets/js/htmx.min.js\""))
		})
	})
	
	ginkgo.Context("Layout Inheritance", func() {
		ginkgo.It("should use consistent layout across pages", func() {
			pages := []string{"/en", "/en/dashboard", "/en/admin"}
			
			for _, pagePath := range pages {
				ginkgo.By(fmt.Sprintf("Checking layout consistency for %s", pagePath))
				gomega.Expect(page.Navigate(baseURL + pagePath)).To(gomega.Succeed())
				
				// Check for common layout elements
				gomega.Eventually(func() string {
					html, _ := page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring("class=\"bg-blue-600 text-white p-4\"")) // Navbar
				gomega.Eventually(func() string {
					html, _ := page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring("<footer")) // Footer
				gomega.Eventually(func() string {
					html, _ := page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring("<main")) // Main content area
			}
		})
		
		ginkgo.It("should have language switcher in layout", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/dashboard")).To(gomega.Succeed())
			
			// Check for language switcher elements
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("EN"))
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("DE"))
		})
		
		ginkgo.It("should have consistent navigation across locales", func() {
			// Check English navigation
			gomega.Expect(page.Navigate(baseURL + "/en/dashboard")).To(gomega.Succeed())
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("href=\"/en/dashboard\""))
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("href=\"/en/admin\""))
			
			// Check German navigation
			gomega.Expect(page.Navigate(baseURL + "/de/dashboard")).To(gomega.Succeed())
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("href=\"/de/dashboard\""))
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("href=\"/de/admin\""))
		})
	})
	
	ginkgo.Context("Theme and Styling", func() {
		ginkgo.It("should have consistent theme across pages", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/dashboard")).To(gomega.Succeed())
			
			// Check for theme attribute
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("data-theme="))
			
			// Check for consistent styling classes
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("bg-gray-50"))
		})
	})
})