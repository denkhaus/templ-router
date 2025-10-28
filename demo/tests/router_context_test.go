package tests

import (
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("RouterContext Query Parameter Integration", func() {
	ginkgo.Context("Query Parameter Demo Route", func() {
		ginkgo.It("should handle query parameters correctly in English", func() {
			testCases := []struct {
				name         string
				queryParams  string
				expectedPage string
				expectedSize string
				expectedSort string
				expectedFilter string
			}{
				{
					name:         "Basic pagination",
					queryParams:  "?page=2&pageSize=25&sort=name",
					expectedPage: "2",
					expectedSize: "25", 
					expectedSort: "name",
					expectedFilter: "",
				},
				{
					name:         "With filter",
					queryParams:  "?page=3&pageSize=50&filter=premium&sort=date",
					expectedPage: "3",
					expectedSize: "50",
					expectedSort: "date", 
					expectedFilter: "premium",
				},
				{
					name:         "Default values",
					queryParams:  "",
					expectedPage: "1",
					expectedSize: "10",
					expectedSort: "name",
					expectedFilter: "",
				},
				{
					name:         "Partial parameters",
					queryParams:  "?page=5&filter=active",
					expectedPage: "5",
					expectedSize: "10", // default
					expectedSort: "name", // default
					expectedFilter: "active",
				},
			}

			for _, tc := range testCases {
				ginkgo.By(fmt.Sprintf("Testing %s with params: %s", tc.name, tc.queryParams))
				
				url := fmt.Sprintf("%s/en/query-demo%s", baseURL, tc.queryParams)
				gomega.Expect(page.Navigate(url)).To(gomega.Succeed())
				
				// Wait for page to load
				var html string
				gomega.Eventually(func() string {
					html, _ = page.HTML()
					return html
				}, PageLoadTimeout).Should(gomega.ContainSubstring("Query Parameter Demo"))
				
				// Check for server errors
				gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
				gomega.Expect(html).ShouldNot(gomega.ContainSubstring("Internal Server Error"))
				
				// Verify query parameters are displayed correctly (simplified check)
				gomega.Expect(html).Should(gomega.ContainSubstring("Page:</span>"))
				gomega.Expect(html).Should(gomega.ContainSubstring(fmt.Sprintf(">%s</span>", tc.expectedPage)))
				gomega.Expect(html).Should(gomega.ContainSubstring("Page Size:</span>"))
				gomega.Expect(html).Should(gomega.ContainSubstring(fmt.Sprintf(">%s</span>", tc.expectedSize)))
				gomega.Expect(html).Should(gomega.ContainSubstring("Sort:</span>"))
				gomega.Expect(html).Should(gomega.ContainSubstring(fmt.Sprintf(">%s</span>", tc.expectedSort)))
				
				// Check filter (can be empty) - simplified
				gomega.Expect(html).Should(gomega.ContainSubstring("Filter:</span>"))
				if tc.expectedFilter != "" {
					gomega.Expect(html).Should(gomega.ContainSubstring(tc.expectedFilter))
				} else {
					gomega.Expect(html).Should(gomega.ContainSubstring("none</em>"))
				}
				
				// Verify user data is still displayed (GetUserData method is called)
				gomega.Expect(html).Should(gomega.ContainSubstring("GetUserData Method Called!"))
				gomega.Expect(html).Should(gomega.ContainSubstring("getuserdata@example.com"))
			}
		})

		ginkgo.It("should handle query parameters correctly in German", func() {
			url := fmt.Sprintf("%s/de/query-demo?page=2&pageSize=20&filter=premium&sort=date", baseURL)
			gomega.Expect(page.Navigate(url)).To(gomega.Succeed())
			
			// Wait for page to load
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Query Parameter Demo"))
			
			// Check for server errors
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
			
			// Verify query parameters
			gomega.Expect(html).Should(gomega.ContainSubstring("Page:</span>\n\t\t\t\t\t\t\t<span class=\"text-green-600 font-mono\">2</span>"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Page Size:</span>\n\t\t\t\t\t\t\t<span class=\"text-green-600 font-mono\">20</span>"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Sort:</span>\n\t\t\t\t\t\t\t<span class=\"text-green-600 font-mono\">date</span>"))
			gomega.Expect(html).Should(gomega.ContainSubstring("premium"))
			
			// Verify German user data (GetUserData method is called)
			gomega.Expect(html).Should(gomega.ContainSubstring("GetUserData Methode aufgerufen!"))
			gomega.Expect(html).Should(gomega.ContainSubstring("getuserdata@beispiel.de"))
		})

		ginkgo.It("should display test links correctly", func() {
			// Start at the demo page
			url := fmt.Sprintf("%s/en/query-demo", baseURL)
			gomega.Expect(page.Navigate(url)).To(gomega.Succeed())
			
			// Wait for page to load
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Test Different Query Parameters"))
			
			// Verify test links are present
			gomega.Expect(html).Should(gomega.ContainSubstring("?page=1&pageSize=5&sort=name"))
			gomega.Expect(html).Should(gomega.ContainSubstring("?page=2&pageSize=20&sort=date&filter=active"))
			gomega.Expect(html).Should(gomega.ContainSubstring("?page=3&pageSize=50&filter=premium&sort=priority"))
			
			// Test direct navigation to one of the parameter combinations
			testURL := fmt.Sprintf("%s/en/query-demo?page=2&pageSize=20&sort=date&filter=active", baseURL)
			gomega.Expect(page.Navigate(testURL)).To(gomega.Succeed())
			
			// Verify the parameters are reflected
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.And(
				gomega.ContainSubstring("Page:</span>\n\t\t\t\t\t\t\t<span class=\"text-green-600 font-mono\">2</span>"),
				gomega.ContainSubstring("Page Size:</span>\n\t\t\t\t\t\t\t<span class=\"text-green-600 font-mono\">20</span>"),
				gomega.ContainSubstring("Sort:</span>\n\t\t\t\t\t\t\t<span class=\"text-green-600 font-mono\">date</span>"),
				gomega.ContainSubstring("active"),
			))
		})
	})

	ginkgo.Context("RouterContext vs Legacy Parameter Comparison", func() {
		ginkgo.It("should demonstrate RouterContext advantages over map[string]string", func() {
			// Test the new query-demo route (uses RouterContext)
			url := fmt.Sprintf("%s/en/query-demo?page=3&pageSize=15&filter=test&sort=priority", baseURL)
			gomega.Expect(page.Navigate(url)).To(gomega.Succeed())
			
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Query Parameter Demo"))
			
			// Verify that RouterContext cleanly separates URL and query parameters
			gomega.Expect(html).Should(gomega.ContainSubstring("URL Parameters")) // Section header
			gomega.Expect(html).Should(gomega.ContainSubstring("Query Parameters")) // Section header
			
			// URL parameters should be in URL section
			gomega.Expect(html).Should(gomega.ContainSubstring("Locale:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("User ID:"))
			
			// Query parameters should be in Query section with correct values
			gomega.Expect(html).Should(gomega.ContainSubstring("Page:</span>\n\t\t\t\t\t\t\t<span class=\"text-green-600 font-mono\">3</span>"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Page Size:</span>\n\t\t\t\t\t\t\t<span class=\"text-green-600 font-mono\">15</span>"))
			gomega.Expect(html).Should(gomega.ContainSubstring("test"))
			gomega.Expect(html).Should(gomega.ContainSubstring("priority"))
			
			// Verify code example is shown
			gomega.Expect(html).Should(gomega.ContainSubstring("RouterContext Code Example"))
			gomega.Expect(html).Should(gomega.ContainSubstring("routerCtx.GetURLParam"))
			gomega.Expect(html).Should(gomega.ContainSubstring("routerCtx.GetQueryParam"))
		})
	})
})