package tests

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Multi-Language Routing", func() {
	ginkgo.Context("English Routes", func() {
		ginkgo.It("should load English homepage", func() {
			gomega.Expect(page.Navigate(baseURL + "/en")).To(gomega.Succeed())
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Dashboard"))
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Admin"))
		})

		ginkgo.It("should load English dashboard", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/dashboard")).To(gomega.Succeed())
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Dashboard"))
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Users"))
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Projects"))
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("$12,345"))
		})

		ginkgo.It("should load English admin page", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/admin")).To(gomega.Succeed())
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Admin"))
		})
	})

	ginkgo.Context("German Routes", func() {
		ginkgo.It("should load German homepage", func() {
			gomega.Expect(page.Navigate(baseURL + "/de")).To(gomega.Succeed())
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Dashboard"))
		})

		ginkgo.It("should load German dashboard with localized content", func() {
			gomega.Expect(page.Navigate(baseURL + "/de/dashboard")).To(gomega.Succeed())
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Dashboard"))
			// Check German currency format
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("12.345"))
		})

		ginkgo.It("should load German admin page", func() {
			gomega.Expect(page.Navigate(baseURL + "/de/admin")).To(gomega.Succeed())
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Admin"))
		})
	})
})

var _ = ginkgo.Describe("Dynamic Routes", func() {
	ginkgo.Context("User Routes", func() {
		ginkgo.It("should handle dynamic user routes in English", func() {
			testUserIDs := []string{"123", "456", "789", "admin"}

			for _, userID := range testUserIDs {
				ginkgo.By(fmt.Sprintf("Testing user ID: %s", userID))
				url := fmt.Sprintf("%s/en/user/%s", baseURL, userID)
				gomega.Expect(page.Navigate(url)).To(gomega.Succeed())

				gomega.Eventually(func() string {
					html, _ := page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring("Profile"))
				gomega.Eventually(func() string {
					html, _ := page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring(userID))
				gomega.Eventually(func() string {
					html, _ := page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring("Profile"))
			}
		})

		ginkgo.It("should handle dynamic user routes in German", func() {
			testUserIDs := []string{"123", "456", "789"}

			for _, userID := range testUserIDs {
				ginkgo.By(fmt.Sprintf("Testing German user ID: %s", userID))
				url := fmt.Sprintf("%s/de/user/%s", baseURL, userID)
				gomega.Expect(page.Navigate(url)).To(gomega.Succeed())

				gomega.Eventually(func() string {
					html, _ := page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring(userID))
				gomega.Eventually(func() string {
					html, _ := page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring("Profile"))
			}
		})
	})

	ginkgo.Context("Product Routes", func() {
		ginkgo.It("should handle dynamic product routes", func() {
			productIDs := []string{"laptop", "phone", "tablet"}

			for _, productID := range productIDs {
				ginkgo.By(fmt.Sprintf("Testing product ID: %s", productID))
				url := fmt.Sprintf("%s/en/product/%s", baseURL, productID)
				gomega.Expect(page.Navigate(url)).To(gomega.Succeed())

				gomega.Eventually(func() string {
					html, _ := page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring("Product"))
				gomega.Eventually(func() string {
					html, _ := page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring(productID))
			}
		})
	})
})