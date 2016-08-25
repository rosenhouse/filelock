package filelock_test

import (
	"github.com/rosenhouse/filelock"

	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Locking using a file", func() {
	var path string
	var tempDir string

	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		path = filepath.Join(tempDir, "dir1", "dir2", "some-file.json")
	})

	AfterEach(func() {
		Expect(os.RemoveAll(tempDir)).To(Succeed())
	})

	AssertBasicThingsWork := func(expectedInitialContents []byte) {
		It("returns a file usable for read and write", func() {
			lock := filelock.Locker{
				Path: path,
			}

			file, err := lock.Open()
			Expect(err).NotTo(HaveOccurred())

			initialContents, err := ioutil.ReadAll(file)
			Expect(err).NotTo(HaveOccurred())

			Expect(initialContents).To(Equal(expectedInitialContents))

			_, err = file.Seek(0, 0)
			Expect(err).NotTo(HaveOccurred())

			_, err = file.Write([]byte("hello"))
			Expect(err).NotTo(HaveOccurred())

			_, err = file.Seek(0, 0)
			Expect(err).NotTo(HaveOccurred())

			allBytes, err := ioutil.ReadAll(file)
			Expect(err).NotTo(HaveOccurred())

			Expect(allBytes).To(Equal([]byte("hello")))

			Expect(file.Close()).To(Succeed())

			allBytes, err = ioutil.ReadFile(path)
			Expect(err).NotTo(HaveOccurred())

			Expect(allBytes).To(Equal([]byte("hello")))
		})
	}

	Context("when the file does not already exist", func() {
		BeforeEach(func() {
			Expect(os.RemoveAll(tempDir)).To(Succeed())
		})
		AssertBasicThingsWork([]byte{})
	})

	Context("when the file already exists", func() {
		var preExistingContents = []byte("foo")

		BeforeEach(func() {
			Expect(os.MkdirAll(filepath.Dir(path), 0700)).To(Succeed())
			Expect(ioutil.WriteFile(path, preExistingContents, 0600)).To(Succeed())
		})
		AssertBasicThingsWork(preExistingContents)
	})

	Context("when the path has already been opened by a locker", func() {
		var otherFileDescriptor filelock.LockedFile

		BeforeEach(func() {
			var err error
			otherFileDescriptor, err = (&filelock.Locker{Path: path}).Open()
			Expect(err).NotTo(HaveOccurred())
		})

		It("blocks the second open until the first one is closed", func(done Done) {
			locker := &filelock.Locker{Path: path}

			openSucceeded := make(chan struct{})
			go func() {
				defer GinkgoRecover()

				_, err := locker.Open()
				Expect(err).NotTo(HaveOccurred())

				close(openSucceeded)
			}()
			Consistently(openSucceeded).ShouldNot(BeClosed())

			Expect(otherFileDescriptor.Close()).To(Succeed())
			Eventually(openSucceeded).Should(BeClosed())

			close(done)
		}, 5 /* max seconds allowed for this spec */)

	})

})
