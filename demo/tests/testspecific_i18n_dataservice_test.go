package tests

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("TestSpecific Route - Specific DataService Method + i18n Integration", func() {
	ginkgo.Context("Specific DataService Method Integration", func() {
		ginkgo.It("should load testspecific page with SpecificOnlyDataService data", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/testspecific")).To(gomega.Succeed())
			
			// Wait for page to load and get HTML once
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("<html"))
			
			// Check for server errors
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("500 Internal Server Error"))
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("Internal Server Error"))
			
			// Check that SpecificOnlyDataService data is displayed
			gomega.Expect(html).Should(gomega.ContainSubstring("specific-only-123"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Specific Method Only Test"))
			gomega.Expect(html).Should(gomega.ContainSubstring("GetSpecificData"))
		})
		
		ginkgo.It("should show that GetSpecificData method was called (not GetData)", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/testspecific")).To(gomega.Succeed())
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("GetSpecificData (no GetData fallback available)"))
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("This data was loaded using ONLY GetSpecificData method"))
		})
	})
	
	ginkgo.Context("i18n Translation Integration - English", func() {
		ginkgo.BeforeEach(func() {
			gomega.Expect(page.Navigate(baseURL + "/en/testspecific")).To(gomega.Succeed())
		})
		
		ginkgo.It("should display English translations from YAML file", func() {
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Specific Method Only Test"))
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Data ID:"))
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Testing GetSpecificData without GetData"))
		})
		
		ginkgo.It("should display English data information labels", func() {
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Data Information"))
			
			gomega.Expect(html).Should(gomega.ContainSubstring("Title"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Description"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Method Called"))
		})
		
		ginkgo.It("should display English test information section", func() {
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("DataService Test Info"))
			
			gomega.Expect(html).Should(gomega.ContainSubstring("Specific Method Only Test"))
			gomega.Expect(html).Should(gomega.ContainSubstring("This template tests a DataService that ONLY implements the specific method"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Service:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Available Methods:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("GetSpecificData() ONLY"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Missing Method:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("GetData() - NOT implemented"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Expected Behavior:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Should call GetSpecificData directly"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Success Indicator:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Welcome to the Specific Method DataService demonstration!"))
		})
		
		ginkgo.It("should display English important note section", func() {
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("📝 Important Note"))
		})
	})
	
	ginkgo.Context("i18n Translation Integration - German", func() {
		ginkgo.BeforeEach(func() {
			gomega.Expect(page.Navigate(baseURL + "/de/testspecific")).To(gomega.Succeed())
		})
		
		ginkgo.It("should display German translations from YAML file", func() {
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Nur Spezifische Methode Test"))
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Daten-ID:"))
			
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Teste GetSpecificData ohne GetData"))
		})
		
		ginkgo.It("should display German data information labels", func() {
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("Dateninformationen"))
			
			gomega.Expect(html).Should(gomega.ContainSubstring("Titel"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Beschreibung"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Aufgerufene Methode"))
		})
		
		ginkgo.It("should display German test information section", func() {
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("DataService Test Info"))
			
			gomega.Expect(html).Should(gomega.ContainSubstring("Nur Spezifische Methode Test"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Dieses Template testet einen DataService, der NUR die spezifische Methode implementiert"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Service:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Verfügbare Methoden:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("NUR GetSpecificData()"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Fehlende Methode:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("GetData() - NICHT implementiert"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Erwartetes Verhalten:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Sollte GetSpecificData direkt aufrufen"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Erfolgsindikator:"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Willkommen zur Spezifische Methode DataService Demonstration!"))
		})
		
		ginkgo.It("should display German important note section", func() {
			gomega.Eventually(func() string {
				html, _ := page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("📝 Wichtiger Hinweis"))
		})
	})
	
	ginkgo.Context("Specific DataService + i18n Combined Functionality", func() {
		ginkgo.It("should show SpecificOnlyDataService data with English translations", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/testspecific")).To(gomega.Succeed())
			
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("specific-only-123"))
			
			// Both SpecificOnlyDataService data AND English translations should be present
			gomega.Expect(html).Should(gomega.ContainSubstring("specific-only-123"))  // From DataService
			gomega.Expect(html).Should(gomega.ContainSubstring("Data Information"))  // From i18n EN
			gomega.Expect(html).Should(gomega.ContainSubstring("GetSpecificData"))  // From DataService
			gomega.Expect(html).Should(gomega.ContainSubstring("Method Called"))  // From i18n EN
		})
		
		ginkgo.It("should show SpecificOnlyDataService data with German translations", func() {
			gomega.Expect(page.Navigate(baseURL + "/de/testspecific")).To(gomega.Succeed())
			
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("specific-only-123"))
			
			// Both SpecificOnlyDataService data AND German translations should be present
			gomega.Expect(html).Should(gomega.ContainSubstring("specific-only-123"))  // From DataService
			gomega.Expect(html).Should(gomega.ContainSubstring("Dateninformationen"))  // From i18n DE
			gomega.Expect(html).Should(gomega.ContainSubstring("GetSpecificData"))  // From DataService
			gomega.Expect(html).Should(gomega.ContainSubstring("Aufgerufene Methode"))  // From i18n DE
		})
	})
	
	ginkgo.Context("Language Switching with Specific DataService", func() {
		ginkgo.It("should maintain SpecificOnlyDataService data when switching languages", func() {
			// Start with English
			gomega.Expect(page.Navigate(baseURL + "/en/testspecific")).To(gomega.Succeed())
			
			var htmlEN string
			gomega.Eventually(func() string {
				htmlEN, _ = page.HTML()
				return htmlEN
			}, PageLoadTimeout).Should(gomega.ContainSubstring("specific-only-123"))
			
			// Switch to German
			gomega.Expect(page.Navigate(baseURL + "/de/testspecific")).To(gomega.Succeed())
			
			var htmlDE string
			gomega.Eventually(func() string {
				htmlDE, _ = page.HTML()
				return htmlDE
			}, PageLoadTimeout).Should(gomega.ContainSubstring("specific-only-123"))
			
			// DataService data should be the same
			gomega.Expect(htmlEN).Should(gomega.ContainSubstring("specific-only-123"))
			gomega.Expect(htmlDE).Should(gomega.ContainSubstring("specific-only-123"))
			gomega.Expect(htmlEN).Should(gomega.ContainSubstring("GetSpecificData"))
			gomega.Expect(htmlDE).Should(gomega.ContainSubstring("GetSpecificData"))
			
			// But translations should be different
			gomega.Expect(htmlEN).Should(gomega.ContainSubstring("Data Information"))
			gomega.Expect(htmlDE).Should(gomega.ContainSubstring("Dateninformationen"))
			gomega.Expect(htmlEN).ShouldNot(gomega.ContainSubstring("Dateninformationen"))
			gomega.Expect(htmlDE).ShouldNot(gomega.ContainSubstring("Data Information"))
		})
	})
	
	ginkgo.Context("Specific Method vs Fallback Method Verification", func() {
		ginkgo.It("should prove that GetSpecificData was called directly (not GetData)", func() {
			gomega.Expect(page.Navigate(baseURL + "/en/testspecific")).To(gomega.Succeed())
			
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("GetSpecificData (no GetData fallback available)"))
			
			// Verify that the specific method was called
			gomega.Expect(html).Should(gomega.ContainSubstring("This data was loaded using ONLY GetSpecificData method"))
			gomega.Expect(html).Should(gomega.ContainSubstring("router successfully called GetSpecificData without needing GetData"))
			
			// Should NOT contain any indication of GetData being called
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("GetData method called"))
			gomega.Expect(html).ShouldNot(gomega.ContainSubstring("fallback method"))
		})
		
		ginkgo.It("should show German-specific data when using German locale", func() {
			gomega.Expect(page.Navigate(baseURL + "/de/testspecific")).To(gomega.Succeed())
			
			var html string
			gomega.Eventually(func() string {
				html, _ = page.HTML()
				return html
			}, PageLoadTimeout).Should(gomega.ContainSubstring("GetSpecificData (kein GetData Fallback verfügbar)"))
			
			// Verify German-specific data from the DataService
			gomega.Expect(html).Should(gomega.ContainSubstring("Diese Daten wurden NUR über GetSpecificData Methode geladen"))
			gomega.Expect(html).Should(gomega.ContainSubstring("Dies beweist, dass der Router die spezifische Methode korrekt aufruft"))
		})
	})
})