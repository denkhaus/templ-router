package tests

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Simple RouterContext Test", func() {
	ginkgo.Context("Basic Query Parameter Functionality", func() {
		ginkgo.It("should load query-demo page without errors", func() {
			url := fmt.Sprintf("%s/en/query-demo", baseURL)
			gomega.Expect(page.Navigate(url)).To(gomega.Succeed())
			
			// Wait for page to load with basic content
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Query Parameter Demo"))
			
			// Get final HTML
			html, err := page.HTML()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			
			// Check for server errors
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("Internal Server Error"))
			
			// Check basic page structure
			gomega.Expect(html).Should(gomega.ContainSubstring("URL Parameters"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Query Parameters"))
			gomega.Expect(html).Should(gomega.ContainSubstring("User Data"))
			
			// Check that default values are shown (no query params)
			gomega.Expect(html).Should(gomega.ContainSubstring("Page:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Page Size:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Sort:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Filter:"))
		})

		ginkgo.It("should handle simple query parameters", func() {
			url := fmt.Sprintf("%s/en/query-demo?page=5&pageSize=15", baseURL)
			gomega.Expect(page.Navigate(url)).To(gomega.Succeed())
			
			// Wait for page to load
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Query Parameter Demo"))
			
			// Get final HTML
			html, err := page.HTML()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			
			// Check for server errors
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
			
			// Check that query parameters are reflected in the page
			gomega.Expect(html).Should(gomega.ContainSubstring(">5</span>")) // page value
			gomega.Expect(html).Should(gomega.ContainSubstring(">15</span>")) // pageSize value
			
			// Check user data is still displayed (GetUserData method is called)
			gomega.Expect(html).Should(gomega.ContainSubstring("GetUserData Method Called!"))
		})

		ginkgo.It("should work with German locale", func() {
			url := fmt.Sprintf("%s/de/query-demo?page=3&filter=test", baseURL)
			gomega.Expect(page.Navigate(url)).To(gomega.Succeed())
			
			// Wait for page to load
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Query Parameter Demo"))
			
			// Get final HTML
			html, err := page.HTML()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			
			// Check for server errors
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
			
			// Check German user data (GetUserData method is called)
			gomega.Expect(html).Should(gomega.ContainSubstring("GetUserData Methode aufgerufen!"))
			gomega.Expect(html).Should(gomega.ContainSubstring("getuserdata@beispiel.de"))
			
			// Check query parameters
			gomega.Expect(html).Should(gomega.ContainSubstring(">3</span>")) // page value
			gomega.Expect(html).Should(gomega.ContainSubstring("test")) // filter value
		})
	})
})