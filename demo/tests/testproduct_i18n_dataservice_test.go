package tests

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("TestProduct Route - DataService + i18n Integration", func() {
	ginkgo.Context("DataService Integration", func() {
		ginkgo.It("should load testproduct page with ProductDataService data", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/testproduct")).To(gomega.Succeed())
			
			// Wait for page to load and get HTML once
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("<html"))
			
			// Check for server errors
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("Internal Server Error"))
			
			// Check that ProductDataService data is displayed
			gomega.Expect(html).Should(gomega.ContainSubstring("Demo Product"))
			gomega.Expect(html).Should(gomega.ContainSubstring("99.99"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Electronics"))
			gomega.Expect(html).Should(gomega.ContainSubstring("demonstration product loaded via GetData method"))
		})
		
		ginkgo.It("should show that GetData method was called (not a specific method)", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/testproduct")).To(gomega.Succeed())
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("GetData method"))
		})
	})
	
	ginkgo.Context("i18n Translation Integration - English", func() {
		ginkgo.BeforeEach(func() {
			gomega.Expect(page.Navigate(baseURL + "/en/testproduct")).To(gomega.Succeed())
		})
		
		ginkgo.It("should display English translations from YAML file", func() {
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Product DataService Test"))
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Product ID:"))
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Testing GetData method ONLY"))
		})
		
		ginkgo.It("should display English product information labels", func() {
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Product Information"))
			
			gomega.Expect(html).Should(gomega.ContainSubstring("Product Name"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Description"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Price"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Category"))
			gomega.Expect(html).Should(gomega.ContainSubstring("In Stock"))
		})
		
		ginkgo.It("should display English test information section", func() {
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("DataService Test Info"))
			
			gomega.Expect(html).Should(gomega.ContainSubstring("GetData Method Only Test"))
			gomega.Expect(html).Should(gomega.ContainSubstring("This template tests a DataService that ONLY implements GetData method"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Service:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Available Methods:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("GetData() only"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Expected Behavior:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Should call GetData method"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Welcome to the Product DataService demonstration!"))
		})
		
		ginkgo.It("should display Yes/No in English for stock status", func() {
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Yes"))
		})
	})
	
	ginkgo.Context("i18n Translation Integration - German", func() {
		ginkgo.BeforeEach(func() {
			gomega.Expect(page.Navigate(baseURL + "/de/testproduct")).To(gomega.Succeed())
		})
		
		ginkgo.It("should display German translations from YAML file", func() {
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Produkt DataService Test"))
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Produkt-ID:"))
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Teste nur GetData Methode"))
		})
		
		ginkgo.It("should display German product information labels", func() {
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Produktinformationen"))
			
			gomega.Expect(html).Should(gomega.ContainSubstring("Produktname"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Beschreibung"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Preis"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Kategorie"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Auf Lager"))
		})
		
		ginkgo.It("should display German test information section", func() {
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("DataService Test Info"))
			
			gomega.Expect(html).Should(gomega.ContainSubstring("GetData Methoden Test"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Dieses Template testet einen DataService, der NUR die GetData Methode implementiert"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Service:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Verfügbare Methoden:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Nur GetData()"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Erwartetes Verhalten:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Sollte GetData Methode aufrufen"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Willkommen zur Produkt DataService Demonstration!"))
		})
		
		ginkgo.It("should display Ja/Nein in German for stock status", func() {
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Ja"))
		})
	})
	
	
	ginkgo.Context("DataService + i18n Combined Functionality", func() {
		ginkgo.It("should show DataService data with English translations", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/testproduct")).To(gomega.Succeed())
			
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Demo Product"))
			
			// Both DataService data AND English translations should be present
			gomega.Expect(html).Should(gomega.ContainSubstring("Demo Product"))  // From DataService
			gomega.Expect(html).Should(gomega.ContainSubstring("Product Information"))  // From i18n EN
			gomega.Expect(html).Should(gomega.ContainSubstring("99.99"))  // From DataService
			gomega.Expect(html).Should(gomega.ContainSubstring("Product Name"))  // From i18n EN
		})
		
		ginkgo.It("should show DataService data with German translations", func() {
			gomega.Expect(page.Navigate(baseURL + "/de/testproduct")).To(gomega.Succeed())
			
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Demo Product"))
			
			// Both DataService data AND German translations should be present
			gomega.Expect(html).Should(gomega.ContainSubstring("Demo Product"))  // From DataService
			gomega.Expect(html).Should(gomega.ContainSubstring("Produktinformationen"))  // From i18n DE
			gomega.Expect(html).Should(gomega.ContainSubstring("99.99"))  // From DataService
			gomega.Expect(html).Should(gomega.ContainSubstring("Produktname"))  // From i18n DE
		})
		
	})
	
	ginkgo.Context("Language Switching with DataService", func() {
		ginkgo.It("should maintain DataService data when switching languages", func() {
			// Start with English
			gomega.Expect(page.Navigate(baseURL + "/en/testproduct")).To(gomega.Succeed())
			
			var htmlEN string
			gomega.Eventually(func() string {
				htmlEN, _ = page.HTML()
				return htmlEN
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Demo Product"))
			
			// Switch to German
			gomega.Expect(page.Navigate(baseURL + "/de/testproduct")).To(gomega.Succeed())
			
			var htmlDE string
			gomega.Eventually(func() string {
				htmlDE, _ = page.HTML()
				return htmlDE
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Demo Product"))
			
			// DataService data should be the same
			gomega.Expect(htmlEN).Should(gomega.ContainSubstring("Demo Product"))
			gomega.Expect(htmlDE).Should(gomega.ContainSubstring("Demo Product"))
			gomega.Expect(htmlEN).Should(gomega.ContainSubstring("99.99"))
			gomega.Expect(htmlDE).Should(gomega.ContainSubstring("99.99"))
			
			// But translations should be different
			gomega.Expect(htmlEN).Should(gomega.ContainSubstring("Product Information"))
			gomega.Expect(htmlDE).Should(gomega.ContainSubstring("Produktinformationen"))
			gomega.Expect(htmlEN).ShouldNot(gomega.ContainSubstring("Produktinformationen"))
			gomega.Expect(htmlDE).ShouldNot(gomega.ContainSubstring("Product Information"))
		})
	})
})