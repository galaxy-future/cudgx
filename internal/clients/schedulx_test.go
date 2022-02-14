package clients_test

import (
	"github.com/galaxy-future/cudgx/internal/predict/xclient"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Xclient", func() {
	ginkgo.Context("SchedulxClient", func() {
		ginkgo.It("CanServiceSchedule", func() {
			_, err := xclient.CanServiceSchedule("gf.cudgx.pi", "gf.cudgx.pi")
			gomega.Expect(err).To(gomega.BeNil())
		})
		ginkgo.It("GetServiceInstanceCount", func() {
			count, err := xclient.GetServiceInstanceCount("gf.cudgx.pi", "gf.cudgx.pi")
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(count > 0).To(gomega.BeTrue())
		})
		ginkgo.It("ExpandService", func() {
			can, err := xclient.CanServiceSchedule("gf.cudgx.pi", "gf.cudgx.pi")
			gomega.Expect(err).To(gomega.BeNil())
			if can {
				err := xclient.ExpandService("gf.cudgx.pi", "gf.cudgx.pi", 1)
				gomega.Expect(err).To(gomega.BeNil())
			}
		})
		ginkgo.It("ShrinkService", func() {
			can, err := xclient.CanServiceSchedule("gf.cudgx.pi", "gf.cudgx.pi")
			gomega.Expect(err).To(gomega.BeNil())
			if can {
				err := xclient.ShrinkService("gf.cudgx.pi", "gf.cudgx.pi", 1)
				gomega.Expect(err).To(gomega.BeNil())
			}
		})
	})
})
