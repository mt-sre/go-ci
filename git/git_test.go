package git

import (
	"context"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _temp string

func TestSuite(t *testing.T) {
	_temp = t.TempDir()

	RegisterFailHandler(Fail)
	RunSpecs(t, "git suite")
}

var _ = BeforeSuite(func() {
	var err error

	cmd := exec.Command("./setup_test.sh", _temp)

	sess, err := Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())

	Eventually(sess).Should(Exit(0))

})

type revParseTestCase struct {
	Options  []RevParseOption
	Expected string
}

var _ = Describe("RevParse", func() {
	temp := _temp

	DescribeTable("Formats",
		func(tc revParseTestCase) {
			ctx := context.Background()

			opts := append(tc.Options, WithWorkingDirectory(temp))

			res, err := RevParse(ctx, opts...)
			Expect(err).ToNot(HaveOccurred())

			Expect(res).To(Equal(tc.Expected))
		},
		Entry("top-level",
			revParseTestCase{
				Options: []RevParseOption{
					WithRevParseFormat(RevParseFormatTopLevel),
				},
				Expected: temp,
			},
		),
		Entry("abbrev-ref",
			revParseTestCase{
				Options: []RevParseOption{
					WithRevParseFormat(RevParseFormatAbbrevRef),
				},
				Expected: "test",
			},
		),
	)
})

var _ = Describe("ListTags", func() {
	It("should list tags", func() {
		ctx := context.Background()

		res, err := ListTags(ctx, WithWorkingDirectory(_temp))
		Expect(err).ToNot(HaveOccurred())

		Expect(res).To(ContainElements([]string{"v1.0.0", "v2.0.0"}))
	})
})

var _ = Describe("LatestTag", func() {
	It("should return the latest tag", func() {
		ctx := context.Background()

		res, err := LatestTag(ctx, WithWorkingDirectory(_temp))
		Expect(err).ToNot(HaveOccurred())

		Expect(res).To(Equal("v2.0.0"))
	})
})

var _ = Describe("LatestVersion", func() {
	It("should return the latest version", func() {
		ctx := context.Background()

		res, err := LatestVersion(ctx, WithWorkingDirectory(_temp))
		Expect(err).ToNot(HaveOccurred())

		Expect(res).To(Equal("v2.0.0"))
	})
})

type statusTestCase struct {
	Options  []StatusOption
	Expected string
}

var _ = DescribeTable("Status",
	func(tc statusTestCase) {
		ctx := context.Background()

		opts := append(tc.Options, WithWorkingDirectory(_temp))

		res, err := Status(ctx, opts...)
		Expect(err).ToNot(HaveOccurred())

		Expect(res).To(Equal(tc.Expected))
	},
	Entry("no format",
		statusTestCase{
			Options:  []StatusOption{},
			Expected: "On branch test\nnothing to commit, working tree clean",
		},
	),
	Entry("porcelain",
		statusTestCase{
			Options: []StatusOption{
				WithStatusFormat(StatusFormatPorcelain),
			},
			Expected: "",
		},
	),
	Entry("long",
		statusTestCase{
			Options: []StatusOption{
				WithStatusFormat(StatusFormatLong),
			},
			Expected: "On branch test\nnothing to commit, working tree clean",
		},
	),
	Entry("short",
		statusTestCase{
			Options: []StatusOption{
				WithStatusFormat(StatusFormatShort),
			},
			Expected: "",
		},
	),
)

type diffTestCase struct {
	Options  []DiffOption
	Expected string
}

var _ = DescribeTable("Diff",
	func(tc diffTestCase) {
		ctx := context.Background()

		opts := append(tc.Options, WithWorkingDirectory(_temp))

		res, err := Diff(ctx, opts...)
		Expect(err).ToNot(HaveOccurred())

		Expect(res).To(Equal(tc.Expected))
	},
	Entry("no format",
		diffTestCase{
			Options: []DiffOption{
				WithDiffFormat(DiffFormatNameOnly),
			},
			Expected: "",
		},
	),
	Entry("name-only",
		diffTestCase{
			Options: []DiffOption{
				WithDiffFormat(DiffFormatNameOnly),
			},
			Expected: "",
		},
	),
	Entry("name-status",
		diffTestCase{
			Options: []DiffOption{
				WithDiffFormat(DiffFormatNameStatus),
			},
			Expected: "",
		},
	),
)
