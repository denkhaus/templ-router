package tests

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Data Service Integration", func() {
	ginkgo.Context("Test Route with Data Service", func() {
		ginkgo.It("should load test page with data service content", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/test")).To(gomega.Succeed())
			
			// Wait for page to load and get HTML once
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("<html"))
			
			// Check for server errors
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("Internal Server Error"))
			
			// Check for data service content - look for any indication the page loaded with data
			gomega.Expect(html).Should(gomega.ContainSubstring("GetUserData"))
		})
		
		ginkgo.It("should load test page in German", func() {
			gomega.Expect(page.Navigate(baseURL + "/de/test")).To(gomega.Succeed())
			
			// Wait for page to load and get HTML once
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("<html"))
			
			// Check for server errors
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("Internal Server Error"))
		})
	})
	
	ginkgo.Context("User Data Service", func() {
		ginkgo.It("should load user data correctly for different IDs", func() {
			testUserIDs := []string{"123", "456", "789", "admin", "test"}
			
			for _, userID := range testUserIDs {
				ginkgo.By(fmt.Sprintf("Testing user data for ID: %s", userID))
				
				// Test English version
				url := fmt.Sprintf("%s/en/user/%s", baseURL, userID)
				gomega.Expect(page.Navigate(url)).To(gomega.Succeed())
				
				// Wait for page to load and get HTML once
				var html string
				gomega.Eventually(func() string {
					html, _ = page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring("<html"))
				
				// Check no server errors
				gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
				
				// Check user ID is displayed
				gomega.Expect(html).Should(gomega.ContainSubstring(userID))
				
				// Check profile content
				gomega.Expect(html).Should(gomega.ContainSubstring("Profile"))
			}
		})
		
		ginkgo.It("should display different user data for different IDs", func() {
			// Test that different user IDs show different content
			gomega.Expect(page.Navigate(baseURL + "/en/user/123")).To(gomega.Succeed())
			
			user123Content, err := page.HTML()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			
			gomega.Expect(page.Navigate(baseURL + "/en/user/456")).To(gomega.Succeed())
			
			user456Content, err := page.HTML()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			
			// Content should be different (at least the user ID)
			gomega.Expect(user123Content).NotTo(gomega.Equal(user456Content))
			gomega.Expect(user123Content).To(gomega.ContainSubstring("123"))
			gomega.Expect(user456Content).To(gomega.ContainSubstring("456"))
		})
	})
	
	ginkgo.Context("DataService Method Resolution Tests", func() {
		ginkgo.It("should call GetUserData preferentially over GetData", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/test")).To(gomega.Succeed())
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("GetUserData Method Called"))
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("specific-method-user"))
		})

		ginkgo.It("should work with ProductDataService that only has GetData", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/testproduct")).To(gomega.Succeed())
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Product DataService Test"))
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("GetData method ONLY"))
		})

		ginkgo.It("should call GetOrderData preferentially over GetData", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/testorder")).To(gomega.Succeed())
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("GetOrderData method called (specific method)"))
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Premium Item"))
		})

		ginkgo.It("should work with SpecificOnlyDataService that has no GetData", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/testspecific")).To(gomega.Succeed())
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Specific Method Only Test"))
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("GetSpecificData (no GetData fallback available)"))
		})

		ginkgo.It("should return error for BrokenDataService", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/testbroken")).To(gomega.Succeed())
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("500"))
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Something went wrong"))
		})
	})
})