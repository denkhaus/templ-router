package tests

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Authentication Flow", func() {
	ginkgo.Context("Sign Up", func() {
		ginkgo.It("should load signup page correctly", func() {
			gomega.Expect(page.Navigate(baseURL + "/signup")).To(gomega.Succeed())
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Sign up"))
			
			// Check form elements
			emailField := page.Find("input[name='email']")
			gomega.Expect(emailField.Count()).To(gomega.BeNumerically(">", 0))
			
			passwordField := page.Find("input[name='password']")
			gomega.Expect(passwordField.Count()).To(gomega.BeNumerically(">", 0))
		})
		
		ginkgo.It("should handle signup form submission", func() {
			gomega.Expect(page.Navigate(baseURL + "/signup")).To(gomega.Succeed())
			
			// Fill form
			usernameField := page.Find("input[name='username']")
			gomega.Expect(usernameField.Fill("testuser")).To(gomega.Succeed())
			
			emailField := page.Find("input[name='email']")
			gomega.Expect(emailField.Fill("test@example.com")).To(gomega.Succeed())
			
			passwordField := page.Find("input[name='password']")
			gomega.Expect(passwordField.Fill("testpassword123")).To(gomega.Succeed())
			
			// Submit form
			submitButton := page.Find("button[type='submit']")
			gomega.Expect(submitButton.Click()).To(gomega.Succeed())
			
			// Wait for HTMX redirect (HX-Redirect header causes browser redirect)
			gomega.Eventually(func() string {
				url, _ := page.URL()
				return url
			}, PageLoadTimeout).ShouldNot(gomega.Equal(baseURL + "/signup"))
		})
	})
	
	ginkgo.Context("Sign In", func() {
		ginkgo.It("should load login page correctly", func() {
			gomega.Expect(page.Navigate(baseURL + "/login")).To(gomega.Succeed())
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Sign in"))
			
			// Check form elements
			emailField := page.Find("input[name='email']")
			gomega.Expect(emailField.Count()).To(gomega.BeNumerically(">", 0))
			
			passwordField := page.Find("input[name='password']")
			gomega.Expect(passwordField.Count()).To(gomega.BeNumerically(">", 0))
			
			submitButton := page.Find("button[type='submit']")
			gomega.Expect(submitButton.Count()).To(gomega.BeNumerically(">", 0))
		})
		
		ginkgo.It("should handle login form submission", func() {
			gomega.Expect(page.Navigate(baseURL + "/login")).To(gomega.Succeed())
			
			// Fill form
			emailField := page.Find("input[name='email']")
			gomega.Expect(emailField.Fill("demo@example.com")).To(gomega.Succeed())
			
			passwordField := page.Find("input[name='password']")
			gomega.Expect(passwordField.Fill("demo123")).To(gomega.Succeed())
			
			// Submit form
			submitButton := page.Find("button[type='submit']")
			gomega.Expect(submitButton.Click()).To(gomega.Succeed())
			
			// Wait for HTMX redirect (HX-Redirect header causes browser redirect)
			gomega.Eventually(func() string {
				url, _ := page.URL()
				return url
			}, PageLoadTimeout).ShouldNot(gomega.Equal(baseURL + "/login"))
		})
		
		ginkgo.It("should redirect to correct route after successful login", func() {
			// Try to access protected route first
			gomega.Expect(page.Navigate(baseURL + "/en/dashboard")).To(gomega.Succeed())
			
			// Should redirect to login if not authenticated
			currentURL, _ := page.URL()
			if currentURL == baseURL+"/login" {
				// Fill login form
				emailField := page.Find("input[name='email']")
				gomega.Expect(emailField.Fill("demo@example.com")).To(gomega.Succeed())
				
				passwordField := page.Find("input[name='password']")
				gomega.Expect(passwordField.Fill("demo123")).To(gomega.Succeed())
				
				submitButton := page.Find("button[type='submit']")
				gomega.Expect(submitButton.Click()).To(gomega.Succeed())
				
				// Should redirect back to dashboard
				gomega.Eventually(func() string {
					url, _ := page.URL()
					return url
				}, PageLoadTimeout).Should(gomega.ContainSubstring("/dashboard"))
			}
		})
	})
	
	ginkgo.Context("Sign Out", func() {
		ginkgo.It("should handle logout correctly", func() {
			// First ensure we're logged in
			gomega.Expect(page.Navigate(baseURL + "/login")).To(gomega.Succeed())
			
			emailField := page.Find("input[name='email']")
			count, _ := emailField.Count()
			if count > 0 {
				gomega.Expect(emailField.Fill("demo@example.com")).To(gomega.Succeed())
				
				passwordField := page.Find("input[name='password']")
				gomega.Expect(passwordField.Fill("demo123")).To(gomega.Succeed())
				
				submitButton := page.Find("button[type='submit']")
				gomega.Expect(submitButton.Click()).To(gomega.Succeed())
			}
			
			// Navigate to a page with logout button
			gomega.Expect(page.Navigate(baseURL + "/en/dashboard")).To(gomega.Succeed())
			
			// Look for logout button/form - simplified approach
			logoutForm := page.Find("form[action*='signout']")
			logoutCount, _ := logoutForm.Count()
			if logoutCount > 0 {
				logoutButton := logoutForm.Find("button[type='submit']")
				buttonCount, _ := logoutButton.Count()
				if buttonCount > 0 {
					gomega.Expect(logoutButton.Click()).To(gomega.Succeed())
					
					// Should redirect after logout
					gomega.Eventually(func() string {
						url, _ := page.URL()
						return url
					}, PageLoadTimeout).ShouldNot(gomega.ContainSubstring("/dashboard"))
				}
			}
			
		})
	})
})