package tests

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/sclevine/agouti"
)

const (
	DefaultTimeout  = 10 * time.Second
	PageLoadTimeout = 5 * time.Second
)

var (
	agoutiDriver *agouti.WebDriver
	page         *agouti.Page
	baseURL      string
)

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Templ Router Demo E2E Test Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	// Setup base URL
	baseURL = os.Getenv("TEST_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8084"
	}

	// Wait for service to be ready
	ginkgo.By("Waiting for Docker service to be ready")
	gomega.Eventually(func() error {
		resp, err := http.Get(baseURL + "/api/health")
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("service not ready, status: %d", resp.StatusCode)
		}
		return nil
	}, 30*time.Second, 2*time.Second).Should(gomega.Succeed())

	// Setup Chrome driver
	agoutiDriver = agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--disable-gpu",
			"--window-size=1920,1080",
		}),
	)

	gomega.Expect(agoutiDriver.Start()).To(gomega.Succeed())
})

var _ = ginkgo.AfterSuite(func() {
	if agoutiDriver != nil {
		gomega.Expect(agoutiDriver.Stop()).To(gomega.Succeed())
	}
})

var _ = ginkgo.BeforeEach(func() {
	var err error
	page, err = agoutiDriver.NewPage()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(page.SetImplicitWait(int(DefaultTimeout.Milliseconds()))).To(gomega.Succeed())
})

var _ = ginkgo.AfterEach(func() {
	if page != nil {
		page.Destroy()
	}
})