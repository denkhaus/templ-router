package tests

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("UserWithId DataService Integration", func() {
	ginkgo.Context("Data Route with UserWithIdDataService", func() {
		ginkgo.It("should load data page with UserWithIdDataService content in English", func() {
			testUserIDs := []string{"123", "456", "789", "admin", "test-user", "abc123"}

			for _, userID := range testUserIDs {
				ginkgo.By(fmt.Sprintf("Testing UserWithId data for ID: %s", userID))
				
				url := fmt.Sprintf("%s/en/data/%s", baseURL, userID)
				gomega.Expect(page.Navigate(url)).To(gomega.Succeed())
				
				// Wait for page to load and get HTML
				var html string
				gomega.Eventually(func() string {
					html, _ = page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring("<html"))
				
				// Check for server errors
				gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
				gomega.Expect(html).ShouldNot(gomega.ContainSubstring("Internal Server Error"))
				
				// Check that the user ID is correctly displayed
				gomega.Expect(html).Should(gomega.ContainSubstring(userID))
				
				// Check for English translations
				gomega.Expect(html).Should(gomega.ContainSubstring("Data Information"))
				gomega.Expect(html).Should(gomega.ContainSubstring("User ID"))
				gomega.Expect(html).Should(gomega.ContainSubstring("Current Language"))
				
				// Check that locale is correctly set to English
				gomega.Expect(html).Should(gomega.ContainSubstring("en"))
			}
		})

		ginkgo.It("should load data page with UserWithIdDataService content in German", func() {
			testUserIDs := []string{"123", "456", "789", "admin", "test-user", "abc123"}

			for _, userID := range testUserIDs {
				ginkgo.By(fmt.Sprintf("Testing German UserWithId data for ID: %s", userID))
				
				url := fmt.Sprintf("%s/de/data/%s", baseURL, userID)
				gomega.Expect(page.Navigate(url)).To(gomega.Succeed())
				
				// Wait for page to load and get HTML
				var html string
				gomega.Eventually(func() string {
					html, _ = page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring("<html"))
				
				// Check for server errors
				gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
				gomega.Expect(html).ShouldNot(gomega.ContainSubstring("Internal Server Error"))
				
				// Check that the user ID is correctly displayed
				gomega.Expect(html).Should(gomega.ContainSubstring(userID))
				
				// Check for German translations
				gomega.Expect(html).Should(gomega.ContainSubstring("Dateninformationen"))
				gomega.Expect(html).Should(gomega.ContainSubstring("Benutzer ID"))
				gomega.Expect(html).Should(gomega.ContainSubstring("Aktuelle Sprache"))
				
				// Check that locale is correctly set to German
				gomega.Expect(html).Should(gomega.ContainSubstring("de"))
			}
		})

		ginkgo.It("should display different user data for different IDs", func() {
			// Test that different user IDs show different content
			gomega.Expect(page.Navigate(baseURL + "/en/data/123")).To(gomega.Succeed())
			
			user123Content, err := page.HTML()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			
			gomega.Expect(page.Navigate(baseURL + "/en/data/456")).To(gomega.Succeed())
			
			user456Content, err := page.HTML()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			
			// Content should be different (at least the user ID)
			gomega.Expect(user123Content).NotTo(gomega.Equal(user456Content))
			gomega.Expect(user123Content).To(gomega.ContainSubstring("123"))
			gomega.Expect(user456Content).To(gomega.ContainSubstring("456"))
		})

		ginkgo.It("should correctly handle locale switching", func() {
			userID := "456"
			
			// Load English version
			gomega.Expect(page.Navigate(fmt.Sprintf("%s/en/data/%s", baseURL, userID))).To(gomega.Succeed())
			
			englishContent, err := page.HTML()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			
			// Load German version
			gomega.Expect(page.Navigate(fmt.Sprintf("%s/de/data/%s", baseURL, userID))).To(gomega.Succeed())
			
			germanContent, err := page.HTML()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			
			// Both should contain the same user ID
			gomega.Expect(englishContent).To(gomega.ContainSubstring(userID))
			gomega.Expect(germanContent).To(gomega.ContainSubstring(userID))
			
			// But different language content
			gomega.Expect(englishContent).To(gomega.ContainSubstring("Data Information"))
			gomega.Expect(germanContent).To(gomega.ContainSubstring("Dateninformationen"))
			
			// And different locale indicators
			gomega.Expect(englishContent).To(gomega.ContainSubstring("en"))
			gomega.Expect(germanContent).To(gomega.ContainSubstring("de"))
		})

		ginkgo.It("should work with valid user IDs and edge cases", func() {
			// Test valid user IDs (based on validation pattern "[0-9a-z]*")
			// Note: Empty string is excluded as it creates invalid URL path
			validIDs := []string{"0", "999", "a", "z", "123abc", "abc123", "000", "test123", "user1"}

			for _, userID := range validIDs {
				ginkgo.By(fmt.Sprintf("Testing valid user ID: '%s'", userID))
				
				url := fmt.Sprintf("%s/en/data/%s", baseURL, userID)
				gomega.Expect(page.Navigate(url)).To(gomega.Succeed())
				
				// Wait for page to load
				var html string
				gomega.Eventually(func() string {
					html, _ = page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring("<html"))
				
				// Should not have server errors
				gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
				
				// Should contain the user ID
				gomega.Expect(html).Should(gomega.ContainSubstring(userID))
				
				// Should always contain the basic page structure
				gomega.Expect(html).Should(gomega.ContainSubstring("Data Information"))
				gomega.Expect(html).Should(gomega.ContainSubstring("User ID"))
			}
		})

		ginkgo.It("should handle undefined user ID gracefully", func() {
			// Test what happens when userID parameter is not provided or invalid
			// This tests the DataService's handling of missing/invalid parameters
			
			url := fmt.Sprintf("%s/en/data/", baseURL) // Note the trailing slash
			gomega.Expect(page.Navigate(url)).To(gomega.Succeed())
			
			// This might result in a 404 or redirect, which is expected behavior
			// We're just ensuring the system doesn't crash
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("<html"))
			
			// The exact behavior depends on router configuration
			// But it should not crash the application
			ginkgo.By("System handles missing userID parameter without crashing")
		})
	})

	ginkgo.Context("UserWithIdDataService Method Resolution", func() {
		ginkgo.It("should call GetUserWithIdData method specifically", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/data/test-method")).To(gomega.Succeed())
			
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("<html"))
			
			// Check that the page loads successfully (indicating the correct method was called)
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
			gomega.Expect(html).Should(gomega.ContainSubstring("test-method"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Data Information"))
		})

		ginkgo.It("should pass correct parameters to UserWithIdDataService", func() {
			userID := "param-test"
			
			// Test English
			gomega.Expect(page.Navigate(fmt.Sprintf("%s/en/data/%s", baseURL, userID))).To(gomega.Succeed())
			
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("<html"))
			
			// Verify that both userId and locale parameters are correctly passed
			gomega.Expect(html).Should(gomega.ContainSubstring(userID)) // userId parameter
			gomega.Expect(html).Should(gomega.ContainSubstring("en"))   // locale parameter
			
			// Test German
			gomega.Expect(page.Navigate(fmt.Sprintf("%s/de/data/%s", baseURL, userID))).To(gomega.Succeed())
			
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("<html"))
			
			// Verify parameters for German version
			gomega.Expect(html).Should(gomega.ContainSubstring(userID)) // userId parameter
			gomega.Expect(html).Should(gomega.ContainSubstring("de"))   // locale parameter
		})
	})

	ginkgo.Context("Template Registry Integration", func() {
		ginkgo.It("should work without YAML dataservice configuration", func() {
			// This test verifies that the route works purely through Template Registry
			// without needing manual YAML dataservice configuration
			
			gomega.Expect(page.Navigate(baseURL + "/en/data/registry-test")).To(gomega.Succeed())
			
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("<html"))
			
			// Should work without errors, proving automatic DataService detection works
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("DataService not found"))
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("Template not found"))
			
			// Should contain expected content
			gomega.Expect(html).Should(gomega.ContainSubstring("registry-test"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Data Information"))
		})
	})
})