package blockcode

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Blockcode", func() {

	Describe(".validateManifest", func() {
		It("should validate correctly", func() {
			cases := map[*Manifest]error{
				&Manifest{}:                                                             fmt.Errorf("manifest error: language is missing"),
				&Manifest{Lang: "c++"}:                                                  fmt.Errorf("manifest error: language {c++} is not supported"),
				&Manifest{Lang: "go", LangVersion: ""}:                                  fmt.Errorf("manifest error: language version is required"),
				&Manifest{Lang: "go", LangVersion: "1.10.2"}:                            fmt.Errorf("manifest error: at least one public function is required"),
				&Manifest{Lang: "go", LangVersion: "1.10.2", PublicFuncs: []string{""}}: fmt.Errorf("manifest error: at least one public function is required"),
			}

			for manifest, err := range cases {
				Expect(validateManifest(manifest)).To(Equal(err))
			}

			validManifest := &Manifest{Lang: "go", LangVersion: "1.10.2", PublicFuncs: []string{"some_func"}}
			Expect(validateManifest(validManifest)).To(BeNil())

		})
	})

	Describe(".FromDir", func() {

		It("should return error = 'project path does not exist' if project path is not found", func() {
			_, err := FromDir("./testdata/unknown")
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("project path does not exist"))
		})

		It("should return error = 'manifest is malformed. invalid character ':' after top-level value' if manifest is not valid JSON", func() {
			_, err := FromDir("./testdata/invalid_manifest")
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("manifest is malformed. invalid character ':' after top-level value"))
		})

		It("should return error = ''package.json' file not found in {./testdata/missing_manifest}' if package.json is missing", func() {
			_, err := FromDir("./testdata/missing_manifest")
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("'package.json' file not found in {./testdata/missing_manifest}"))
		})

		It("should successfully create a Blockcode", func() {
			bc, err := FromDir("./testdata/blockcode_example")
			Expect(err).To(BeNil())
			Expect(bc.code).ToNot(BeEmpty())
			Expect(bc.Manifest).ToNot(BeNil())
			Expect(bc.Manifest.Lang).ToNot(BeEmpty())
			Expect(bc.Manifest.LangVersion).ToNot(BeEmpty())
			Expect(bc.Manifest.PublicFuncs).ToNot(BeEmpty())
		})
	})

	Describe(".Len", func() {
		It("should return 3072", func() {
			bc, err := FromDir("./testdata/blockcode_example")
			Expect(err).To(BeNil())
			Expect(bc.Len()).To(Equal(3072))
		})
	})

	Describe(".Bytes", func() {
		It("should return bytes", func() {
			bc, err := FromDir("./testdata/blockcode_example")
			Expect(err).To(BeNil())
			Expect(bc.Bytes()).ToNot(BeEmpty())
		})
	})

	Describe(".Hash", func() {
		It("should return Hash", func() {
			bc, err := FromDir("./testdata/blockcode_example")
			Expect(err).To(BeNil())
			Expect(bc.Hash()).To(Equal([]byte{148, 25, 211, 248, 52, 247, 87, 232, 85, 122, 54, 128, 104, 171, 53, 68, 243, 150, 28, 66, 217, 58, 223, 47, 111, 118, 36, 95, 219, 145, 83, 185}))
		})
	})

	Describe(".ID", func() {
		It("should return ID", func() {
			bc, err := FromDir("./testdata/blockcode_example")
			Expect(err).To(BeNil())
			Expect(bc.ID()).To(Equal("9419d3f834f757e8557a368068ab3544f3961c42d93adf2f6f76245fdb9153b9"))
		})
	})

	Describe(".FromBytes", func() {
		It("should return ID", func() {
			bc, err := FromDir("./testdata/blockcode_example")
			Expect(err).To(BeNil())
			bs := bc.Bytes()
			blockcode, err := FromBytes(bs)
			Expect(err).To(BeNil())
			Expect(blockcode).To(Equal(bc))
		})
	})

	Describe(".Read", func() {
		It("should return err = 'destination path does not exist' if destination path does not exist", func() {
			bc, err := FromDir("./testdata/blockcode_example")
			Expect(err).To(BeNil())
			err = bc.Read("./unknown/path")
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("destination path does not exist"))
		})

		It("should successfully un-tar to destination", func() {

			destination := "/tmp/blockcode_example_untar"
			err := os.Mkdir(destination, 0700)
			Expect(err).To(BeNil())
			defer os.RemoveAll(destination)

			bc, err := FromDir("./testdata/blockcode_example")
			Expect(err).To(BeNil())

			err = bc.Read(destination)
			Expect(err).To(BeNil())
		})
	})
})