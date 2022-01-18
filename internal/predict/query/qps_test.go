package query_test

import (
	"github.com/galaxy-future/cudgx/internal/predict/query"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const serviceName = "test-service"
const cluster1Name = "cluster1"
const cluster2Name = "cluster2"
const begin int64 = 1639642577
const end int64 = 1639642583

var _ = ginkgo.Describe("Qps", func() {
	ginkgo.Context("AverageQPS", func() {
		ginkgo.It("specified cluster", func() {
			samples, err := query.AverageQPS(serviceName, cluster1Name, begin, end)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(len(samples)).To(gomega.Equal(int(end - begin)))
			for i, sample := range samples {
				gomega.Expect(begin + int64(i)).To(gomega.Equal(sample.Timestamp))
				if sample.ClusterName == cluster1Name {
					gomega.Expect(sample.Value).To(gomega.Equal(float64(2)))
				} else {
					ginkgo.Fail("cluster field wrong value")
				}
			}
		})
		ginkgo.It("all cluster", func() {
			samples, err := query.AverageQPS(serviceName, "", begin, end)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(len(samples)).To(gomega.Equal(int(end-begin) * 2))
			for i, sample := range samples {
				gomega.Expect(begin + int64(i/2)).To(gomega.Equal(sample.Timestamp))
				if sample.ClusterName == cluster1Name {
					gomega.Expect(sample.Value).To(gomega.Equal(float64(2)))
				} else if sample.ClusterName == cluster2Name {
					gomega.Expect(sample.Value).To(gomega.Equal(2.5))
				} else {
					ginkgo.Fail("cluster field wrong value")
				}
			}
		})
	})

	ginkgo.Context("TotalQPS", func() {
		ginkgo.It("specified cluster", func() {
			samples, err := query.TotalQPS(serviceName, cluster1Name, begin, end)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(len(samples)).To(gomega.Equal(int(end - begin)))
			for i, sample := range samples {
				gomega.Expect(begin + int64(i)).To(gomega.Equal(sample.Timestamp))
				if sample.ClusterName == cluster1Name {
					gomega.Expect(sample.Value).To(gomega.Equal(float64(6)))
				} else {
					ginkgo.Fail("cluster field wrong value")
				}
			}
		})
		ginkgo.It("all cluster", func() {
			samples, err := query.TotalQPS(serviceName, "", begin, end)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(len(samples)).To(gomega.Equal(int(end-begin) * 2))
			for i, sample := range samples {
				gomega.Expect(begin + int64(i/2)).To(gomega.Equal(sample.Timestamp))
				if sample.ClusterName == cluster1Name {
					gomega.Expect(sample.Value).To(gomega.Equal(float64(6)))
				} else if sample.ClusterName == cluster2Name {
					gomega.Expect(sample.Value).To(gomega.Equal(float64(5)))
				} else {
					ginkgo.Fail("cluster field wrong value")
				}
			}
		})
	})

	ginkgo.Context("InstanceCountByQPS", func() {
		ginkgo.It("specified cluster", func() {
			samples, err := query.InstanceCountByQPS(serviceName, cluster1Name, begin, end)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(len(samples)).To(gomega.Equal(int(end - begin)))
			for i, sample := range samples {
				gomega.Expect(begin + int64(i)).To(gomega.Equal(sample.Timestamp))
				if sample.ClusterName == cluster1Name {
					gomega.Expect(sample.Value).To(gomega.Equal(float64(3)))
				} else {
					ginkgo.Fail("cluster field wrong value")
				}
			}
		})
		ginkgo.It("all cluster", func() {
			samples, err := query.InstanceCountByQPS(serviceName, "", begin, end)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(len(samples)).To(gomega.Equal(int(end-begin) * 2))
			for i, sample := range samples {
				gomega.Expect(begin + int64(i/2)).To(gomega.Equal(sample.Timestamp))
				if sample.ClusterName == cluster1Name {
					gomega.Expect(sample.Value).To(gomega.Equal(float64(3)))
				} else if sample.ClusterName == cluster2Name {
					gomega.Expect(sample.Value).To(gomega.Equal(float64(2)))
				} else {
					ginkgo.Fail("cluster field wrong value")
				}
			}
		})
	})
})
