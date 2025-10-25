package tests

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Content Validation", func() {
	ginkgo.Context("Server Error Prevention", func() {
		routes := []string{
			"/",
			"/en",
			"/de", 
			"/en/dashboard",
			"/de/dashboard",
			"/en/admin",
			"/de/admin",
			"/en/user/123",
			"/de/user/123",
			"/en/product/laptop",
			"/de/product/laptop",
			"/en/test",
			"/de/test",
			"/login",
			"/signup",
		}
		
		for _, route := range routes {
			route := route // capture loop variable
			ginkgo.It(fmt.Sprintf("should load %s without server errors", route), func() {
				gomega.Expect(page.Navigate(baseURL + route)).To(gomega.Succeed())
				
				// Wait for page to load and get HTML once
				var html string
				gomega.Eventually(func() string {
					html, _ = page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring("<html"))
				
				// Check for server errors in the loaded HTML
				gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
				gomega.Expect(html).ShouldNot(gomega.ContainSubstring("Internal Server Error"))
				gomega.Expect(html).ShouldNot(gomega.ContainSubstring("panic:"))
				gomega.Expect(html).ShouldNot(gomega.ContainSubstring("runtime error:"))
				
				// Check that page has meaningful content
				gomega.Expect(html).Should(gomega.ContainSubstring("</html>"))
			})
		}
	})
	
	ginkgo.Context("Content Quality", func() {
		ginkgo.It("should have proper HTML structure on all pages", func() {
			routes := []string{"/en", "/de", "/en/dashboard", "/de/dashboard"}
			
			for _, route := range routes {
				ginkgo.By(fmt.Sprintf("Checking HTML structure for %s", route))
				gomega.Expect(page.Navigate(baseURL + route)).To(gomega.Succeed())
				
				// Wait for page to load and get HTML once
				var html string
				gomega.Eventually(func() string {
					html, _ = page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring("<html"))
				
				// Check all HTML structure elements in the loaded HTML
				gomega.Expect(html).Should(gomega.ContainSubstring("<head>"))
				gomega.Expect(html).Should(gomega.ContainSubstring("<body"))
				gomega.Expect(html).Should(gomega.ContainSubstring("</body>"))
			}
		})
		
		ginkgo.It("should have navigation elements on main pages", func() {
			routes := []string{"/en", "/de", "/en/dashboard", "/de/dashboard"}
			
			for _, route := range routes {
				ginkgo.By(fmt.Sprintf("Checking navigation for %s", route))
				gomega.Expect(page.Navigate(baseURL + route)).To(gomega.Succeed())
				
				// Wait for page to load and get HTML once
				var html string
				gomega.Eventually(func() string {
					html, _ = page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring("<nav"))
				
				// Check for footer in the loaded HTML
				gomega.Expect(html).Should(gomega.ContainSubstring("<footer"))
			}
		})
	})
})