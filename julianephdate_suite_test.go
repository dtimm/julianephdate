package julianephdate_test

import (
	"testing"
	"time"

	"github.com/dtimm/julianephdate"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJulianephdate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Julianephdate Suite")
}

var _ = Describe("julianephdate", func() {
	Describe("Date", func() {
		Context("given a stdlib time.Time", func() {
			type testEntry struct {
				Time     time.Time
				Expected float64
			}
			DescribeTable("return the float64 of the Julian ephemeris date",
				func(t time.Time, jed float64) {
					Expect(julianephdate.Date(t)).To(BeNumerically("~", 0.01, jed))
					Expect(julianephdate.StdTime(jed)).To(BeTemporally("~", t, time.Millisecond))
				},
				Entry("J2000.0", time.Date(2000, 1, 1, 11, 58, 55, 816e+6, time.UTC), 2451545.0),
			)
		})
	})
})
